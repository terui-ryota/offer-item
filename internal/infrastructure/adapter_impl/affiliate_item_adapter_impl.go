package adapter_impl

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/dgraph-io/ristretto"
	"github.com/terui-ryota/offer-item/internal/domain/model"
	"github.com/terui-ryota/offer-item/internal/infrastructure/component/rakuten"
	"github.com/terui-ryota/offer-item/pkg/logger"
	"golang.org/x/sync/errgroup"
	"golang.org/x/sync/singleflight"

	"github.com/terui-ryota/offer-item/internal/domain/adapter"
	"github.com/terui-ryota/offer-item/pkg/apperr"
	"go.opencensus.io/trace"
	"go.uber.org/zap"
)

func NewAffiliateItemAdapterImpl(rakutenClient *rakuten.RakutenIchibaClient) adapter.AffiliateItemAdapter {
	// コネクションを生成する
	//conn, err := grpc.Dial(
	//config.AffiliateItem.Target(),
	//	grpc.WithTransportCredentials(insecure.NewCredentials()),
	//	grpc.WithUnaryInterceptor(
	//		grpc_middleware.UnaryClientInterceptor(consts.ContextName),
	//	),
	//	grpc.WithStatsHandler(&ocgrpc.ClientHandler{
	//		StartOptions: trace.StartOptions{
	//			Sampler: trace.ProbabilitySampler(0.5),
	//		},
	//	}),
	//)
	//if err != nil {
	//	// 接続できない場合
	//	panic(err)
	//}
	// クライアントを生成する
	//affiliateItemAdminClient := affiliate_item.NewAffiliateItemAdminHandlerClient(conn)
	//affiliateItemClient := affiliate_item.NewAffiliateItemHandlerClient(conn)
	//
	//cache, err := ristretto.NewCache(&ristretto.Config{
	//	NumCounters: 10000,
	//	MaxCost:     10 * 1024 * 1024, // 10MB
	//	BufferItems: 64,
	//})
	//if err != nil {
	//	panic(err)
	//}

	cache, err := ristretto.NewCache(&ristretto.Config{
		NumCounters: 10000,
		MaxCost:     10 * 1024 * 1024, // 10MB
		BufferItems: 64,
	})
	if err != nil {
		panic(err)
	}

	return &AffiliateItemAdapterImpl{
		rakutenClient: rakutenClient,
		localCache:    cache,
	}
}

type AffiliateItemAdapterImpl struct {
	rakutenClient *rakuten.RakutenIchibaClient
	sg            singleflight.Group
	localCache    *ristretto.Cache
}

// 案件ID、DF案件IDを指定して、案件情報を取得する
func (a *AffiliateItemAdapterImpl) GetItems(ctx context.Context, itemIdentifier model.ItemIdentifier) (*model.Items, error) {
	ctx, span := trace.StartSpan(ctx, "AffiliateItemAdapterImpl.GetItems")
	defer span.End()

	itemMap, err := a.bulkGetItems(ctx, []model.ItemIdentifier{itemIdentifier}, true)
	if err != nil {
		return nil, fmt.Errorf("a.bulkGetItems: %w", err)
	}
	items, ok := itemMap[itemIdentifier]
	if !ok {
		return nil, apperr.OfferItemAffiliateItemNotFoundError.Wrap(fmt.Errorf("item not found: %v", itemIdentifier))
	}

	return &items, nil
}

// itemIdentifiersを指定して、案件情報マップを取得する
func (a *AffiliateItemAdapterImpl) BulkGetItems(ctx context.Context, itemIdentifiers model.ItemIdentifiers) (map[model.ItemIdentifier]model.Items, error) {
	ctx, span := trace.StartSpan(ctx, "AffiliateItemAdapterImpl.BulkGetItems")
	defer span.End()

	return a.bulkGetItems(ctx, itemIdentifiers, true)
}

