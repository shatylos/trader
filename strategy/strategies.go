package strategy

import (
	"bitbucket.org/shatylos/trader/strategy/buyCheapSellHigh"
	"bitbucket.org/shatylos/trader/strategy/struct"
	"bitbucket.org/shatylos/trader/utils"
	"fmt"
)

func GetStrategyByCode(code string) (_struct.StrategyInterface, error) {
	switch code {
	case "buy_cheap_sell_high":
		return &buyCheapSellHigh.BuyCheapSellHigh{}, nil
	}
	return nil, utils.AppError{
		Message: fmt.Sprintf("strategy with code \"%s\" not implemented", code),
	}
}
