package remote

import (
	"log"
	"net/http"
	"testing"
	"time"
)

func TestNewReader(t *testing.T) {
	NewReader(
		Retry(3),
		Timeout(time.Second),
		SkipTLSVerify(),
		UserAgent("test"))
}

func TestReader_Read(t *testing.T) {
	_, err := NewReader().Read("https://google.com")
	if err != nil {
		t.Error(err)
	}
}

func TestReader_Bytes(t *testing.T) {
	content, err := NewReader().Bytes("https://google.com")
	if err != nil {
		t.Error(err)
	}
	if len(content) == 0 {
		t.Error("Bytes return empty content")
	}
}

func TestReader_JSON(t *testing.T) {
	// start a small http server
	http.HandleFunc("/json/valid", func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte("{\"content\": 1}"))
		if err != nil {
			log.Fatal("fail to write back")
		}
	})
	http.HandleFunc("/json/invalid", func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte("boom"))
		if err != nil {
			log.Fatal("fail to write back")
		}
	})
	go func(t *testing.T) {
		log.Fatal(http.ListenAndServe(":8080", nil))
	}(t)

	// try to read valid and invalid json
	url := "http://localhost:8080/json"
	type testData struct {Content int `json:"content"`}
	result := &testData{}
	if err := NewReader().JSON(url + "/invalid",  &testData{}); err == nil {
		t.Error("read invalid json response")
	}
	if err := NewReader().JSON(url + "/valid", result); err != nil {
		t.Error("fail to read json response")
	}
	if result.Content != 1 {
		t.Error("invalid result Json", result.Content)
	}
}
