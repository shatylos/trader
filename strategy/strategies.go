package strategy

import (
	"fmt"
	"github.com/shatylos/trader/strategy/buyCheapSellHigh"
	"github.com/shatylos/trader/strategy/struct"
	"github.com/shatylos/trader/utils"
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
