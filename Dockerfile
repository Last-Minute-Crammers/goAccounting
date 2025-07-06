FROM golang:1.24-alpine AS builder


ENV GOPROXY=https://goproxy.cn,direct
ENV GOSUMDB=sum.golang.google.cn

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o goAccounting main.go

# 第二阶段：精简运行镜像
FROM alpine:latest

# 安装必要的系统依赖
RUN apk --no-cache add tzdata ca-certificates && \
    cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime

WORKDIR /app

# 拷贝编译好的二进制和配置
COPY --from=builder /app/goAccounting .
COPY config.yaml .

# 需要其它静态文件/证书等，也在这里 COPY

EXPOSE 8080 

CMD ["./goAccounting"]