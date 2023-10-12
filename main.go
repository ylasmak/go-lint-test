package main

import (
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/proxy/lib"
)

//go:generate oapi-codegen --package=lib -generate=types -o ./lib/model.gen.go ./swagger.yaml

const RESPONSE = "%s Service %s "

type Server struct {
	Color       string
	APIType     string
	Certificate string
	PrivateKey  string
	CA          string
	Secure      bool
	Port        string
}

func main() {

	s := NewServer()

	router := gin.Default()
	path := "/" + strings.ToLower(s.Color)
	router.GET(path, s.getIdentity)
	router.POST(path, s.callAPI)
	router.POST(path+"/selfSigned-ca", s.callAPIWithSelfSignedCert)

	if s.Secure {
		s.startTLSServer(router)
	} else {
		s.startServer(router)
	}

}

func (s Server) startTLSServer(router *gin.Engine) {
	tlsCert, err := tls.X509KeyPair([]byte(s.Certificate), []byte(s.PrivateKey))
	if err != nil {
		panic(err)
	}

	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{tlsCert},
	}

	server := &http.Server{
		Addr:      fmt.Sprintf(":%s", s.Port),
		Handler:   router,
		TLSConfig: tlsConfig,
	}

	err = server.ListenAndServeTLS("", "")
	if err != nil {
		panic(err)
	}
}

func (s Server) startServer(router *gin.Engine) {

	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", s.Port),
		Handler: router,
	}

	err := server.ListenAndServe()
	if err != nil {
		panic(err)
	}

}

func (s Server) getIdentity(c *gin.Context) {
	ips, err := lib.GetPrivateIpAddress()
	
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
	}
	
	name := fmt.Sprintf(RESPONSE, s.Color, s.APIType)
	service := lib.ServiceName{
		Name:       &name,
		PrivateIp:  &ips,
		RequestURI: &c.Request.RequestURI,
		Host:       &c.Request.Host,
	}
	c.JSON(http.StatusOK, service)
}

func (s Server) callAPIWithSelfSignedCert(c *gin.Context) {
	s.invokeAPI(c, true)
}

func (s Server) callAPI(c *gin.Context) {
	s.invokeAPI(c, false)
}

func (s Server) invokeAPI(c *gin.Context, checkSelfSignedCA bool) {

	ips, err := lib.GetPrivateIpAddress()
	if err != nil {
		
		c.JSON(http.StatusInternalServerError, err)
	}

	var service lib.ExtrenalService
	err = json.NewDecoder(c.Request.Body).Decode(&service)
	
	if err != nil {
		
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	var response *http.Response
	if checkSelfSignedCA {
		response, err = lib.HttpsClient(s.CA).Get(*service.Url)
	} else {
		response, err = lib.HttpClient().Get(*service.Url)
	}

	if err != nil {
		
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	var result interface{}
	err = json.NewDecoder(response.Body).Decode(&result)
	fmt.Print(result)
	if err != nil {
	
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	postResponse := lib.Response{
		ExternalServiceName: &result,
		MyIP:                &ips,
	}
	c.JSON(http.StatusOK, postResponse)

}

func NewServer() Server {

	secure := os.Getenv("SECURE")
	port := os.Getenv("PORT")
	if len(port) == 0 {
		port = "8080"
	}

	if secure == "TRUE" {
		certificate, err := base64.StdEncoding.Strict().DecodeString(os.Getenv("SSL_CERTIFICATE"))
		if err != nil {
			panic(err)
		}

		privateKey, err := base64.StdEncoding.Strict().DecodeString(os.Getenv("SSL_PRIVATE_KEY"))
		if err != nil {
			panic(err)
		}
		caCert, err := base64.StdEncoding.Strict().DecodeString(os.Getenv("SSL_CA"))
		if err != nil {
			panic(err)
		}
		return Server{
			Color:       os.Getenv("COLOR"),
			APIType:     os.Getenv("API_TYPE"),
			Certificate: string(certificate),
			PrivateKey:  string(privateKey),
			CA:          string(caCert),
			Port:        port,
			Secure:      true,
		}

	} else {
		return Server{
			Color:   os.Getenv("COLOR"),
			APIType: os.Getenv("API_TYPE"),
			Port:    port,
			Secure:  false,
		}
	}

}
