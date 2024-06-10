FROM golang:1.22.1-alpine as build

ENV BIN_FILE /opt/calendar/calendar-app
ENV CODE_DIR /go/src/app

WORKDIR ${CODE_DIR}

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY ../../ ${CODE_DIR}

ARG VERSION
ARG BUILD_TIME

RUN CGO_ENABLED=0 go build -ldflags "-X main.Version=${VERSION} -X main.BuildTime=${BUILD_TIME}" -o ${BIN_FILE} ./cmd/calendar

FROM alpine:3.9

LABEL ORGANIZATION="OTUS Online Education"
LABEL SERVICE="calendar"
LABEL MAINTAINERS="j.zadorozhnaya@gmail.com"

ENV BIN_FILE "/opt/calendar/calendar-app"
COPY --from=build ${BIN_FILE} ${BIN_FILE}

ENV CONFIG_FILE /etc/calendar/calendar_config.toml
COPY ../../config/calendar_config.toml ${CONFIG_FILE}

CMD ${BIN_FILE} -config ${CONFIG_FILE}
