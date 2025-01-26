package api

import (
	"firebase.google.com/go/v4/auth"
	"github.com/gin-gonic/gin"
	db "github.com/imrishuroy/legal-referral/db/sqlc"
	"net/http"
	"strconv"
)

type createDiscussionRequest struct {
	AuthorID        string   `json:"author_id"`
	Topic           string   `json:"topic"`
	InvitedUsersIDs []string `json:"invited_users_ids"`
}

func (s *Server) createDiscussion(ctx *gin.Context) {

	var req *createDiscussionRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(400, errorResponse(err))
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*auth.Token)
	if authPayload.UID != req.AuthorID {
		ctx.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}

	// create discussion
	arg := db.CreateDiscussionParams{
		AuthorID: req.AuthorID,
		Topic:    req.Topic,
	}

	discussion, err := s.Store.CreateDiscussion(ctx, arg)
	if err != nil {
		ctx.JSON(400, errorResponse(err))
		return
	}

	// invite users
	for _, invitedUserID := range req.InvitedUsersIDs {
		arg := db.InviteUserToDiscussionParams{
			DiscussionID:  discussion.DiscussionID,
			InviteeUserID: req.AuthorID,
			InvitedUserID: invitedUserID,
		}

		err = s.Store.InviteUserToDiscussion(ctx, arg)
		if err != nil {
			errorCode := db.ErrorCode(err)
			if errorCode == db.UniqueViolation {
				ctx.JSON(400, gin.H{"message": "Already invited"})
				return
			}
			ctx.JSON(400, errorResponse(err))
			return
		}
	}
	ctx.JSON(200, discussion)

}

func (s *Server) updateDiscussionTopic(ctx *gin.Context) {

	var req db.UpdateDiscussionTopicParams
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(400, errorResponse(err))
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*auth.Token)

	if authPayload.UID == "" {
		ctx.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}

	err := s.Store.UpdateDiscussionTopic(ctx, req)
	if err != nil {
		ctx.JSON(400, errorResponse(err))
		return
	}

	ctx.JSON(200, gin.H{"message": "Updated"})
}

type inviteUserToDiscussionRequest struct {
	InviteeUserID string `json:"invitee_user_id"`
	InvitedUserID string `json:"invited_user_id"`
}

func (s *Server) inviteUserToDiscussion(ctx *gin.Context) {
	discussionIDStr := ctx.Param("discussion_id")
	discussionID, err := strconv.Atoi(discussionIDStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Invalid discussion ID"})
		return
	}

	var req inviteUserToDiscussionRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(400, errorResponse(err))
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*auth.Token)
	if authPayload.UID != req.InviteeUserID {
		ctx.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}

	arg := db.InviteUserToDiscussionParams{
		DiscussionID:  int32(discussionID),
		InviteeUserID: req.InviteeUserID,
		InvitedUserID: req.InvitedUserID,
	}

	err = s.Store.InviteUserToDiscussion(ctx, arg)
	if err != nil {
		errorCode := db.ErrorCode(err)
		if errorCode == db.UniqueViolation {
			ctx.JSON(400, gin.H{"message": "Already invited"})
			return
		}

		ctx.JSON(400, errorResponse(err))
		return
	}
	ctx.JSON(200, gin.H{"message": "Invited"})
}

func (s *Server) joinDiscussion(ctx *gin.Context) {
	discussionIDStr := ctx.Param("discussion_id")
	discussionID, err := strconv.Atoi(discussionIDStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Invalid discussion ID"})
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*auth.Token)
	if authPayload.UID == "" {
		ctx.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}

	arg := db.JoinDiscussionParams{
		DiscussionID:  int32(discussionID),
		InvitedUserID: authPayload.UID,
	}

	err = s.Store.JoinDiscussion(ctx, arg)
	if err != nil {
		errorCode := db.ErrorCode(err)
		if errorCode == db.UniqueViolation {
			ctx.JSON(400, gin.H{"message": "Already joined"})
			return
		}

		ctx.JSON(400, errorResponse(err))
		return
	}
	ctx.JSON(200, gin.H{"message": "Joined"})
}

