package sqlite

import (
	"bitbucket.org/shatylos/trader/strategy/buyCheapSellHigh/storage/structs"
	"fmt"
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

func (s *SqliteStorage) GetNotCalculatedDomainOrders() ([]structs.HistoryOrder, error) {
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
			WHERE AveragePrice IS NULL
				OR AveragePrice = 0`,
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
