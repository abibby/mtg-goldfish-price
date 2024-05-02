package main

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"regexp"
	"strconv"
)

type Result struct {
	Name   string   `json:"name"`
	Set    string   `json:"set"`
	Prices []*Price `json:"prices"`
}

type Price struct {
	Date  string  `json:"date"`
	Price float64 `json:"price"`
}

func main() {
	http.HandleFunc("GET /price/{set}/{name}", handle)
	http.ListenAndServe(":9800", nil)
}

func handle(w http.ResponseWriter, r *http.Request) {
	resp, err := http.Get("https://www.mtggoldfish.com/" + r.URL.Path)
	if err != nil {
		sendError(w, err)
		return
	}
	defer resp.Body.Close()

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		sendError(w, err)
		return
	}

	re := regexp.MustCompile(`\s+d \+= "([^"]+)"`)

	matches := re.FindAllSubmatch(b, -1)

	prices := make([]*Price, len(matches))
	for i, match := range matches {
		line := bytes.TrimPrefix(match[1], []byte("\\n"))
		parts := bytes.Split(line, []byte(", "))

		price, err := strconv.ParseFloat(string(parts[1]), 64)
		if err != nil {
			sendError(w, err)
			return
		}
		prices[i] = &Price{
			Date:  string(parts[0]),
			Price: price,
		}
	}

	err = json.NewEncoder(w).Encode(&Result{
		Name:   r.PathValue("name"),
		Set:    r.PathValue("set"),
		Prices: prices,
	})
	if err != nil {
		log.Print(err)
	}
}

func sendError(w http.ResponseWriter, err error) {
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(err.Error()))
	}
}
