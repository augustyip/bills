package services

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"strings"

	"golang.org/x/net/html"
)

// Wsd www.esd.wsd.gov.hk account details
type Wsd struct {
	Username string
	Password string
}

// ElectronicBill get latest info
func ElectronicBill(c Wsd) string {

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

	var loginPageUrl = "https://www.esd.wsd.gov.hk/esd/login.do"

	loginPageResp, err := http.Get(loginPageUrl)
	if err != nil {
		// handle error
	}
	defer loginPageResp.Body.Close()
	loginPageBody, err := ioutil.ReadAll(loginPageResp.Body)

	loginPageDoc, err := html.Parse(strings.NewReader(string(loginPageBody[:])))
	if err != nil {
	}

	// login url https://www.esd.wsd.gov.hk/esd/login.do

	var loginBody = "userID=" + c.Username + "&password=" + c.Password
	loginReq, err := http.NewRequest("POST", loginPageUrl, strings.NewReader(loginBody))
	if err != nil {
		// handle error
	}

	// loginReq.Header.Set("X-CSRFToken", csrfToken)
	// loginReq.Header.Set("X-Requested-With", "XMLHttpRequest")
	loginReq.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	loginResp, err := client.Do(loginReq)
	fmt.Println(loginResp.StatusCode)

	return "xxx"
	// body, err := ioutil.ReadAll(loginResp.Body)
	// return string(body[:])
}
