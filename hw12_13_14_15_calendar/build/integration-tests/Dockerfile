FROM golang:1.19
ENV CODE_DIR /go/src/

WORKDIR ${CODE_DIR}

LABEL ORGANIZATION="OTUS Online Education"
LABEL SERVICE="intergration_tests"
LABEL MAINTAINERS="frenatz@gmail.com"

# Кэшируем слои с модулями
#COPY go.mod .
#COPY go.sum .
COPY . ${CODE_DIR}
RUN go mod download

CMD go test -run Integration ./internal/integration-tests/


