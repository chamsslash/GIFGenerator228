package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
)

type answerStruct struct {
	Data []struct {
		Images struct {
			Original struct {
				URL string `json:"url"`
			} `json:"original"`
		} `json:"images"`
	} `json:"data"`
}

func listen() {
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("Cant build server")
		log.Fatal(err)
	}

}
func register(url string) {
	http.HandleFunc("/get", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			log.Println("Bad request type")
			http.Error(w, "Bad type of request", http.StatusBadRequest)
			return
		}
		if url == "" {
			log.Println("No URL provided for the GIF.")
			http.Error(w, "No GIF found", http.StatusInternalServerError)
			return
		}

		cntx := map[string]string{"URL": url}
		tmpl, err := template.ParseFiles("templates/index.html")
		if err != nil {
			log.Println("Can't Parse Template:", err)
			http.Error(w, "Can't Parse Template", http.StatusInternalServerError)
			return
		}
		err = tmpl.Execute(w, cntx)
		if err != nil {
			log.Println("Template Execution Error:", err)
			http.Error(w, "Can't render template", http.StatusInternalServerError)
		}
	})
}
func main() {
	var stop string
	go listen()
	key := "vZSvA48xzSN6OzC4b4RhKrFkOCLFW6B7"
	query := "dog"
	url := fmt.Sprintf("https://api.giphy.com/v1/gifs/search?q=%s&api_key=%s&limit=1", query, key)
	fmt.Println(url)
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Problems with GIPHY")
		log.Fatal("Problems with GIPHY", err)

	}
	reader, readerr := io.ReadAll(resp.Body)
	defer resp.Body.Close()
	if readerr != nil {
		fmt.Println("Problems with GIPHY")
		log.Fatal("Problems with GIPHY", err)

	}

	var answer answerStruct
	jserr := json.Unmarshal(reader, &answer)
	if jserr != nil {
		log.Fatalf("Ошибка при парсинге JSON: %v", err)
	}

	if answer.Data[0].Images.Original.URL != "" {
		gif := answer.Data[0].Images.Original.URL
		register(gif)

	} else {
		log.Println("No GIF found.")
	}
	fmt.Scan(&stop)

}
