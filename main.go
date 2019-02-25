package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gocolly/colly"
)

// constants that do not change
const (
	username = "dzone-refcardz@mailcatch.com" // created just for this purpose
	password = "password123456"               // created just for this purpose
	baseURL  = "https://dzone.com/services/widget/assets-listV2/DEFAULT?hidefeat=true&page="
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
	login := colly.NewCollector()

	// null the counters/pagers
	counter := 1
	stop := false

	// authenticate, otherwise you won't be able to download refcardz later
	login.OnHTML("form[role=form] input[type=hidden][name=TH_CSRF]", func(e *colly.HTMLElement) {
		thCsrf := e.Attr("value")
		err := login.Post("https://dzone.com/j_spring_security_check", map[string]string{"TH_CSRF": thCsrf, "_spring_security_remember_me": "true", "j_username": "dzone-refcardz@mailcatch.com", "j_password": "password123456"})
		if err != nil {
			log.Fatal(err)
		}
		return
	})

	// visit the users-login page
	login.Visit("https://dzone.com/users/login.html")
	fmt.Println("Logging in...")

	// instantiate second collector for assets-list
	assets := colly.NewCollector()

	// clone the logging collector for the actual downloads
	downloader := login.Clone()
	downloader.SetRequestTimeout(180 * time.Second)

	// is there any error in the assets-list response?
	assets.OnError(func(r *colly.Response, e error) {
		log.Println("Error:", e, r.Request.URL, string(r.Body))
	})

	// perform the following block on each assets-list page
	assets.OnResponse(func(r *colly.Response) {

		// unmarshal JSON based on the struct
		err := json.Unmarshal([]byte(r.Body), &data)
		if err != nil { // eventually error out
			log.Fatal(err)
		}

		// loop through the data assets containing Titles and incomplete Pdf links
		for _, obj := range data.Result.Data.Assets {
			fileName := strings.Replace(obj.Title+".pdf", " ", "_", -1) // title is going to be filename, eplace the " " with a "_" in the filename and add extension
			link := "https://dzone.com" + obj.Pdf                       // complete the HTTP link

			downloader.OnResponse(func(q *colly.Response) {
				q.Save("downloads/" + fileName) // save in the downloads directory
			})

			downloader.Visit(link)
			fmt.Println("#", counter, "Downloaded", fileName, "from", link) // show verbose progress
			link = ""
			fileName = ""
			counter++
		}

		// when the response is the last page, stop continuing with other pages and exit out
		if string(r.Body) == `{"success":true,"result":{"data":{"assets":[],"sort":"downloads"}},"status":200}` {
			defer fmt.Println("Last page reached!")
			stop = true
			return
		}

	}) // end of the assets-list response block

	// visit assets-list until the last page is reached
	for page := 1; page < 999; page++ {
		if stop {
			break
		} else {
			assets.Visit(baseURL + strconv.Itoa(page) + "&sort=downloads&type=refcard")
		}
	}

} // end of the main function
