FROM golang:1.21-bullseye as build

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY producer ./producer
COPY consumer ./consumer

RUN mkdir -p /out
RUN go build -o /out/producer ./producer
RUN go build -o /out/consumer ./consumer

#

FROM debian:bullseye-slim

COPY --from=build /out/producer /app/producer
COPY --from=build /out/consumer /app/consumer

EXPOSE 8080