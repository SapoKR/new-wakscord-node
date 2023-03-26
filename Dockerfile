FROM golang:1.20.2-alpine3.17 AS builder

WORKDIR /app

COPY . .

RUN go build -o app

FROM scratch

COPY --from=builder /app/app .

ENV KEY HOST PORT ID OWNER MAX_CONCURRENT WAIT_CONCURRENT

ENTRYPOINT [ "./app" ]