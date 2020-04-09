package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
	"github.com/joho/godotenv"

	tb "gopkg.in/tucnak/telebot.v2"
)

func main() {
	godotenv.Load()
	var (
		port      = os.Getenv("PORT")       // sets automatically
		publicURL = os.Getenv("PUBLIC_URL") // you must add it to your config vars
		token     = os.Getenv("TOKEN")      // you must add it to your config vars
	)

	webhook := &tb.Webhook{
		Listen:   ":" + port,
		Endpoint: &tb.WebhookEndpoint{PublicURL: publicURL},
	}
	pref := tb.Settings{
		Token:  token,
		Poller: webhook,
	}

	b, err := tb.NewBot(pref)
	if err != nil {
		log.Fatal(err)
	}
	b.Handle("/hello", func(m *tb.Message) {
		b.Send(m.Sender, "Hi!")
	})
	crawl("https://acgn-stock.com/company/1")
}

func crawl(url string) {
	c := colly.NewCollector(colly.UserAgent("Mozilla/5.0 (compatible; Googlebot/2.1; +http://www.google.com/bot.html)"))
	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})
	c.OnError(func(_ *colly.Response, err error) {
		fmt.Println("Something went wrong:", err)
	})

	c.OnResponse(func(r *colly.Response) {
		fmt.Println("Visited", r.Request.URL)
	})

	c.OnHTML(".media", func(e *colly.HTMLElement) {
		e.DOM.Find("div.title").Each(func(i int, s *goquery.Selection) {
			fmt.Println(strings.TrimSpace(s.Text()))
		})
	})

	c.OnHTML(".page-item a[aria-label='下一頁']", func(e *colly.HTMLElement) {
		e.Request.Visit(e.Attr("href"))
	})

	c.OnScraped(func(r *colly.Response) {
		fmt.Println("Finished", r.Request.URL)
	})

	c.Visit(url)
}
