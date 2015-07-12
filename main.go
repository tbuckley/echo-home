package main

import (
	"net/http"
	"os"

	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
	"github.com/tbuckley/go-alexa"
)

func main() {
	skill := alexa.New(os.Getenv("AMAZON_APPID"))

	router := mux.NewRouter()
	n := negroni.Classic()

	router.Handle("/echo/skill", negroni.New(
		negroni.HandlerFunc(skill.HandlerFuncWithNext),
		negroni.Wrap(new(EchoRequestHandler)),
	))

	n.UseHandler(router)
	n.Run(":8081")
}

type EchoRequestHandler int

func (h *EchoRequestHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	req := alexa.GetEchoRequest(r)

	if req.GetRequestType() == "IntentRequest" || req.GetRequestType() == "LaunchRequest" {
		switch req.GetIntentName() {
		case "AddMessage":
			AddMessageHandler(w, r)
		case "ListMessages":
			ListMessagesHandler(w, r)
		case "ClearMessages":
			ClearMessagesHandler(w, r)
		}
	}
}

func AddMessageHandler(w http.ResponseWriter, r *http.Request) {
	res := alexa.NewResponse()
	res.OutputSpeech("Hello world from my new Echo test app!")
	res.Card("Hello World", "This is a test card.")

	json, _ := res.ToJSON()
	w.Header().Set("Content-Type", "application/json;charset=UTF-8")
	w.Write(json)
}

func ListMessagesHandler(w http.ResponseWriter, r *http.Request) {
	res := alexa.NewResponse()
	res.OutputSpeech("Hello world from my new Echo test app!")
	res.Card("Hello World", "This is a test card.")

	json, _ := res.ToJSON()
	w.Header().Set("Content-Type", "application/json;charset=UTF-8")
	w.Write(json)
}

func ClearMessagesHandler(w http.ResponseWriter, r *http.Request) {
	res := alexa.NewResponse()
	res.OutputSpeech("Hello world from my new Echo test app!")
	res.Card("Hello World", "This is a test card.")

	json, _ := res.ToJSON()
	w.Header().Set("Content-Type", "application/json;charset=UTF-8")
	w.Write(json)
}
