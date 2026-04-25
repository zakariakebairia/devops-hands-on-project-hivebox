FROM golang:1.25-alpine AS builder
WORKDIR /app
COPY go.mod ./
RUN go mod download
COPY main.go box.go ./
RUN CGO_ENABLED=0 GOOS=linux go build -o hivebox .

# I used scratch, but I got issues with tls certificates
# Since scratch doesn't have anything in it
# SOLUTION: COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
# FROM scratch
FROM gcr.io/distroless/static:nonroot
COPY --from=builder /app/hivebox /hivebox
ENTRYPOINT [ "/hivebox" ]
