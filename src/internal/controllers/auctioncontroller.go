package controllers

import (
	"encoding/json"
	maxheap "internal/collections"
	"internal/models"
	"internal/utils"
	"math"
	"net/http"
	"time"
)

func GetLocations(w http.ResponseWriter, r *http.Request) {
	auctionLocations := make([]*models.AuctionLocation, 0)
	err := database.Table("auction_db.public.\"Auction Center Locations\"").Scan(&auctionLocations).Error
	if err != nil {
		resp := utils.Message(false, "Connection error. Please retry")
		utils.Respond(w, resp)
		return
	}
	resp := utils.Message(true, "success")
    resp["locations"] = auctionLocations
	utils.Respond(w, resp)
}

func GetProductCategories(w http.ResponseWriter, r *http.Request) {
	categories := make([]*models.Category, 0)
	err := database.Table("auction_db.public.\"Product Categories\"").Scan(&categories).Error
	if err != nil {
		resp := utils.Message(false, "Unable to retrieve categories " + err.Error())
		utils.Respond(w, resp)
		return
	}
	resp := utils.Message(true, "success")
	resp["categories"] = categories
	utils.Respond(w, resp)
}

//Can only be called by broker
func AddLot(w http.ResponseWriter, r *http.Request) {
	lot := &struct {models.Lot; RemainingWeight float64}{}
	err := json.NewDecoder(r.Body).Decode(lot)
	if err != nil {
		resp := utils.Message(false,  "Unable to unmarshal lot " + err.Error() + ".")
		utils.Respond(w, resp)
		return
	}

	lot.RemainingWeight = lot.OriginalWeight
	isNewRecord := database.Table("auction_db.public.Lots").NewRecord(lot)
	if isNewRecord {
		err = database.Table("auction_db.public.Lots").Create(&lot).Error
		if err != nil {
			resp := utils.Message(false, "Unable to add lot" + err.Error())
			utils.Respond(w, resp)
			return
		}
		resp := utils.Message(true, "success")
		utils.Respond(w, resp)
	} else {
		resp := utils.Message(false, "Unable to add lot, already exits.")
		utils.Respond(w, resp)
	}
}

func GetLot(w http.ResponseWriter, r *http.Request) {
	lotId := r.URL.Query().Get("lot_id")
	lot:= &models.LotDetailed{}
	err := database.Table("auction_db.public.Lots").Where("lot_id = ?", lotId).
		First(&lot).Error
	if err != nil {
		resp := utils.Message(false, "Unable to retrieve lot " + err.Error() + "." )
		utils.Respond(w, resp)
		return
	}

	resp := utils.Message(true, "success")
	resp["lot"] = lot

	utils.Respond(w, resp)
}

//Can only be called by broker
func RemoveLot(w http.ResponseWriter, r *http.Request) {
	lot := &models.LotDetailed{}
	err := json.NewDecoder(r.Body).Decode(lot)
	if err != nil {
		resp := utils.Message(false,  "Unable to unmarshal lot " + err.Error() + ".")
		utils.Respond(w, resp)
		return
	}

	err = database.Table("auction_db.public.Lots").Where("lot_id = ?", lot.LotId).Delete(models.LotDetailed{}).Error
	if err != nil {
		resp := utils.Message(false, "Unable to remove lot " + err.Error() + ".")
		utils.Respond(w, resp)
		return
	}

	resp := utils.Message(true, "success")
	utils.Respond(w, resp)
}

func UpdateLot(w http.ResponseWriter, r *http.Request) {
	lot := &models.LotDetailed{}
	err := json.NewDecoder(r.Body).Decode(lot)
	if err != nil {
		resp := utils.Message(false,  "Unable to unmarshal lot " + err.Error() + ".")
		utils.Respond(w, resp)
		return
	}

	err = database.Table("auction_db.public.Lots").Where("lot_id = ?", lot.LotId).
		Update("product_id", lot.ProductId).
		Update("original_weight", lot.OriginalWeight).
		Update("awr_no", lot.AwrNo).
		Update("garden_report_no", lot.GardenReportNo).
		Update("date_arrived", lot.DateArrived).
		Update("remaining_weight", lot.RemainingWeight).Error
	if err != nil {
		resp := utils.Message(false, "Unable to update lot " + err.Error() + ".")
		utils.Respond(w, resp)
		return
	}

	resp := utils.Message(true, "Success")
	utils.Respond(w, resp)
}

