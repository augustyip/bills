package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/augustyip/bills/services"
	log "github.com/sirupsen/logrus"
)

// Certification struct
type Certification struct {
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

	file, _ := os.Open("cert.json")
	decoder := json.NewDecoder(file)
	certs := make([]Certification, 0)
	err := decoder.Decode(&certs)
	if err != nil {
		fmt.Println("error:", err)
	}

	c := make(chan string, len(certs))

	for _, cert := range certs {

		switch s := cert.Service; s {
		case "towngas":
			log.Info("Starting to run Towngas service...")
			towngas := services.Towngas{cert.Username, cert.Password}
			go services.GetNewsNoticeAsync(towngas, c)
			// fmt.Printf(r)

		case "clp":
			log.Info("Starting to run CLP service...")
			clp := services.Clp{cert.Username, cert.Password}
			go services.GetServiceDashboard(clp, c)

		case "wsd":
			log.Info("Starting to run WSD service...")
			wsd := services.Wsd{cert.Username, cert.Password}
			go services.ElectronicBill(wsd, c)

		}
	}
	for i := range c {
		log.Info(i)
	}

}
