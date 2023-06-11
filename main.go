package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

type UCBRate struct {
	BidRate float32 `json:"rate"`
}

type ErsteCurrencyRates struct {
	Buying      float32 `json:"buying"`
	BuyingCash  float32 `json:"buyingCash"`
	Middle      float32 `json:"middle"`
	Selling     float32 `json:"selling"`
	SellingCash float32 `json:"sellingCash"`
	Hnb         float32 `json:"hnb"`
}
type ErsteCurrency struct {
	Name  string             `json:"name"`
	Rates ErsteCurrencyRates `json:"rates"`
}
type ErsteRate struct {
	Date       string          `json:"date"`
	Currencies []ErsteCurrency `json:"currencies"`
}

type ErsteResponseRate struct {
	Date string  `json:"date"`
	Rate float32 `json:"rate"`
}

func GetUCBExchange() {
	const TIME_FORMAT = "20060102T03:04:05.0-0700"
	dateTo := time.Now().Format(TIME_FORMAT)
	dateFrom := time.Now().AddDate(0, 0, -7).Format(TIME_FORMAT)
	// Set the User-Agent header
	req, err := http.NewRequest("POST", "https://www.unicredit.ba/cwa/GetExchangeRates", strings.NewReader(`{"Currency": "GBP", "DateFrom": "`+dateFrom+`", "DateTo": "`+dateTo+`"}`))
	if err != nil {
		fmt.Println(err)
		return
	}

	req.Header.Set("authority", "www.unicredit.ba")
	req.Header.Set("accept", "*/*")
	req.Header.Set("accept-language", "en-US,en;q=0.6")
	req.Header.Set("cache-control", "no-cache")
	req.Header.Set("content-type", "application/json")
	req.Header.Add("content-encoding", "gzip")
	req.Header.Set("entitycode", "BH")
	req.Header.Set("language", "BS")
	req.Header.Set("origin", "https://www.unicredit.ba")
	req.Header.Set("pragma", "no-cache")
	req.Header.Set("referer", "https://www.unicredit.ba/ba/stanovnistvo/tecajna_lista.html")
	req.Header.Set("product", "PWS")
	req.Header.Set("sourcesystem", "PWS")
	req.Header.Set("user-agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/114.0.0.0 Safari/537.36")

	// Get the exchange rate page from Unicredit.ba
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Check if response is valid
	if resp.StatusCode != http.StatusOK {
		fmt.Println("Error: ", resp)
		return
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Print the response body
	fmt.Println(string(body))
}

func GetErsteExchange() {
	const TIME_FORMAT = "2006-01-02"
	dateTo := time.Now().Format(TIME_FORMAT)
	dateFrom := time.Now().AddDate(0, 0, -7).Format(TIME_FORMAT)
	req, err := http.NewRequest("GET", `https://local.erstebank.hr/api/v1/fx?dateFrom=`+dateFrom+`&dateThru=`+dateTo+``, nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}

	if resp.StatusCode != http.StatusOK {
		fmt.Println("Error: ", resp)
		return
	}

	var rates []ErsteRate
	var rspRates []ErsteResponseRate

	err = json.NewDecoder(resp.Body).Decode(&rates)
	if err != nil {
		fmt.Println(err)
		return
	}

	for _, r := range rates {
		var gbpCurrency ErsteCurrency
		var bamCurrency ErsteCurrency
		for _, c := range r.Currencies {
			if c.Name == "GBP" {
				gbpCurrency = c
			}
			if c.Name == "BAM" {
				bamCurrency = c
			}
		}
		rspRates = append(rspRates, ErsteResponseRate{
			Date: r.Date,
			Rate: bamCurrency.Rates.Selling / gbpCurrency.Rates.Buying},
		)
	}

	jsonResp, err := json.Marshal(rspRates)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(string(jsonResp))
}

func main() {
	GetErsteExchange()
}
