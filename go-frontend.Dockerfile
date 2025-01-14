FROM golang:1.23-alpine AS builder
WORKDIR /src
COPY . .
RUN go build -o /usr/bin/frontend .

FROM alpine
COPY --from=builder /usr/bin/frontend /usr/bin/frontend
ENTRYPOINT ["/usr/bin/frontend"]
