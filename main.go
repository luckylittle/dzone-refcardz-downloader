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
	"github.com/gocolly/colly/debug"
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
	f := colly.NewCollector(colly.Debugger(&debug.LogDebugger{}))

	// instantiate second collector for assets-list
	g := colly.NewCollector()

	// null the counters/pagers
	page := 0
	counter := 1

	// authenticate, otherwise you won't be able to download refcardz later
	f.OnHTML("form[role=form] input[type=hidden][name=TH_CSRF]", func(e *colly.HTMLElement) {
		thCsrf := e.Attr("value")
		err := f.Post("https://dzone.com/j_spring_security_check", map[string]string{"TH_CSRF": thCsrf, "_spring_security_remember_me": "true", "j_username": "dzone-refcardz@mailcatch.com", "j_password": "password123456"})
		if err != nil {
			log.Fatal(err)
		}
		return
	})

	// visit the users-login page
	f.Visit("https://dzone.com/users/login.html")

	// clone the logging collector
	h := f.Clone()
	h.SetRequestTimeout(180 * time.Second)

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
			//return

		}) // end of response block

		// loop through the data assets containing Titles and incomplete Pdf links
		for _, obj := range data.Result.Data.Assets {
			fileName := strings.Replace(obj.Title+".pdf", " ", "_", -1) // title is going to be filename
			link := "https://dzone.com" + obj.Pdf                       // complete the HTTP link
			fmt.Println("link after assignment:", link)

			h.OnResponse(func(q *colly.Response) {
				// replace the " " with a "_" in the filename and add extension
				q.Save("downloads/" + fileName)                                 // save
				fmt.Println("#", counter, "Downloaded", fileName, "from", link) // show verbose progress
			})

			// visit
			fmt.Println("link before h.visit:", link)
			h.Visit(link)
			link = ""
			counter++
		}
		/*
		   Debug log:
		   [000001] 1 [     1 - request] map["url":"https://dzone.com/users/login.html"] (36.392µs)
		   [000002] 1 [     1 - response] map["url":"https://dzone.com/users/login.html" "status":"OK"] (1.535579866s)
		   [000003] 1 [     1 - html] map["selector":"form[role=form] input[type=hidden][name=TH_CSRF]" "url":"https://dzone.com/users/login.html"] (1.537939626s)
		   [000004] 1 [     2 - request] map["url":"https://dzone.com/j_spring_security_check"] (1.538068512s)
		   [000005] 1 [     2 - response] map["url":"https://dzone.com/index.html" "status":"OK"] (2.565529343s)
		   [000006] 1 [     2 - scraped] map["url":"https://dzone.com/index.html"] (2.570756573s)
		   [000007] 1 [     1 - scraped] map["url":"https://dzone.com/users/login.html"] (2.570774088s)
		   link after assignment: https://dzone.com/asset/download/280333
		   link before h.visit: https://dzone.com/asset/download/280333
		   [000008] 3 [     1 - request] map["url":"https://dzone.com/asset/download/280333"] (3.377048525s)
		   [000009] 3 [     1 - response] map["url":"https://dzone.com/storage/assets/11325551-dzone-refcard288-gettingstartedwithgit0221.pdf" "status":"OK"] (11.466351159s)
		   # 1 Downloaded Getting_Started_With_Git.pdf from https://dzone.com/asset/download/280333
		   [000010] 3 [     1 - scraped] map["url":"https://dzone.com/storage/assets/11325551-dzone-refcard288-gettingstartedwithgit0221.pdf"] (11.468748885s)
		   link after assignment: https://dzone.com/asset/download/279342
		   link before h.visit: https://dzone.com/asset/download/279342
		   [000011] 3 [     2 - request] map["url":"https://dzone.com/asset/download/279342"] (11.469000995s)
		   [000012] 3 [     2 - response] map["status":"OK" "url":"https://dzone.com/storage/assets/11283656-dzone-refcard288-timeseriesdata.pdf"] (1m15.373064768s)
		   # 2 Downloaded Getting_Started_With_Git.pdf from
		   # 2 Downloaded Working_With_Time_Series_Data.pdf from https://dzone.com/asset/download/279342
		   [000013] 3 [     2 - scraped] map["url":"https://dzone.com/storage/assets/11283656-dzone-refcard288-timeseriesdata.pdf"] (1m15.385235892s)
		   link after assignment: https://dzone.com/asset/download/278339
		   link before h.visit: https://dzone.com/asset/download/278339
		   [000014] 3 [     3 - request] map["url":"https://dzone.com/asset/download/278339"] (1m15.385308434s)
		   [000015] 3 [     3 - response] map["url":"https://dzone.com/storage/assets/11231342-dzone-refcard288-introtolowcode.pdf" "status":"OK"] (1m22.233645502s)
		   # 3 Downloaded Getting_Started_With_Git.pdf from
		   # 3 Downloaded Working_With_Time_Series_Data.pdf from
		   # 3 Downloaded Low-Code_Application_Development.pdf from https://dzone.com/asset/download/278339
		   [000016] 3 [     3 - scraped] map["url":"https://dzone.com/storage/assets/11231342-dzone-refcard288-introtolowcode.pdf"] (1m22.257504381s)
		   link after assignment: https://dzone.com/asset/download/279336
		   link before h.visit: https://dzone.com/asset/download/279336
		   [000017] 3 [     4 - request] map["url":"https://dzone.com/asset/download/279336"] (1m22.25785153s)
		   ^Csignal: interrupt
		*/
		// next page
		page++

		// visitassets-list website construction
		g.Visit("https://dzone.com/services/widget/assets-listV2/DEFAULT?hidefeat=true&page=" + strconv.Itoa(page) + "&sort=downloads&type=refcard")

	} // end of the for loop

} // end of the main function
