FROM golang:latest AS build

WORKDIR /build

COPY . .

RUN go build

FROM golang:latest

WORKDIR /app

COPY --from=BUILD /build/http_server .

CMD ["./http_server"]
