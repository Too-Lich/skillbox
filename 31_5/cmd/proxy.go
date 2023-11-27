package main

import (
	"bytes"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"io"
	"log"
	"net/http"
)

const proxyAddr string = "localhost:9000"

var (
	counter            int    = 0
	firstInstanceHost  string = "http://localhost:8080"
	secondInstanceHost string = "http://localhost:8082"
	//customTransport           = http.DefaultTransport
)

func handlerProxy(w http.ResponseWriter, r *http.Request) {

	textBytes, err := io.ReadAll(r.Body)
	if err != nil {
		log.Fatalln(err)
	}
	//w.Write(textBytes)
	text := string(textBytes)

	if counter == 0 {
		counter++
		resp, err := http.Post(firstInstanceHost+r.URL.Path, "text/json", bytes.NewBuffer([]byte(text)))
		if err != nil {
			log.Fatal(err)
		}

		textBytes, err = io.ReadAll(resp.Body)
		if err != nil {
			log.Fatalln(err)
		}
		defer resp.Body.Close()

		fmt.Print(string(textBytes))
		w.Write(textBytes)

		return
	}

	if counter == 1 {
		counter--
		resp, err := http.Post(secondInstanceHost+r.URL.Path, "text/json", bytes.NewBuffer([]byte(text)))
		if err != nil {
			log.Fatal(err)
		}

		textBytes, err = io.ReadAll(resp.Body)
		if err != nil {
			log.Fatalln(err)
		}
		defer resp.Body.Close()

		fmt.Print(string(textBytes))
		w.Write(textBytes)

		return

	}

}

func main() {
	fmt.Println("\033[32m<<-- Proxy. Перенаправление соединений -->>\033[0m")
	fmt.Println("proxy: слушаю порт ", proxyAddr)

	var r = chi.NewRouter()
	r.Use(middleware.Logger)
	http.HandleFunc("/", handlerProxy)
	log.Fatalln(http.ListenAndServe(proxyAddr, nil))

}
