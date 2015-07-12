package main

import (
	"fmt"
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
		case "RecordScore":
			RecordScoreHandler(w, r)
		case "GetScore":
			GetScoreHandler(w, r)
		case "GetLiteral":
			GetLiteralHandler(w, r)
		}
	}
}

func GetNextQuestion(res *alexa.EchoResponse) (string, bool) {
	_, ok := res.SessionAttributes["game"]
	if !ok {
		return "game", true
	}

	_, ok = res.SessionAttributes["score"]
	if !ok {
		return "score", true
	}

	_, ok = res.SessionAttributes["players"]
	if !ok {
		return "players", true
	}

	return "", false
}

func AskNextQuestion(w http.ResponseWriter, req *alexa.EchoRequest, res *alexa.EchoResponse) {
	prompt, complete := GetNextQuestion(res)

	if !complete {
		for _, prop := range []string{"game", "score", "players"} {
			val, ok := req.Session.Attributes.String[prop]
			if !ok {
				res.SessionAttributes[prop] = val
			}
		}
		res.SessionAttributes["prompt"] = prompt
		res.EndSession(false)
		switch {
		case prompt == "game":
			res.OutputSpeech("What game did you play?")
		case prompt == "score":
			res.OutputSpeech("What score did you get?")
		case prompt == "score":
			res.OutputSpeech("Who got that score?")
		}
	} else {
		res.OutputSpeech("Score recorded!")
		content := fmt.Sprintf("%v scored %v while playing %v", res.SessionAttributes["players"], res.SessionAttributes["score"], res.SessionAttributes["game"])
		res.Card("Score recorded", content)
		res.EndSession(true)
	}

	json, _ := res.ToJSON()
	w.Header().Set("Content-Type", "application/json;charset=UTF-8")
	w.Write(json)
}

func RecordScoreHandler(w http.ResponseWriter, r *http.Request) {
	req := alexa.GetEchoRequest(r)
	res := alexa.NewResponse()

	game, err := req.GetSlotValue("Game")
	if err == nil {
		res.SessionAttributes["game"] = game
	}

	score, err := req.GetSlotValue("Score")
	if err == nil {
		res.SessionAttributes["score"] = score
	}

	AskNextQuestion(w, req, res)
}

func GetScoreHandler(w http.ResponseWriter, r *http.Request) {
	req := alexa.GetEchoRequest(r)
	res := alexa.NewResponse()

	prompt, ok := req.Session.Attributes.String["prompt"]
	if !ok || prompt != "score" {
		res.OutputSpeech("That's not what I was expecting")
		res.EndSession(true)
		json, _ := res.ToJSON()
		w.Header().Set("Content-Type", "application/json;charset=UTF-8")
		w.Write(json)
		return
	}

	score, _ := req.GetSlotValue("Score")
	res.SessionAttributes["score"] = score

	AskNextQuestion(w, req, res)
}

func GetLiteralHandler(w http.ResponseWriter, r *http.Request) {
	req := alexa.GetEchoRequest(r)
	res := alexa.NewResponse()

	prompt, ok := req.Session.Attributes.String["prompt"].(string)
	if !ok {
		res.OutputSpeech("That's not what I was expecting")
		res.EndSession(true)
		json, _ := res.ToJSON()
		w.Header().Set("Content-Type", "application/json;charset=UTF-8")
		w.Write(json)
		return
	}

	literal, _ := req.GetSlotValue("Literal")
	res.SessionAttributes[prompt] = literal

	AskNextQuestion(w, req, res)
}
