package models

import "time"


type Auction struct {
	UserId int       `json:"user_id"`
	LotId        int       `json:"lot_id"`
	ReservePrice float64   `json:"reserve_price"`
	StartTime    time.Time `json:"start_time"`
	OfferedWeight float64  `json:"offered_weight"`
}

type AuctionDetailed struct {
	Auction
	AuctionId int `json:"auction_id"`
	LocationId   int       `json:"location_id"`
	StateId      int       `json:"state_id"`
	SoldPrice    float64   `json:"sold_price"`
	SoldWeight   float64   `json:"sold_weight"`
}

type AuctionState int
const (
	Open   AuctionState = 1
	Clearing  AuctionState = 2
	Closed   AuctionState = 3
)

type Bid struct {
	BidId int `gorm:"primary_key;auto_increment" json:"bid_id"`
	UserId int ` json:"user_id"`
	Value float64 `json:"value"`
	Weight float64 `json:"weight"`
	Time time.Time `json:"time"`
	Budget float64 `json:"budget"`
	AuctionId int `json:"auction_id"`
}

type BidDetailed struct {
	Bid
	IsWinningBid bool `json:"is_winning_bid"`
	UserTypeId int `json:"user_type_id"`
}

type AuctionLocation struct {
	City		string `json:"city"`
	Country 	string `json:"country"`
	LocationId 	string `json:"location_id"`
}

type Lot struct {
	ProductId        int        `json:"product_id"`
	OriginalWeight   float64    `json:"original_weight"`
	AwrNo            int        `json:"awr_no"`
	GardenReportNo   int        `json:"garden_report_no"`
	DateArrived      time.Time  `json:"date_arrived"`
}

type LotDetailed struct {
	Lot
	LotId            int  `json:"lot_id"`
	RemainingWeight  int  `json:"remaining_weight"`
}