func (a *AffiliateItemAdapterImpl) getAffiliateItems(ctx context.Context, pairs model.ItemIdentifiers) (model.AffiliateItemPairList, error) {
	mutex := &sync.Mutex{}
	// Context は評価後に Cancel されるため別名で定義する
	eg, egCtx := errgroup.WithContext(ctx)

	var ids []string
	for _, v := range pairs {
		ids = append(ids, v.DFItemID().String())
	}

	var affiliateItemDFItemList model.AffiliateItemDFItemList
	//var err error

	uncachedList := make(model.AffiliateItemDFItemList, 0, len(ids))
	for _, itemPair := range pairs {
		id := itemPair.DFItemID().String()
		eg.Go(func() error {
			v, err, _ := a.sg.Do(id, func() (interface{}, error) {
				return a.rakutenClient.GetItemsByItemId(egCtx, id)
			})
			if err != nil {
				switch err := err.(type) {
				case rakuten.TooManyRequestsErr:
					// 楽天API でリクエスト数超過が出た場合は400系で返却するためにTooManyRequestErrorで返却する
					return rakuten.NewASPTooManyRequestError(itemPair.ItemID())
				default:
					// 正常に取得できた item は返却する必要があるので、リクエスト数超過以外のエラーはここでエラーログで出力し関数自体は正常終了させる。
					logger.Default().Error("Failed to get rakuten items by ids", zap.String("id", id), zap.Error(err))
					return nil
				}
			}
			result, ok := v.(*rakuten.ItemResult)
			if !ok {
				return apperr.AffiliateItemInternalError.Wrap(errors.New("failed convert to rakuten.ItemResult"))
			}
			// model変換
			affiliateItemDFItemList = a.convertAffiliateItemDFItems(egCtx, result, itemPair.ItemID())
			mutex.Lock()
			defer mutex.Unlock()
			uncachedList = append(uncachedList, affiliateItemDFItemList...)
			return nil
		})
	}
	if err := eg.Wait(); err != nil {
		return nil, err
	}

	var affiliateItemPairList model.AffiliateItemPairList
	affiliateItemRate := float32(4)
	affiliateItemStaticRate := float32(4)
	affiliateItemCommission := model.NewAffiliateItemCommission(&affiliateItemRate, nil, &affiliateItemStaticRate, "")
	affiliateItemItem := model.NewAffiliateItemItem("RK000001", "https://stat.amebame.com/pub/content/5164757434/amebapick/item/rakuten/logo.png", "楽天市場", affiliateItemCommission, "楽天株式会社", model.URLMap{model.PlatformTypeAll: "https://www.rakuten.co.jp/"}, true)
	for _, affiliateItemDFItem := range uncachedList {
		affiliateItemPair := model.NewAffiliateItemPair(affiliateItemItem, affiliateItemDFItem)
		affiliateItemPairList = append(affiliateItemPairList, affiliateItemPair)
	}
	return affiliateItemPairList, nil
}

