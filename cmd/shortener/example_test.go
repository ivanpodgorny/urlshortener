package main

import (
	"fmt"
	"log"
	"net/http/cookiejar"
	"time"

	"github.com/imroc/req/v3"
)

func Example() {
	go func() {
		log.Fatal(Execute())
	}()
	time.Sleep(1000 * time.Millisecond)

	jar, _ := cookiejar.New(nil)
	c := req.C().
		SetCookieJar(jar).
		SetRedirectPolicy(req.NoRedirectPolicy())

	// Создание сокращенного URL.
	respCreate, _ := c.R().
		SetBodyString("https://ya.ru").
		Post("http://localhost:8080/")

	// Вызов GET /{id}. Происходит редирект на оригинальный URL.
	respGet, _ := c.R().Get(respCreate.String())

	fmt.Println(respGet.Header.Get("Location"))

	// Output:
	// Build version: N/A
	// Build date: N/A
	// Build commit: N/A
	// https://ya.ru
}

func Example_shorten() {
	reqJSON := struct {
		URL string `json:"url"`
	}{
		URL: "https://ya.ru",
	}
	respJSON := struct {
		URL string `json:"result"`
	}{}
	jar, _ := cookiejar.New(nil)
	resp, _ := req.C().SetCookieJar(jar).R().
		SetBodyJsonMarshal(&reqJSON).
		SetSuccessResult(&respJSON).
		Post("http://localhost:8080/api/shorten")

	if resp.IsSuccessState() {
		fmt.Println(respJSON.URL)
	}
}

func Example_shortenBatch() {
	reqJSON := []struct {
		ID  string `json:"correlation_id"`
		URL string `json:"original_url"`
	}{
		{
			ID:  "1",
			URL: "https://ya.ru",
		},
		{
			ID:  "2",
			URL: "https://google.com",
		},
	}
	respJSON := make([]struct {
		ID  string `json:"correlation_id"`
		URL string `json:"short_url"`
	}, 0, len(reqJSON))
	jar, _ := cookiejar.New(nil)
	resp, _ := req.C().SetCookieJar(jar).R().
		SetBodyJsonMarshal(&reqJSON).
		SetSuccessResult(&respJSON).
		Post("http://localhost:8080/api/shorten/batch")

	if resp.IsSuccessState() {
		for _, u := range respJSON {
			fmt.Printf("ID: %q. Short URL: %q.\n", u.ID, u.URL)
		}
	}
}

func Example_getUserURLs() {
	respJSON := make([]struct {
		ShortURL    string `json:"short_url"`
		OriginalURL string `json:"original_url"`
	}, 0)
	jar, _ := cookiejar.New(nil)
	resp, _ := req.C().SetCookieJar(jar).R().
		SetSuccessResult(&respJSON).
		Get("http://localhost:8080/api/user/urls")

	if resp.IsSuccessState() {
		for _, u := range respJSON {
			fmt.Printf("Short URL: %q. Original URL: %q.\n", u.ShortURL, u.OriginalURL)
		}
	}
}

func Example_deleteUserURLs() {
	reqJSON := []string{"zUTzokUGe9KDmhoL", "pMx8KJdDbnJWmMsk"}
	jar, _ := cookiejar.New(nil)
	resp, _ := req.C().SetCookieJar(jar).R().
		SetBodyJsonMarshal(&reqJSON).
		Delete("http://localhost:8080/api/user/urls")

	if resp.IsSuccessState() {
		fmt.Println("Удаление выполнено.")
	}
}
