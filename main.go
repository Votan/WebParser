package main

import (
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"
	"sync"

	"golang.org/x/net/html"
)

var total int

func main() {
	var wg sync.WaitGroup

	urls := []string{
		"https://golang.org/",
		"https://golang.org/doc/",
		"https://golang.org/pkg/compress/",
		"https://golang.org/pkg/compress/gzip/",
		"https://golang.org/pkg/crypto/md5/",
	}

	for _, url := range urls {
		wg.Add(1)

		go getString(url, &wg)
	}

	wg.Wait()

	fmt.Println("Total:", total)
}

func getString(url string, wg *sync.WaitGroup) {
	defer wg.Done()

	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	body := string(bodyBytes)

	bodyWithoutTags := removeHTMLTags(body)

	pattern := "\\bGo\\b"

	// Count the occurrences of the word "Go"
	r, err := regexp.Compile(pattern)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	count := len(r.FindAllString(bodyWithoutTags, -1))

	total += count

	fmt.Printf("Count for %s: %d\n", url, count)
}

func removeHTMLTags(text string) string {
	tokenizer := html.NewTokenizer(strings.NewReader(text))
	var result strings.Builder

	for {
		tokenType := tokenizer.Next()
		if tokenType == html.ErrorToken {
			break
		}

		if tokenType == html.TextToken {
			token := tokenizer.Token()
			result.WriteString(token.Data)
		}
	}

	return result.String()
}
