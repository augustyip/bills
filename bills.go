package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/augustyip/bills/model"

	"github.com/augustyip/bills/services"
	log "github.com/sirupsen/logrus"
)

// Bill struct
type Bill struct {
	AccountNo       string
	Balance         string
	LastPaymentDate string
}

func init() {
	log.SetLevel(log.DebugLevel)
}

func main() {

	file, _ := os.Open("accounts.json")
	decoder := json.NewDecoder(file)
	accounts := make([]model.Account, 0)
	err := decoder.Decode(&accounts)
	if err != nil {
		fmt.Println("error:", err)
	}

	c := make(chan string, 3)

	for _, acc := range accounts {
		services.GetSummy(acc, c)
	}
	for i := range c {
		log.Info(i)
	}

}
