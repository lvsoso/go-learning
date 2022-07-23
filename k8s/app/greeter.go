package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		metaName := os.Getenv("META_NAME")
		fmt.Fprint(w, fmt.Sprint("Hello, World! Here is ", metaName, "\n"))
	})
	http.HandleFunc("/_status/healthz", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "alive!")
	})

	log.Fatalln(http.ListenAndServe(":80", nil))
}
