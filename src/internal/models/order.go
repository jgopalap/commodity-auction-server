package models

import "time"

type Order struct {
	BidId int `json:"bid_id"`
	OrderCreationDate time.Time `json:"order_creation_date"`
	Weight float64 `json:"weight"`
	Price float64 `json:"price"`
	LotId int `json:"lot_id"`
}

type OrderDetailed struct {
	Order
	OrderId int `json:"order_id"`
	OrderStatusId int `json:"order_status_id"`
	OrderDeliveredDate time.Time `json:"order_delivered_date"`
}

type OrderStatus int
const (
	Created   OrderStatus = 1
	Shipped  OrderStatus = 2
	Delivered   OrderStatus = 3
)