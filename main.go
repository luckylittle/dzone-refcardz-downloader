package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/gocolly/colly"
)

// constants that do not change
const (
	username = "dzone-refcardz@mailcatch.com" // created just for this purpose
	password = "password123456"               // created just for this purpose
)

func main() { // main function start

	// instantiate first collector
	f := colly.NewCollector()

	// always start with page zero
	page := 0

	// construct the JSON structure
	data := struct {
		Result struct {
			Data struct {
				Assets []struct {
					Title string `json:"title"`
					Pdf   string `json:"pdf"`
				}
			}
		}
	}{}

	// keep visiting the assets-list websites until the specific response indicating last page is returned
	for {

		// is there any error in the response?
		f.OnError(func(r *colly.Response, e error) {
			log.Println("Error:", e, r.Request.URL, string(r.Body))
		})

		// perform the following block on each response
		f.OnResponse(func(r *colly.Response) {

			// when the response is the last page, stop continuing with other pages
			if string(r.Body) == `{"success":true,"result":{"data":{"assets":[],"sort":"downloads"}},"status":200}` {
				fmt.Println("Last page reached")
				defer fmt.Println("!")
				os.Exit(0)
			}

			// unmarshal JSON based on the struct
			err := json.Unmarshal([]byte(r.Body), &data)
			if err != nil { // eventually error out
				log.Fatal(err)
			}

			// naked return; returns the current values in the return arguments local variables
			return
		}) // end of response block

		fmt.Println("Struct #", page, ":") // display the current page
		fmt.Println(data)                  // only for testing, display the struct of the current page

		// TODO: Perform something like below for each struct data
		// instantiate new collector
		// p := colly.NewCollector()

		// authenticate, otherwise you won't be able to download refcardz
		// err := p.Post("https://dzone.com/services/internal/action/users-login", map[string]string{"username": username, "password": password})
		// if err != nil {
		//	log.Fatal(err)
		// }

		// attach callbacks after login
		// p.OnResponse(func(q *colly.Response) {
		//	log.Println("Login response received:", q.StatusCode)
		//	replace the " " with a "_" in the filename
		//	result := strings.Replace(title, " ", "_", -1)
		//	q.Save(result)
		// })

		// download link example: https://dzone.com/asset/download/279342
		// p.Visit("https://dzone.com" + pdf)

		page++ // increase counter after visiting the assets-list webpage

		// assets-list website construction
		f.Visit("https://dzone.com/services/widget/assets-listV2/DEFAULT?hidefeat=true&page=" + strconv.Itoa(page) + "&sort=downloads&type=refcard")

	} // end of the for loop

} // end of the main function
