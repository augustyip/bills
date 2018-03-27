package services

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"strings"
)

// Clp clp.com.hk account details
type Clp struct {
	Username string
	Password string
}

// Bill bill details
type Bill struct {
	AccountNo string `json:"caNo"`
	Balance   string `json:"LastBillAmount"`
}

// GetServiceDashboard get latest info
func (b *Bill) GetServiceDashboard(c Clp) {
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

	cookiesResp, err := client.Get("https://services.clp.com.hk/zh/login/index.aspx")
	if err != nil {
		// handle error
	}
	defer cookiesResp.Body.Close()
	for _, cookie := range cookiesResp.Cookies() {
		if cookie.Name == "K2Cie90hi___AntiXsrfToken" {
			csrfToken = cookie.Value
		}
	}

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
	data, err := ioutil.ReadAll(resp.Body)
	json.Unmarshal(data, &b)
}
