package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/PuerkitoBio/goquery"
)

func ThreadReaper(url string, pageTitle string) {
	res, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf(`<!doctype html><html>
 
  <head>
  <meta charset="utf-8">
  <meta http-equiv="X-UA-Compatible" content="IE=edge">
    <meta name="viewport" content="width=device-width, initial-scale=1">
  <link rel="stylesheet" href="https://maxcdn.bootstrapcdn.com/bootstrap/3.3.7/css/bootstrap.min.css" integrity="sha384-BVYiiSIFeK1dGmJRAkycuHAHRg32OmUcww7on3RYdg4Va+PmSTsz/K68vbdEjh4u" crossorigin="anonymous">
  <style type="text/css">
  .row {
    margin-top: 1.5em;
  }
  </style>
  </head>
  <body>
  <div class="container">
  <div class="page-header">
  `)
	t := time.Now()
	fmt.Printf("<h1>%s</h1>", pageTitle)
	fmt.Printf(`<p><a href="%v">Original Thread</a></p>`, url)
	fmt.Printf("<time>%s</time>", t.Format("2006-01-02T15:04:05-07:00"))
	fmt.Println(`</div><!-- page-header -->`)

	doc.Find(".content-tweet").Each(func(i int, s *goquery.Selection) {
		screenName, e1 := s.Attr("data-screenname")
		tweet, e2 := s.Attr("data-tweet")
		s.Find("img").Each(func(ii int, ss *goquery.Selection) {
			imgURL, imgExists := ss.Attr("data-src")
			if imgExists {
				ss.SetAttr("src", imgURL)
			}
		})
		s.Find(".tw-permalink").Remove()
		res, err := s.Html()
		if err != nil {
			log.Fatal(err)
			return
		}
		if e1 && e2 {
			fmt.Printf(`<div class="row tweet" data-screenname="%s" data-tweet="%s">%s</div>`, screenName, tweet, res)
		}
	})

}

func main() {
	url := os.Args[1]
	pageTitle := ""
	if len(os.Args) > 2 {
		pageTitle = os.Args[2]
	}
	ThreadReaper(url, pageTitle)
}
