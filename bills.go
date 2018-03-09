package main

import (
	"fmt"

	"github.com/augustyip/bills/services"
)

func main() {
	r := services.GetNewsNoticeAsync()
	fmt.Printf(r)
}
