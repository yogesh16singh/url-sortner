package main

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"
)

type URL struct {
	ID           string    `json:"id"`
	OriginalURL  string    `json:"original_url"`
	ShortUrl     string    `json:"short_url"`
	CreationDate time.Time `json:"creation_date"`
}

var urlDb = make(map[string]URL)

func generateShortURL(url string) string {
	fmt.Println("url", url)
	hasher := md5.New()
	hasher.Write([]byte(url))
	fmt.Println("hasher", hasher)
	data := hasher.Sum(nil)
	fmt.Println("data", data)
	hash := hex.EncodeToString(data)
	fmt.Println("hash", hash)
	fmt.Println("hash8", hash[:8])
	return hash[:8]
}

func createURL(url string) string {
	shortURL := generateShortURL(url)
	id := shortURL
	urlDb[id] = URL{ID: id, OriginalURL: url, ShortUrl: shortURL, CreationDate: time.Now()}
	return shortURL
}

func getUrl(id string) (URL, error) {
	_, ok := urlDb[id]
	if !ok {
		return URL{}, errors.New("url not found")
	}
	return urlDb[id], nil
}
func main() {
	fmt.Println("starting main")

	fmt.Println("starting the server")
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, %s!", r.URL.Path[1:])
	})

	http.HandleFunc("/shorten", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		var data struct {
			URL string `json:"url"`
		}
		json.NewDecoder(r.Body).Decode(&data)
		shortURL := createURL(data.URL) // url := ()data.URL
		response := struct {
			ShortUrl string `json:"short_url"`
		}{
			ShortUrl: shortURL,
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	})

	http.HandleFunc("/url/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		id := r.FormValue("id")
		url, err := getUrl(id)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.Write([]byte(fmt.Sprintf(`{"original_url": "%s"}`, url.OriginalURL)))
	})
	http.HandleFunc("/redirect/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		fmt.Println(r.URL.Path[len("/redirect/"):])
		id := r.URL.Path[len("/redirect/"):]
		url, err := getUrl(id)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		http.Redirect(w, r, url.OriginalURL, http.StatusFound)
	})
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err)
	}

}
