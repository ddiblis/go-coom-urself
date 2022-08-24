package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"regexp"
	"strings"
)

type Artists []struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Service string `json:"service"`
}

type Sites struct {
	Coomer Artists `json:"coomer_artists"`
	Kemono Artists `json:"kemono_artists"`
}

func (s Sites) getArtists() error {
	coomerList := Artists{}
	kemonoList := Artists{}
	sites := Sites{}
	json.NewDecoder(getBody("https://coomer.party/api/creators")).Decode(&coomerList)
	sites.Coomer = coomerList
	json.NewDecoder(getBody("https://kemono.party/api/creators")).Decode(&kemonoList)
	sites.Kemono = kemonoList
	file, _ := json.MarshalIndent(sites, "", " ")
	return ioutil.WriteFile("creators.json", file, 0644)
}

func comparinator(a string, b string) bool {
	reg, err := regexp.Compile("[^a-zA-Z0-9]+")
	if err != nil {
		log.Fatal(err)
	}
	pa := strings.ToLower(reg.ReplaceAllString(a, ""))
	pb := strings.ToLower(reg.ReplaceAllString(b, ""))
	return strings.Contains(pa, pb)
}

func searchArtists(input string) (Artists, Artists) {
	file, _ := ioutil.ReadFile("creators.json")
	data := Sites{}
	coomerMatches := Artists{}
	kemonoMatches := Artists{}

	_ = json.Unmarshal([]byte(file), &data)

	for _, user := range data.Coomer {
		if comparinator(user.Name, input) {
			coomerMatches = append(coomerMatches, user)
		}
	}

	for _, user := range data.Kemono {
		if comparinator(user.Name, input) {
			kemonoMatches = append(kemonoMatches, user)
		}
	}

	return coomerMatches, kemonoMatches
}
