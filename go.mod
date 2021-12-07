module github.com/payfazz/fz-sentry

go 1.13

require (
	github.com/go-kit/kit v0.10.0
	github.com/go-redis/redis v6.15.9+incompatible
	github.com/go-sql-driver/mysql v1.4.1
	github.com/gofrs/uuid v3.3.0+incompatible
	github.com/grpc-ecosystem/go-grpc-middleware v1.3.0
	github.com/grpc-ecosystem/go-grpc-prometheus v1.2.0
	github.com/jmoiron/sqlx v1.2.0
	github.com/lib/pq v1.8.0
	github.com/payfazz/fz-router v0.0.0-20200807154353-fb9dfc3429a4
	github.com/prometheus/client_golang v1.3.0
	github.com/slack-go/slack v0.7.2
	go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc v0.27.0
	go.opentelemetry.io/otel/trace v1.2.0
	go.uber.org/zap v1.15.0
	golang.org/x/net v0.0.0-20210726213435-c6fcb2dbf985 // indirect
	golang.org/x/sys v0.0.0-20210630005230-0f9fa26af87c // indirect
	google.golang.org/genproto v0.0.0-20210729151513-df9385d47c1b // indirect
	google.golang.org/grpc v1.42.0
)
