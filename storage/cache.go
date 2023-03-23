package storage

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/plitn/wb_school_l0/models"
	"log"
)

type CacheStruct struct {
	cacheMap map[string]models.Order
}

func NewCache() *CacheStruct {
	return &CacheStruct{}
}

func (c *CacheStruct) Init() error {
	c.cacheMap = make(map[string]models.Order)
	postgre := NewConn()
	err := postgre.Conn()
	if err != nil {
		return err
	}

	// позорный костыль для скана
	var orderU string

	// тут восстанавливаем все из бд в кеш
	orderRows := postgre.SelectOrder()
	defer func(orderRows *sql.Rows) {
		err := orderRows.Close()
		if err != nil {
			log.Println("conn closing error")
			return
		}
	}(orderRows)
	for orderRows.Next() {
		var order models.Order
		err := orderRows.Scan(&order.OrderUid, &order.TrackNumber, &order.Entry,
			&order.Locale, &order.InternalSignature, &order.CustomerId, &order.DeliveryService,
			&order.Shardkey, &order.SmId, &order.DateCreated, &order.OofShard)
		if err != nil {
			log.Fatal(err)
		}
		deliveryRows := postgre.SelectDelivery(order.OrderUid)
		var delivery models.Delivery
		fmt.Println(order)
		err = deliveryRows.Scan(&delivery.Name, &delivery.Phone, &delivery.Zip,
			&delivery.City, &delivery.Address, &delivery.Region,
			&delivery.Email, &orderU)
		if err != nil {
			fmt.Println(err)
			fmt.Println("scan error")
			return err
		}
		order.Delivery = delivery
		fmt.Println(order)
		paymentRows := postgre.SelectPayment(order.OrderUid)

		err = paymentRows.Scan(&order.Payment.Transaction, &order.Payment.RequestId,
			&order.Payment.Currency, &order.Payment.Provider, &order.Payment.Amount,
			&order.Payment.PaymentDt, &order.Payment.Bank, &order.Payment.DeliveryCost,
			&order.Payment.GoodsTotal, &order.Payment.CustomFee, &orderU)
		if err != nil {
			return err
		}

		itemsRows := postgre.SelectItems(order.TrackNumber)

		for itemsRows.Next() {
			var items models.Items
			err := itemsRows.Scan(&items.ChrtId, &items.TrackNumber, &items.Price,
				&items.Rid, &items.Name, &items.Sale, &items.Size, &items.TotalPrice,
				&items.NmId, &items.Brand, &items.Status)
			if err != nil {
				return err
			}
			order.Items = append(order.Items, items)
		}
		fmt.Println(order)
	}
	return nil
}

func (c *CacheStruct) GetData(id string) (*models.Order, error) {
	value, exists := c.cacheMap[id]
	if !exists {
		fmt.Printf("order with id %s does not exist in cache", id)
		return nil, errors.New("order with id does not exist")
	}
	return &value, nil
}

func (c *CacheStruct) SetData(order models.Order) error {
	c.cacheMap[order.OrderUid] = order
	return nil
}
