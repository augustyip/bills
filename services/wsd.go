package services

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"strings"

	"github.com/PuerkitoBio/goquery"
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
	var token string

	// Get prelogin page for token
	preLoginPageResp, err := http.Get("https://www.esd.wsd.gov.hk/esd/preLogin.do?pageFlag=1")
	if err != nil {
		// handle error
	}
	defer preLoginPageResp.Body.Close()
	preLoginPageDoc, err := goquery.NewDocumentFromReader(preLoginPageResp.Body)
	if err != nil {
		// handle error
	}
	// For token hidden field and get toekn value
	preLoginPageDoc.Find("input[name='org.apache.struts.taglib.html.TOKEN']").Each(func(i int, s *goquery.Selection) {
		token, _ = s.Attr("value")
	})

	var cookies = preLoginPageResp.Cookies()

	var loginBody = "org.apache.struts.taglib.html.TOKEN=" + token + "&userID=" + c.Username + "&password=" + c.Password + "&submit=%E9%81%9E%E4%BA%A4"
	fmt.Println(loginBody)
	loginReq, err := http.NewRequest("POST", "https://www.esd.wsd.gov.hk/esd/login.do", strings.NewReader(loginBody))
	if err != nil {
		// handle error
	}
	cookieJar.SetCookies(loginReq.URL, cookies)

	// loginReq.Header.Set("X-CSRFToken", csrfToken)
	// loginReq.Header.Set("X-Requested-With", "XMLHttpRequest")
	loginReq.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	loginResp, err := client.Do(loginReq)
	fmt.Println(loginResp.StatusCode)

	// return "xxx"
	body, err := ioutil.ReadAll(loginResp.Body)
	return string(body[:])
}
