package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/PuerkitoBio/goquery"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/gocolly/colly"
	"github.com/joho/godotenv"
)

type waifustock struct {
	name    string
	price   string //股價
	capital string //資本額
	value   string //市值
	release string //釋股量
	surplus string //盈餘
}

func main() {
	godotenv.Load()
	bot, err := tgbotapi.NewBotAPI(os.Getenv("TOKEN"))
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil { // ignore any non-Message Updates
			continue
		}
		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

		if update.Message.IsCommand() {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
			switch update.Message.Command() {
			case "help":
				msg.Text = "type /sayhi or /status."
				bot.Send(msg)
			case "colly":
				wife := crawl("https://acgn-stock.com/company/1")
				// msg.Text = wife[0]
				for _, wifu := range wife {
					msg.Text = wifu
					bot.Send(msg)
				}
			case "status":
				msg.Text = "I'm not ok."
				bot.Send(msg)
			default:
				msg.Text = "I don't know that command"
				bot.Send(msg)
			}
			//bot.Send(msg)
		}
	}
}

func crawl(url string) []string {
	var wifus []string
	c := colly.NewCollector(colly.UserAgent("Mozilla/5.0 (compatible; Googlebot/2.1; +http://www.google.com/bot.html)"))
	c.OnRequest(func(r *colly.Request) {
		log.Println("Visiting", r.URL)
	})
	c.OnError(func(_ *colly.Response, err error) {
		log.Println("Something went wrong:", err)
	})

	c.OnResponse(func(r *colly.Response) {
		log.Println("Visited", r.Request.URL)
	})

	c.OnHTML(".media", func(e *colly.HTMLElement) {
		e.DOM.Find("div.title").Each(func(i int, s *goquery.Selection) {
			log.Println(strings.TrimSpace(s.Text()))
			wifus = append(wifus, strings.TrimSpace(s.Text()))
		})
	})

	// c.OnHTML(".page-item a[aria-label='下一頁']", func(e *colly.HTMLElement) {
	// 	e.Request.Visit(e.Attr("href"))
	// })

	c.OnScraped(func(r *colly.Response) {
		fmt.Println("Finished", r.Request.URL)
	})

	c.Visit(url)
	return wifus
}
