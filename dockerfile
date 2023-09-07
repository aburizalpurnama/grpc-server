FROM golang:1.21-alpine AS build
ENV APP_PORT=8081
WORKDIR /app

COPY . /app/
RUN go mod tidy
RUN go build -o /app/main

FROM alpine:3
WORKDIR /app
COPY --from=build /app/main .
EXPOSE ${APP_PORT}
CMD /app/main