package lib

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net"
	"net/http"
)

func GetPrivateIpAddress() ([]string, error) {
	// get list of available addresses

	var ipAddresses []string
	addr, err := net.InterfaceAddrs()
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	for _, addr := range addr {
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() && ipnet.IP.To4() != nil {
			// print available addresses
			fmt.Println(ipnet.IP.String())
			ipAddresses = append(ipAddresses, ipnet.IP.String())

		}
	}
	return ipAddresses, nil
}

func HttpsClient(cert string) *http.Client {
	// Create a custom TLS configuration
	tlsConfig := &tls.Config{
		InsecureSkipVerify: false, // Ignore verification of the server's certificate
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

func HttpClient() *http.Client {
	// Create an HTTP client with the custom TLS configuration
	client := &http.Client{}
	return client

}
