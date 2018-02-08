package main

import (
  "fmt"
  "net/http"
  "net/url"
  "io/ioutil"
  "encoding/json"
  "os"
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

  resp, err := http.PostForm("https://eservice.towngas.com/EAccount/Login/SignIn", url.Values{"LoginID": {"username"},"password": {"password"}})
  if err != nil {
    // handle error
  }
  defer resp.Body.Close()
  body, err := ioutil.ReadAll(resp.Body)
  bodyContent := string(body[:])
  fmt.Printf(bodyContent)
}
