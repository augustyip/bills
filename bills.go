package main

import (
  "fmt"
  "net/http"
  "net/url"
  "io/ioutil"
  )
func main() {
  resp, err := http.PostForm("https://eservice.towngas.com/EAccount/Login/SignIn", url.Values{"LoginID": {"username"},"password": {"password"}})
  if err != nil {
    // handle error
  }
  defer resp.Body.Close()
  body, err := ioutil.ReadAll(resp.Body)
  bodyContent := string(body[:])
  fmt.Printf(bodyContent)
}
