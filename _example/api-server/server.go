package main

import (
	"log"
	"net/http"
	"strconv"

	rent "github.com/neighborhood999/fiveN1-rent-scraper"
)

func callbackHandler(w http.ResponseWriter, r *http.Request) {
	options := rent.NewOptions()
	r.ParseForm()

	if r.Method == "GET" {
		options.Area = r.FormValue("area")
		options.Role = r.FormValue("role")
		options.Other = r.FormValue("other")
		options.Order = r.FormValue("order")
		options.Shape = r.FormValue("shape")
		options.Floor = r.FormValue("floor")
		options.Option = r.FormValue("option")
		options.HasImg = r.FormValue("hasImage")
		options.Section = r.FormValue("section")
		options.NotCover = r.FormValue("notCover")
		options.RentPrice = r.FormValue("rentPrice")
		options.OrderType = r.FormValue("orderType")
		options.Sex, _ = strconv.Atoi(r.FormValue("sex"))
		options.Kind, _ = strconv.Atoi(r.FormValue("kind"))
		options.FirstRow, _ = strconv.Atoi(r.FormValue("firstRow"))

		url, err := rent.GenerateURL(options)
		if err != nil {
			log.Fatal(err)
		}

		f := rent.NewFiveN1(url)
		f.SetReqCookie(r.FormValue("urlJump"))

		if err := f.Scrape(1); err != nil {
			log.Fatal(err)
		}
		json := rent.ConvertToJSON(f.RentList)

		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token")

		w.Write(json)
	}
}

func main() {
	http.HandleFunc("/", callbackHandler)
	http.ListenAndServe(":8080", nil)
}
