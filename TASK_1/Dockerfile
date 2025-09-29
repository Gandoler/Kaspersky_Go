

FROM golang:1.24-alpine AS builder
WORKDIR /src


COPY go.mod ./
RUN go mod download

COPY . .

ENV CGO_ENABLED=0
RUN go build -o /out/app .


FROM scratch AS runner
WORKDIR /

ENV WORKERS=4
ENV QUEUE_SIZE=64

EXPOSE 8080
COPY --from=builder /out/app /app

ENTRYPOINT ["/app"]
