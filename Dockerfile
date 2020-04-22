FROM golang:1.14

ENV APP_NAME mock-oauth2
ENV PORT 8090

COPY . /app/${APP_NAME}
WORKDIR /app/${APP_NAME}

RUN go build -o ${APP_NAME} ./cmd/server/

CMD ./${APP_NAME}

EXPOSE ${PORT}
