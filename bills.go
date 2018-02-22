package main

import (
  "fmt"
  "net/http"
  "net/url"
  "io/ioutil"
  "encoding/json"
  "os"
  "net/http/cookiejar"
)

type Certification struct {
  Service  string
  Username  string
  Password  string
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
    
    loginResp, err := client.PostForm("https://eservice.towngas.com/EAccount/Login/SignIn", url.Values{"LoginID": {cert.Username},"password": {cert.Password}})
    if err != nil {
      // handle error
    }
    defer loginResp.Body.Close()

    resp, err := http.Get("https://eservice.towngas.com/en/BillingUsage/NewsNotices")
    if err != nil {
      // handle error
    }
    defer resp.Body.Close()
    body, err := ioutil.ReadAll(resp.Body)
    bodyContent := string(body[:])
    fmt.Printf(bodyContent)
  }

}
