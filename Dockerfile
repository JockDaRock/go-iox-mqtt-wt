FROM golang:alpine as go-build

RUN apk update
RUN apk add git
RUN go get gopkg.in/ini.v1
RUN go get github.com/eclipse/paho.mqtt.golang

WORKDIR /

COPY main.go .
COPY package_config.ini .

RUN go build -o goMQ main.go

FROM alpine

COPY --from=go-build /goMQ /usr/local/bin/goMQ
RUN chmod +x /usr/local/bin/goMQ
COPY package_config.ini /usr/local/bin/package_config.ini

CMD ["goMQ"]

LABEL cisco.info.name="DevNet Create IOx Demo" \
      cisco.info.description="Simple App to relay IoT Data" \
      cisco.info.version="0.02" \
      cisco.info.author-link="" \
      cisco.info.author-name="" \
      cisco.type=docker \
      cisco.cpuarch=x86_64 \
      cisco.resources.profile=custom \
      cisco.resources.cpu=100 \
      cisco.resources.memory=50 \
      cisco.resources.disk=10 \
      cisco.resources.network.0.interface-name=eth0


