package services

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
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

	cookiesResp, err := client.Get("https://services.clp.com.hk/zh/login/index.aspx")
	if err != nil {
		// handle error
	}
	defer cookiesResp.Body.Close()
	// fmt.Println(loginPageResp.Cookies()[0].String())
	// for _, cookie := range cookiesResp.Cookies() {
	// 	fmt.Println(cookie.String())
	// }

	var clpCookies = cookiesResp.Cookies()

	// loginPageBody, err := ioutil.ReadAll(loginPageResp.Body)
	// fmt.Println(string(loginPageBody[:]))

	var loginBody = "username=" + c.Username + "&password=" + c.Password
	loginReq, err := http.NewRequest("POST", "https://services.clp.com.hk/Service/ServiceLogin.ashx", strings.NewReader(loginBody))
	if err != nil {
		// handle error
	}
	cookieJar.SetCookies(loginReq.URL, clpCookies)
	loginResp, err := client.Do(loginReq)
	if err != nil {
		// handle error
	}
	loginRespBody, err := ioutil.ReadAll(loginResp.Body)
	fmt.Println(string(loginRespBody[:]))

	defer loginResp.Body.Close()

	var loginedCookies = loginResp.Cookies()

	// fmt.Println(loginCookies[0].String())
	req, err := http.NewRequest("POST", "https://services.clp.com.hk/Service/ServiceDashboard.ashx", strings.NewReader("assCA="))

	cookieJar.SetCookies(req.URL, loginedCookies)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := client.Do(req)

	if err != nil {
		// handle error
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	return string(body[:])
}
