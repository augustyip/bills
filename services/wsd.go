package services

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"strings"
)

// Wsd www.esd.wsd.gov.hk account details
type Wsd struct {
	Username string
	Password string
}

// ElectronicBill get latest info
func ElectronicBill(c Clp) string {
	var csrfToken string
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

	// login url https://www.esd.wsd.gov.hk/esd/login.do

	var loginBody = "username=" + c.Username + "&password=" + c.Password
	loginReq, err := http.NewRequest("POST", "https://services.clp.com.hk/Service/ServiceLogin.ashx", strings.NewReader(loginBody))
	if err != nil {
		// handle error
	}

	loginReq.Header.Set("X-CSRFToken", csrfToken)
	loginReq.Header.Set("X-Requested-With", "XMLHttpRequest")
	loginReq.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")

	loginResp, err := client.Do(loginReq)
	if err != nil {
		// handle error
	}
	defer loginResp.Body.Close()

	var loginedCookies = loginResp.Cookies()

	req, err := http.NewRequest("POST", "https://services.clp.com.hk/Service/ServiceDashboard.ashx", strings.NewReader("assCA="))

	cookieJar.SetCookies(req.URL, loginedCookies)
	req.Header.Set("X-CSRFToken", csrfToken)
	req.Header.Set("X-Requested-With", "XMLHttpRequest")
	resp, err := client.Do(req)

	if err != nil {
		// handle error
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	return string(body[:])
}
