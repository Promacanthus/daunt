package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/gocolly/colly"
	"github.com/gocolly/colly/proxy"
)

const scholar = "https://scholar.google.com/scholar?hl=zh-CN&as_sdt=0,5&q="

func main() {

	ID := 0
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

			authorAndPress := make([]string, 2)
			pressAndDate := make([]string, 2)

			tmpAP := e.DOM.Find(".gs_a").Text()
			if strings.Contains(tmpAP, "-") {
				authorAndPress = strings.Split(tmpAP, "-")
			} else {
				authorAndPress[0] = tmpAP
				authorAndPress[1] = ""
			}

			if authorAndPress[1] != "" {
				if strings.Contains(authorAndPress[1], ",") {
					pressAndDate = strings.Split(authorAndPress[1], ",")
				} else {
					pressAndDate[0] = authorAndPress[1]
					pressAndDate[1] = ""
				}
			} else {
				pressAndDate[0] = ""
				pressAndDate[1] = ""
			}

			ID += 1
			writer.Write([]string{
				strconv.Itoa(ID),                                                    // ID
				e.DOM.Find(".gs_rt").Text(),                                         // Title
				strings.TrimSpace(authorAndPress[0]),                                // Author
				strings.TrimSpace(pressAndDate[0]),                                  // Press
				strings.TrimSpace(pressAndDate[1]),                                  // Date
				e.DOM.Find(".gs_rs").Text(),                                         // Abstract
				strings.Split(e.DOM.Find(".gs_fl>a:nth-of-type(3)").Text(), "：")[1], // Reference Number
			})
		})
	})

	c.OnHTML("#gs_nml", func(e *colly.HTMLElement) {
		e.ForEach("a[href]", func(i int, e *colly.HTMLElement) {
			link := e.Attr("href")
			c.Visit(e.Request.AbsoluteURL(link))
		})
	})

	c.OnRequest(func(r *colly.Request) {
		log.Printf("Visiting : %s\n", r.URL.String())
	})

	c.OnError(func(r *colly.Response, err error) {
		fmt.Println("Request URL:", r.Request.URL, "failed with response:", r, "\nError:", err)
	})

	err = c.Visit(scholar + "acute+pancreatitis")
	if err != nil {
		log.Fatalln(err)
	}
}
