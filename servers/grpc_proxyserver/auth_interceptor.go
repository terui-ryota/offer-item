package grpc_proxyserver

type CheckPermissionFunc = func(fullMethodName string, permissions []string) bool

const (
	keyPermissions  = "x-odessa-permission"
	keyMediaType    = "x-odessa-media-type"
	keyMediaUserID  = "x-odessa-media-user-id"
	keyApplicantID  = "x-odessa-applicant-id"
	keyAffiliatorID = "x-odessa-affiliator-id"

	fieldAffiliatorID = "AffiliatorId"
	fieldApplicantID  = "ApplicantId"
	fieldMediaType    = "MediaType"
	fieldMediaUserId  = "MediaUserId"
)

//func authUnaryServerInterceptor(checkPermissionFuncs ...func(fullMethodName string, permissions []string) bool) grpc.UnaryClientInterceptor {
//	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
//		md, ok := metadata.FromOutgoingContext(ctx)
//		if !ok {
//			md = metadata.Pairs()
//		}
//		var permissions []string
//		// 認可を行う
//		if v, ok := getMetadataValue(md, keyPermissions); ok {
//			permissions = strings.Split(v, ",")
//		}
//		hasAuthenticated := false
//		for _, f := range checkPermissionFuncs {
//			if f(method, permissions) {
//				hasAuthenticated = true
//				break
//			}
//		}
//		if !hasAuthenticated {
//			return apperr.PermissionDenied
//		}
//		// ID変換
//		rv := reflect.ValueOf(req)
//		if rv.Kind() == reflect.Ptr {
//			rv = rv.Elem()
//		}
//		// アフェリエイターID変換
//		if v, ok := getMetadataValue(md, keyAffiliatorID); ok {
//			if f := rv.FieldByName(fieldAffiliatorID); f.IsValid() && f.Kind() == reflect.String {
//				f.SetString(v)
//			}
//		}
//
//		// アフェリエイターID変換
//		if v, ok := getMetadataValue(md, keyApplicantID); ok {
//			if f := rv.FieldByName(fieldApplicantID); f.IsValid() && f.Kind() == reflect.String {
//				f.SetString(v)
//			}
//		}
//		// メディアタイプ変換
//		if v, ok := getMetadataValue(md, keyMediaType); ok {
//			if f := rv.FieldByName(fieldMediaType); f.IsValid() && f.Kind() == reflect.Int32 {
//				if mediaType, ok := common.MediaType_value[v]; ok {
//					f.SetInt(int64(mediaType))
//				}
//			}
//		}
//		// メディアユーザーID変換
//		if v, ok := getMetadataValue(md, keyMediaUserID); ok {
//			if f := rv.FieldByName(fieldMediaUserId); f.IsValid() && f.Kind() == reflect.String {
//				f.SetString(v)
//			}
//		}
//		return invoker(ctx, method, req, reply, cc, opts...)
//	}
//}

//func getMetadataValue(md metadata.MD, key string) (string, bool) {
//	v := md.Get(key)
//	if len(v) > 0 {
//		return v[0], true
//	}
//	return "", false
//}
