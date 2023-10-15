 #!/bin/bash

export COLOR="GREEN"
export API_TYPE="BUSINESS"
export SSL_CERTIFICATE=$(cat ./cert/jwt.gateway.k8s.local.fullchain.crt | base64)
export SSL_PRIVATE_KEY=$(cat ./cert/jwt.gateway.k8s.local.key | base64)
export SSL_CA=$(cat ./cert/server_rootCA.crt | base64)
export PORT="8080"
export SECURE="FALSE"

go run main.go

# go install golang.org/x/vuln/cmd/govulncheck@latest
# go install github.com/securego/gosec/v2/cmd/gosec@latest
