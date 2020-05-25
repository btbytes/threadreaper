package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/gosimple/slug"
	"github.com/manifoldco/promptui"
)

func Check(e error) {
	if e != nil {
		panic(e)
	}
}

func GetTheadId(url string) string {

	parts := strings.Split(url, "/")
	fname := strings.Split(parts[len(parts)-1], ".")[0]
	return fname
}

func main() {
	var author string
	var title string
	var css string
	var ask bool
	var outFile string
	flag.StringVar(&author, "author", "", "Name of the author")
	flag.StringVar(&title, "title", "", "Page title")
	flag.StringVar(&css, "css", "https://cdn.jsdelivr.net/gh/kognise/water.css@latest/dist/light.min.css", "CSS URL")
	flag.StringVar(&outFile, "out", "", "Output filename")
	flag.BoolVar(&ask, "p", true, "prompt me for options")
	flag.Parse()
	url := flag.Arg(0)
	if ask {
		// Author
		p1 := promptui.Prompt{
			Label: "Author",
		}
		result, err1 := p1.Run()
		Check(err1)
		author = result
		// Title
		p2 := promptui.Prompt{
			Label: "Title",
		}
		result, err2 := p2.Run()
		Check(err2)
		title = result

		// outFile
		p3 := promptui.Prompt{
			Label: "Output File",
		}
		result, err3 := p3.Run()
		Check(err3)
		outFile = result
	}

	res, err := http.Get(url)
	Check(err)
	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}
	if outFile == "" {
		values := []string{}
		values = append(values, slug.Make(author))
		values = append(values, slug.Make(title))
		outFile = strings.Join(values, "-")
		if author == "" && title == "" {
			outFile = GetTheadId(url)
		}
		sLen := len(outFile)
		if sLen > 40 {
			sLen = 40
		}
		outFile = outFile[:sLen] + ".html"
	}
	f, err := os.Create(outFile)
	Check(err)
	defer f.Close()
	doc, err := goquery.NewDocumentFromReader(res.Body)
	Check(err)
	fmt.Fprintf(f, `<!doctype html>
  <html>
  <head>
  <meta charset="utf-8">
  <meta http-equiv="X-UA-Compatible" content="IE=edge">
  <meta name="viewport" content="width=device-width, initial-scale=1">
  <link rel="stylesheet" href="%s">
  <style type="text/css">
  .row {
    margin-top: 1.5em;
  }
  </style>
  </head>
  <body>
  <div class="container">
  <div class="page-header">
  `, css)
	t := time.Now()
	fmt.Fprintf(f, `<h1 class="title">%s</h1>`, title)
	fmt.Fprintf(f, `<p>by <span class="author"><a href="https://twitter.com/%s">@%s</a></span></p>`, author, author)
	fmt.Fprintf(f, `</div><!-- page-header -->`)
	//  doc.Find(".thread-info .time").Each()
	doc.Find(".t-main .content-tweet").Each(func(i int, s *goquery.Selection) {
		screenName, _ := s.Attr("data-screenname")
		tweet, _ := s.Attr("data-tweet")
		s.Find("img").Each(func(ii int, ss *goquery.Selection) {
			imgURL, imgExists := ss.Attr("data-src")
			if imgExists {
				ss.SetAttr("src", imgURL)
			}
		})
		s.Find(".tw-permalink").Remove()
		res, err := s.Html()
		Check(err)
		fmt.Fprintf(f, `<div class="row tweet" data-screenname="%s" data-tweet="%s">%s</div>`, screenName, tweet, res)
	})
	fmt.Fprintf(f, `</div><footer>Retrieved on <span class="time"><time>%s</time></span>`, t.Format("2006-01-02T15:04:05-07:00"))
	fmt.Fprintf(f, `, from <a href="%s">threadreaderapp page</a></footer></body></html>`, url)

}
