package main

import (
	"crypto/tls"
	"io"
	"net"
	"net/http"
	"strconv"
	"testing"
	"time"
)

func Test_startServer(t *testing.T) {
	var err error
	port := freeTCPPort()
	go func() {
		err = startServer(":"+strconv.Itoa(port), "")
	}()
	time.Sleep(2 * time.Second)
	if err != nil {
		t.Fatal(err)
	}
	clientCertificate, err := issueCertificate("localhost")
	if err != nil {
		t.Fatal(err)
	}
	httpClient := http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
				Certificates:       []tls.Certificate{*clientCertificate},
			},
		},
	}
	httpResponse, err := httpClient.Get("https://localhost:" + strconv.Itoa(port) + "/subpath")
	if err != nil {
		t.Fatal(err)
	}
	if httpResponse.StatusCode != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, httpResponse.StatusCode)
	}
	responseData, err := io.ReadAll(httpResponse.Body)
	if err != nil {
		t.Fatal(err)
	}
	println(string(responseData))
}

// freeTCPPort asks the kernel for a free open port that is ready to use.
// Taken from https://gist.github.com/sevkin/96bdae9274465b2d09191384f86ef39d
func freeTCPPort() (port int) {
	if a, err := net.ResolveTCPAddr("tcp", "localhost:0"); err == nil {
		var l *net.TCPListener
		if l, err = net.ListenTCP("tcp", a); err == nil {
			defer l.Close()
			return l.Addr().(*net.TCPAddr).Port
		} else {
			panic(err)
		}
	} else {
		panic(err)
	}
}
