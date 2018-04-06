package services

import (
	"fmt"
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
	preLoginPageDoc, _ := goquery.NewDocumentFromReader(preLoginPageResp.Body)
	// For token hidden field and get toekn value
	preLoginPageDoc.Find("input[name='org.apache.struts.taglib.html.TOKEN']").Each(func(i int, s *goquery.Selection) {
		token, _ = s.Attr("value")
	})

	// Login action
	var preLoginCookies = preLoginPageResp.Cookies()
	var loginBody = "org.apache.struts.taglib.html.TOKEN=" + token + "&userID=" + c.Username + "&password=" + c.Password + "&submit=%E9%81%9E%E4%BA%A4"
	loginReq, err := http.NewRequest("POST", "https://www.esd.wsd.gov.hk/esd/login.do", strings.NewReader(loginBody))
	if err != nil {
		// handle error
	}
	cookieJar.SetCookies(loginReq.URL, preLoginCookies)
	loginReq.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	loginResp, err := client.Do(loginReq)
	if err != nil {
		// handle error
	}

	// Get electronicBill page
	var cookies = loginResp.Cookies()
	electronicBillInitReq, err := http.NewRequest("GET", "https://www.esd.wsd.gov.hk/esd/bnc/electronicBill/init.do?pageFlag=1", nil)
	cookieJar.SetCookies(electronicBillInitReq.URL, cookies)
	electronicBillInitResp, err := client.Do(electronicBillInitReq)
	if err != nil {
		// handle error
	}
	defer electronicBillInitResp.Body.Close()

	var accountID string
	electronicBillInitDoc, _ := goquery.NewDocumentFromReader(electronicBillInitResp.Body)

	// Get account ID from the form
	// TODO support multiple account
	electronicBillInitDoc.Find("input[name='accountID']").Each(func(i int, s *goquery.Selection) {
		accountID, _ = s.Attr("value")
	})

	// Submit to processSelectAccount.do page
	var selectAccount = "org.apache.struts.taglib.html.TOKEN=" + token + "&accountID=" + accountID + "&page=2&submit=Next"
	processSelectAccountReq, err := http.NewRequest("POST", "https://www.esd.wsd.gov.hk/esd/bnc/electronicBill/processSelectAccount.do", strings.NewReader(selectAccount))

	cookieJar.SetCookies(processSelectAccountReq.URL, cookies)
	processSelectAccountReq.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	processSelectAccountResp, err := client.Do(processSelectAccountReq)
	if err != nil {
		// handle error
	}
	defer processSelectAccountResp.Body.Close()

	// Submit to processSelectBillServices.do page
	var selectBillServices = "org.apache.struts.taglib.html.TOKEN=" + token + "&services=summary&submit=Next"
	processSelectBillServicesReq, err := http.NewRequest("POST", "https://www.esd.wsd.gov.hk/esd/bnc/electronicBill/processSelectBillServices.do", strings.NewReader(selectBillServices))

	cookieJar.SetCookies(processSelectBillServicesReq.URL, cookies)
	processSelectBillServicesReq.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	processSelectBillServicesResp, err := client.Do(processSelectBillServicesReq)
	if err != nil {
		// handle error
	}
	defer processSelectBillServicesResp.Body.Close()
	var billTable string
	processSelectBillServicesDoc, _ := goquery.NewDocumentFromReader(processSelectBillServicesResp.Body)
	processSelectBillServicesDoc.Find("table.style_table").Each(func(i int, s *goquery.Selection) {
		billTable, _ = s.Html()
	})
	return billTable
}
