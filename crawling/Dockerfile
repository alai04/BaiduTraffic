FROM golang:1.14-alpine

RUN apk add --no-cache tzdata \
    && ln -snf /usr/share/zoneinfo/Asia/Shanghai /etc/localtime \
    && echo "Asia/Shanghai" > /etc/timezone
ENV TZ Asia/Shanghai

RUN mkdir -p /opt/crawling
WORKDIR /opt/crawling
ADD *.go /opt/crawling/
ADD go.* /opt/crawling/
RUN go build

CMD ["./crawling", "-t"]
