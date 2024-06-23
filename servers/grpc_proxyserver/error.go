package grpc_proxyserver

import (
	"context"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/status"

	"github.com/terui-ryota/offer-item/pkg/apperr"
	commonpb "github.com/terui-ryota/protofiles/go/common"
)

// HttpErrorHandler grpc server側でエラーが発生した時用のハンドラー
// 現状不要なので、エラー時にはgrpcのheader、trailerはここで握り潰しています。
func HttpErrorHandler(_ context.Context, _ *runtime.ServeMux, marshaler runtime.Marshaler, w http.ResponseWriter, _ *http.Request, err error) {
	s := status.Convert(err)
	w.Header().Del("Trailer")
	w.Header().Set("Content-Type", marshaler.ContentType(s))

	ae := apperr.UnknownError.Proto()
	for _, d := range s.Details() {
		if t, ok := d.(*commonpb.ApplicationError); ok {
			ae = t
			break
		}
	}
	buf, merr := marshaler.Marshal(ae)
	if merr != nil {
		grpclog.Infof("Failed to marshal error message %q: %v", ae, merr)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	st := runtime.HTTPStatusFromCode(s.Code())
	w.WriteHeader(st)
	if _, err := w.Write(buf); err != nil {
		grpclog.Infof("Failed to write response: %v", err)
	}
}
