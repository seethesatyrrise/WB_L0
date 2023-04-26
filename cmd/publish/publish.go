package main

import (
	"encoding/json"
	"fmt"
	"github.com/nats-io/nats.go"
	"http-nats-psql/internal/models"
	"os"
	"time"
)

func main() {
	url := os.Getenv("NATS_URL")
	if url == "" {
		url = nats.DefaultURL
	}

	nc, _ := nats.Connect(url)
	defer nc.Drain()

	js, _ := nc.JetStream()

	//deleteStream(js, "MSG")
	//addStream(js, "MSG")

	publish(js, getOrder1())
}

func publish(js nats.JetStreamContext, order *models.Order) {
	message, err := json.Marshal(order)
	if err != nil {
		fmt.Println(err)
	}

	js.Publish("msg", message)
}

func addStream(js nats.JetStreamContext, streamName string) {
	js.AddStream(&nats.StreamConfig{
		Name:     streamName,
		Subjects: []string{"msg"},
		MaxMsgs:  100,
		MaxAge:   time.Hour,
	})
}

func deleteStream(js nats.JetStreamContext, streamName string) {
	js.DeleteStream(streamName)
}

func getOrder1() *models.Order {
	return &models.Order{
		OrderUid:    "ba8d924e8c93ca35test",
		TrackNumber: "WBILMTESTTRACK",
		Entry:       "WBIL",
		Delivery: models.Delivery{
			Name:    "Hazel Lapina",
			Phone:   "+9720000000",
			Zip:     "2639809",
			City:    "Moscow",
			Address: "Ploshad Mira 15",
			Region:  "Moscow",
			Email:   "test@gmail.com",
		},
		Payment: models.Payment{
			Transaction:  "ba8d924e8c93ca35test",
			RequestId:    "",
			Currency:     "USD",
			Provider:     "wbpay",
			Amount:       1294,
			PaymentDt:    1637907727,
			Bank:         "alpha",
			DeliveryCost: 0,
			GoodsTotal:   144,
			CustomFee:    0,
		},
		Items: []models.Items{{
			ChrtId:      9934930,
			TrackNumber: "WBILMTESTTRACK",
			Price:       453,
			Rid:         "ab4219087a764ae0btest",
			Name:        "Mascaras",
			Sale:        30,
			Size:        "0",
			TotalPrice:  317,
			NmId:        2389212,
			Brand:       "Vivienne Sabo",
			Status:      202,
		}, {
			ChrtId:      9934931,
			TrackNumber: "WBILMTESTTRACK",
			Price:       368,
			Rid:         "ab4219087a764ae0btest",
			Name:        "Mascaras",
			Sale:        30,
			Size:        "0",
			TotalPrice:  317,
			NmId:        2389212,
			Brand:       "Revolution Makeup",
			Status:      202,
		}},
		Locale:            "en",
		InternalSignature: "",
		CustomerId:        "test",
		DeliveryService:   "meest",
		Shardkey:          "9",
		SmId:              99,
		DateCreated: func() time.Time {
			t, _ := time.Parse(time.RFC3339, "2021-11-26T06:22:19Z")
			return t
		}(),
		OofShard: "1",
	}
}