func (s *Server) rejectDiscussion(ctx *gin.Context) {

	discussionIDStr := ctx.Param("discussion_id")
	discussionID, err := strconv.Atoi(discussionIDStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Invalid discussion ID"})
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*auth.Token)
	if authPayload.UID == "" {
		ctx.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}

	arg := db.RejectDiscussionParams{
		DiscussionID:  int32(discussionID),
		InvitedUserID: authPayload.UID,
	}

	err = s.Store.RejectDiscussion(ctx, arg)
	if err != nil {
		ctx.JSON(400, errorResponse(err))
		return
	}
	ctx.JSON(200, gin.H{"message": "Rejected"})
}

func (s *Server) listDiscussionInvites(ctx *gin.Context) {
	userID := ctx.Param("user_id")

	authPayload := ctx.MustGet(authorizationPayloadKey).(*auth.Token)
	if authPayload.UID != userID {
		ctx.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}

	invites, err := s.Store.ListDiscussionInvites(ctx, userID)
	if err != nil {
		ctx.JSON(400, errorResponse(err))
		return
	}

	ctx.JSON(200, invites)
}

func (s *Server) listActiveDiscussions(ctx *gin.Context) {
	userID := ctx.Param("user_id")

	authPayload := ctx.MustGet(authorizationPayloadKey).(*auth.Token)
	if authPayload.UID != userID {
		ctx.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}

	discussions, err := s.Store.ListActiveDiscussions(ctx, userID)
	if err != nil {
		ctx.JSON(400, errorResponse(err))
		return
	}

	ctx.JSON(200, discussions)
}

func (s *Server) listDiscussionParticipants(ctx *gin.Context) {
	discussionIDStr := ctx.Param("discussion_id")
	discussionID, err := strconv.Atoi(discussionIDStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Invalid discussion ID"})
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*auth.Token)
	if authPayload.UID == "" {
		ctx.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}

	participants, err := s.Store.ListDiscussionParticipants(ctx, int32(discussionID))
	if err != nil {
		ctx.JSON(400, errorResponse(err))
		return
	}

	ctx.JSON(200, participants)
}

func (s *Server) listUninvitedParticipants(ctx *gin.Context) {
	discussionIDStr := ctx.Param("discussion_id")
	discussionID, err := strconv.Atoi(discussionIDStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Invalid discussion ID"})
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*auth.Token)
	if authPayload.UID == "" {
		ctx.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}

	participants, err := s.Store.ListUninvitedParticipants(ctx, int32(discussionID))
	if err != nil {
		ctx.JSON(400, errorResponse(err))
		return
	}

	ctx.JSON(200, participants)
}

type invitedUsersToDiscussionReq struct {
	InvitedUserIDs []string `json:"invited_user_ids"`
}

func (s *Server) inviteUsersToDiscussion(ctx *gin.Context) {
	discussionIDStr := ctx.Param("discussion_id")
	discussionID, err := strconv.Atoi(discussionIDStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Invalid discussion ID"})
		return
	}

	var req invitedUsersToDiscussionReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(400, errorResponse(err))
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*auth.Token)
	if authPayload.UID == "" {
		ctx.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}

	for _, invitedUserID := range req.InvitedUserIDs {
		arg := db.InviteUserToDiscussionParams{
			DiscussionID:  int32(discussionID),
			InviteeUserID: authPayload.UID,
			InvitedUserID: invitedUserID,
		}

		err = s.Store.InviteUserToDiscussion(ctx, arg)
		if err != nil {
			errorCode := db.ErrorCode(err)
			if errorCode == db.UniqueViolation {
				ctx.JSON(400, gin.H{"message": "Already invited"})
				return
			}
			ctx.JSON(400, errorResponse(err))
			return
		}

	}

	ctx.JSON(200, gin.H{"message": "Invited"})
}
