## Build
FROM golang:1.21 AS build

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY *.go ./

RUN CGO_ENABLED=0 go build -o /notify-address

## Deploy
FROM gcr.io/distroless/static

WORKDIR /

COPY --from=build /notify-address /notify-address

ENTRYPOINT ["/notify-address"]