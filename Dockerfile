FROM golang:1.15-alpine AS build

WORKDIR /go/src/app

COPY go .

RUN go get -d -v ./...
RUN go build -v .

FROM alpine:3 as prod
WORKDIR /app
COPY --from=build /go/src/app/airelibre-bot /app

ENV ACCESS_TOKEN="" \
    ACCESS_TOKEN_SECRET="" \
    CONSUMER_KEY="" \
    CONSUMER_SECRET="" \
    API_URL=https://rald-dev.greenbeep.com/api/v1/aqi    

CMD ["./airelibre-bot"]
