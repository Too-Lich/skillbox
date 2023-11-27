package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
)

const addr string = "localhost:8081"

func main() {
	http.HandleFunc("/", handle)
	log.Fatalln(http.ListenAndServe(addr, nil))
}

func handle(w http.ResponseWriter, r *http.Request) {

	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err)
	}
	defer r.Body.Close()

	text := string(bodyBytes)
	response := "1 instance: " + r.Method + text

	if _, err := w.Write([]byte(response + r.Host)); err != nil {
		log.Fatal(err)
	}
	fmt.Println(r.Host)
}
