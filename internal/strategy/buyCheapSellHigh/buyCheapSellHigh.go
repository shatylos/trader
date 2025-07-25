package buyCheapSellHigh

import (
	"fmt"
	"github.com/shatylos/trader/internal/domain"
	"github.com/shatylos/trader/internal/strategy/struct"
	tradeConst "github.com/shatylos/trader/internal/trading/constant"
	"github.com/shatylos/trader/tools/logger"
	"time"
)

type BuyCheapSellHigh struct {
	Id                        string
	CoinPare                  string
	Domain                    domain.SpotDomainInterface
	MainCurrency              string
	TradeCurrency             string
	Resolution                string
	TimeoutSeconds            time.Duration
	CostRanges                []int64
	PercentRanges             []int64
	LongTermMaxPrice          float64
	LongTermMinPrice          float64
	LongTermPercentBuffer     float64
	MaxQtyDiff                float64
	MinQty                    float64
	MainCurrencyPrecision     int64
	PurchaseVolumePrecision   int64
	PurchasePricePrecision    int64
	CommissionPercent         float64
	MinutesToReducePriceRange int64
	AvgPriceCandleLimit       int64
	AvgPriceCandleOffset      int64
	ManualBuyPriceBeforeStart float64
	WithdrawPercent           float64
	Enabled                   bool
}

var _ _struct.StrategyInterface = (*BuyCheapSellHigh)(nil)

func (s *BuyCheapSellHigh) GetId() string {
	return s.Id
}

func (s *BuyCheapSellHigh) GetTitle() string {
	if !s.Enabled {
		return fmt.Sprintf("Buy Cheap Sell High: %s (%s) (DISABLED)", s.Id, s.CoinPare)
	}
	return fmt.Sprintf("Buy Cheap Sell High: %s (%s)", s.Id, s.CoinPare)
}

func (s *BuyCheapSellHigh) DoAction() error {
	if !s.Enabled {
		return nil
	}

	openOrders, err := s.getOpenOrders()
	if err != nil {
		return err
	}

	if len(openOrders) >= 2 {
		err := s.cancelOldOrdersWithBigRanges(openOrders)
		if err != nil {
			return err
		}
		return nil
	}

	if len(openOrders) == 1 {
		err = s.cancelAllOrder(openOrders[0])
		if err != nil {
			return err
		}
		err = s.fillPrices()
		if err != nil {
			return err
		}
		err = s.calculateHistoryOrderValues()
		if err != nil {
			return err
		}
	}

	historyOrders, err := s.getHistoryOrders()
	if err != nil {
		return err
	}

	baseCurrencyBalance, tradeCurrencyBalance, err := s.getBalances()
	if err != nil {
		return err
	}

	buyPrice, buyQty, sellPrice, sellQty, err := s.getPricesAndQtysToNewOrders(historyOrders, baseCurrencyBalance, tradeCurrencyBalance)
	if err != nil {
		return err
	}

	if buyPrice > 0 && buyQty > 0 {
		orderId, err := s.setLimitOrder(buyPrice, buyQty, tradeConst.SideBuy)
		if err != nil {
			return err
		}
		err = s.setOrderToStorage(orderId, baseCurrencyBalance, tradeCurrencyBalance)
		if err != nil {
			return err
		}
	} else {
		logger.Info(fmt.Sprintf("[Buy Cheap Sell High] unexpected values for buy orders: buyPrice: %f, buyQty: %f", buyPrice, buyQty))
	}
	if sellPrice > 0 && sellQty > 0 {
		orderId, err := s.setLimitOrder(sellPrice, sellQty, tradeConst.SideSell)
		if err != nil {
			return err
		}
		err = s.setOrderToStorage(orderId, baseCurrencyBalance, tradeCurrencyBalance)
		if err != nil {
			return err
		}
	} else {
		logger.Info(fmt.Sprintf("[Buy Cheap Sell High] unexpected values for sell orders: sellPrice: %f, sellQty: %f", sellPrice, sellQty))
	}

	return nil
}

func (s *BuyCheapSellHigh) Wait() {
	time.Sleep(time.Second * s.TimeoutSeconds)
}
