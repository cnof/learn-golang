// HTTPS 双向认证时使用 client.crt 及 client.key 请求接口

package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

func main() {

}

func request(url string) {
	pool := x509.NewCertPool()
	caCrt, err := ioutil.ReadFile("client.crt")
	if err != nil {
		log.Fatal("read ca.crt file error:", err.Error())
	}
	pool.AppendCertsFromPEM(caCrt)
	cliCrt, err := tls.LoadX509KeyPair("client.crt", "client.key")
	if err != nil {
		log.Fatalln("LoadX509KeyPair error:", err.Error())
	}
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{
			ClientCAs:          pool,
			ClientAuth:         tls.RequireAndVerifyClientCert,
			MinVersion:         tls.VersionTLS12,
			InsecureSkipVerify: true,
			RootCAs:            pool,
			Certificates:       []tls.Certificate{cliCrt},
		},
	}

	client := &http.Client{
		Transport: tr,
	}
	req, _ := http.NewRequest("GET", url, nil)

	res, err := client.Do(req)
	if err != nil {
		log.Fatalln(err)
	}

	defer func() {
		if res.Body != nil {
			err = res.Body.Close()
		}
	}()
	//var ret map[string]interface{}
	bytes, _ := ioutil.ReadAll(res.Body)
	fmt.Println(string(bytes))
}
