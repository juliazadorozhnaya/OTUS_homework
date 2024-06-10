FROM golang:1.22.1-alpine as build

ENV BIN_FILE /bin/app
ENV CODE_DIR /go/src/

WORKDIR ${CODE_DIR}
COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . ${CODE_DIR}

ARG LDFLAGS
RUN CGO_ENABLED=0 go build -ldflags "$LDFLAGS" -o ${BIN_FILE} ./cmd/calendar_scheduler

FROM alpine:3.16

LABEL ORGANIZATION="OTUS Online Education"
LABEL SERVICE="calendar"
LABEL MAINTAINERS="j.zadorozhnaya@gmail.com"

ENV BIN_FILE /bin/app
ENV CONFIG_FILE config.toml

COPY --from=build ${BIN_FILE} ${BIN_FILE}
COPY config/scheduler_config.toml ${CONFIG_FILE}

CMD sleep 15 && ${BIN_FILE} -config ${CONFIG_FILE}
