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
			r := services.GetServiceDashboard(clp)
			fmt.Printf(r)
		}
	}
}
