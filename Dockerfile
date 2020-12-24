FROM alpine:3.12

RUN mkdir /lib64
RUN ln -s /lib/libc.musl-x86_64.so.1 /lib64/ld-linux-x86-64.so.2

ENV PYTHONUNBUFFERED=1
RUN apk add --update --no-cache python3 && ln -sf python3 /usr/bin/python
RUN python3 -m ensurepip
RUN pip3 install --no-cache --upgrade pip setuptools requests

RUN apk add --no-cache ca-certificates curl bash \
    && mkdir /data /data/bin

WORKDIR /data

ADD ./cmd/prometheus-bm/prometheus-bm /data/bin/

ENV PATH=${PATH}:/data/bin
