package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type createUserReq struct {
	ID              string `json:"id"`
	FirstName       string `json:"first_name"`
	LastName        string `json:"last_name"`
	MobileNumber    string `json:"mobile_number"`
	Email           string `json:"email"`
	BarLicenceNo    string `json:"bar_licence_no"`
	PracticingField string `json:"practicing_field"`
	Experience      int32  `json:"experience"`
}

func (server *Server) createUser(ctx *gin.Context) {
	var req createUserReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	//if req.Email == "" || req.Name == "" {
	//	ctx.JSON(http.StatusBadRequest, gin.H{"error": "Email and Name are required"})
	//	return
	//}

	//	authPayload := ctx.MustGet(authorizationPayloadKey).(*auth.Token)

	// search if req email already exists in db
	//dbUser, err := server.store.GetUserByEmail(ctx, req.Email)
	//if err != nil {
	//	if !errors.Is(err, db.ErrRecordNotFound) {
	//		ctx.JSON(http.StatusBadRequest, errorResponse(err))
	//		return
	//	}
	//}
	//
	//// found the user with req email
	//if dbUser.ID != "" {
	//	ctx.JSON(http.StatusBadRequest, gin.H{"error": "user with email already exists"})
	//	return
	//}

	// check if email is verified
	//u, e := server.auth.GetUserByEmail(ctx, req.Email)
	//if e != nil {
	//	ctx.JSON(http.StatusBadRequest, errorResponse(e))
	//	return
	//}
	//
	//if !u.EmailVerified {
	//	ctx.JSON(http.StatusBadRequest, gin.H{"error": "Email is not verified"})
	//	return
	//}

	// create user

	//arg := db.CreateUserParams{
	//	ID:              authPayload.UID,
	//	FirstName:       req.FirstName,
	//	LastName:        req.LastName,
	//	MobileNumber:    req.MobileNumber,
	//	Email:           req.Email,
	//	BarLicenceNo:    req.BarLicenceNo,
	//	PracticingField: req.PracticingField,
	//	Experience:      req.Experience,
	//}

	//user, err := server.store.CreateUser(ctx, arg)
	//if err != nil {
	//	errCode := db.ErrorCode(err)
	//	if errCode == db.ForeignKeyViolation || errCode == db.UniqueViolation {
	//		ctx.JSON(http.StatusForbidden, errorResponse(err))
	//		return
	//	}
	//	ctx.JSON(http.StatusInternalServerError, errorResponse(err))
	//	return
	//
	//}

	//ctx.JSON(http.StatusOK, user)
	ctx.JSON(http.StatusOK, gin.H{"message": "user created successfully"})

}
