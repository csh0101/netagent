package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/csh0101/netagent.git/demo/mtls_v2/constant"
)

func main() {
	clientCert, err := tls.LoadX509KeyPair(constant.ClientPem, constant.ClientKey)
	if err != nil {
		fmt.Println("Error loading client certificate:", err)
		return
	}

	caCert, err := ioutil.ReadFile(constant.CA)
	if err != nil {
		fmt.Println("Error loading CA certificate:", err)
		return
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	config := &tls.Config{
		Certificates: []tls.Certificate{clientCert},
		RootCAs:      caCertPool,
	}

	tr := &http.Transport{
		TLSClientConfig: config,
	}
	client := &http.Client{Transport: tr}

	resp, err := client.Get("https://localhost:8443")
	if err != nil {
		fmt.Println("Error connecting to server:", err)
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return
	}

	fmt.Println("Response:", string(body))
}
