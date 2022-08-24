package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

var client http.Client

func getBody(link string) io.ReadCloser {
	req, err := http.NewRequest("GET", link, nil)
	if err != nil {
		log.Fatal(err)
	}

	res, err := client.Do(req)
	if err != nil {
		log.Fatalf("Status code error: %d: %s", res.StatusCode, res.Status)
	}

	return res.Body
}

func getPage(link string) Page {
	page := Page{}
	json.NewDecoder(getBody(link)).Decode(&page)
	return page
}

func getArtist(service string, artistName string) {
	i := 0

	for empty := false; !empty; {
		pagenum := i * 25
		url := fmt.Sprintf("https://coomer.party/api/%s/user/%s?o=%d", service, artistName, pagenum)
		fmt.Println(url)
		page := getPage(url)
		if len(page) == 0 {
			empty = true
		}
		page.getPosts(service, artistName)
		i += 1
	}
}

func downloadContent(path string, text string) {
	file, err := os.Create(fmt.Sprintf("%s/post.txt", path))
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	_, err = file.WriteString(text)
	if err != nil {
		log.Fatal(err)
	}
}

func downloadAttachments(filename string, url string) {
	file, err := os.Create(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	_, err = io.Copy(file, getBody(url))
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	// sites := Sites{}
	// sites.getArtists()
	// getArtist("onlyfans", "belledelphine")
	cm, km := searchArtists("belle delphine")
	fmt.Println(cm)
	fmt.Println(km)
}
