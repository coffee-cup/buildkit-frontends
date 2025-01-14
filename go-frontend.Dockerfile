FROM golang:1.23-alpine AS builder

WORKDIR /src

# Copy go.mod and go.sum first to leverage Docker cache
COPY go.mod go.sum ./
RUN go mod download

# Now copy the rest of the source code
COPY . .
RUN go build -o /usr/bin/frontend .

FROM alpine

LABEL org.opencontainers.image.source=https://github.com/coffee-cup/buildkit-frontends
LABEL org.opencontainers.image.description="Go Frontend for Buildkit"
LABEL org.opencontainers.image.licenses=MIT

COPY --from=builder /usr/bin/frontend /usr/bin/frontend
ENTRYPOINT ["/usr/bin/frontend"]
