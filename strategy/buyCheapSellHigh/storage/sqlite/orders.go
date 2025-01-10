package sqlite

import (
	"fmt"
	"github.com/shatylos/trader/strategy/buyCheapSellHigh/storage/structs"
	"time"
)

// AddDomainOrderOnce add order to storage.
func (s *SqliteStorage) AddDomainOrderOnce(order structs.HistoryOrder) (bool, error) {

	isAdded := false

	orderTableName := getOrderTableName(s.setupCode)

	var count int
	q := `select count(*) from ` + orderTableName + ` where DomainOrderId = ?`
	err := s.db.QueryRow(q, order.DomainOrderId).Scan(&count)
	if err != nil {
		return false, err
	}

	if count == 0 {
		q = `INSERT INTO ` + orderTableName + ` 
			(DomainOrderId, FilledPrice, FilledQty, Side, CreatedTime, UpdatedTime, MainCurrencyAmountBefore, TradeCurrencyAmountBefore, AveragePrice)
			VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`

		_, err = s.db.Exec(
			q,
			order.DomainOrderId,
			order.FilledPrice,
			order.FilledQty,
			order.Side,
			order.CreatedTime,
			order.UpdatedTime,
			order.MainCurrencyAmountBefore,
			order.TradeCurrencyAmountBefore,
			order.AveragePrice,
		)

		if err != nil {
			return false, err
		}
		isAdded = true
	}

	return isAdded, nil
}

func (s *SqliteStorage) GetNotFilledHistoryOrders() ([]structs.HistoryOrder, error) {
	orderTableName := getOrderTableName(s.setupCode)

	q := fmt.Sprintf(`SELECT 
				DomainOrderId, 
				FilledPrice, 
				FilledQty, 
				Side, 
				CreatedTime, 
				UpdatedTime, 
				MainCurrencyAmountBefore, 
				TradeCurrencyAmountBefore, 
				AveragePrice 
			FROM %s 
			WHERE FilledPrice = 0 OR FilledQty = 0 OR Side = '' OR UpdatedTime = 0`,
		orderTableName)

	rows, err := s.db.Query(q)
	if err != nil {
		return nil, err
	}

	var orders []structs.HistoryOrder

	for rows.Next() {
		var order structs.HistoryOrder
		err := rows.Scan(
			&order.DomainOrderId,
			&order.FilledPrice,
			&order.FilledQty,
			&order.Side,
			&order.CreatedTime,
			&order.UpdatedTime,
			&order.MainCurrencyAmountBefore,
			&order.TradeCurrencyAmountBefore,
			&order.AveragePrice,
		)
		if err != nil {
			return nil, err
		}
		orders = append(orders, order)
	}

	return orders, nil
}

func (s *SqliteStorage) GetNotCalculatedHistoryOrders() ([]structs.HistoryOrder, error) {
	panic("not implemented")
}

func (s *SqliteStorage) GetLastCalculatedOrder() (*structs.HistoryOrder, error) {
	panic("not implemented")
}

func (s *SqliteStorage) GetCalculatedHistoryOrders(from time.Time, to time.Time) ([]structs.HistoryOrder, error) {
	panic("not implemented")
}

func (s *SqliteStorage) RemoveOrder(domainOrderId string) error {
	panic("not implemented")
}

func (s *SqliteStorage) UpdateOrder(order structs.HistoryOrder) error {
	panic("not implemented")
}
func (s *SqliteStorage) ResetHistoryOrderData() error {
	panic("not implemented")
}
