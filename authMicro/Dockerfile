FROM golang:1.25-alpine AS builder-debug
WORKDIR /app
RUN apk add --no-cache git
COPY common/go.mod ./common/
COPY authMicro/go.mod authMicro/go.sum ./authMicro/
COPY authMicro/grpcApi/go.mod ./authMicro/grpcApi/
RUN go work init ./common ./authMicro ./authMicro/grpcApi
RUN go work sync
RUN go mod download
WORKDIR /app/authMicro
RUN go install github.com/go-delve/delve/cmd/dlv@latest
COPY common/ /app/common/
COPY authMicro/cmd/ ./cmd/
COPY authMicro/internal/ ./internal/
COPY authMicro/utlis/ ./utlis/
COPY authMicro/grpcApi/ ./grpcApi/
COPY authMicro/docs/ ./docs/
RUN CGO_ENABLED=0 GOOS=linux go build -gcflags "all=-N -l" -o auth-app ./cmd/app

FROM golang:1.25-alpine AS builder-release
WORKDIR /app
RUN apk add --no-cache git
COPY common/go.mod ./common/
COPY authMicro/go.mod authMicro/go.sum ./authMicro/
COPY authMicro/grpcApi/go.mod ./authMicro/grpcApi/
RUN go work init ./common ./authMicro ./authMicro/grpcApi
RUN go work sync
RUN go mod download
WORKDIR /app/authMicro
COPY common/ /app/common/
COPY authMicro/cmd/ ./cmd/
COPY authMicro/internal/ ./internal/
COPY authMicro/utlis/ ./utlis/
COPY authMicro/grpcApi/ ./grpcApi/
COPY authMicro/docs/ ./docs/
RUN CGO_ENABLED=0 GOOS=linux go build -o auth-app ./cmd/app

FROM alpine:3.19 AS run-debug
WORKDIR /app
COPY --from=builder-debug /app/authMicro/auth-app .
COPY authMicro/db/migrations /migrations
COPY authMicro/locales ./locales
COPY --from=builder-debug /go/bin/dlv /usr/local/bin/dlv
EXPOSE 40000
EXPOSE 80
CMD dlv exec ./auth-app --headless --listen=:40000 --api-version=2 --accept-multiclient --continue --log

FROM alpine:3.19 AS run-release
WORKDIR /app
COPY --from=builder-release /app/authMicro/auth-app .
COPY authMicro/db/migrations /migrations
COPY authMicro/locales ./locales
EXPOSE 80
CMD ./auth-app
