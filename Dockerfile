FROM golang:1.22-bullseye AS deploy-builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -trimpath -ldflags "-w -s" -o app

# ---------------------------------------------------

FROM debian:bullseye-slim AS deploy

RUN apt-get update
# deploy-builderの/appディレクトリ上のappという実行ファイルをdeployのカレントへ配置
COPY --from=deploy-builder /app/app .
# 配置したapp実行ファイルを実行
CMD ["./app"]

# ---------------------------------------------------

FROM golang:1.22 AS dev
WORKDIR /app
RUN go install github.com/air-verse/air@latest
CMD ["air"]