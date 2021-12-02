package common

import (
	"strconv"

	"github.com/FreemanFeng/dragon/deps/calc/compute"
	"github.com/shopspring/decimal"
)

func Calculate(input string) string {
	res, err := compute.Evaluate(input)
	if err != nil {
		return input
	}
	s := strconv.FormatFloat(res, 'G', -1, 64)
	num, err := decimal.NewFromString(s)
	if err != nil {
		Log(err)
		return input
	}
	return num.String()
}
