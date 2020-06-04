FROM golang:1.14-alpine

RUN apk add --no-cache tzdata \
    && ln -snf /usr/share/zoneinfo/Asia/Shanghai /etc/localtime \
    && echo "Asia/Shanghai" > /etc/timezone
ENV TZ Asia/Shanghai

RUN mkdir -p /opt/adjust
WORKDIR /opt/adjust
ADD *.go /opt/adjust/
ADD go.* /opt/adjust/
RUN go build

EXPOSE 8080
CMD ["./adjust", "-h"]
