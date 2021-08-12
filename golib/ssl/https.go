//Package ssl HTTPS/SSL (both one way and mutual)
package ssl

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"golang.org/x/crypto/pkcs12"
)

var onewayClient *http.Client
var mutualClient *http.Client

func getCACerts() *x509.CertPool {
	certlocation := os.Getenv("GLOBAL_TRUSTED_CA_STORE")
	if len(certlocation) == 0 {
		log.Println("ENV variable GLOBAL_TRUSTED_CA_STORE is empty , using current directory to search for certificates ...... ")
		certlocation = "."
	}
	files, err := ioutil.ReadDir(certlocation)
	if err != nil {
		log.Fatal(err)
	}
	caCertPool := x509.NewCertPool()
	for _, file := range files {
		filepath := []string{certlocation, file.Name()}
		fullname := strings.Join(filepath, "/")
		if strings.HasSuffix(file.Name(), "cer") {
			log.Println("Loading certitifate ", file.Name(), " .....")
			caCertbytes, err := ioutil.ReadFile(fullname)
			caCert, err := x509.ParseCertificate(caCertbytes)
			if err != nil {
				log.Fatal(err)
			}
			caCertPool.AddCert(caCert)
		} else if strings.HasSuffix(file.Name(), "pem") {
			log.Println("Loading certitifate ", file.Name(), " .....")
			caCertbytes, err := ioutil.ReadFile(fullname)
			if err != nil {
				log.Fatal(err)
			}
			caCertPool.AppendCertsFromPEM(caCertbytes)
		}
	}
	return caCertPool
}

func getClientCert(pfxFile string, password string) *tls.Certificate {
	certlocation := os.Getenv("GLOBAL_TRUSTED_CA_STORE")
	if len(certlocation) == 0 {
		log.Println("ENV variable GLOBAL_TRUSTED_CA_STORE is empty , using current directory to search for certificates ...... ")
		certlocation = "."
	}
	filepath := []string{certlocation, pfxFile}
	fullname := strings.Join(filepath, "/")
	p12file, err := ioutil.ReadFile(fullname)
	p12pemArray, err := pkcs12.ToPEM(p12file, password)
	if err != nil {
		log.Fatal(err)
	}
	encoded := pem.EncodeToMemory(p12pemArray[0])
	p12pemArray[0].Bytes = encoded
	encodedkey := pem.EncodeToMemory(p12pemArray[1])
	p12pemArray[1].Bytes = encodedkey
	clientcert, err := tls.X509KeyPair(encoded, encodedkey)
	log.Println("Loading certitifate ", pfxFile, " .....")
	if err != nil {
		log.Fatal(err)
	}
	return &clientcert
}

//GetOneWaySSLClient - create one way SSL client
func GetOneWaySSLClient() *http.Client {
	if onewayClient != nil {
		return onewayClient
	}
	caCertPool := getCACerts()
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				RootCAs: caCertPool,
			},
		},
	}
	onewayClient = client
	return client
}

//GetMutualSSLClient - create mutual SSL client with given pfx fime and password
func GetMutualSSLClient(pfxFileName string, pfxPassword string) *http.Client {
	if mutualClient != nil {
		return mutualClient
	}
	caCertPool := getCACerts()
	clientcert := getClientCert(pfxFileName, pfxPassword)
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				RootCAs:      caCertPool,
				Certificates: []tls.Certificate{*clientcert},
			},
		},
	}
	mutualClient = client
	return client
}

func example() {
	client := GetMutualSSLClient("y0319t1783.p12", "TDPCmp7ynG23Cv7oxLLYW")
	var jsonStr = []byte(`
		{
		  "EventType":"Effective",
		  "Location":"835",
		  "AsOf":"2017-07-26T00:11:05",
		  "Price":[
			{
			  "RmsSkuId":"56650527",
			  "CurrentPrice":{
				"PriceType":"R",
				"Price":500,
				"Reason":"Error in Retail"
			  },
			  "OwnershipPrice":{
				"PriceType":"R",
				"Price":500,
				"Reason":"Error in Retail"
			  },
			  "RegularPrices":[
				{
				  "PriceTypeId":"R-1501367",
				  "PriceType":"R",
				  "Price":500,
				  "StartDateTime":"2017-06-08T07:00:00",
				  "Reason":"Error in Retail"
				}
			  ]
			}
		  ],
		  "PublishedDateTime":"2017-07-26T00:11:05.422"
		}`)
	req, err := http.NewRequest("POST", "https://sxgtest.nordstrom.net:4439/NLS/api/PoAllocation", bytes.NewBuffer(jsonStr))
	// ...
	req.Header.Add("jwn-event-creation-time", "2017-07-26T00,11,05.403Z")
	req.Header.Add("jwn-event-subscription-id", "3b3711e5-680b-43be-8016-de40769387f2")
	req.Header.Add("nordapiversion", "1.0")
	req.Header.Add("jwn-event-name", "Price.Effective")
	req.Header.Add("jwn-subscription-id", "e29b8575-4607-4fe8-ba5d-99e3b3ed41a7")
	req.Header.Add("correlationid", "a8403636-3383-4e2c-b84d-6d8df80a0bcc")
	req.Header.Add("jwn-event-id", "a8403636-3383-4e2c-b84d-6d8df80a0bcc")
	req.Header.Add("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(resp)
	body, err := ioutil.ReadAll(resp.Body)
	fmt.Println("get:\n", string(body))
}
