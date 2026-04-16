FROM golang:1.25-alpine AS builder-debug
WORKDIR /app
RUN apk add --no-cache git
COPY common/go.mod ./common/
COPY sample/userMicro/go.mod sample/userMicro/go.sum ./sample/userMicro/
COPY authMicro/grpcApi/go.mod ./authMicro/grpcApi/
RUN go work init ./common ./sample/userMicro ./authMicro/grpcApi
RUN go work sync
RUN go mod download
WORKDIR /app/sample/userMicro
RUN go install github.com/go-delve/delve/cmd/dlv@latest
COPY common/ /app/common/
COPY authMicro/grpcApi/ /app/authMicro/grpcApi/
COPY sample/userMicro/ ./
RUN CGO_ENABLED=0 GOOS=linux go build -gcflags "all=-N -l" -o user-app ./cmd/app

FROM golang:1.25-alpine AS builder-release
WORKDIR /app
RUN apk add --no-cache git
COPY common/go.mod ./common/
COPY sample/userMicro/go.mod sample/userMicro/go.sum ./sample/userMicro/
COPY authMicro/grpcApi/go.mod ./authMicro/grpcApi/
RUN go work init ./common ./sample/userMicro ./authMicro/grpcApi
RUN go work sync
RUN go mod download
WORKDIR /app/sample/userMicro
COPY common/ /app/common/
COPY authMicro/grpcApi/ /app/authMicro/grpcApi/
COPY sample/userMicro/ ./
RUN CGO_ENABLED=0 GOOS=linux go build -o user-app ./cmd/app

FROM alpine:3.19 AS run-debug
WORKDIR /app
COPY --from=builder-debug /app/sample/userMicro/user-app .
COPY --from=builder-debug /app/sample/userMicro/internal/infrastructure/outbound/repository/database/migrations /migrations
COPY --from=builder-debug /app/sample/userMicro/locales ./locales
COPY --from=builder-debug /go/bin/dlv /usr/local/bin/dlv
EXPOSE 40001
EXPOSE 80
CMD dlv exec ./user-app --headless --listen=:40001 --api-version=2 --accept-multiclient --continue --log

FROM alpine:3.19 AS run-release
WORKDIR /app
COPY --from=builder-release /app/sample/userMicro/user-app .
COPY --from=builder-release /app/sample/userMicro/internal/infrastructure/outbound/repository/database/migrations /migrations
COPY --from=builder-release /app/sample/userMicro/locales ./locales
EXPOSE 80
CMD ["./user-app"]