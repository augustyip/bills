package services

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
)

// Clp towngas sercice account details
type Clp struct {
	Username string
	Password string
}

// GetServiceDashboard get latest info
func GetServiceDashboard(c Clp) string {

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

	loginResp, err := client.PostForm("https://services.clp.com.hk/Service/ServiceLogin.ashx", url.Values{"LoginID": {c.Username}, "password": {c.Password}})
	if err != nil {
		// handle error
	}
	defer loginResp.Body.Close()

	var loginCookies = loginResp.Cookies()
	fmt.Println(loginCookies)
	req, err := http.NewRequest("POST", "https://services.clp.com.hk/Service/ServiceDashboard.ashx", strings.NewReader("assCA="))

	cookieJar.SetCookies(req.URL, loginCookies)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := client.Do(req)

	if err != nil {
		// handle error
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	return string(body[:])
}
