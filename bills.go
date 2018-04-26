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
			towngasAcc := services.Towngas{
				Username: acc.Username,
				Password: acc.Password,
			}
			go services.GetNewsNoticeAsync(towngasAcc, c)

		case "clp":
			clpAcc := services.Clp{
				Username: acc.Username,
				Password: acc.Password,
			}
			go services.GetServiceDashboard(clpAcc, c)

		case "wsd":
			wsdAcc := services.Wsd{
				Username: acc.Username,
				Password: acc.Password,
			}
			go services.ElectronicBill(wsdAcc, c)

		}
	}
	for i := range c {
		log.Info(i)
	}

}