func GetAuctions(w http.ResponseWriter, r *http.Request) {
	auctions := make([]*models.AuctionDetailed, 0)
	err := database.Table("auction_db.public.Auctions").Order("start_time asc").Limit(50).Scan(&auctions).Error
	if err != nil {
		resp := utils.Message(false, "Unable to get auctions.")
		utils.Respond(w, resp)
		return
	}
	resp := utils.Message(true, "success")
	resp["auctions"] = auctions
	utils.Respond(w, resp)
}

//Can only be called by broker
func AddAuction(w http.ResponseWriter, r *http.Request) {
	auction := &models.Auction{}
	err := json.NewDecoder(r.Body).Decode(auction)
	if err != nil {
		resp := utils.Message(false,  "Unable to unmarshal auction " + err.Error() + ".")
		utils.Respond(w, resp)
		return
	}

	user := &models.User{}
	err = database.Table("auction_db.public.\"User Accounts\"").Where("user_id = ?", auction.UserId).First(&user).Error
	if err != nil {
		resp := utils.Message(false, "Unable to add auction " + err.Error())
		utils.Respond(w, resp)
		return
	} else if user.UserTypeId == int(models.Buyer) || user.UserTypeId == int(models.Admin) {
		resp := utils.Message(false, "Unable to add auction, wrong user type.")
		utils.Respond(w, resp)
		return
	}

	isNewRecord := database.Table("auction_db.public.Auctions").NewRecord(auction)
	if isNewRecord {
		err = database.Table("auction_db.public.Auctions").Create(&auction).Error
		if err != nil {
			resp := utils.Message(false, "Unable to add auction " + err.Error())
			utils.Respond(w, resp)
			return
		}
		resp := utils.Message(true, "success")
		utils.Respond(w, resp)
	} else {
		resp := utils.Message(false, "Unable to add auction, already exits.")
		utils.Respond(w, resp)
	}
}

//Can only be called by broker
func RemoveAuction(w http.ResponseWriter, r *http.Request) {
	auction := &models.AuctionDetailed{}
	err := json.NewDecoder(r.Body).Decode(auction)
	if err != nil {
		resp := utils.Message(false,  "Unable to unmarshal auction " + err.Error() + ".")
		utils.Respond(w, resp)
		return
	}

	err = database.Table("auction_db.public.Auctions").Where("auction_id = ?", auction.AuctionId).Delete(models.AuctionDetailed{}).Error
	if err != nil {
		resp := utils.Message(false, "Unable to remove auction " + err.Error() + ".")
		utils.Respond(w, resp)
		return
	}

	resp := utils.Message(true, "success")
	utils.Respond(w, resp)
}

//Can only be called by broker or service when the auction ends
func UpdateAuction(w http.ResponseWriter, r *http.Request) {
	auctionDb := &models.AuctionDetailed{}
	auction := &models.AuctionDetailed{}
	err := json.NewDecoder(r.Body).Decode(auction)
	if err != nil {
		resp := utils.Message(false,  "Unable to unmarshal auction " + err.Error() + ".")
		utils.Respond(w, resp)
		return
	}

	err = database.Table("auction_db.public.Auctions").Where("auction_id = ?",
		auction.AuctionId).First(&auctionDb).Error
	if err != nil {
		resp := utils.Message(false, "Unable to update auction "+err.Error()+".")
		utils.Respond(w, resp)
		return
	}

	if auctionDb.StateId > 1 {
		resp := utils.Message(false, "Unable to update auction, the auction is in progress.")
		utils.Respond(w, resp)
		return
	}

	err = database.Table("auction_db.public.Auctions").Where("auction_id = ?", auction.AuctionId).
		Update("reserve_price", auction.ReservePrice).
		Update("sold_price", auction.SoldPrice).
		Update("sold_weight", auction.SoldWeight).
		Update("state_id", auction.StateId).
		Update("lot_id", auction.LotId).
		Update("start_time", auction.StartTime).
		Update("offered_weight", auction.OfferedWeight).Error
	if err != nil {
		resp := utils.Message(false, "Unable to update auction " + err.Error() + ".")
		utils.Respond(w, resp)
		return
	}

	resp := utils.Message(true, "Success")
	utils.Respond(w, resp)
}

