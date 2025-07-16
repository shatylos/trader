package main

import (
	"fmt"
	"github.com/shatylos/trader/tools/math"
)

func main() {

	//Не виводячи 13 міс (1 рік 1 міс) до 2071
	//Не виводячи 23 міс (1 рік 11 міс) до 5372
	//
	//Не виводячи			37 міс (3 роки  2 міс.) до 20402
	//Виводячи 50% з 500	44 міс (3 роки  9 міс.) до 20729
	//Виводячи 50% з 200	60 міс (5 років 1 міс.) до 20518

	totalAmount := 2500.0
	months := 12 * 12
	percentToWidthdraw := 0.0
	minAmountToWidthdraw := 0.0
	percentRevenue := 2.5

	for month := 1; month <= months; month++ {
		revenue := math.Mul(math.Div(totalAmount, 100), percentRevenue)
		withdraw := math.Mul(math.Div(revenue, 100), percentToWidthdraw)
		if withdraw < minAmountToWidthdraw {
			withdraw = 0
		}
		increaseDepo := revenue - withdraw
		totalAmount += increaseDepo
		year := month/12 + 1
		fmt.Printf("Year: %d, Month: %d, Withdraw: %.2f, Total Amount after: %.2f \n", year, month, withdraw, totalAmount)
	}
}
