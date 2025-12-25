FROM golang:1.25-alpine AS builder-debug
WORKDIR /workspace
RUN apk add --no-cache git
COPY common/go.mod common/go.sum ./common/
COPY userMicro/go.mod userMicro/go.sum ./userMicro/
RUN go work init ./common ./userMicro
RUN go work sync
WORKDIR /workspace/userMicro
RUN --mount=type=cache,target=/go/pkg/mod go mod download
RUN go install github.com/go-delve/delve/cmd/dlv@latest
COPY common/ /workspace/common/
COPY userMicro/ ./
RUN --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 GOOS=linux go build -gcflags "all=-N -l" -o user-app ./cmd/app

FROM alpine:3.19 AS run-debug
WORKDIR /app
COPY --from=builder-debug /workspace/userMicro/user-app .
COPY --from=builder-debug /workspace/userMicro/internal/repository/database/migrations /migrations
COPY --from=builder-debug /workspace/userMicro/locales ./locales
COPY --from=builder-debug /go/bin/dlv /usr/local/bin/dlv
EXPOSE 80
CMD dlv exec ./user-app --headless --listen=:${DEBUG_PORT:-40001} --api-version=2 --accept-multiclient --continue --log
