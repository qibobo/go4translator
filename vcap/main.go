package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.ibm.com/qiyangbj/go4translator/handler"
)

type Credential struct {
	ApiKey string `json:"apikey"`
	URL    string `json:"url"`
}
type Service struct {
	Name        string     `json:"name"`
	Plan        string     `json:"plan"`
	Credentials Credential `json:"credentials"`
}

func GetCredentialFromVcapServices() *Credential {
	serviceClassName := "language-translator"
	cfInstanceIndex := os.Getenv("CF_INSTANCE_INDEX")
	if cfInstanceIndex != "" {
		serviceClassName = "language_translator"
	}
	vcapStr := os.Getenv("VCAP_SERVICES")
	log.Printf("----vcap_services string: %v\n", vcapStr)
	log.Printf("----serviceClassName: %v\n", serviceClassName)
	var vcapServices map[string][]Service
	err := json.Unmarshal([]byte(vcapStr), &vcapServices)
	if err != nil {
		return nil
	}
	if vcapServices[serviceClassName] != nil && len(vcapServices[serviceClassName]) > 0 {
		return &(vcapServices[serviceClassName][0].Credentials)
	}
	return nil
}

func main() {
	log.Print("go4translator started.")

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	credential := GetCredentialFromVcapServices()
	log.Printf("The credential is %v\n", credential)

	apikey := credential.ApiKey
	url := credential.URL
	log.Printf("----apikey: %s, url: %s\n", apikey, url)
	handler := handler.NewHandler(apikey, url)
	http.HandleFunc("/", handler.Handle)

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}
