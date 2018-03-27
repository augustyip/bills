package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/augustyip/bills/services"
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

func main() {

	file, _ := os.Open("cert.json")
	decoder := json.NewDecoder(file)
	certs := make([]Certification, 0)
	err := decoder.Decode(&certs)
	if err != nil {
		fmt.Println("error:", err)
	}
	for _, cert := range certs {

		switch s := cert.Service; s {
		case "towngas":
			towngas := services.Towngas{cert.Username, cert.Password}
			r := services.GetNewsNoticeAsync(towngas)
			fmt.Printf(r)

		case "clp":
			clp := services.Clp{cert.Username, cert.Password}
			var clpBill services.Bill
			clpBill.GetServiceDashboard(clp)
			fmt.Printf("%+v\n", clpBill)

			// case "wsd":
			// 	wsd := services.Wsd{cert.Username, cert.Password}
			// 	r := services.ElectronicBill(wsd)
			// 	fmt.Printf(r)
		}
	}
}
