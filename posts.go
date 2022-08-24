package main

import (
	"fmt"
	"os"
	"strings"
	"sync"
)

type Attachments []struct {
	Path string `json:"path"`
}

func (A Attachments) getMaterial(path string) {
	var wg sync.WaitGroup
	comp := make(chan bool, 3)
	for j, img := range A {
		fmt.Println(img)
		extension := strings.Split(img.Path, ".")[1]
		imgurl := fmt.Sprintf("https://data16.coomer.party/data/%s", img.Path)
		imgpath := fmt.Sprintf("%s/%d.%s", path, j, extension)
		wg.Add(1)
		comp <- true
		go func() {
			downloadAttachments(imgpath, imgurl)
			<-comp
			wg.Done()
		}()
	}
	wg.Wait()
}

type Page []struct {
	Attachments Attachments `json:"attachments"`
	Content     string      `json:"content"`
	ID          string      `json:"id"`
	Published   string      `json:"published"`
	Title       string      `json:"title"`
}

func (p Page) getPosts(service string, artistName string) {
	for _, post := range p {
		path := fmt.Sprintf("%s/%s/%s", service, artistName, post.Published)
		os.MkdirAll(path, os.ModePerm)
		downloadContent(path, post.Content)
		post.Attachments.getMaterial(path)
	}
}
