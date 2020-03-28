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
	vcapStr := os.Getenv("VCAP_SERVICES")
	if vcapStr == "" {
		vcapStr = jsonStr
	}
	var vcapServices map[string][]Service
	err := json.Unmarshal([]byte(vcapStr), &vcapServices)
	if err != nil {
		return nil
	}
	for _, v := range vcapServices["language-translator"] {
		if v.Name == name {
			return &v.Credentials
		}
	}
	return nil
}

func main() {
	log.Print("Hello world sample started.")

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	name := os.Getenv("TRANSLATOR_NAME")
	if name == "" {
		name = "thetranslator"
	}
	credential := GetCredentialFromVcapServicesByName(name)
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
		  "iam_apikey_description": "Auto-generated for key dca63618-4af0-416c-9508-03752e10e242",
		  "iam_apikey_name": "thetranslator",
		  "iam_role_crn": "crn:v1:bluemix:public:iam::::serviceRole:Manager",
		  "iam_serviceid_crn": "crn:v1:bluemix:public:iam-identity::a/ef6a34810cbcd892507d3ebe01e3d95a::serviceid:ServiceId-5bfd5402-0655-4592-a25b-9012f0335d23",
		  "url": "https://api.us-south.language-translator.watson.cloud.ibm.com/instances/c7a49da9-c44d-4dd4-a0f3-8714fb15f4c0"
		},
		"name": "thetranslator",
		"plan": "standard"
	  }
	],
	"cloud-object-storage": [
	  {
		"credentials": {
		  "apikey": "xxx",
		  "cos_hmac_keys": "{\"access_key_id\":\"xxx\",\"secret_access_key\":\"yyy\"}",
		  "endpoints": "https://control.cloud-object-storage.cloud.ibm.com/v2/endpoints",
		  "iam_apikey_description": "Auto-generated for key ff229a84-ac88-4290-8e66-568859cf84af",
		  "iam_apikey_name": "thecos",
		  "iam_role_crn": "crn:v1:bluemix:public:iam::::serviceRole:Manager",
		  "iam_serviceid_crn": "crn:v1:bluemix:public:iam-identity::a/ef6a34810cbcd892507d3ebe01e3d95a::serviceid:ServiceId-1cfcbb24-7bfe-407e-b001-4045188136a2",
		  "resource_instance_id": "crn:v1:bluemix:public:cloud-object-storage:global:a/ef6a34810cbcd892507d3ebe01e3d95a:a1016433-81fe-47a0-a043-0396321e73e8::"
		},
		"name": "thecos",
		"plan": "lite"
	  }
	]
  }`
