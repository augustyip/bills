package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/augustyip/bills/services"
	log "github.com/sirupsen/logrus"
)

type account struct {
	Service  string
	Username string
	Password string
}

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
	accounts := make([]account, 0)
	err := decoder.Decode(&accounts)
	if err != nil {
		fmt.Println("error:", err)
	}

	c := make(chan string, 3)

	for _, acc := range accounts {

		switch s := acc.Service; s {
		case "towngas":
			log.Info("Starting to run Towngas service...")
			towngas := services.Towngas{acc.Username, acc.Password}
			go services.GetNewsNoticeAsync(towngas, c)

		case "clp":
			log.Info("Starting to run CLP service...")
			clp := services.Clp{acc.Username, acc.Password}
			go services.GetServiceDashboard(clp, c)

		case "wsd":
			log.Info("Starting to run WSD service...")
			wsd := services.Wsd{acc.Username, acc.Password}
			go services.ElectronicBill(wsd, c)

		}
	}
	for i := range c {
		log.Info(i)
	}

}
