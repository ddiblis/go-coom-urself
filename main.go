package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

var client http.Client

type Attachments []struct {
	Name string `json:"name"`
	Path string `json:"path"`
}

type Page []struct {
	Attachments Attachments `json:"attachments"`
	Content     string      `json:"content"`
	File        Attachments `json:"file"`
	ID          string      `json:"id"`
	Service     string      `json:"service"`
	Title       string      `json:"title"`
}

type Artist []struct {
	Links Page
}

func getData(link string) {

	req, err := http.NewRequest("GET", link, nil)
	if err != nil {
		log.Fatal(err)
	}

	res, err := client.Do(req)
	if err != nil {
		log.Fatalf("Status code error: %d: %s", res.StatusCode, res.Status)
	}

	defer res.Body.Close()

	artist := Page{}
	json.NewDecoder(res.Body).Decode(&artist)
	fmt.Print(artist)
}

func main() {
	getData("https://coomer.party/api/onlyfans/user/belledelphine")
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
