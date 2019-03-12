package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/gocolly/colly"
	"github.com/gocolly/colly/debug"
)

// constants that do not change
const (
	username = "dzone-refcardz@mailcatch.com" // created just for this purpose
	password = "password123456"               // created just for this purpose
	baseURL  = "https://dzone.com/services/widget/assets-listV2/DEFAULT?hidefeat=true&page="
	user     = 3590306 // user ID can be found https://dzone.com/users/3590306/dzone-refcardz.html"
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

// {"user":3590306}
type Payload0 struct {
	User int `json:"user"`
}

// {"item":333545,"type":"rejected","referral":"Web"}
type Payload1 struct {
	Item     int    `json:"item"`
	Type     string `json:"type"`
	Referral string `json:"referral"`
}

// {"item":325625,"user":3146177,"collectNow":true}
type Payload2 struct {
	Item       int  `json:"item"`
	User       int  `json:"user"`
	CollectNow bool `json:"collectNow"`
}

// global variables
var data RefcardzData
var thCsrf string

func main() { // main function start

	// null the counters & stopper
	stop := false

	// create folder for downloads
	os.Mkdir("downloads", 0700)

	// instantiate first collector for log in
	login := colly.NewCollector(colly.Debugger(&debug.LogDebugger{}), colly.UserAgent("Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/535.11 (KHTML, like Gecko) Chrome/69.0.3497.92 Safari/537.36"))

	// authenticate, otherwise you won't be able to download refcardz later
	login.OnHTML("form[role=form] input[type=hidden][name=TH_CSRF]", func(e *colly.HTMLElement) {
		thCsrf = e.Attr("value")
		err := e.Request.Post("https://dzone.com/j_spring_security_check", map[string]string{"TH_CSRF": thCsrf, "_spring_security_remember_me": "true", "j_username": username, "j_password": password})
		if err != nil {
			log.Fatal(err)
		}
		return
	})

	// visit the users-login page
	login.Visit("https://dzone.com/users/login.html")

	login.OnScraped(func(response *colly.Response) {
		fmt.Println("Logged in...")
	})

	// instantiate second collector for assets-list
	assets := login.Clone()

	// is there any error in the assets-list response?
	assets.OnError(func(r *colly.Response, e error) {
		log.Println("Error:", e, r.Request.URL, string(r.Body))
	})

	// perform the following block on each assets-list page
	assets.OnResponse(func(r *colly.Response) {

		var link string
		var fileName string

		// clone the login collector for the actual downloads
		downloader := login.Clone()
		downloader.SetRequestTimeout(180 * time.Second) // DZone can be really slow sometimes

		// unmarshal JSON based on the struct
		err := json.Unmarshal([]byte(r.Body), &data)
		if err != nil { // eventually error out
			log.Fatal(err)
		}

		// do this block on the response of Pdf website, e.g. https://dzone.com/asset/download/202338
		downloader.OnResponse(func(q *colly.Response) {
			if strings.Contains(q.Request.URL.String(), "https://dzone.com/interstitial?asset=") { // if response is interstitial website

				// clone the downloader collector
				i := login.Clone()

				interstitialLink := q.Request.URL.String() + "#"
				linkSplit := strings.Split(q.Request.URL.String(), "=")
				itemStr := linkSplit[len(linkSplit)-1] // last element called item
				fmt.Println(itemStr)
				item, errr := strconv.Atoi(itemStr)
				if errr != nil {
					log.Fatal(errr)
				}
				fmt.Println("Interstitial page encountered:", interstitialLink, "item #", item)

				// post the JSON payload to the canDownloadMembership page
				bytes, err00 := json.Marshal(Payload0{
					User: user,
				})
				if err00 != nil {
					panic(err00)
				}
				payload0 := []byte(bytes)

				i.PostRaw("https://dzone.com/services/internal/action/dzoneUsers-canDownloadMembership", payload0)

				// post the JSON payload to the trackClick page
				bytes, err0 := json.Marshal(Payload1{
					Item:     item,
					Type:     "accepted",
					Referral: "Web",
				})
				if err0 != nil {
					panic(err0)
				}
				payload1 := []byte(bytes)
				err1 := i.PostRaw("https://dzone.com/services/internal/action/campaigns-trackClick", payload1)
				if err1 != nil {
					log.Fatal(err1)
				}

				// post the JSON payload to the deliverLead page
				bytes, err := json.Marshal(Payload2{
					Item:       item,
					User:       user,
					CollectNow: true,
				})
				if err != nil {
					panic(err)
				}
				payload2 := []byte(bytes)
				err2 := i.PostRaw("https://dzone.com/services/internal/action/leadgen-deliverLead", payload2)
				if err2 != nil {
					log.Fatal(err2)
				}

				// set special headers for the interstitial pages
				i.OnRequest(func(t *colly.Request) {
					t.Headers.Set("X-TH-CSRF", thCsrf)
					t.Headers.Set("Content-Type", "application/json;charset=UTF-8")
				})

				// save the PDF file
				i.OnResponse(func(in *colly.Response) {
					in.Save("downloads/" + fileName)
				})

				// visit the link, e.g. https://dzone.com/interstitial?asset=1976511&item=313585#
				i.Visit(interstitialLink)

				i.OnScraped(func(url *colly.Response) {
					fmt.Println("Visited", url.Request.URL)
				})

			} else {
				// save the PDF file
				q.Save("downloads/" + fileName)
			}
		})

		downloader.OnScraped(func(url *colly.Response) {
			fmt.Println("Visited", url.Request.URL)
		})

		// loop through the data assets containing Titles and incomplete Pdf links
		for _, obj := range data.Result.Data.Assets {
			fileName = strings.Replace(obj.Title+".pdf", " ", "_", -1) // title is going to be filename, replace the " " with a "_" in the filename and add extension
			reg, err := regexp.Compile("[^a-zA-Z0-9._]+")              // only keep alphanumerics, dots and underscores
			if err != nil {
				log.Fatal(err)
			}
			fileName = reg.ReplaceAllString(fileName, "")
			link = "https://dzone.com" + obj.Pdf // complete the HTTP link

			// visit the Pdf website, e.g. https://dzone.com/asset/download/202338
			downloader.Visit(link)
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
			break // if stop = true
		} else {
			if err := assets.Visit(baseURL + strconv.Itoa(page) + "&sort=downloads&type=refcard"); err != nil { // visit assets-list
				fmt.Println("Error:", err)
				break
			}
		}
	}

} // end of the main function
