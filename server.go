package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

func helloHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, mTLS!")
}

func main() {
	// Load server's certificate and private key
	serverCert, err := tls.LoadX509KeyPair("/server.crt", "/server.key")
	if err != nil {
		log.Fatalf("server: loadkeys: %s", err)
	}

	// Load CA certificate
	caCert, err := ioutil.ReadFile("ca.crt")
	if err != nil {
		log.Fatalf("server: readca: %s", err)
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	// Setup mTLS configuration
	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{serverCert},
		ClientCAs:    caCertPool,
		ClientAuth:   tls.RequireAndVerifyClientCert,
	}

	server := &http.Server{
		Addr:      ":8443",
		TLSConfig: tlsConfig,
		Handler:   http.HandlerFunc(helloHandler),
	}

	log.Println("server: listening on https://localhost:8443")
	log.Fatal(server.ListenAndServeTLS("", "")) // Cert and key already loaded via tlsConfig
}
