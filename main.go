// package main

// import (
// 	"fmt"
// 	"net/http"
// 	"os"
// )

// func handler(w http.ResponseWriter, r *http.Request) {
// 	fmt.Fprintf(w, "Challenge4 is running 🚀")
// }

// func main() {
// 	port := os.Getenv("PORT")
// 	if port == "" {
// 		port = "8080"
// 	}

// 	http.HandleFunc("/", handler)

// 	fmt.Println("Running on port", port)
// 	http.ListenAndServe(":"+port, nil)
// }


package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
)

type homeResponse struct {
	Name        string
	Description string
	Url         string
}

type statusResponse struct {
	Status string `json:"status"`
}

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, homeResponse{
			Name:        "Hello",
			Description: "World",
			Url:         r.Host,
		})
	})

	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, statusResponse{Status: "ok"})
	})

	addr := ":4444"
	if value := os.Getenv("PORT"); value != "" {
		addr = ":" + value
	}

	log.Printf("listening on %s", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatal(err)
	}
}

func writeJSON(w http.ResponseWriter, statusCode int, payload any) {
	body, err := json.Marshal(payload)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		log.Printf("encode response: %v", err)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(statusCode)

	if _, err := w.Write(body); err != nil {
		log.Printf("write response: %v", err)
	}
}