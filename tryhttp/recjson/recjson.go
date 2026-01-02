package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

type Request struct {
	Format string `json:"format"`
	TZ     string `json:"tz"`
}

type Response struct {
	Time string `json:"time"`
}

type Error struct {
	Error string `json:"error"`
}

func ReadJSON[T any](r io.ReadCloser) (T, error) {
	var v T
	err := json.NewDecoder(r).Decode(&v)
	return v, errors.Join(err, r.Close())
}

func WriteJSON(w http.ResponseWriter, v any) error {
	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(v)
}

func WriteError(w http.ResponseWriter, err error, code int) {
	log.Printf("%d %v: %v", code, http.StatusText(code), err)
	w.Header().Set("Content-Type", "encondig/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(Error{err.Error()})
}

func getTimeV2(w http.ResponseWriter, r *http.Request) {

	req, err := ReadJSON[Request](r.Body)
	if err != nil {
		WriteError(w, err, 400)
		return
	}

	var tz *time.Location = time.Local
	if req.TZ != "" {
		var err error
		tz, err = time.LoadLocation(req.TZ)
		if err != nil || tz == nil {
			WriteError(w, err, 400)
			return
		}
	}

	format := time.RFC3339
	if req.Format != "" {
		format = req.Format
	}

	WriteJSON(w, Response{time.Now().In(tz).Format(format)})

}

func getTime(w http.ResponseWriter, r *http.Request) {
	var req Request

	w.Header().Set("Content-Type", "encoding/json")
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(400)
		json.NewEncoder(w).Encode(Error{err.Error()})
		return
	}

	r.Body.Close()

	var tz *time.Location = time.Local
	if req.TZ != "" {
		var err error
		tz, err = time.LoadLocation(req.TZ)
		if err != nil || tz == nil {
			w.WriteHeader(400)
			json.NewEncoder(w).Encode(Error{err.Error()})
			return
		}
	}

	format := time.RFC3339
	if req.Format != "" {
		format = req.Format
	}

	resp := Response{time.Now().In(tz).Format(format)}
	json.NewEncoder(w).Encode(resp)

}

var client = &http.Client{Timeout: 2 * time.Second}

func sendRequest(tz, format string) {
	body := new(bytes.Buffer)
	json.NewEncoder(body).Encode(Request{TZ: tz, Format: format})
	log.Printf("request body: %v", body)
	req, err := http.NewRequestWithContext(context.TODO(), "GET", "http://localhost:8080", body)
	if err != nil {
		panic(err)
	}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	resp.Write(os.Stdout)
	resp.Body.Close()
}

func main() {
	server := http.Server{Addr: ":8080", Handler: http.HandlerFunc(getTimeV2)}
	go server.ListenAndServe()

	sendRequest("", "")
	sendRequest("America/Los_Angeles", time.RFC3339)
	sendRequest("Asia/Bangkok", time.RFC822Z)
	sendRequest("faketz", "")
}
