package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"
)

type UCBRate struct {
	BidRate                      float32 `json:"bidRate"`
	ExchangeRateUpdatedTimestamp string  `json:"exchangeRateUpdatedTimestamp"`
}

type ErsteRate struct {
	Date       string          `json:"date"`
	Currencies []ErsteCurrency `json:"currencies"`
}

type ErsteCurrency struct {
	Name  string             `json:"name"`
	Rates ErsteCurrencyRates `json:"rates"`
}

type ErsteCurrencyRates struct {
	Buying  float32 `json:"buying"`
	Selling float32 `json:"selling"`
}

type ResponseRate struct {
	Date string  `json:"date"`
	Rate float32 `json:"rate"`
}

type ApiResponse struct {
	Erste []ResponseRate `json:"erste"`
	UCB   []ResponseRate `json:"ucb"`
}

func getUCBCookie() string {
	req, err := http.NewRequest("GET", "https://www.unicredit.ba/ba/stanovnistvo/tecajna_lista.html", nil)
	if err != nil {
		fmt.Println(err)
		return ""
	}

	req.Header.Set("authority", "www.unicredit.ba")
	req.Header.Set("accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8")
	req.Header.Set("accept-language", "en-US,en;q=0.6")
	req.Header.Set("cache-control", "no-cache")
	req.Header.Set("pragma", "no-cache")
	req.Header.Set("referer", "https://www.unicredit.ba/ba/stanovnistvo.html")
	req.Header.Set("upgrade-insecure-requests", "1")
	req.Header.Set("user-agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/114.0.0.0 Safari/537.36")

	rsp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println(err)
		return ""
	}

	var sb strings.Builder

	for _, c := range rsp.Cookies() {
		sb.WriteString(c.Name + `=` + c.Value + `;`)
	}

	return sb.String()
}

func getUCBExchange() ([]ResponseRate, error) {
	const TIME_FORMAT = "20060102T03:04:05.0-0700"
	dateTo := time.Now().Format(TIME_FORMAT)
	dateFrom := time.Now().AddDate(0, 0, -7).Format(TIME_FORMAT)

	req, err := http.NewRequest("POST", "https://www.unicredit.ba/cwa/GetExchangeRates", strings.NewReader(`{"Currency": "GBP", "DateFrom": "`+dateFrom+`", "DateTo": "`+dateTo+`"}`))
	if err != nil {
		fmt.Println(err)
		return nil, err
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
	req.Header.Set("cookie", getUCBCookie())

	// Get the exchange rate page from Unicredit.ba
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	// Check if response is valid
	if resp.StatusCode != http.StatusOK {
		fmt.Println("Error: ", resp)
		return nil, errors.New("status code not OK")
	}
	defer resp.Body.Close()

	var rates []UCBRate
	var rspRates []ResponseRate

	err = json.NewDecoder(resp.Body).Decode(&rates)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	for _, r := range rates {
		rspRates = append(rspRates, ResponseRate{Date: r.ExchangeRateUpdatedTimestamp, Rate: r.BidRate})
	}

	return rspRates, nil
}

func getErsteExchange() ([]ResponseRate, error) {
	const TIME_FORMAT = "2006-01-02"
	dateTo := time.Now().Format(TIME_FORMAT)
	dateFrom := time.Now().AddDate(0, 0, -8).Format(TIME_FORMAT)
	req, err := http.NewRequest("GET", `https://local.erstebank.hr/api/v1/fx?dateFrom=`+dateFrom+`&dateThru=`+dateTo+``, nil)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		fmt.Println("Error: ", resp)
		return nil, errors.New("status code not OK")
	}
	defer resp.Body.Close()

	var rates []ErsteRate
	var rspRates []ResponseRate

	err = json.NewDecoder(resp.Body).Decode(&rates)
	if err != nil {
		fmt.Println(err)
		return nil, err
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
		rspRates = append(rspRates, ResponseRate{
			Date: r.Date,
			Rate: bamCurrency.Rates.Selling / gbpCurrency.Rates.Buying},
		)
	}

	return rspRates, nil
}

func Exchange(w http.ResponseWriter, r *http.Request) {
	ratesUcb, _ := getUCBExchange()
	ratesErste, _ := getErsteExchange()

	jsonResp, err := json.Marshal(ApiResponse{Erste: ratesErste, UCB: ratesUcb})
	if err != nil {
		fmt.Println(err)
		return
	}

	w.Write(jsonResp)
}
