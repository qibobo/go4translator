package main

import (
	"encoding/json"
	"fmt"
	"github.com/IBM/go-sdk-core/core"
	"github.com/watson-developer-cloud/go-sdk/languagetranslatorv3"
	"log"
	"net/http"
	"os"
)

type Secret struct {
	APIKey string `json:"apikey"`
	URL    string `json:"url"`
}
type Status struct {
	SecretName Secret `json:"secretName"`
}
type Credential struct {
	BackendServiceIdentifier string `json:"id"`
	Status                   Status `json:"status"`
}
type Handler struct {
	service *languagetranslatorv3.LanguageTranslatorV3
}

func GetVCAPSERVICES(name string) *Credential {
	vcapStr := os.Getenv("VCAP_SERVICES")
	if vcapStr == "" {
		vcapStr = jsonStr
	}
	var credentials []Credential
	err := json.Unmarshal([]byte(vcapStr), &credentials)
	if err != nil {
		return nil
	}
	for _, v := range credentials {
		if v.BackendServiceIdentifier == name {
			return &v
		}
	}
	return nil
}
func NewHandler(apiKey, url string) *Handler {
	authenticator := &core.IamAuthenticator{
		ApiKey: apiKey,
	}
	service, _ := languagetranslatorv3.
		NewLanguageTranslatorV3(&languagetranslatorv3.LanguageTranslatorV3Options{
			URL:           url,
			Version:       "2018-02-16",
			Authenticator: authenticator,
		})

	return &Handler{
		service: service,
	}
}

func (h *Handler) Handle(w http.ResponseWriter, r *http.Request) {

	/* TRANSLATE */
	model := r.URL.Query().Get("model")
	text := r.URL.Query().Get("text")
	textToTranslate := []string{
		text,
	}

	translateOptions := h.service.NewTranslateOptions(textToTranslate).
		SetModelID(model)

	// Call the languageTranslator Translate method
	translateResult, response, responseErr := h.service.Translate(translateOptions)

	// Check successful call
	if responseErr != nil {
		panic(responseErr)
	}

	fmt.Println(response)
	fmt.Println(translateResult)
	fmt.Fprintf(w, "result is: %s\n", *(translateResult.Translations[0].Translation))
}

func main() {
	log.Print("Hello world sample started.")

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// vcapStr := os.Getenv("VCAP_SERVICES")
	// if vcapStr == "" {
	// 	vcapStr = jsonStr
	// }
	// var vcap map[string]map[string]map[string]interface{}
	// _ = json.Unmarshal([]byte(vcapStr), &vcap)
	// fmt.Printf("%v", vcap)
	// secret := (vcap["mytranslator-binding"]["status"]["secretName"]) //interface
	// secrets := secret.(map[string]interface{})
	// url := (secrets["url"]).(string)
	// apikey := (secrets["apikey"]).(string)
	name := os.Getenv("TRANSLATOR_NAME")
	if name == "" {
		name = "mytranslator-binding"
	}
	credential := GetCredential(name)
	apikey := credential.Status.SecretName.APIKey
	url := credential.Status.SecretName.URL

	handler := NewHandler(apikey, url)
	http.HandleFunc("/", handler.Handle)

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}

const jsonStr = `[{"backendServiceIdentifier":"mytranslator-binding","spec":{"serviceName":"my-translator"},"status":{"instanceId":"crn:v1:bluemix:public:language-translator:us-south:a/ef6a34810cbcd892507d3ebe01e3d95a:8b9b6da2-b594-47f2-8f33-1907f6797034::","keyInstanceId":"crn:v1:bluemix:public:language-translator:us-south:a/ef6a34810cbcd892507d3ebe01e3d95a:8b9b6da2-b594-47f2-8f33-1907f6797034:resource-key:485b1173-8be1-4561-9243-817bae7f4452","message":"Online","secretName":{"apikey":"xxx","iam_apikey_description":"Auto-generated for key 485b1173-8be1-4561-9243-817bae7f4452","iam_apikey_name":"mytranslator-binding","iam_role_crn":"crn:v1:bluemix:public:iam::::serviceRole:Manager","iam_serviceid_crn":"crn:v1:bluemix:public:iam-identity::a/ef6a34810cbcd892507d3ebe01e3d95a::serviceid:ServiceId-3270a074-2c25-44d8-98c1-846cd8dcb312","url":"https://api.us-south.language-translator.watson.cloud.ibm.com/instances/8b9b6da2-b594-47f2-8f33-1907f6797034"},"state":"Online","status.secretName":{}}},{"backendServiceIdentifier":"mytranslator-binding2","spec":{"serviceName":"my-translator2"},"status":{"instanceId":"crn:v1:bluemix:public:language-translator:us-south:a/ef6a34810cbcd892507d3ebe01e3d95a:c6c9be8e-68fe-4a17-8cfe-1c2c5f69dabf::","keyInstanceId":"crn:v1:bluemix:public:language-translator:us-south:a/ef6a34810cbcd892507d3ebe01e3d95a:c6c9be8e-68fe-4a17-8cfe-1c2c5f69dabf:resource-key:00899b75-fa24-4eaa-b35d-f267c5f0c7df","message":"Online","secretName":{"apikey":"xxx","iam_apikey_description":"Auto-generated for key 00899b75-fa24-4eaa-b35d-f267c5f0c7df","iam_apikey_name":"mytranslator-binding2","iam_role_crn":"crn:v1:bluemix:public:iam::::serviceRole:Manager","iam_serviceid_crn":"crn:v1:bluemix:public:iam-identity::a/ef6a34810cbcd892507d3ebe01e3d95a::serviceid:ServiceId-74aad50c-8e8d-4255-90ae-589900cf30f4","url":"https://api.us-south.language-translator.watson.cloud.ibm.com/instances/c6c9be8e-68fe-4a17-8cfe-1c2c5f69dabf"},"state":"Online","status.secretName":{}}}]`
