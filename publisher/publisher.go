package main

import (
	"L0/internal"
	"log"
	"time"

	json "github.com/mailru/easyjson"
	"github.com/nats-io/stan.go"
)

func publishEvent(data *internal.Order) {
	sc, err := stan.Connect(
		"test-cluster",
		"delivery-service-publisher",
		stan.NatsURL("nats://127.0.0.1:4222"),
	)
	if err != nil {
		log.Println(err)
		return
	}
	defer sc.Close()

	channel := "order-channel"
	byteData, _ := json.Marshal(data)

	err = sc.Publish(channel, byteData)
	if err != nil {
		return
	}
	log.Println("Published message on channel: " + channel)
}

func main() {
	order := internal.Order{
		OrderUid:    "3",
		TrackNumber: "111",
		Entry:       "awdaa",
		Delivery: struct {
			Name    string `json:"name"`
			Phone   string `json:"phone"`
			Zip     string `json:"zip"`
			City    string `json:"city"`
			Address string `json:"address"`
			Region  string `json:"region"`
			Email   string `json:"email"`
		}{"asd", "1", "asf", "asf", "asf", "dfh", "dfh"},
		Payment: struct {
			Transaction  string    `json:"transaction"`
			RequestId    string    `json:"request_id"`
			Currency     string    `json:"currency"`
			Provider     string    `json:"provider"`
			Amount       int       `json:"amount"`
			PaymentDt    time.Time `json:"payment_dt"`
			Bank         string    `json:"bank"`
			DeliveryCost int       `json:"delivery_cost"`
			GoodsTotal   int       `json:"goods_total"`
			CustomFee    int       `json:"custom_fee"`
		}{"qwer", "zsxdc", "afaf", "dhg", 1, time.Now(), "fthdf", 4, 1, 1},
		Items: []internal.Item{
			internal.Item{
				ChrtId:      1,
				TrackNumber: "a",
				Price:       1,
				Rid:         "a",
				Name:        "a",
				Sale:        1,
				Size:        "a",
				TotalPrice:  1,
				NmId:        1,
				Brand:       "a",
				Status:      1,
			},
			internal.Item{
				ChrtId:      3,
				TrackNumber: "c",
				Price:       3,
				Rid:         "c",
				Name:        "c",
				Sale:        3,
				Size:        "c",
				TotalPrice:  3,
				NmId:        3,
				Brand:       "c",
				Status:      3,
			},
		},
	}
	publishEvent(&order)
}
