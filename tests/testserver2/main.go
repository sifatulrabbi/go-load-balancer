package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/*", func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		b, _ := json.Marshal(map[string]any{
			"message":    "API 2 is up and running",
			"statusCode": 200,
			"success":    true,
		})

		w.WriteHeader(200)
		w.Write(b)
	})

	if err := http.ListenAndServe("127.0.0.1:5002", mux); err != nil {
		log.Panicln(err)
	}
}
