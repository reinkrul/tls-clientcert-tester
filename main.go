package main

import (
	"bufio"
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"math/big"
	"net"
	"net/http"
	"os"
	"strings"
	"sync"
)

// main starts the application. Arguments: <port>, <PEM file of trusted certs>
func main() {
	if len(os.Args) < 2 {
		panic("Usage: main <addr> <optional PEM file with server certificate>")
	}
	pemFile := ""
	if len(os.Args) > 2 {
		pemFile = os.Args[2]
	}
	println("Starting server on " + os.Args[1])
	err := startServer(os.Args[1], pemFile)
	if err != nil {
		panic(err)
	}
}

func startServer(addr string, pemFile string) error {
	if !strings.Contains(addr, ":") {
		addr = ":" + addr
	}
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	defer ln.Close()
	http.DefaultServeMux.HandleFunc("GET /", func(writer http.ResponseWriter, request *http.Request) {
		writer.WriteHeader(200)
		_, _ = writer.Write([]byte("OK"))
	})

	var cert *tls.Certificate
	if pemFile == "" {
		hostname, err := os.Hostname()
		if err != nil {
			return err
		}
		if cert, err = issueCertificate(hostname); err != nil {
			return err
		}
	} else {
		lCert, err := tls.LoadX509KeyPair(pemFile, pemFile)
		if err != nil {
			return err
		}
		cert = &lCert
	}

	printLock := &sync.Mutex{}
	stdoutWriter := bufio.NewWriter(os.Stdout)
	echo := func(msg string, args ...interface{}) {
		printLock.Lock()
		defer printLock.Unlock()
		_, _ = stdoutWriter.Write([]byte(fmt.Sprintf(msg, args...) + "\n"))
		_ = stdoutWriter.Flush()
	}

	srv := http.Server{
		Addr: addr,
		TLSConfig: &tls.Config{
			ClientAuth:   tls.RequireAnyClientCert,
			Certificates: []tls.Certificate{*cert},
			VerifyConnection: func(state tls.ConnectionState) error {
				var chainPEM []string
				for _, cert := range state.PeerCertificates {
					chainPEM = append(chainPEM, pemEncode(cert.Raw))
				}
				echo("client cert issuer: %s, subject: %s, chain (leaf first):\n%s", state.PeerCertificates[0].Issuer.String(), state.PeerCertificates[0].Subject.String(), strings.Join(chainPEM, "\n"))
				return nil
			},
		},
	}

	return srv.ServeTLS(ln, "", "")
}

func pemEncode(derBytes []byte) string {
	return string(pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: derBytes}))
}

func issueCertificate(hostname string) (*tls.Certificate, error) {
	priv, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, fmt.Errorf("failed to generate private key: %w", err)
	}
	// self-sign certificate
	template := &x509.Certificate{
		Subject: pkix.Name{
			CommonName: hostname,
		},
		DNSNames:     []string{hostname},
		SerialNumber: big.NewInt(1),
	}
	derBytes, err := x509.CreateCertificate(rand.Reader, template, template, priv.Public(), priv)
	if err != nil {
		return nil, fmt.Errorf("failed to create certificate: %w", err)
	}
	certificate, err := x509.ParseCertificate(derBytes)
	if err != nil {
		return nil, err
	}
	return &tls.Certificate{
		Certificate: [][]byte{derBytes},
		PrivateKey:  priv,
		Leaf:        certificate,
	}, nil
}
