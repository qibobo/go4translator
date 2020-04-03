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

func GetCredentialFromVcapServicesByName(name string) *Credential {
	serviceClassName := "language-translator"
	cfInstanceIndex := os.Getenv("CF_INSTANCE_INDEX")
	if cfInstanceIndex != "" {
		serviceClassName = "language_translator"
	}
	vcapStr := os.Getenv("VCAP_SERVICES")
	if vcapStr == "" {
		vcapStr = jsonStr
	}
	log.Printf("----vcap_services string: %v\n", vcapStr)
	log.Printf("----name: %v\n", name)
	log.Printf("----serviceClassName: %v\n", serviceClassName)
	var vcapServices map[string][]Service
	err := json.Unmarshal([]byte(vcapStr), &vcapServices)
	if err != nil {
		return nil
	}
	for _, v := range vcapServices[serviceClassName] {
		if v.Name == name {
			return &v.Credentials
		}
	}
	return nil
}

func main() {
	log.Print("go4translator started.")

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	name := os.Getenv("TRANSLATOR_NAME")
	if name == "" {
		name = "thetranslator"
	}
	credential := GetCredentialFromVcapServicesByName(name)
	log.Printf("The credential is %v\n", credential)

	apikey := credential.ApiKey
	url := credential.URL
	log.Printf("-----------apikey: %s, url: %s\n", apikey, url)
	handler := handler.NewHandler(apikey, url)
	http.HandleFunc("/", handler.Handle)

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}

const jsonStr = `{
	"language-translator": [
	  {
		"credentials": {
		  "apikey": "xxx",
		  "iam_apikey_description": "Auto-generated for key 07f29313-3c4d-4dea-a431-f6caae6b03d9",
		  "iam_apikey_name": "ying-translator-binding",
		  "iam_role_crn": "crn:v1:bluemix:public:iam::::serviceRole:Manager",
		  "iam_serviceid_crn": "crn:v1:bluemix:public:iam-identity::a/8d63fb1cc5e99e86dd7229dddffc05a5::serviceid:ServiceId-5c954360-76cd-47c7-ad04-660782f68947",
		  "url": "https://api.us-south.language-translator.watson.cloud.ibm.com/instances/271b6ff6-4e0f-4673-a389-f7b7fbbe3afa"
		},
		"name": "thetranslator",
		"plan": "standard"
	  }
	]
  }`
