package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
)

func main() {
	parseurls("https://acgn-stock.com/company/1")
}

func crawl(url string) string {
	clinet := &http.Client{}
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("User-Agent", "Mozilla/5.0 (compatible; Googlebot/2.1; +http://www.google.com/bot.html)")
	resp, err := clinet.Do(req)
	if err != nil {
		fmt.Println("Http get err:", err)
		return ""
	}
	if resp.StatusCode != 200 {
		fmt.Println("Http status code:", resp.StatusCode)
		return ""
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("err:", err)
		return ""
	}
	return string(body)
}

func parseurls(url string) {
	body := crawl(url)
	body = strings.Replace(body, "\n", "", -1)
	// rp := regexp.MustCompile(`<div class= company-card company-card-default>(.*?)</div>`)
	title := regexp.MustCompile(`<a href="/company/detail/(.*?)</a>`)
	company := title.FindAllStringSubmatch(body, -1)
	// k := title.FindAllStringSubmatch(body, -1)
	// fmt.Println(title.MatchString(k[0][0]))
	for _, item := range company {
		fmt.Println(item[0])
	}

}
