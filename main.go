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

type VCAP_SERVICE struct {
	content map[string]map[string]interface{}
}

type Handler struct {
	service *languagetranslatorv3.LanguageTranslatorV3
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
	fmt.Fprintf(w, " %s\n", *(translateResult.Translations[0].Translation))
}

func main() {
	log.Print("Hello world sample started.")

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	vcapStr := os.Getenv("VCAP_SERVICES")
	// vcapStr = jsonStr
	var vcap map[string]map[string]map[string]map[string]string
	_ = json.Unmarshal([]byte(vcapStr), &vcap)
	// fmt.Printf("%v", vcap["mytranslator-binding"])
	secret := vcap["mytranslator-binding"]["status"]["secretName"]
	url := secret["url"]
	apikey := secret["apikey"]
	handler := NewHandler(apikey, url)
	http.HandleFunc("/", handler.Handle)

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}