func (a *AffiliateItemAdapterImpl) bulkGetItems(ctx context.Context, pairs model.ItemIdentifiers, useCache bool) (map[model.ItemIdentifier]model.Items, error) {
	res := make(map[model.ItemIdentifier]model.Items)
	afRequestPairs := make(model.ItemIdentifiers, 0, len(pairs))

	// cacheから取得できるアイテムと取得できないアイテムを分ける
	for _, pair := range pairs {
		if useCache {
			if items, found := a.getListItemsFromLocalCache(ctx, pair); found {
				res[pair] = *items
				continue
			}
		}
		// cacheから取得できないアイテムはaffiliateItemから取得する為、afRequestPairsに追加
		afRequestPairs = append(afRequestPairs, pair)
	}

	logger.FromContext(ctx).Debugf("Number of cache hits : %d, Number of cache misses: %d", len(res), len(afRequestPairs))
	if len(afRequestPairs) == 0 {
		return res, nil
	}

	lir, err := a.getAffiliateItems(ctx, afRequestPairs)
	if err != nil {
		return nil, fmt.Errorf("a.getAffiliateItems: %w", err)
	}

	if len(lir) == 0 {
		return nil, apperr.OfferItemAffiliateItemNotFoundError.Wrap(fmt.Errorf("item not found: %v", afRequestPairs))
	}

	items := make(map[string]model.Item)
	dfItems := make(map[string]model.DFItem)

	for i := range lir {
		if affiliateItemItem := lir[i].AffiliateItemItem(); affiliateItemItem != nil {
			item, err := model.AffiliateItemItemToItem(affiliateItemItem)
			if err != nil {
				return nil, fmt.Errorf("converter.AffiliateItemItemToItem: %w", err)
			}
			items[item.ID().String()] = *item
		}
		if affiliateItemDFItem := lir[i].AffiliateItemDFItem(); affiliateItemDFItem != nil {
			dfItem, err := model.AffiliateItemDFItemToDFItem(affiliateItemDFItem)
			if err != nil {
				return nil, fmt.Errorf("converter.AffiliateItemDFItemToDFItem: %w", err)
			}
			dfItems[dfItem.ID().String()] = *dfItem
		}
	}

	for i := range afRequestPairs {
		_, isItemExist := items[afRequestPairs[i].ItemID().String()]
		_, isDFItemExist := dfItems[afRequestPairs[i].DFItemID().String()]

		// DfItemが存在するが、DFItemを取得出来ない場合はログに出力して次の処理へいく
		if afRequestPairs[i].DFItemID().String() != "" && !isDFItemExist {
			logger.FromContext(ctx).Warn("DFItem not found", zap.String("DFItemID", afRequestPairs[i].DFItemID().String()))
			continue
		}

		// ItemIDとDFItemIDの両方が非空であり、それぞれのアイテムがitemsとdfItemsに存在する場合
		if afRequestPairs[i].ItemID().String() != "" && afRequestPairs[i].DFItemID().String() != "" && isItemExist && isDFItemExist {
			res[afRequestPairs[i]] = model.Items{
				Item:   items[afRequestPairs[i].ItemID().String()],
				DFItem: dfItems[afRequestPairs[i].DFItemID().String()],
			}
			// キャッシュに保存する
			a.saveListItemsToLocalCache(afRequestPairs[i], res[afRequestPairs[i]])
			continue
		}

		if afRequestPairs[i].ItemID().String() != "" && isItemExist {
			res[afRequestPairs[i]] = model.Items{
				Item: items[afRequestPairs[i].ItemID().String()],
			}
			a.saveListItemsToLocalCache(afRequestPairs[i], res[afRequestPairs[i]])
			continue
		}
	}

	return res, nil
}

//func (a *AffiliateItemAdapterImpl) bulkGetItems(ctx context.Context, pairs model.ItemIdentifiers, useCache bool) (map[model.ItemIdentifier]model.Items, error) {
//	res := make(map[model.ItemIdentifier]model.Items)
//
//	lir, err := a.getAffiliateItems(ctx, pairs)
//	if err != nil {
//		return nil, fmt.Errorf("a.getAffiliateItems: %w", err)
//	}
//
//	if len(lir) == 0 {
//		return nil, apperr.OfferItemAffiliateItemNotFoundError.Wrap(fmt.Errorf("item not found: %v", pairs))
//	}
//
//	items := make(map[string]model.Item)
//	dfItems := make(map[string]model.DFItem)
//
//	for i := range lir {
//		if affiliateItemItem := lir[i].AffiliateItemItem(); affiliateItemItem != nil {
//			item, err := model.AffiliateItemItemToItem(affiliateItemItem)
//			if err != nil {
//				return nil, fmt.Errorf("converter.AffiliateItemItemToItem: %w", err)
//			}
//			items[item.ID().String()] = *item
//		}
//		if affiliateItemDFItem := lir[i].AffiliateItemDFItem(); affiliateItemDFItem != nil {
//			dfItem, err := model.AffiliateItemDFItemToDFItem(affiliateItemDFItem)
//			if err != nil {
//				return nil, fmt.Errorf("converter.AffiliateItemDFItemToDFItem: %w", err)
//			}
//			dfItems[dfItem.ID().String()] = *dfItem
//		}
//	}
//
//	for i := range pairs {
//		_, isItemExist := items[pairs[i].ItemID().String()]
//		_, isDFItemExist := dfItems[pairs[i].DFItemID().String()]
//		// DfItemが存在するが、DFItemを取得出来ない場合はログに出力して次の処理へいく
//		if pairs[i].DFItemID().String() != "" && !isDFItemExist {
//			logger.FromContext(ctx).Warn("DFItem not found", zap.String("DFItemID", pairs[i].DFItemID().String()))
//			continue
//		}
//		// ItemIDとDFItemIDの両方が非空であり、それぞれのアイテムがitemsとdfItemsに存在する場合
//		if pairs[i].ItemID().String() != "" && pairs[i].DFItemID().String() != "" && isItemExist && isDFItemExist {
//			res[pairs[i]] = model.Items{
//				Item:   items[pairs[i].ItemID().String()],
//				DFItem: dfItems[pairs[i].DFItemID().String()],
//			}
//			continue
//		}
//
//		if pairs[i].ItemID().String() != "" && isItemExist {
//			res[pairs[i]] = model.Items{
//				Item: items[pairs[i].ItemID().String()],
//			}
//			continue
//		}
//	}
//
//	return res, nil
//}

