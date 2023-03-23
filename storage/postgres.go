package storage

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
	"github.com/plitn/wb_school_l0/models"
)

const (
	host        = "localhost"
	port        = 5432
	user        = "postgres"
	password    = "admin"
	dbname      = "postgres"
	insertOrder = `insert into "orders"("order_uid", "track_number",
                     "entry", "locale", "internal_signature", "customer_id",
                     "delivery_service", "shardkey", "sm_id", "date_created",
                     "oof_shard") values($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`
	insertDeliveries = `insert into "deliveries"("name", "phone", 
                     "zip", "city", "address", "region", "email",
                     "order_uid") values($1, $2, $3, $4, $5, $6, $7, $8)`
	insertItems = `insert into "items"("chrt_id", "track_number",
                    "price", "rid", "name", "sale", "size", "total_price",
                    "nm_id", "brand", "status") values($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`
	insertPayments = `insert into "payments"("order_id", "transaction",
                         "request_id", "currency", "provider",
                         "amount", "payment_dt", "bank", "delivery_cost",
                         "goods_total", "custom_fee") values($1, $2, $3, $4, $5, $6, $7, $8,
                                                             $9, $10, $11)`
	selectOrder    = `select * from "orders";`
	selectDelivery = `select * from "deliveries" where order_uid = $1;`
	selectPayment  = `select * from "payments" where order_id=$1;`
	selectItems    = `select * from "items" where track_number=$1;`
)

type PostgreConn struct {
	db *sql.DB
}

func (p *PostgreConn) CloseDb() {
	err := p.db.Close()
	if err != nil {
		fmt.Println(err)
		return
	}

	if p.db == nil {
		fmt.Println("db is nil")
		return
	}
}

func (p *PostgreConn) SelectOrder() *sql.Rows {
	rows, err := p.db.Query(selectOrder)
	if err != nil {
		fmt.Println("select orders error")
	}
	return rows
}

func (p *PostgreConn) SelectDelivery(orderUid string) *sql.Row {
	rows := p.db.QueryRow(selectDelivery, orderUid)
	return rows
}

func (p *PostgreConn) SelectPayment(orderUid string) *sql.Row {
	rows := p.db.QueryRow(selectPayment, orderUid)
	return rows
}

func (p *PostgreConn) SelectItems(orderUid string) *sql.Rows {
	rows, err := p.db.Query(selectItems, orderUid)
	if err != nil {
		fmt.Println("select items error")
	}
	return rows
}

func (p *PostgreConn) InsertOrder(order models.Order) error {
	_, err := p.db.Exec(insertOrder, order.OrderUid, order.TrackNumber, order.Entry, order.Locale,
		order.InternalSignature, order.CustomerId, order.DeliveryService, order.Shardkey, order.SmId,
		order.DateCreated, order.OofShard)
	if err != nil {
		fmt.Println(err)
	}
	return err
}

func (p *PostgreConn) InsertDelivery(order models.Order) error {
	_, err := p.db.Exec(insertDeliveries, order.Delivery.Name, order.Delivery.Phone,
		order.Delivery.Zip, order.Delivery.City, order.Delivery.Address,
		order.Delivery.Region, order.Delivery.Email, order.OrderUid)
	if err != nil {
		fmt.Println(err)
	}
	return err
}

func (p *PostgreConn) InsertPayment(order models.Order) error {
	_, err := p.db.Exec(insertPayments, order.OrderUid, order.Payment.Transaction,
		order.Payment.RequestId, order.Payment.Currency, order.Payment.Provider,
		order.Payment.Amount, order.Payment.PaymentDt, order.Payment.Bank,
		order.Payment.DeliveryCost, order.Payment.GoodsTotal, order.Payment.CustomFee)
	if err != nil {
		fmt.Println(err)
	}
	return err
}

func (p *PostgreConn) InsertItems(order models.Order) error {
	for _, item := range order.Items {
		_, err := p.db.Exec(insertItems, item.ChrtId, item.TrackNumber, item.Price,
			item.Rid, item.Name, item.Sale, item.Size, item.TotalPrice, item.NmId, item.Brand, item.Status)
		if err != nil {
			fmt.Println(err)
			return err
		}
	}
	return nil
}

func NewConn() *PostgreConn {
	return &PostgreConn{}
}

func (p *PostgreConn) Conn() error {
	connString := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	var err error
	p.db, err = sql.Open("postgres", connString)
	if err != nil {
		return err
	}

	err = p.db.Ping()
	if err != nil {
		return err
	}

	return nil
}
