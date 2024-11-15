FROM golang:1.23.3-alpine as builder

RUN mkdir /app

COPY . /app

WORKDIR /app

RUN CGO_ENABLED=0 go build -o karaBank .

RUN chmod +x /app/karaBank

# build a tiny docker image
FROM alpine:latest

RUN mkdir /app

COPY --from=builder /app/karaBank /app

CMD [ "/app/karaBank" ]