FROM golang:1.20 as builder
WORKDIR "/builder"
COPY . ./
RUN CGO_ENABLED=1 GOOS=linux go build -o appp ./cmd

FROM debian:stable-slim

WORKDIR /app/

COPY --from=builder /builder/appp  .
COPY --from=builder /builder/.env  .
COPY --from=builder /builder/config.yml  .
COPY --from=builder /builder/certificate.crt  .
COPY --from=builder /builder/private.pem  .
COPY --from=builder /builder/public.pem  .
COPY --from=builder /builder/blunder.html  .
CMD ["./appp"]
