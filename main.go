package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/urfave/cli/v2"
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
	var name string
	var cm, km Artists

	app := &cli.App{
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:    "coomer",
				Aliases: []string{"c"},
				Value:   false,
				Usage:   "search coomer.party",
			},
			&cli.BoolFlag{
				Name:    "kemono",
				Aliases: []string{"k"},
				Value:   false,
				Usage:   "search kemono.party",
			},
			&cli.StringFlag{
				Name:        "name",
				Aliases:     []string{"n"},
				Usage:       "name of the artist",
				Destination: &name,
			},
		},
		Commands: []*cli.Command{
			{
				Name:    "download",
				Aliases: []string{"d"},
				Usage:   "Download all media of certain artist",
				Action: func(cli *cli.Context) error {
					cm, km := searchArtists(name)

					displayTable(cm, "Coomer")
					displayTable(km, "Kemono")

					var i int
					_, err := fmt.Scanf("%d", &i)
					if err != nil {
						log.Fatal("err")
					}

					if cli.FlagNames()[0] == "c" {
						getArtist(cm[i].Service, cm[i].ID, "https://coomer.party/api/%s/user/%s?o=%d", "c")
					} else if cli.FlagNames()[0] == "k" {
						getArtist(km[i].Service, km[i].ID, "https://kemono.party/api/%s/user/%s?o=%d", "k")
					}

					return nil
				},
			},
			{
				Name:    "update",
				Aliases: []string{"u"},
				Usage:   "Update artist list",
				Action: func(*cli.Context) error {
					fmt.Println("updating artist list")
					sites := Sites{}
					sites.getArtists()
					return nil
				},
			},
			{
				Name:    "table",
				Aliases: []string{"t"},
				Usage:   "Search artists list for artist Name",
				Action: func(cli *cli.Context) error {
					cm, km = searchArtists(cli.Args().First())

					displayTable(cm, "Coomer")
					displayTable(km, "Kemono")
					return nil
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

// {
// 	Name:    "coomer",
// 	Aliases: []string{"c"},
// 	Usage:   "Search coomer.party",
// 	Action: func(cli *cli.Context) error {
// 		// intVar, _ := strconv.Atoi(cli.Args().Slice()[1])
// 		// cm, _ := searchArtists(cli.Args().First())
// 		// getArtist(cm[intVar].Service, cm[intVar].ID, "https://coomer.party/api/%s/user/%s?o=%d", "c")
// 		return nil
// 	},
// },
// {
// 	Name:    "kemono",
// 	Aliases: []string{"k"},
// 	Usage:   "Search kemono.party",
// 	Action: func(cli *cli.Context) error {
// 		intVar, _ := strconv.Atoi(cli.Args().Slice()[1])
// 		_, km := searchArtists(cli.Args().First())
// 		getArtist(km[intVar].Service, km[intVar].ID, "https://kemono.party/api/%s/user/%s?o=%d", "k")
// 		return nil
// 	},
// },
