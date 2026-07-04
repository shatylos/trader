package strategy

import (
	"fmt"
	"github.com/shatylos/trader/internal/strategy/buyCheapSellHigh"
	"github.com/shatylos/trader/internal/strategy/fibonacci"
	"github.com/shatylos/trader/internal/strategy/investor"
	"github.com/shatylos/trader/internal/strategy/scalper"
	"github.com/shatylos/trader/internal/strategy/struct"
	"github.com/shatylos/trader/internal/strategy/vwapReversion"
	"github.com/shatylos/trader/tools"
)

func GetStrategyByCode(code string) (_struct.StrategyInterface, error) {
	switch code {
	case "buy_cheap_sell_high":
		return &buyCheapSellHigh.BuyCheapSellHigh{}, nil
	case "fibonacci":
		return &fibonacci.Fibonacci{}, nil
	case "investor":
		return &investor.Investor{}, nil
	case "vwap_reversion":
		return &vwapReversion.VwapReversion{}, nil
	case "scalper":
		return &scalper.Scalper{}, nil
	}
	return nil, tools.AppError{
		Message: fmt.Sprintf("strategy with code \"%s\" not implemented", code),
	}
}
