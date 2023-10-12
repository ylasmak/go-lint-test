## Build
FROM golang:1.20.4-buster AS build

# Install necessary packages
RUN apt-get update && apt-get install -y ca-certificates

# # Copy root CA certificate
# COPY ./cert/server_rootCA.crt /usr/local/share/ca-certificates/
# # Update the certificate store
# RUN update-ca-certificates
# RUN   ls /etc/ssl/certs/

WORKDIR /app

COPY . .

RUN go mod tidy
RUN go mod vendor
RUN go build -o  /go-proxy -ldflags="-s -w" main.go 

# Deploy
FROM gcr.io/distroless/base-debian11
WORKDIR /

# COPY --from=build /etc/ssl/certs/ /etc/ssl/certs/


COPY --from=build /go-proxy /go-proxy
USER nonroot:nonroot
EXPOSE 8080
ENTRYPOINT ["/go-proxy"]