//Can only be called by broker
func AddBid(w http.ResponseWriter, r *http.Request) {
	bid := &models.BidDetailed{}
	err := json.NewDecoder(r.Body).Decode(bid)
	if err != nil {
		resp := utils.Message(false,  "Unable to unmarshal bid " + err.Error() + ".")
		utils.Respond(w, resp)
		return
	}
	if bid.UserTypeId != int(models.Buyer) {
		resp := utils.Message(false, "Unable to add bid, wrong user type.")
		utils.Respond(w, resp)
		return
	}
	auctionInfo := struct{StartTime time.Time; StateId int}{}
	err = database.Table("auction_db.public.Auctions").Where("auction_id = ?", bid.AuctionId).First(&auctionInfo).Error
	if err != nil {
		utils.Respond(w, utils.Message(false, "Unable to add bid, auction doesn't exist"))
		return
	}
	if bid.Time.After(auctionInfo.StartTime) || auctionInfo.StateId > int(models.Open) {
		utils.Respond(w, utils.Message(false, "Bid is inadmissible."))
		return
	}
	isNewRecord := database.Table("auction_db.public.Bids").NewRecord(bid)
	if isNewRecord {
		err = database.Table("auction_db.public.Bids").Create(&bid.Bid).Error
		if err != nil {
			resp := utils.Message(false, "Unable to add bid " + err.Error())
			utils.Respond(w, resp)
			return
		}
		resp := utils.Message(true, "success")
		resp["bid"] = bid
		utils.Respond(w, resp)
	} else {
		resp := utils.Message(false, "Unable to add bid, already exits.")
		utils.Respond(w, resp)
	}
}

func RemoveBid(w http.ResponseWriter, r *http.Request) {
	bid := &models.BidDetailed{}
	err := json.NewDecoder(r.Body).Decode(bid)
	if err != nil {
		resp := utils.Message(false,  "Unable to unmarshal bid " + err.Error() + ".")
		utils.Respond(w, resp)
		return
	}

	if bid.UserTypeId != int(models.Buyer) {
		resp := utils.Message(false, "Unable to remove bid, wrong user type.")
		utils.Respond(w, resp)
		return
	}

	err = database.Table("auction_db.public.Bids").Where("bid_id = ?", bid.BidId).
		Delete(models.BidDetailed{}).Error
	if err != nil {
		resp := utils.Message(false, "Unable to remove bid " + err.Error() + ".")
		utils.Respond(w, resp)
		return
	}

	resp := utils.Message(true, "success")
	utils.Respond(w, resp)
}

func UpdateBid(w http.ResponseWriter, r *http.Request) {
	bid := &models.BidDetailed{}
	err := json.NewDecoder(r.Body).Decode(bid)
	if err != nil {
		resp := utils.Message(false,  "Unable to unmarshal bid " + err.Error() + ".")
		utils.Respond(w, resp)
		return
	}

	err = database.Table("auction_db.public.Bids").Where("bid_id = ?", bid.BidId).
		Update("value", bid.Value).
		Update("weight", bid.Weight).
		Update("time", bid.Time).
		Update("budget", bid.Budget).
		Update("auction_id", bid.AuctionId).Error
	if err != nil {
		resp := utils.Message(false, "Unable to update bid " + err.Error() + ".")
		utils.Respond(w, resp)
		return
	}

	resp := utils.Message(true, "Success")
	utils.Respond(w, resp)
}

func GetBidsForUser(w http.ResponseWriter, r *http.Request) {
	userId := r.URL.Query().Get("user_id")

	bids := make([]*models.BidDetailed, 0)
	err := database.Table("auction_db.public.Bids").Where("user_id = ?", userId).Scan(&bids).Error
	if err != nil {
		resp := utils.Message(false, "Unable to find bids for user " + err.Error() + ".")
		utils.Respond(w, resp)
		return
	}

	resp := utils.Message(true, "Success")
	resp["bids"] = bids
	utils.Respond(w, resp)
}

