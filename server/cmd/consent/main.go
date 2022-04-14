package main

import (
	"fmt"
	"net/http"
	"os"
)

func main() {
	port := "8000"
	if override, ok := os.LookupEnv("PORT"); ok {
		port = override
	}
	http.ListenAndServe(fmt.Sprintf(":%s", port), http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	}))
}
