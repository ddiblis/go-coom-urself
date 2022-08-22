package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
)

var client http.Client

type Attachments []struct {
	Path string `json:"path"`
}

type Page []struct {
	Attachments Attachments `json:"attachments"`
	Content     string      `json:"content"`
	ID          string      `json:"id"`
	Published   string      `json:"published"`
	Title       string      `json:"title"`
}

type Artist struct {
	Pages []Page
}

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
	var wg sync.WaitGroup

	for empty := false; !empty; {
		page := i * 25
		url := fmt.Sprintf("https://coomer.party/api/%s/user/%s?o=%d", service, artistName, page)
		pages := getPage(url)
		if len(pages) == 0 {
			empty = true
		}
		for _, page := range pages {
			wg.Add(len(page.Attachments))
			path := fmt.Sprintf("%s/%s/%s", service, artistName, page.Published)
			os.MkdirAll(path, os.ModePerm)
			getText(path, page.Content)
			for j, img := range page.Attachments {
				go func(j int) {
					fmt.Println(img)
					defer wg.Done()
					extension := strings.Split(img.Path, ".")[1]
					imgurl := fmt.Sprintf("https://data16.coomer.party/data/%s", img.Path)
					imgpath := fmt.Sprintf("%s/%d.%s", path, j, extension)
					getPosts(imgpath, imgurl)
				}(j)
			}
			wg.Wait()
		}
		i += 1
	}
}

// file, _ := json.MarshalIndent(artist, "", " ")
// _ = ioutil.WriteFile("test.json", file, 0644)

func getText(path string, text string) {
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

func getPosts(filename string, url string) {
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
	getArtist("onlyfans", "belledelphine")
}

// func init() {
// 	jar, err := cookiejar.New(nil)
// 	if err != nil {
// 		log.Fatalf("Got error %s", err)
// 	}

// 	client = http.Client{
// 		Jar: jar,
// 	}
// }

// func GetDoc(link string) *goquery.Document {

// 	req, err := http.NewRequest("GET", link, nil)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	// req.AddCookie(cookie)

// 	res, err := client.Do(req)
// 	if err != nil {
// 		log.Fatalf("Error at request: %s", err)
// 	}

// 	defer res.Body.Close()

// 	if res.StatusCode != 200 {
// 		log.Fatalf("Status code error: %d, %s", res.StatusCode, res.Status)
// 	}

// 	doc, err := goquery.NewDocumentFromReader(res.Body)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	return doc

// }

// func site(link string) map[string]map[string]string {

// 	doc := GetDoc(link)

// 	m := make(map[string]map[string]string)

// 	doc.Find(".card__headlines").Each(func(i int, s *goquery.Selection) {
// 		info := make(map[string]string)
// 		l := s.Find("a.card__headline--long")
// 		link, _ := l.Attr("href")
// 		title := l.Text()
// 		if link != "" {
// 			info["link"] = link
// 			m[title] = info
// 		}
// 	})
// 	return m
// }

// func article(site map[string]map[string]string) {

// 	for name, link := range site {

// 		// link["test"] = "test"

// 		// site[key] = link

// 		doc := GetDoc(link["link"])

// 		doc.Find("section.js-main-content-list").Each(func(i int, s *goquery.Selection) {

// 			s.Find("div.cli-text").Each(func(j int, d *goquery.Selection) {

// 				site[name]["article"] += fmt.Sprintf("\n %s \n", d.Find("p").Text())

// 			})

// 			s.Find("figure.cli-image--header-media").Each(func(j int, d *goquery.Selection) {

// 				img, _ := d.Find("img.landscape").Attr("src")

// 				imgname := fmt.Sprintf("Image_%v", j+1)

// 				site[name][imgname] = img

// 			})

// 		})

// 		// fmt.Print(site[name])
// 		// fmt.Print("\n\n------------------------------------------------------------------------------\n\n")
// 	}

// 	fmt.Print(site)

// }

// func main() {
// 	site := site("https://www.huffpost.com/")
// 	article(site)
// }