func GetBidsForAuction(w http.ResponseWriter, r *http.Request) {
	auctionId := r.URL.Query().Get("auction_id")

	bids := make([]*models.BidDetailed, 0)
	err := database.Table("auction_db.public.Bids").Where("auction_id = ?", auctionId).Scan(&bids).Error
	if err != nil {
		resp := utils.Message(false, "Unable to find bids for auction " + err.Error() + ".")
		utils.Respond(w, resp)
		return
	}

	resp := utils.Message(true, "Success")
	resp["bids"] = bids
	utils.Respond(w, resp)
}

func RunAuction(w http.ResponseWriter, r *http.Request) {
	auction := &models.AuctionDetailed{}
	err := json.NewDecoder(r.Body).Decode(auction)
	if err != nil {
		resp := utils.Message(false,  "Unable to unmarshal auction " + err.Error() + ".")
		utils.Respond(w, resp)
		return
	}

	err = database.Table("auction_db.public.Auctions").Where("auction_id = ?",
		auction.AuctionId).First(&auction).Update("state_id", int(models.Clearing)).Error
	if err != nil {
		resp := utils.Message(false, "Unable to run auction, doesn't exist.")
		utils.Respond(w, resp)
		return
	}

	bids := make([]*models.BidDetailed, 0)
	err = database.Table("auction_db.public.Bids").Where("auction_id = ?", auction.AuctionId).Scan(&bids).Error
	if err != nil || len(bids) == 0 {
		resp := utils.Message(false, "Unable to run auction, retrieving bids failed.")
		utils.Respond(w, resp)
		return
	}

	bidHeap := &maxheap.BidMaxHeap{}
	for _, bid := range bids {
		maxheap.Push(bidHeap, bid)
	}

	winningBid := maxheap.Pop(bidHeap).(*models.BidDetailed)
	winningPrice := winningBid.Value
    createOrder(winningPrice, winningBid.BidId, auction.LotId, winningBid.Weight, w)

	lot := &models.Lot{}
	database.Table("auction_db.public.Lots").Where("lot_id = ?", auction.LotId).First(&lot)
	lotAvailableWeight := lot.OriginalWeight - winningBid.Weight
	//minimum shipping size is 10kg
	for lotAvailableWeight >= float64(10) && bidHeap.Len() > 0 {
		winningBid = maxheap.Pop(bidHeap).(*models.BidDetailed)
	    weight := math.Floor(winningBid.Budget / winningPrice)
	    if weight >= winningBid.Weight {
	    	weight = winningBid.Weight
		}

	    createOrder(winningPrice, winningBid.BidId, auction.LotId, weight, w)
	    lotAvailableWeight -= weight
	}

	err = database.Table("auction_db.public.Auctions").Where("auction_id = ?",
		auction.AuctionId).First(&auction).Update("state_id", int(models.Closed)).Error
	if err != nil {
		resp := utils.Message(false, "Unable to close auction.")
		utils.Respond(w, resp)
		return
	}
}

func createOrder(winningPrice float64, bidId int, lotId int, weight float64, w http.ResponseWriter) {
	//update bid to reflect win
	err := database.Table("auction_db.public.Bids").Where("bid_id = ?", bidId).
		Update("is_winning_bid", true).Error
	if err != nil {
		resp := utils.Message(false, "Unable to update winning bid.")
		utils.Respond(w, resp)
		return
	}

	//create order
	order := &models.Order{}
	order.OrderCreationDate = time.Now()
	order.BidId = bidId
	order.LotId = lotId
	order.Weight = weight
	order.Price = winningPrice

	isNewRecord := database.Table("auction_db.public.Orders").NewRecord(order)
	if isNewRecord {
		err := database.Table("auction_db.public.Orders").Create(&order).Error
		if err != nil {
			resp := utils.Message(false, "Unable to create order " + err.Error())
			utils.Respond(w, resp)
			return
		}
		resp := utils.Message(true, "success")
		utils.Respond(w, resp)
	} else {
		resp := utils.Message(false, "Unable to create order, already exits.")
		utils.Respond(w, resp)
	}
}
