FROM swr.cn-east-3.myhuaweicloud.com/cocos-paas/golang:tzdata AS builder

ARG UID=1000
ARG GID=1000

ADD ./bin /opt/logservice2

ADD ./infrastructure/i18n/locales/ /opt/logservice2/infrastructure/i18n/locales/

ADD ./server/conf/logservice2.yaml /opt/logservice2/

RUN groupadd -r -g $GID cocos \
    && useradd -r -u $UID -g cocos cocos \
    && chown -R $GID:$UID /opt/logservice2

# 指定工作目录
WORKDIR /opt/logservice2/

# Change user
USER 1000

# Expose port
EXPOSE 13020

# 指定启动命令
ENTRYPOINT [ "./logservice2"]
