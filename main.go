package main

import (
	"log"

	"github.com/gocolly/colly"
)

/*
Simple workflow for testing the procedure of logging on to Dzone.com and downloading the refcardz PDF file.
*/

func main() {
	// create a new collector
	c := colly.NewCollector()

	// authenticate on validateCredentials (POST method)
	err := c.Post("https://dzone.com/services/internal/action/dzoneUsers-validateCredentials", map[string]string{"username": "dzone-refcardz@mailcatch.com", "password": "password123456"})
	if err != nil {
		log.Fatal(err)
	}

	// visit the validateCredentials
	c.Visit("https://dzone.com/services/internal/action/dzoneUsers-validateCredentials")

	// attach callbacks after login to validateCredentials
	c.OnResponse(func(r *colly.Response) {
		log.Println("response received", r.StatusCode, r.Request.URL)
	})

	// clone the "c" collector
	f := c.Clone()

	// authenticate to users-login
	err2 := f.Post("https://dzone.com/services/internal/action/users-login", map[string]string{"username": "dzone-refcardz@mailcatch.com", "password": "password123456"})
	if err2 != nil {
		log.Fatal(err)
	}

	// visit the users-login
	f.Visit("https://dzone.com/services/internal/action/users-login")

	// attach callback after login to users-login
	f.OnResponse(func(s *colly.Response) {
		log.Println("response received", s.StatusCode, s.Request.URL)
	})

	// finally clone the users-login "f" collector
	g := f.Clone()

	// visit random refcard, link that redirects to the PDF file
	g.Visit("https://dzone.com/asset/download/279342")

	// save the PDF file, which is only possible when you are logged in ("c" and "f" collectors)
	g.OnResponse(func(t *colly.Response) {
		log.Println("response received", t.StatusCode, t.Request.URL)
		t.Save("filename.pdf")
	})

}
