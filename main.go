package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/mjwood10/avwx"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s [port number]\n", os.Args[0])
		os.Exit(1)
	}

	port := os.Args[1]

	http.HandleFunc("/metar/", handleMetar)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func handleMetar(w http.ResponseWriter, r *http.Request) {
	tokens := strings.Split(r.URL.Path, "/")
	arg := tokens[len(tokens)-1]
	icao, err := avwx.FormatICAO(arg)
	if err != nil {
		var errorResp ErrorResponse
		errorResp.Error = fmt.Sprintf("%v", err)
		b, _ := json.Marshal(errorResp)
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "%s", string(b))
		return
	}

	resp := avwx.FetchMetar(icao)
	b, _ := json.Marshal(resp)

	fmt.Fprintf(w, "%s", string(b))
}

type ErrorResponse struct {
	Error string `json:"error"`
}
