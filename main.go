package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/jedib0t/go-pretty/v6/table"
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

func getArtist(service string, artistName string, baseURL string, site string) {
	i := 0

	for empty := false; !empty; {
		pagenum := i * 25
		url := fmt.Sprintf(baseURL, service, artistName, pagenum)
		fmt.Println(url)
		page := getPage(url)
		if len(page) == 0 {
			empty = true
		}
		page.getPosts(service, artistName, site)
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

func displayTable(artists Artists, service string) {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{service})
	t.AppendHeader(table.Row{"IDX", "username", "service"})
	for i, user := range artists {
		t.AppendRow(table.Row{i, user.Name, user.Service})
		t.AppendSeparator()
	}
	t.Render()
}

func main() {
	// sites := Sites{}
	// sites.getArtists()
	// getArtist("onlyfans", "belledelphine")
	cm, km := searchArtists("belledelphine")

	displayTable(cm, "Coomer")
	displayTable(km, "Kemono")

	useKemono := flag.Bool("k", false, "Search kemono.party")
	useCoomer := flag.Bool("c", false, "Search coomer.party")
	idx := flag.Int("i", 0, "index of desired artist")

	flag.Parse()

	if *useCoomer {
		getArtist(cm[*idx].Service, cm[*idx].ID, "https://coomer.party/api/%s/user/%s?o=%d", "c")
		return
	}

	if *useKemono {
		getArtist(km[*idx].Service, km[*idx].ID, "https://kemono.party/api/%s/user/%s?o=%d", "k")
	}
}
