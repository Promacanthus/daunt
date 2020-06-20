package main

import (
	"encoding/csv"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/gocolly/colly"
	"github.com/gocolly/colly/proxy"
)

const scholar = "https://scholar.google.com/scholar?hl=zh-CN&as_sdt=0,5&q="

func main() {

	fileName := "articles.csv"
	file, err := os.Create(fileName)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	err = writer.Write([]string{"ID", "Title", "Author", "Press", "Date", "Abstract", "Reference"})
	if err != nil {
		log.Fatalln(err)
	}

	c := colly.NewCollector(
		colly.UserAgent(
			"Mozilla/5.0 (X11; Linux x86_64) " +
				"AppleWebKit/537.36 (KHTML, like Gecko) " +
				"Chrome/81.0.4044.122 Safari/537.36"),
	)

	rp, err := proxy.RoundRobinProxySwitcher("socks5://127.0.0.1:1080")
	if err != nil {
		log.Fatalln(err)
	}
	c.SetProxyFunc(rp)

	c.OnHTML("#gs_res_ccl", func(e *colly.HTMLElement) {
		e.ForEach(".gs_ri", func(i int, e *colly.HTMLElement) {
			authorAndPress := strings.Split(e.DOM.Find(".gs_a").Text(), "-")
			pressAndDate := strings.Split(authorAndPress[1], ",")
			err := writer.Write([]string{
				strconv.Itoa(i),                                                     // ID
				e.DOM.Find(".gs_rt").Text(),                                         // Title
				strings.TrimSpace(authorAndPress[0]),                                // Author
				strings.TrimSpace(pressAndDate[0]),                                  // Press
				strings.TrimSpace(pressAndDate[1]),                                  // Date
				e.DOM.Find(".gs_rs").Text(),                                         // Abstract
				strings.Split(e.DOM.Find(".gs_fl>a:nth-of-type(3)").Text(), "：")[1], // Reference Number
			})
			if err != nil {
				log.Fatalln(err)
			}
		})
	})

	c.OnHTML("#gs_nml", func(e *colly.HTMLElement) {
		e.ForEach("a[href]", func(i int, e *colly.HTMLElement) {
			link := e.Text
			_ = link
			_ = link
		})
	})

	err = c.Visit(scholar + "acute+pancreatitis")
	if err != nil {
		log.Fatalln(err)
	}
}
