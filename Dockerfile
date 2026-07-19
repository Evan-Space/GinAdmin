FROM golang:1.26-alpine AS builder

WORKDIR /app

# 第一阶段
# 先复制依赖描述文件，利用 Docker 缓存
COPY go.mod go.sum ./
RUN go mod download

# 复制全部源码，编译成二进制文件
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o gin-admin .

# 第二阶段
FROM alpine:3.21

WORKDIR /app

# 从构建阶段拿编译好的二进制
COPY --from=builder /app/gin-admin .

# 把 Docker 专用配置覆盖为容器内的默认配置
COPY --from=builder /app/config/config.docker.yaml ./config/config.yaml

# Casbin 需要的权限模型文件
COPY --from=builder /app/config/rbac_model.conf .

# 声明容器监听的端口
EXPOSE 8080

# 启动服务
CMD ["./gin-admin", 'server']