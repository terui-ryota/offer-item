# Application Error

# 目的

* アプリケーションのエラーを一元管理する
* アプリケーション外（ライブラリーなど）で発生したエラーをラップする
* エラーの原因によって処理を分岐させる

# 使用方法

## エラーのラップ

```go
a, err := entity.FindAffiliator(ctx, tx, string(id))
if err != nil {
    return nil, apperr.AffiliatorNotFound.Wrap(err)
}
```

## エラーによる処理分岐

xerrors を使用する

```go
a, err := affiliatorRepository.Fetch(id)
if err != nil {
    if xerrors.Is(err, apperr.AffiliatorNotFound) {
    	...
    }
}
```

# RPC 間のエラーの伝播

interceptor を実装する

## gRPC Server

ApplicationErrorStreamServerInterceptor()、ApplicationErrorUnaryServerInterceptor() を実装する

```go
server := grpc.NewServer(
    grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(
        grpc_zap.StreamServerInterceptor(logger),
        grpc_recovery.StreamServerInterceptor(),
        common_metadata.CommonMetadataStreamServerInterceptor(),
        custom_log.CustomLogStreamServerInterceptor(config.ContextName),
        application_error.ApplicationErrorStreamServerInterceptor(),
    )),
    grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
        grpc_zap.UnaryServerInterceptor(logger),
        grpc_recovery.UnaryServerInterceptor(),
        common_metadata.CommonMetadataUnaryServerInterceptor(),
        custom_log.CustomLogUnaryServerInterceptor(config.ContextName),
        application_error.ApplicationErrorUnaryServerInterceptor(),
    )),
)
```

## gRPC Client (unary)

ApplicationErrorUnaryClientInterceptor() を実装する

```go
conn, err := grpc.Dial(
    target,
    grpc.WithInsecure(),
    grpc.WithUnaryInterceptor(
        grpc_middleware.ChainUnaryClient(
            common_metadata.CommonMetadataUnaryClientInterceptor(),
            application_error.ApplicationErrorUnaryClientInterceptor(),
        ),
    ),
)
```

## gRPC Client（stream）
ApplicationErrorStreamClientInterceptor() を実装する

```go
conn, err := grpc.Dial(
    target,
    grpc.WithInsecure(),
    grpc.WithStreamInterceptor(
        grpc_middleware.ChainStreamClient(
            common_metadata.CommonMetadataStreamClientInterceptor(),
            application_error.ApplicationErrorStreamClientInterceptor(),
        ),
    ),
)
```