func (a *AffiliateItemAdapterImpl) saveListItemsToLocalCache(itemIdentifier model.ItemIdentifier, items model.Items) {
	ttl := 24 * time.Hour
	key := itemIdentifier.ItemID().String()
	if itemIdentifier.DFItemID().String() != "" {
		key = key + "#" + itemIdentifier.DFItemID().String()
	}
	a.localCache.SetWithTTL(key, items, 1, ttl)
}

func (r *AffiliateItemAdapterImpl) getListItemsFromLocalCache(ctx context.Context, itemIdentifier model.ItemIdentifier) (*model.Items, bool) {
	key := itemIdentifier.ItemID().String()
	if itemIdentifier.DFItemID().String() != "" {
		key = key + "#" + itemIdentifier.DFItemID().String()
	}
	if value, isFound := r.localCache.Get(key); isFound {
		if items, ok := value.(model.Items); ok {
			return &items, true
		} else {
			logger.FromContext(ctx).Warnf("Failed to parse cached. key: %s, value: %d", key, items)
		}
	}

	return nil, false
}

// convertAffiliateItemDFItems はrakuten商品検索APIの検索結果をAffiliateItemDFItemModelに変換します。
func (r *AffiliateItemAdapterImpl) convertAffiliateItemDFItems(ctx context.Context, result *rakuten.ItemResult, itemId model.ItemID) model.AffiliateItemDFItemList {
	// 楽天のDF案件
	affiliateItemDFItemList := make(model.AffiliateItemDFItemList, 0, len(result.Items))
	for _, resultItem := range result.Items {
		affiliateRate := float32(resultItem.Item.AffiliateRate)
		commissionRate := model.NewAffiliateItemCommission(
			&affiliateRate,
			&affiliateRate,
			&affiliateRate,
			"",
		)

		// 商品URL
		itemUrl, err := resultItem.Item.ParsedItemUrl()
		if err != nil {
			logger.Default().Warn("Failed to parse url.", zap.String("url", resultItem.Item.AffiliateURL))
		}

		// 金額
		price := model.AffiliateItemDFItemPriceList{
			model.AffiliateItemDFItemPrice{
				SalePrice:   model.Price(resultItem.Item.ItemPrice),
				RetailPrice: model.Price(resultItem.Item.ItemPrice),
			},
		}
		affiliateItemDFItem := model.NewAffiliateItemDFItem(
			resultItem.Item.ItemCode,
			resultItem.Item.ItemName, // CatchcopyはCarrierによって変わる（PC、スマホ、などなど）
			itemUrl,
			resultItem.MakeMaterial(),
			resultItem.Item.ShopName,
			resultItem.Item.ShopCode,
			itemId,
			price,
			model.AffiliateItemCommissionList{commissionRate},
		)

		affiliateItemDFItemList.Append(affiliateItemDFItem)
	}
	return affiliateItemDFItemList
}
