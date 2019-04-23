FROM golang:alpine as go-build

RUN apk update
RUN apk add git
RUN go get gopkg.in/ini.v1
RUN go get github.com/eclipse/paho.mqtt.golang

WORKDIR /

COPY main.go .
COPY msg.ini .

RUN go build -0 goMQ main.go

FROM alpine

COPY --from=go-build /goMQ /usr/local/bin/goMQ
RUN chmod +x /usr/local/bin/goMQ
COPY msg.ini /usr/local/bin/msg.ini

CMD ["goMQ"]


