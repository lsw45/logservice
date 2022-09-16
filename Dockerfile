FROM swr.cn-east-3.myhuaweicloud.com/cocos-paas/golang:tzdata AS builder
LABEL stage=gobuilder

ENV CGO_ENABLED 0
ENV GOPROXY https://goproxy.cn,direct

WORKDIR /build

ADD go.mod .
ADD go.sum .
RUN go mod download
COPY . .
COPY server/conf /app/conf

RUN go build -ldflags="-s -w" -o /app/logservice2 server/main.go

FROM swr.cn-east-3.myhuaweicloud.com/cocos-paas/alpine:tzdata

COPY --from=builder /usr/share/zoneinfo/Asia/Shanghai /usr/share/zoneinfo/Asia/Shanghai
ENV TZ Asia/Shanghai

WORKDIR /app
COPY --from=builder /app/logservice2 /app/logservice2
COPY --from=builder /app/conf /app/conf

# Expose port
EXPOSE 38080

# 指定启动命令
ENTRYPOINT ["/opt/logservice2"]
