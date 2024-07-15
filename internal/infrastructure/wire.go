package infrastructure

import (
	"github.com/google/wire"
	"github.com/terui-ryota/offer-item/internal/infrastructure/component/rakuten"

	"github.com/terui-ryota/offer-item/internal/infrastructure/adapter_impl"
	"github.com/terui-ryota/offer-item/internal/infrastructure/repository_impl"
)

var WireSet = wire.NewSet(
	repository_impl.NewOfferItemRepositoryImpl,
	repository_impl.NewAssigneeRepositoryImpl,
	repository_impl.NewExaminationRepositoryImpl,
	repository_impl.NewQuestionnaireRepositoryImpl,
	repository_impl.NewQuestionnaireQuestionAnswerRepositoryImpl,
	adapter_impl.NewAffiliateItemAdapterImpl,
	rakuten.NewRakutenIchibaClient,
	rakuten.NewApplicationIDHelper,
)
