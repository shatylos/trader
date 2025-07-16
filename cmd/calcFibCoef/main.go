package main

import "fmt"

func main() {
	minPrice := 70000.0
	maxPrice := 80000.0
	qty := 0.01

	qty1Percent := 20.0
	qty2Percent := 35.0
	qty3Percent := 45.0
	ep1Coef := 0.5
	ep2Coef := 0.26
	ep3Coef := 0.05
	tp1Coef := 1.25
	tp2Coef := 1.1
	tp3Coef := 0.95
	slCoef := -0.1

	qty1 := qty / 100.0 * qty1Percent
	qty2 := qty / 100.0 * qty2Percent
	qty3 := qty / 100.0 * qty3Percent

	tp1Price := (maxPrice-minPrice)*tp1Coef + minPrice
	tp2Price := (maxPrice-minPrice)*tp2Coef + minPrice
	tp3Price := (maxPrice-minPrice)*tp3Coef + minPrice

	ep1Price := (maxPrice-minPrice)*ep1Coef + minPrice
	ep2Price := (maxPrice-minPrice)*ep2Coef + minPrice
	ep3Price := (maxPrice-minPrice)*ep3Coef + minPrice

	slPrice := (maxPrice-minPrice)*slCoef + minPrice

	tp1Order := (tp1Price * qty1) - (ep1Price * qty1)

	tp2Order1 := (tp2Price * qty1) - (ep1Price * qty1)
	tp2Order2 := (tp2Price * qty2) - (ep2Price * qty2)

	tp3Order1 := (tp3Price * qty1) - (ep1Price * qty1)
	tp3Order2 := (tp3Price * qty2) - (ep2Price * qty2)
	tp3Order3 := (tp3Price * qty3) - (ep3Price * qty3)

	sl1 := (slPrice * qty1) - (ep1Price * qty1)
	sl2 := (slPrice * qty2) - (ep2Price * qty2)
	sl3 := (slPrice * qty3) - (ep3Price * qty3)

	fmt.Printf("TP 1 order: %f\n", tp1Order)
	fmt.Printf("TP 2 orders: %f\n", tp2Order1+tp2Order2)
	fmt.Printf("TP 3 orders: %f\n", tp3Order1+tp3Order2+tp3Order3)
	fmt.Printf("SL: %f\n", sl1+sl2+sl3)
}
