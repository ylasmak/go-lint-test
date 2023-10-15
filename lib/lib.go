package lib

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"net"
	"net/http"
)

func GetPrivateIPAddress() ([]string, error) {
	// get list of available addresses
	var ipAddresses []string
	addr, err := net.InterfaceAddrs()
	if err != nil {
		return nil, err
	}
	for _, addr := range addr {
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() && ipnet.IP.To4() != nil {
			ipAddresses = append(ipAddresses, ipnet.IP.String())
		}
	}
	return ipAddresses, nil
}

func HTTPSClient(cert string) *http.Client {
	// Create a custom TLS configuration
	tlsConfig := &tls.Config{
		InsecureSkipVerify: false, // Ignore verification of the server's certificate
		MinVersion:         tls.VersionTLS12,
	}

	// Load the certificate and key into a certificate pool
	certPool := x509.NewCertPool()
	certPool.AppendCertsFromPEM([]byte(cert))

	// Add the certificate pool to the TLS configuration
	tlsConfig.RootCAs = certPool
	// Create an HTTP client with the custom TLS configuration
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: tlsConfig,
		},
	}
	return client
}

func HTTPClient() *http.Client {
	// Create an HTTP client with the custom TLS configuration
	client := &http.Client{}
	return client
}

func GetRequestWithContext(url string) (*http.Request, error) {
	return http.NewRequestWithContext(context.TODO(), http.MethodGet, url, nil)
}
