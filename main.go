package main

import (
	"fmt"
	"log"

	"github.com/gocolly/colly"
)

// variables
var (
	username = "XXX" // !!! change to your e-mail address
	password = "XXX" // !!! enter your password here
)

func main() {
	// instantiate default collector
	c := colly.NewCollector()

	// authenticate
	err := c.Post("https://dzone.com/services/internal/action/users-login", map[string]string{"username": username, "password": password})
	if err != nil {
		log.Fatal(err)
	}

	// attach callbacks after login
	c.OnResponse(func(r *colly.Response) {
		log.Println("response received", r.StatusCode)
	})

	//
	var counter = 0
	c.OnHTML("div[class=asset-subtitle] a[href]", func(f *colly.HTMLElement) {
		title := f.Text
		fmt.Println("Book title is:", title)
		counter++
	})

	fmt.Println("Amount of Refcardz found:", counter)

	c.OnHTML("dz-download", func(d *colly.HTMLElement) {
		download := d.Attr("asset")
		// Download link example: https://dzone.com/asset/download/279342
		fmt.Println("Download link: " + "https://dzone.com" + download)
	})

	// start scraping Refcardz
	c.Visit("https://dzone.com/refcardz")

} // end of func main()
