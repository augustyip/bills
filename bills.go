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
	}

	for _, cert := range certs {

		loginResp, err := client.PostForm("https://eservice.towngas.com/EAccount/Login/SignIn", url.Values{"LoginID": {cert.Username}, "password": {cert.Password}})
		if err != nil {
			// handle error
		}
		defer loginResp.Body.Close()

		var loginCookies = loginResp.Cookies()

		// https://eservice.towngas.com/Common/GetMeterReadingHistoryAsync accountNo:7220678095
		// req, err := http.NewRequest("GET", "https://eservice.towngas.com/en/BillingUsage/NewsNotices", nil)
		req, err := http.NewRequest("PostForm", "https://eservice.towngas.com/NewsNotices/GetNewsNoticeAsync", strings.NewReader("accountNo=7220678095"))

		cookieJar.SetCookies(req.URL, loginCookies)

		req.Header.Add("origin", "https://eservice.towngas.com")
		req.Header.Add("referer", "https://eservice.towngas.com/en/BillingUsage/NewsNotices")

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
