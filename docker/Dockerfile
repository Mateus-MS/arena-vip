FROM golang:1.23-alpine AS builder

ARG TAILWIND_VERSION=v3.4.17

RUN go install github.com/a-h/templ/cmd/templ@v0.3.1001

RUN wget -qO /usr/local/bin/tailwindcss \
    "https://github.com/tailwindlabs/tailwindcss/releases/download/${TAILWIND_VERSION}/tailwindcss-linux-x64" \
    && chmod +x /usr/local/bin/tailwindcss

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN templ generate
RUN tailwindcss -i ./static/css/input.css -o ./static/css/output.css --minify
RUN CGO_ENABLED=0 GOOS=linux go build -o arena-vip .

FROM alpine:3.21

WORKDIR /app

COPY --from=builder /build/arena-vip .
COPY --from=builder /build/static ./static

EXPOSE 8080

ENTRYPOINT ["./arena-vip"]
