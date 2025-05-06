# Build 專案
FROM    codingxiang/go_vc AS builder
ENV     RUN_PATH=/audit PROJ_PATH=/build
RUN     mkdir -p $RUN_PATH
USER    root
ADD     . ${PROJ_PATH}
WORKDIR ${PROJ_PATH}
RUN     make build pack \
        && tar -zxf audit-v*.tar.gz -C ${RUN_PATH} \
        && rm -rf ${PROJ_PATH}
# 打包 Image
FROM    registry.digiwincloud.com.cn/base/alpine:3.1.0.0
USER    root
ENV     RUN_PATH=/audit app=audit
RUN     mkdir -p $RUN_PATH && apk add --no-cache ca-certificates bash
#安裝 apk
ENV     TZ=Asia/Taipei
RUN     ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone
COPY    --from=builder ${RUN_PATH} ${RUN_PATH}
WORKDIR ${RUN_PATH}
ENTRYPOINT ["/setup/start.sh"]