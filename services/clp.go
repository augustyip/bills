package services

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"strings"

	"github.com/augustyip/bills/model"

	log "github.com/sirupsen/logrus"
)

// GetServiceDashboard get latest info
func GetServiceDashboard(acc model.Account, channel chan string) {
	log.Debug("[CLP] Starting to run CLP service...")

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

	log.Debug("[CLP] Get login page for the CSRF token.")
	cookiesResp, err := client.Get("https://services.clp.com.hk/zh/login/index.aspx")
	if err != nil {
		log.Error(err)
	}
	defer cookiesResp.Body.Close()
	for _, cookie := range cookiesResp.Cookies() {
		if cookie.Name == "K2Cie90hi___AntiXsrfToken" {
			csrfToken = cookie.Value
		}
	}

	log.Debug("[CLP] Logging into...")
	var loginBody = "username=" + acc.Username + "&password=" + acc.Password
	loginReq, err := http.NewRequest("POST", "https://services.clp.com.hk/Service/ServiceLogin.ashx", strings.NewReader(loginBody))
	if err != nil {
		log.Error(err)
	}

	loginReq.Header.Set("X-CSRFToken", csrfToken)
	loginReq.Header.Set("X-Requested-With", "XMLHttpRequest")
	loginReq.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")

	loginResp, err := client.Do(loginReq)
	if err != nil {
		log.Error(err)
	}
	defer loginResp.Body.Close()

	var loginedCookies = loginResp.Cookies()

	log.Debug("[CLP] Getting service dashboard info...")
	req, err := http.NewRequest("POST", "https://services.clp.com.hk/Service/ServiceDashboard.ashx", strings.NewReader("assCA="))

	cookieJar.SetCookies(req.URL, loginedCookies)
	req.Header.Set("X-CSRFToken", csrfToken)
	req.Header.Set("X-Requested-With", "XMLHttpRequest")
	resp, err := client.Do(req)

	if err != nil {
		log.Error(err)
	}
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	channel <- string(data[:])
}
