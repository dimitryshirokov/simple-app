FROM golang:1.20.8-alpine3.17 as builder
ENV GO111MODULE=on
WORKDIR /build
COPY . .
ENV CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64
RUN go build -mod vendor -installsuffix cgo -o bin/app cmd/main.go && \
    go build -mod vendor -installsuffix cgo -o bin/migrator cmd/migrator.go

FROM alpine:3.17
ENV TZ=Europe/Moscow \
    LANG=en_US.UTF-8 \
    LANGUAGE=en_US.UTF-8 \
    LC_CTYPE=en_US.UTF-8 \
    LC_ALL=en_US.UTF-8
RUN apk add --no-cache tzdata ca-certificates; \
    cp /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone; \
    mkdir -p /usr/local/bin
COPY --from=builder /build/bin/* /usr/local/bin/
