package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"strings"
)

// Certification struct
type Certification struct {
	Service  string
	Username string
	Password string
}

func main() {

	file, _ := os.Open("cert.json")
	decoder := json.NewDecoder(file)
	certs := make([]Certification, 0)
	err := decoder.Decode(&certs)
	if err != nil {
		fmt.Println("error:", err)
	}

	cookieJar, _ := cookiejar.New(nil)
	client := &http.Client{
		Jar: cookieJar,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			fmt.Println(req.URL)
			for _, r := range via {
				fmt.Println(r.URL)
			}
			return http.ErrUseLastResponse
		},
	}

	for _, cert := range certs {

		loginResp, err := client.PostForm("https://eservice.towngas.com/EAccount/Login/SignIn", url.Values{"LoginID": {cert.Username}, "password": {cert.Password}})
		if err != nil {
			// handle error
		}
		defer loginResp.Body.Close()

		var loginCookies = loginResp.Cookies()

		// https://eservice.towngas.com/Common/GetMeterReadingHistoryAsync accountNo:7220678095
		req, err := http.NewRequest("POST", "https://eservice.towngas.com/NewsNotices/GetNewsNoticeAsync", strings.NewReader("accountNo=7220678095"))

		cookieJar.SetCookies(req.URL, loginCookies)

		req.Header.Set("origin", "https://eservice.towngas.com")
		req.Header.Set("referer", "https://eservice.towngas.com/en/BillingUsage/NewsNotices")
		req.Header.Set("x-requested-with", "XMLHttpRequest")
		req.Header.Set("user-agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_3) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/64.0.3282.186 Safari/537.36")
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		fmt.Println(req.Header)
		resp, err := client.Do(req)

		if err != nil {
			// handle error
		}
		fmt.Println(resp.StatusCode)
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		bodyContent := string(body[:])
		fmt.Printf(bodyContent)
	}

}
