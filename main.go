package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/gocolly/colly"
	"github.com/joho/godotenv"
)

type waifustock struct {
	Name    string `json:"name,omitempty"`
	Price   string `json:"price,omitempty"`   //股價
	Capital string `json:"capital,omitempty"` //資本額
	Value   string `json:"value,omitempty"`   //市值
	Release string `json:"release,omitempty"` //釋股量
	Surplus string `json:"surplus,omitempty"` //盈餘
}

type wifusArray struct {
	wifus []waifustock
}

func main() {
	// crawl("https://acgn-stock.com/company/1")
	tele()
}

func crawl(url string) {

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

	c.OnHTML("div.company-card.company-card-default", func(e *colly.HTMLElement) {
		jr := &waifustock{
			Name:    strings.TrimSpace(e.DOM.Find("div.title").Text()),
			Price:   strings.TrimSpace(e.DOM.Find("div.row.row-info.d-flex.justify-content-between").Eq(2).Text()),
			Capital: strings.TrimSpace(e.DOM.Find("div.row.row-info.d-flex.justify-content-between").Eq(3).Text()),
			Value:   strings.TrimSpace(e.DOM.Find("div.row.row-info.d-flex.justify-content-between").Eq(4).Text()),
			Release: strings.TrimSpace(e.DOM.Find("div.row.row-info.d-flex.justify-content-between").Eq(5).Text()),
			Surplus: strings.TrimSpace(e.DOM.Find("div.row.row-info.d-flex.justify-content-between").Eq(6).Text()),
		}

		jsondata, _ := json.Marshal(jr)
		fmt.Println(string(jsondata))
	})

	// c.OnHTML(".page-item a[aria-label='下一頁']", func(e *colly.HTMLElement) {
	// 	e.Request.Visit(e.Attr("href"))
	// })

	c.OnScraped(func(r *colly.Response) {
		fmt.Println("Finished", r.Request.URL)
	})

	c.Visit(url)
}

func tele() {
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
				crawl("https://acgn-stock.com/company/1")
				msg.Text = "Done!"
				bot.Send(msg)
				// msg.Text = wife[0]
				// for _, wifu := range wife {
				// 	msg.Text = wifu
				// 	bot.Send(msg)
				// }
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
