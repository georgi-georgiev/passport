FROM golang:1.20 as builder
WORKDIR "/builder"
COPY . ./
RUN CGO_ENABLED=1 GOOS=linux go build -o appp ./

FROM debian:stable-slim

WORKDIR /app/

COPY --from=builder /builder/appp  .
COPY --from=builder /builder/.env  .
COPY --from=builder /builder/config.yml  .
CMD ["./appp"]
