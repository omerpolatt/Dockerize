package main

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
	"strings"
)

func ReadFile(choice string) (string, error) {
	filename := "templates/stil/" + choice + ".txt"
	file, err := os.ReadFile(filename)
	if err != nil {
		return "", err
	}
	return string(file), err
}

func ReadAndFind(rows []string, word string, output http.ResponseWriter) {
	for h := 1; h < 9; h++ {
		for i := 0; i < len(word); i++ {
			for satirIndex, row := range rows {
				if satirIndex == (int(word[i])-32)*9+h {
					fmt.Fprint(output, row)
				}
			}
		}
		fmt.Fprintln(output)
	}
}

func asciiHandler(output http.ResponseWriter, request *http.Request) {
	if request.Method == http.MethodGet {
		tmpl, err := template.ParseFiles("templates/index.html")
		if err != nil {
			http.Error(output, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		tmpl.Execute(output, nil)
		return
	}

	if request.Method == http.MethodPost {
		text := request.FormValue("metin")
		choice := request.FormValue("secim")

		if choice == "" && text == "" {
			http.Error(output, "Text and Type Cannot Be Blank", http.StatusBadRequest)
			return
		}

		if text == "" {
			http.Error(output, "Text Cannot be left blank", http.StatusBadRequest)
			return
		}
		if choice == "" {
			http.Error(output, "Type Cannot be left blank", http.StatusBadRequest)
			return
		}

		for _, letter := range text {
			if letter >= 1 && letter <= 127 {
				continue
			} else {
				http.Error(output, "Invalid character in text", http.StatusBadRequest)
				return
			}
		}

		statement := strings.Split(text, "\n")

		rows, err := ReadFile(choice)
		if err != nil {
			http.Error(output, "File Read Error", http.StatusInternalServerError)
			return
		}

		for i, word := range statement {
			if word == "" {
				if i != 0 {
					fmt.Fprintln(output)
				}
				continue
			}

			ReadAndFind(strings.Split(rows, "\n"), word, output)

			if err != nil {
				http.Error(output, fmt.Sprintf("ERROR: %v", err), http.StatusInternalServerError)
				return
			}
		}
	}
}

func main() {
	http.HandleFunc("/", asciiHandler)

	fs := http.FileServer(http.Dir("./templates/static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	port := ":8081"
	fmt.Printf("The server is listening on port %s...\n", port)
	err := http.ListenAndServe(port, nil)
	if err != nil {
		fmt.Println("Server error:", err)
	}
}
