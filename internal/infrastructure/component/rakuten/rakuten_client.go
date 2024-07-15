package rakuten

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"

	//"net/url"
	//"time"

	"github.com/pkg/errors"
	"github.com/terui-ryota/offer-item/internal/app/grpcserver/config"
	"github.com/terui-ryota/offer-item/pkg/logger"
	"go.opencensus.io/trace"
	"go.uber.org/zap"
)

func NewRakutenIchibaClient(
	config *config.RakutenConfig,
	client *http.Client,
	helper *ApplicationIDHelper,
) *RakutenIchibaClient {
	return &RakutenIchibaClient{
		client: client,
		config: *config,
		helper: helper,
	}
}

// RakutenIchibaClient は楽天商品検索API用Clientです。
// URL:https://webservice.rakuten.co.jp/api/ichibaitemsearch/
type RakutenIchibaClient struct {
	client *http.Client
	config config.RakutenConfig
	helper *ApplicationIDHelper
}

type TooManyRequestsErr struct{}

func (TooManyRequestsErr) Error() string {
	return ""
}

const (
	rakutenIchibaItemURL = "https://app.rakuten.co.jp/services/api/IchibaItem/Search/20170706"
)

// GetItemsByItemId は指定の商品コードの商品情報を取得します。
func (r *RakutenIchibaClient) GetItemsByItemId(ctx context.Context, itemCode string) (*ItemResult, error) {
	ctx, span := trace.StartSpan(ctx, "RakutenIchibaClient#GetItemsByItemId")
	defer span.End()

	applicationId := r.helper.GetApplicationID(ctx)
	// パラメーター生成
	param := NewItemSearchByItemCodeParam(
		applicationId,
		r.config.RakutenIchiba.Format,
		itemCode,
	)

	// リクエスト生成
	req, err := r.createRequest(ctx, rakutenIchibaItemURL, param)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Failed to generate request rakuten url:%s param:%+v", rakutenIchibaItemURL, param))
	}
	logger.Default().Debug("Rakuten request.", zap.String("url", req.URL.String()))

	// リクエスト実行
	resp, err := r.client.Do(req)
	if err != nil {
		return nil, errors.Wrapf(err, "Failed to Rakuten http request")
	}
	defer func() {
		_, _ = io.Copy(ioutil.Discard, resp.Body)
		_ = resp.Body.Close()
	}()

	var reader io.Reader = resp.Body
	// debug用
	// r = io.TeeReader(r, os.Stderr)

	if resp.StatusCode >= 400 {
		// エラーメッセージをデコード
		var rakutenError RakutenError
		if err := json.NewDecoder(reader).Decode(&rakutenError); err != nil {
			b, _ := ioutil.ReadAll(reader)
			logger.Default().Warn("Failed to error decode.", zap.String("data", string(b)))
			return nil, errors.Wrap(err, fmt.Sprintf("Failed to decode rakuten error. url:%s param:%+v", rakutenIchibaItemURL, param))
		}
		logger.Default().Warn("Failed to rakuten request.", zap.Reflect("error", rakutenError))

		if rakutenError.IsInValidItemCode() {
			// 存在しないItemIDの場合は空で返却する。
			return &ItemResult{}, nil
		} else if rakutenError.IsTooManyRequest() {
			// 429が返ってきている場合はInternalServerErrorではなく400系で返却するためTooManyRequestsErrで返却する
			return nil, TooManyRequestsErr{}
		}

		return nil, errors.New(fmt.Sprintf("Failed to rakuten request. %+v", rakutenError))
	}

	// レスポンスをデコード
	var data ItemResult
	if err := json.NewDecoder(reader).Decode(&data); err != nil {
		return nil, errors.Wrapf(err, "Failed to decode rakuten api. url:%s param:%+v", rakutenIchibaItemURL, param)
	}

	return &data, nil
}

// doRequest は楽天商品検索API用のリクエストURLを生成します。
func (r *RakutenIchibaClient) createRequest(ctx context.Context, targetUrl string, param RakutenSearchParam) (*http.Request, error) {
	u, err := url.Parse(targetUrl)
	if err != nil {
		return nil, fmt.Errorf("url.Parse: %w", err)
	}
	req, err := NewRakutenRequest(
		targetUrl,
		param.UrlValues(),
	).GetRequest()
	if err != nil {
		return nil, fmt.Errorf("NewRakutenRequest: %w", err)
	}
	// targetURL にパス変数が含まれないため Path をそのまま設定します
	// パス (URL) を設定する箇所とヘッダーを設定する場所を分離しないように修正してください
	req.Header.Add("X-Path-Pattern", u.Path)
	return req, nil
}
