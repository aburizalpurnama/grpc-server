FROM golang:1.21-alpine AS build

ENV APP_PORT=8081
ENV DB_USERNAME=test
ENV DB_PASSWORD=test
ENV DB_HOST=localhost
ENV DB_PORT=3435
ENV DB_NAME=test

WORKDIR /app

COPY . /app/
RUN go mod tidy
RUN go build -o /app/main

FROM alpine:3
WORKDIR /app
COPY --from=build /app/main .
EXPOSE ${APP_PORT}
CMD /app/main