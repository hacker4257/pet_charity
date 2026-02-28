# ---- 构建阶段 ----
FROM golang:1.25-alpine AS builder

WORKDIR /app

# 先复制依赖文件，利用 Docker 缓存
COPY go.mod go.sum ./
RUN go mod download

# 复制源码并编译
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o server ./cmd/server

# ---- 运行阶段 ----
FROM alpine:3.19

WORKDIR /app

# 安装时区数据（Go 的 time 包需要）
RUN apk add --no-cache tzdata ca-certificates

COPY --from=builder /app/server .
COPY --from=builder /app/configs ./configs

# 创建上传目录
RUN mkdir -p uploads

EXPOSE 8080

CMD ["./server"]