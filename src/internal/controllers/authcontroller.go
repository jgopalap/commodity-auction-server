package controllers

import (
	"encoding/json"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"internal/models"
	"internal/utils"
)

func GetUserTypes(w http.ResponseWriter, r *http.Request) {
	userTypes := make([]*models.UserType, 0)
	err := database.Table("auction_db.public.\"User Types\"").Scan(&userTypes).Error
	if err != nil {
		resp := utils.Message(false, "Connection error " + err.Error() + ".  Please try again." )
		utils.Respond(w, resp)
		return
	}
	resp := utils.Message(true, "success")
	resp["user_types"] = userTypes
	utils.Respond(w, resp)
}

func AddUser(w http.ResponseWriter, r *http.Request) {
	user := &models.User{}
	err := json.NewDecoder(r.Body).Decode(user)
	if err != nil {
		resp := utils.Message(false,  "Unable to unmarshal user " + err.Error() + ".")
		utils.Respond(w, resp)
		return
	}
	pass, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		resp := utils.Message(false, "Password encryption failed.")
		utils.Respond(w, resp)
		return
	}

	user.Password = string(pass)
	isNewRecord := database.Table("auction_db.public.\"User Accounts\"").NewRecord(user)
	if isNewRecord {
		userWithId := &models.UserDetailed{}
		err = database.Table("auction_db.public.\"User Accounts\"").Create(&user).Scan(&userWithId).Error
	    if err != nil {
			resp := utils.Message(false, "Unable to add user " + err.Error())
			utils.Respond(w, resp)
			return
		}
		resp := utils.Message(true, "success")
		resp["user_id"] = userWithId.UserId
		utils.Respond(w, resp)
	} else {
		resp := utils.Message(false, "Unable to add user, user already exits.")
		utils.Respond(w, resp)
	}
}

func GetUsers(w http.ResponseWriter, r *http.Request) {
	users := make([]*models.UserDetailed, 0)
	err := database.Table("auction_db.public.\"User Accounts\"").Scan(&users).Error
	if err != nil {
		resp := utils.Message(false, "Unable to get users " + err.Error() + ".  Please try again." )
		utils.Respond(w, resp)
		return
	}
	resp := utils.Message(true, "success")
	resp["users"] = users
	utils.Respond(w, resp)
}

func GetUser(w http.ResponseWriter, r *http.Request) {
	userId := r.URL.Query().Get("user_id")
	user:= &models.UserDetailed{}
	err := database.Table("auction_db.public.\"User Accounts\"").Where("user_id = ?", userId).
		First(&user).Error
	if err != nil {
		resp := utils.Message(false, "Unable to retrieve user " + err.Error() + "." )
		utils.Respond(w, resp)
		return
	}

	resp := utils.Message(true, "success")
	resp["user"] = user
	utils.Respond(w, resp)
}

/*only admins can call this API*/
func UpdateUserStatus(w http.ResponseWriter, r *http.Request) {
	user := &models.UserDetailed{}
	err := json.NewDecoder(r.Body).Decode(user)
	if err != nil {
		resp := utils.Message(false,  "Unable to unmarshal user status " + err.Error() + ".")
		utils.Respond(w, resp)
		return
	}

	err = database.Table("auction_db.public.\"User Accounts\"").Where("user_id = ?",
		user.UserId).Update("account_status_id", user.AccountStatusId).Error
    if err != nil {
		resp := utils.Message(false, "Connection error " + err.Error() + ". Please retry")
		utils.Respond(w, resp)
		return
	}

	resp := utils.Message(true, "Success")
	utils.Respond(w, resp)
}

func UpdateUser(w http.ResponseWriter, r *http.Request) {
	user := &models.UserDetailed{}
	err := json.NewDecoder(r.Body).Decode(user)
	if err != nil {
		resp := utils.Message(false,  "Unable to unmarshal user " + err.Error() + ".")
		utils.Respond(w, resp)
		return
	}

	if user.Password != "" {
		pass, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		if err != nil {
			resp := utils.Message(false, "Password encryption failed.")
			utils.Respond(w, resp)
			return
		}

		user.Password = string(pass)
	}

	err = database.Table("auction_db.public.\"User Accounts\"").Where("user_id = ?", user.UserId).
		Update("email_address", user.EmailAddress).
		Update("password", user.Password).
		Update("first_name", user.FirstName).
		Update("last_name", user.LastName).
		Update("user_type_id",user.UserTypeId).
		Update("business_address", user.BusinessAddress).
		Update("business_name", user.BusinessName).
		Update("shipping_address", user.BusinessName).Error

	if err != nil {
		resp := utils.Message(false, "Unable to update user " + err.Error() + ".")
		utils.Respond(w, resp)
		return
	}

	resp := utils.Message(true, "Success")
	utils.Respond(w, resp)
}

func Login(w http.ResponseWriter, r *http.Request) {
	user := &models.UserDetailed{}
	login := &models.LoginCredentials{}
	err := json.NewDecoder(r.Body).Decode(login)
	if err != nil {
		resp := utils.Message(false,  "Unable to unmarshal login credentials " + err.Error() + ".")
		utils.Respond(w, resp)
		return
	}

	err = database.Table("auction_db.public.\"User Accounts\"").Where("user_name = ?",
		login.UserName).First(&user).Error
	if err != nil {
		resp := utils.Message(false, "Unable to retrieve login creds " + err.Error() + ".")
		utils.Respond(w, resp)
		return
	} else if user.AccountStatusId != int(models.Active) {
		resp := utils.Message(false, "Unable to login, user not active. ")
		utils.Respond(w, resp)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte (user.Password), []byte (login.Password))
	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword { //Password does not match!
		resp := utils.Message(false, "Invalid login credentials. Please try again")
		utils.Respond(w, resp)
		return
	}

	resp := utils.Message(true, "success")
	resp["user"] = user
	utils.Respond(w, resp)

    return
}
