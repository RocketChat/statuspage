FROM golang:1.12.8 as backend

RUN go get -u github.com/golang/dep/cmd/dep

WORKDIR /go/src/github.com/RocketChat/statuscentral

COPY . .

RUN dep ensure && \
    GOOS=linux && \
    go build

FROM alpine:latest

WORKDIR /root/

RUN apk --no-cache add ca-certificates libc6-compat && \ 
    mkdir app

ENV GIN_MODE=release

COPY --from=backend /go/src/github.com/RocketChat/statuscentral/statuscentral .
COPY --from=backend /go/src/github.com/RocketChat/statuscentral/templates templates
COPY --from=backend /go/src/github.com/RocketChat/statuscentral/static static

EXPOSE 5000

CMD ["/root/statuscentral"]