package main

import (
	"fmt"
	"log"
	"time"

	"github.com/gocolly/colly"
)

var thCsrf string

func main() {
	c := colly.NewCollector()
	c.OnHTML("form[role=form] input[type=hidden][name=TH_CSRF]", func(e *colly.HTMLElement) {
		thCsrf := e.Attr("value")
		err := c.Post("https://dzone.com/j_spring_security_check", map[string]string{"TH_CSRF": thCsrf, "_spring_security_remember_me": "true", "j_username": "dzone-refcardz@mailcatch.com", "j_password": "password123456"})
		if err != nil {
			log.Fatal(err)
		}
		return
	})
	c.Visit("https://dzone.com/users/login.html")
	// ---------------------------------------------------------------------------------------------
	d := c.Clone()
	d.SetRequestTimeout(180 * time.Second)

	d.OnResponse(func(r *colly.Response) {
		fmt.Println(r.Request.URL.String())
		r.Save("279342.pdf")
	})

	d.Visit("https://dzone.com/asset/download/279342")
}
