FROM alpine

ARG ARCH

RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.tuna.tsinghua.edu.cn/g' /etc/apk/repositories && \
	apk update --no-cache && apk --no-cache add curl bash sqlite

ENV TZ Asia/Shanghai

WORKDIR /app
COPY ./frp-panel-${ARCH} /app/frp-panel
COPY ./etc /app/etc

RUN ln -sf /app/etc/Shanghai /etc/localtime && mv /app/etc/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt && mkdir -p /data

# web port
EXPOSE 9000

# rpc port
EXPOSE 9001

ENV DB_DSN /data/data.db

ENTRYPOINT [ "/app/frp-panel" ]

CMD [ "master" ]