FROM golang:1.14-alpine AS build

RUN apk add --no-cache ca-certificates git
WORKDIR /go/src/github.com/RocketChat/statuscentral
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a

FROM scratch as runtime

ARG GIN_MODE=release
ARG PORT=5000
ARG CONF_FILE=statuscentral.yaml
ENV GIN_MODE=${GIN_MODE}
ENV PORT=${PORT}
ENV CONF_FILE=${CONF_FILE}

WORKDIR /usr/local/statuscentral

COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=build /go/src/github.com/RocketChat/statuscentral/${CONF_FILE} ./statuscentral.yaml
COPY --from=build /go/src/github.com/RocketChat/statuscentral/statuscentral .
COPY --from=build /go/src/github.com/RocketChat/statuscentral/templates templates
COPY --from=build /go/src/github.com/RocketChat/statuscentral/static static

EXPOSE ${PORT}

ENTRYPOINT ["./statuscentral"]