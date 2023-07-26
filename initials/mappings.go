package initials

import (
	"log"
	"github.com/shopspring/decimal"
)

type ListingData struct {
	Name				string
	Price				decimal.Decimal
}

func convertToDecimal(value string) decimal.Decimal {
	n, err := decimal.NewFromString(value)
	if err != nil {
		log.Println(err)
	}
	return n
}

func ListingsMapping() []ListingData {
	return []ListingData{
		{Name: "Brand New royal Enfield 250 CC For special Sale", Price: convertToDecimal("6000.00")},
		{Name: "Wedding wow Exclusive Cupple Ring (S2022)", Price: convertToDecimal("3000.00")},
		{Name: "Toyota AIGID A Class Hatchback Sale", Price: convertToDecimal("2000.00")},
		{Name: "Havit HV-G61 USB Black Double Game With Vibrat", Price: convertToDecimal("5000.85")},
		{Name: "Brand New Honda CBR 600 RR For Sale (2022)", Price: convertToDecimal("9000.00")},
		{Name: "IPhone 11 Pro Max All Variants Available For Sale", Price: convertToDecimal("4000.00")},
	}
}