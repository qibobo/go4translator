package handler

import (
	"fmt"
	"log"
	"net/http"

	"github.com/IBM/go-sdk-core/core"
	"github.com/watson-developer-cloud/go-sdk/languagetranslatorv3"
)

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
	log.Printf("translateOptions is : %v\n", translateOptions)
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
