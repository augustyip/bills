package services

import "github.com/augustyip/bills/model"

// GetSummy return the service summy
func GetSummy(acc model.Account, c chan string) {
	switch acc.Service {
	case "towngas":
		go GetNewsNoticeAsync(acc, c)

	case "clp":
		go GetServiceDashboard(acc, c)

	case "wsd":
		go ElectronicBill(acc, c)

	}
}
