package app

import (
	"github.com/gorilla/mux"
	"internal/controllers"
	"net/http"
)

func Init() {
	controllers.Init()
	router := mux.NewRouter()

	//User
	router.HandleFunc("/getUserTypes", controllers.GetUserTypes).Methods("GET")
	router.HandleFunc("/addUser", controllers.AddUser).Methods("POST")
	router.HandleFunc("/getUsers", controllers.GetUsers).Methods("GET")
	router.HandleFunc("/getUser", controllers.GetUser).Methods("GET")
	router.HandleFunc("/updateUserStatus", controllers.UpdateUserStatus).Methods("PATCH")
	router.HandleFunc("/updateUser", controllers.UpdateUser).Methods("PUT")
	router.HandleFunc("/login", controllers.Login).Methods("POST")

	//Auction
	router.HandleFunc("/getAuctionLocations", controllers.GetLocations).Methods("GET")
	router.HandleFunc("/getProductCategories", controllers.GetProductCategories).Methods("GET")
	router.HandleFunc("/addAuction", controllers.AddAuction).Methods("POST")
	router.HandleFunc("/removeAuction", controllers.RemoveAuction).Methods("DELETE")
	router.HandleFunc("/updateAuction", controllers.UpdateAuction).Methods("PUT")
	router.HandleFunc("/runAuction", controllers.RunAuction).Methods("POST")
	router.HandleFunc("/getAuctions", controllers.GetAuctions).Methods("GET")
	router.HandleFunc("/addLot", controllers.AddLot).Methods("POST")
	router.HandleFunc("/removeLot", controllers.RemoveLot).Methods("DELETE")
	router.HandleFunc("/updateLot", controllers.UpdateLot).Methods("PUT")
	router.HandleFunc("/getLot", controllers.GetLot).Methods("GET")
	router.HandleFunc("/addBid", controllers.AddBid).Methods("POST")
	router.HandleFunc("/removeBid", controllers.RemoveBid).Methods("DELETE")
	router.HandleFunc("/updateBid", controllers.UpdateBid).Methods("PUT")
	router.HandleFunc("/getBidsForUser", controllers.GetBidsForUser).Methods("GET")
	router.HandleFunc("/getBidsForAuction", controllers.GetBidsForAuction).Methods("GET")
	router.HandleFunc("/removeLot", controllers.RemoveLot).Methods("DELETE")
	router.HandleFunc("/updateLot", controllers.UpdateLot).Methods("PUT")


	//Product
	router.HandleFunc("/addProduct", controllers.AddProduct).Methods("POST")
	router.HandleFunc("/removeProduct", controllers.RemoveProduct).Methods("DELETE")
	router.HandleFunc("/updateProduct", controllers.UpdateProduct).Methods("PUT")
	router.HandleFunc("/getProduct", controllers.GetProduct).Methods("GET")

	http.ListenAndServe(":8000", router)
}



