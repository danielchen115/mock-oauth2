FROM golang:1.14

ENV APP_NAME mock-oauth2
ENV PORT 8090

COPY . /go/src/${APP_NAME}
WORKDIR /go/src/${APP_NAME}

RUN go get ./
RUN go build ./cmd/server/ -o ${APP_NAME}

CMD ./${APP_NAME}

EXPOSE ${PORT}
