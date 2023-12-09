# ビルドステージ
FROM golang:1.21 AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o server .

# 実行ステージ
#FROM jrottenberg/ffmpeg:4.4-alpine

FROM ubuntu:20.04
RUN apt-get update && apt-get install -y ffmpeg ca-certificates && rm -rf /var/lib/apt/lists/*

RUN apt-get update && apt install -y golang && rm -rf /var/lib/apt/lists/*

WORKDIR /bin

# アプリケーションの実行可能ファイルをコピー
COPY --from=builder /app/server .

# ポートのエクスポーズ
EXPOSE 8080

# コンテナ起動時にアプリケーションを実行
CMD ["/app/server"]
