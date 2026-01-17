FROM golang:1.25-alpine AS builder-debug
WORKDIR /workspace
RUN apk add --no-cache git
COPY ./common/ ./common/
COPY ./userMicro/go.mod ./userMicro/go.sum ./userMicro/
COPY ./authMicro/grpcApi/go.mod ./authMicro/grpcApi/go.sum ./authMicro/grpcApi/
RUN go work init ./common ./userMicro ./authMicro/grpcApi
RUN go work sync
RUN go mod download
RUN go install github.com/go-delve/delve/cmd/dlv@latest
COPY ./common/ /workspace/common/
COPY ./authMicro/grpcApi/ /workspace/authMicro/grpcApi/
COPY ./userMicro/ /workspace/userMicro/
RUN CGO_ENABLED=0 GOOS=linux go build -gcflags "all=-N -l" -o user-app ./userMicro/cmd/app

FROM alpine:3.19 AS run-debug
WORKDIR /app
COPY --from=builder-debug /workspace/user-app .
COPY --from=builder-debug /workspace/userMicro/internal/infrastructure/outbound/repository/database/migrations /migrations
COPY --from=builder-debug /workspace/userMicro/locales ./locales
COPY --from=builder-debug /go/bin/dlv /usr/local/bin/dlv
EXPOSE 80
CMD dlv exec ./user-app --headless --listen=:${DEBUG_PORT:-40001} --api-version=2 --accept-multiclient --continue --log
