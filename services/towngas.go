package services

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
)

// Towngas towngas sercice account details
type Towngas struct {
	Username string
	Password string
}

// GetNewsNoticeAsync get latest info
func GetNewsNoticeAsync(c Towngas) (r string) {

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

	var result string

	loginResp, err := client.PostForm("https://eservice.towngas.com/EAccount/Login/SignIn", url.Values{"LoginID": {c.Username}, "password": {c.Password}})
	if err != nil {
		// handle error
	}
	defer loginResp.Body.Close()

	var loginCookies = loginResp.Cookies()

	// Get Gosted Tg Account Number
	getAccNumReq, _ := http.NewRequest("POST", "https://eservice.towngas.com/Common/GetHostedTGAccountAsync", nil)
	cookieJar.SetCookies(getAccNumReq.URL, loginCookies)
	getAccNumReq.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	getAccNumResp, err := client.Do(getAccNumReq)

	if err != nil {
		// handle error
	}
	defer getAccNumResp.Body.Close()

	getAccNumRespBody, _ := ioutil.ReadAll(getAccNumResp.Body)
	var accountNum []string
	dec := json.NewDecoder(strings.NewReader(string(getAccNumRespBody[:])))
	decErr := dec.Decode(&accountNum)
	if decErr != nil {
		// handle error
	}

	var formBody = "accountNo=" + accountNum[0]
	req, err := http.NewRequest("POST", "https://eservice.towngas.com/NewsNotices/GetNewsNoticeAsync", strings.NewReader(formBody))

	cookieJar.SetCookies(req.URL, loginCookies)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := client.Do(req)

	if err != nil {
		// handle error
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	result = string(body[:])

	return result
}
