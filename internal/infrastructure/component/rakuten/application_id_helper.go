package rakuten

import (
	"context"
	"math/rand"
	"time"

	"github.com/terui-ryota/offer-item/internal/app/grpcserver/config"
	"go.opencensus.io/trace"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

type ApplicationIDHelper struct {
	applicationIDs []string
}

func NewApplicationIDHelper(conf *config.RakutenConfig) *ApplicationIDHelper {
	return &ApplicationIDHelper{
		applicationIDs: conf.ApplicationID,
	}
}

func (h *ApplicationIDHelper) GetApplicationID(ctx context.Context) string {
	ctx, span := trace.StartSpan(ctx, "ApplicationIDHelper#GetApplicationID")
	defer span.End()

	applicationID := h.applicationIDs[rand.Intn(len(h.applicationIDs))]

	return applicationID
}
