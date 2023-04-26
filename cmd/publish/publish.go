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

	streamName := "MSG"

	js.AddStream(&nats.StreamConfig{
		Name: streamName,
		//Retention: nats.WorkQueuePolicy,
		Subjects: []string{"msg"},
	})

	message, err := json.Marshal(getOrder1())
	if err != nil {
		fmt.Println(err)
	}

	js.Publish("msg", message)

	//input := bufio.NewScanner(os.Stdin)
	//input.Scan()
}

func getOrder1() *models.Order {
	return &models.Order{
		OrderUid:    "b563feb7b2b84b6te",
		TrackNumber: "WBILMTESTTRACK",
		Entry:       "WBIL",
		Delivery: models.Delivery{
			Name:    "Test Testov",
			Phone:   "+9720000000",
			Zip:     "2639809",
			City:    "Kiryat Mozkin",
			Address: "Ploshad Mira 15",
			Region:  "Kraiot",
			Email:   "test@gmail.com",
		},
		Payment: models.Payment{
			Transaction:  "b563feb7b2b84b6test",
			RequestId:    "",
			Currency:     "USD",
			Provider:     "wbpay",
			Amount:       1817,
			PaymentDt:    1637907727,
			Bank:         "alpha",
			DeliveryCost: 1500,
			GoodsTotal:   317,
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
