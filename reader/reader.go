package reader

import (
	"encoding/json"
	_ "github.com/lib/pq"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/stan.go"
	"github.com/plitn/wb_school_l0/models"
	"github.com/plitn/wb_school_l0/storage"
	"log"
)

type NutsReader struct {
	nc    *nats.Conn
	sc    stan.Conn
	cache *storage.CacheStruct
}

// NewNatsReader конструктор
func NewNatsReader(c *storage.CacheStruct) *NutsReader {
	nr := &NutsReader{
		cache: c,
	}
	return nr

}

// Init инит соединения и подписка на слушание
func (r *NutsReader) Init() error {
	host := "localhost:4223"
	natsconn, err := nats.Connect(host)

	if err != nil {
		return err
	}

	stanconn, err := stan.Connect("test-cluster", "reader", stan.NatsConn(natsconn))
	if err != nil {
		return err
	}

	r.nc = natsconn
	r.sc = stanconn
	_, err = r.sc.Subscribe("orders_data", r.readData, stan.DeliverAllAvailable())
	log.Printf("Reading on host %s, cluster: %s", host, "test-cluster")
	return nil
}

// тут слушаем дату и обрабатываем, запихиваем в бд и в кеш
func (r *NutsReader) readData(message *stan.Msg) {
	var order models.Order

	// намаршалим полученные данные
	err := json.Unmarshal(message.Data, &order)
	if err != nil {
		log.Printf("unmarshal error: %v", err)
		return
	}

	// записываем в кеш полученную модель
	err = r.cache.SetData(order)
	if err != nil {
		log.Printf("cache insert error: %v", err)
	}

	// подкл. к бд
	postgres := storage.NewConn()
	err = postgres.Conn()
	if err != nil {
		log.Printf("conn error: %v", err)
		return
	}

	// записываем в бд все поля модели
	err = postgres.InsertOrder(order)
	if err != nil {
		log.Printf("order insert error: %v", err)
		return
	}
	err = postgres.InsertDelivery(order)
	if err != nil {
		log.Printf("delivery insert error : %v", err)
		return
	}

	err = postgres.InsertPayment(order)
	if err != nil {
		log.Printf("payment insert error: %v", err)
		return
	}
	err = postgres.InsertItems(order)
	if err != nil {
		log.Printf("items insert error: %v", err)
		return
	}

	// закрываем соединение с бд
	postgres.CloseDb()
}
