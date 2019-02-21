package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/gocolly/colly"
)

// constants that do not change
const (
	username = "dzone-refcardz@mailcatch.com" // created just for this purpose
	password = "password123456"               // created just for this purpose
)

// RefcardzData is an exported title and pdf link
type RefcardzData struct {
	Result struct {
		Data struct {
			Assets []struct {
				Title string //`json:"title"`
				Pdf   string //`json:"pdf"`
			} `json:"assets"`
		} `json:"data"`
	} `json:"result"`
}

// global variables
var data RefcardzData

func main() { // main function start

	// create folder for downloads
	os.Mkdir("downloads", 0700)

	// instantiate first collector for log in
	f := colly.NewCollector()

	// instantiate second collector for assets-list
	g := colly.NewCollector()

	// null the counters/pagers
	page := 0
	counter := 1

	// authenticate, otherwise you won't be able to download refcardz later
	err := f.Post("https://dzone.com/services/internal/action/dzoneUsers-validateCredentials", map[string]string{"username": username, "password": password})
	if err != nil {
		log.Fatal(err)
	}

	// visit the users-login page
	f.Visit("https://dzone.com/services/internal/action/users-login")

	// keep visiting the assets-list websites until the specific response indicating last page is returned
	for {

		// is there any error in the response?
		g.OnError(func(r *colly.Response, e error) {
			log.Println("Error:", e, r.Request.URL, string(r.Body))
		})

		// perform the following block on each assets-page response
		g.OnResponse(func(r *colly.Response) {

			// unmarshal JSON based on the struct
			err := json.Unmarshal([]byte(r.Body), &data)
			if err != nil { // eventually error out
				log.Fatal(err)
			}

			// when the response is the last page, stop continuing with other pages and exit
			if string(r.Body) == `{"success":true,"result":{"data":{"assets":[],"sort":"downloads"}},"status":200}` {
				fmt.Println("Last page reached")
				defer fmt.Println("!")
				os.Exit(0)
			}

			// naked return; returns the current values in the return arguments local variables
			return

		}) // end of response block

		// clone the logging collector
		h := f.Clone()

		// loop through the data assets containing Titles and incomplete Pdf links
		for _, obj := range data.Result.Data.Assets {
			filename := obj.Title                                            // title is going to be filename
			link := "https://dzone.com" + obj.Pdf                            // complete the HTTP link
			fmt.Println("#", counter, "Downloading", filename, "from", link) // show verbose progress

			h.OnResponse(func(q *colly.Response) {
				// replace the " " with a "_" in the filename and add extension
				result := strings.Replace("downloads/"+filename+".pdf", " ", "_", -1)
				q.Save(result) // save
			})
			// visit
			h.Visit(link)
			counter++
		}

		// next page
		page++

		// visitassets-list website construction
		g.Visit("https://dzone.com/services/widget/assets-listV2/DEFAULT?hidefeat=true&page=" + strconv.Itoa(page) + "&sort=downloads&type=refcard")

	} // end of the for loop

} // end of the main function
