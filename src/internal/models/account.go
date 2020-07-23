package models

type User struct {
	UserName string `json:"username"`
	Password string `json:"password"`
	EmailAddress string `json:"email_address"`
	FirstName string `json:"first_name"`
	LastName string `json:"last_name"`
	UserTypeId int `json:"user_type_id"`
	BusinessAddress string `json:"business_address"`
	BusinessName string `json:"business_name"`
	ShippingAddress string `json:"shipping_address"`
}

type UserDetailed struct {
	User
	UserId int `json:"user_id"`
	AccountStatusId int `json:"account_status_id"`
}

type UserType struct {
	UserType string `json:"user_type"`
	UserTypeId string `json:"id"`
}

type AccountStatus int
const (
	Pending AccountStatus = 1
	Active  AccountStatus = 2
	Inactive AccountStatus = 3
)

type AccountType int
const (
	Admin AccountType = 1
	Broker  AccountType = 2
	Buyer AccountType = 3
	Producer AccountType = 4
)

type LoginCredentials struct {
	UserName string `json:"username"`
	Password string `json:"password"`
}