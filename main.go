package main

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

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

func GetSparkasseExchange() {
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

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(string(body))
}

func main() {
	GetUCBExchange()
}
