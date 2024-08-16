FROM golang:latest AS build

WORKDIR /build

COPY . .

RUN go build

FROM golang:latest

WORKDIR /app

COPY --from=BUILD /build/http_server .

EXPOSE 8000

CMD ["./http_server", "-addr", "0.0.0.0:8000"]
