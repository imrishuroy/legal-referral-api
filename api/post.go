package api

import (
	"context"
	"encoding/json"
	"errors"
	"firebase.google.com/go/v4/auth"
	"fmt"
	"github.com/rs/zerolog/log"
	"mime/multipart"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	db "github.com/imrishuroy/legal-referral/db/sqlc"
)

// PostType represents the type of post
type PostType string

const (
	PostTypeText     PostType = "text"
	PostTypeImage    PostType = "image"
	PostTypeVideo    PostType = "video"
	PostTypeAudio    PostType = "audio"
	PostTypeLink     PostType = "link"
	PostTypeDocument PostType = "document"
	PostTypePoll     PostType = "poll"
	PostTypeOther    PostType = "other"
)

type createPollReq struct {
	OwnerID   string     `form:"owner_id" binding:"required"`
	PollTitle string     `form:"poll_title"`
	Options   []string   `form:"options"`
	EndTime   *time.Time `form:"end_time"`
}

type createPostReq struct {
	OwnerID   string                  `form:"owner_id" binding:"required"`
	Content   string                  `form:"content"`
	Files     []*multipart.FileHeader `form:"files"`
	PostType  PostType                `form:"post_type" binding:"required"`
	PollTitle string                  `form:"poll_title"`
	Options   []string                `form:"options"`
	EndTime   *time.Time              `form:"end_time"`
}

func (srv *Server) CreatePost(ctx *gin.Context) {

	var req createPostReq

	// Bind the form fields to the struct
	if err := ctx.ShouldBind(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	log.Info().Msgf("Request: %+v", req)

	// Check if the authenticated user is the owner
	authPayload := ctx.MustGet(authorizationPayloadKey).(*auth.Token)
	if authPayload.UID != req.OwnerID {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	imageUrls := make([]string, 0)

	if req.PostType == PostTypeImage || req.PostType == PostTypeVideo || req.PostType == PostTypeDocument {
		urls, err := srv.handleFilesUpload(ctx, req.Files)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}
		log.Info().Msgf("URLs: %+v", urls)
		imageUrls = append(imageUrls, urls...)

	}

	var pollID *int32
	if req.PostType == PostTypePoll {
		createPollReq := createPollReq{
			OwnerID:   req.OwnerID,
			PollTitle: req.PollTitle,
			Options:   req.Options,
			EndTime:   req.EndTime,
		}

		poll, err := srv.createPoll(ctx, &createPollReq)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}
		pollID = &poll.PollID
	}

	arg := db.CreatePostParams{
		OwnerID:  req.OwnerID,
		Content:  &req.Content,
		Media:    imageUrls,
		PostType: db.PostType(req.PostType),
		PollID:   pollID,
	}

	dbPost, err := srv.Store.CreatePost(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	post := post{
		PostID:    dbPost.PostID,
		OwnerID:   dbPost.OwnerID,
		Content:   dbPost.Content,
		Media:     dbPost.Media,
		PostType:  PostType(dbPost.PostType),
		PollID:    dbPost.PollID,
		CreatedAt: dbPost.CreatedAt,
	}

	// cache the post
	postKey := fmt.Sprintf("post:%d", post.PostID)
	log.Info().Msg(postKey)

	if err := srv.cachePost(ctx, postKey, post, 12*time.Hour); err != nil {
		log.Error().Err(err).Msg("Failed to cache post")
	}

	srv.publishToKafka("publish-feed", req.OwnerID, string(post.PostID))

	ctx.JSON(http.StatusOK, gin.H{"success": "Post created successfully"})
}

func (srv *Server) IsPostFeatured(ctx *gin.Context) {
	postIDStr := ctx.Param("post_id")

	postID, err := strconv.Atoi(postIDStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	int32PostID := int32(postID)

	featured, err := srv.Store.IsPostFeatured(ctx, int32PostID)
	if err != nil {
		if errors.Is(err, db.ErrRecordNotFound) {
			ctx.JSON(http.StatusOK, false)
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, featured)
}

func (srv *Server) cachePost(ctx context.Context, key string, post post, expDuration time.Duration) error {
	data, err := json.Marshal(post)
	if err != nil {
		return fmt.Errorf("error serializing post data: %v", err)
	}
	// convert data to string
	dataStr := string(data)
	return srv.ValkeyClient.Do(ctx, srv.ValkeyClient.B().Set().Key(key).Value(dataStr).Ex(expDuration).Build()).Error()
}

func (srv *Server) createPoll(ctx *gin.Context, req *createPollReq) (*db.Poll, error) {
	if req == nil {
		return nil, errors.New("poll request is nil")
	}

	arg := db.CreatePollParams{
		OwnerID: req.OwnerID,
		Title:   req.PollTitle,
		Options: req.Options,
		//EndDate: pgtype.Timestamptz{Time: *req.EndDate, Valid: req.EndDate != nil},
	}

	poll, err := srv.Store.CreatePoll(ctx, arg)
	if err != nil {
		return nil, err

	}
	return &poll, nil
}

func (srv *Server) postToNewsFeed(ctx *gin.Context, userID string, postID int32) error {
	userIDs, err := srv.Store.ListConnectedUserIDs(ctx, userID)
	if err != nil {
		return err
	}

	userIDs = append(userIDs, userID)
	for _, id := range userIDs {
		srv.publishToKafka("publish-feed", id.(string), string(postID))
	}
	return nil
}

func (srv *Server) PostLikesAndCommentsCount(ctx *gin.Context) {
	postIDStr := ctx.Param("post_id")

	authPayload := ctx.MustGet(authorizationPayloadKey).(*auth.Token)
	if authPayload.UID == "" {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	postID, err := strconv.Atoi(postIDStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	int32PostID := int32(postID)

	likes, err := srv.Store.GetPostLikesCount(ctx, &int32PostID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	comments, err := srv.Store.GetPostCommentsCount(ctx, int32(postID))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"likes":    likes,
		"comments": comments,
	})
}

func (srv *Server) IsPostLiked(ctx *gin.Context) {
	postIDStr := ctx.Param("post_id")

	authPayload := ctx.MustGet(authorizationPayloadKey).(*auth.Token)
	if authPayload.UID == "" {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	postID, err := strconv.Atoi(postIDStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	int32PostID := int32(postID)

	arg := db.GetPosIsLikedByCurrentUserParams{
		UserID: authPayload.UID,
		PostID: &int32PostID,
	}

	liked, err := srv.Store.GetPosIsLikedByCurrentUser(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, liked)
}

func (srv *Server) DeletePost(ctx *gin.Context) {
	postIDStr := ctx.Param("post_id")
	postID, err := strconv.Atoi(postIDStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*auth.Token)
	if authPayload.UID == "" {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	int32PostID := int32(postID)

	arg := db.DeletePostParams{
		PostID:  int32PostID,
		OwnerID: authPayload.UID,
	}

	err = srv.Store.DeletePost(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Post deleted successfully"})
}

func (srv *Server) GetPost(ctx *gin.Context) {
	postIDStr := ctx.Param("post_id")
	postID, err := strconv.Atoi(postIDStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	int32PostID := int32(postID)

	post, err := srv.Store.GetPost(ctx, int32PostID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, post)
}

type searchPostsReq struct {
	Limit       int32  `form:"limit" binding:"required"`
	Offset      int32  `form:"offset" binding:"required"`
	SearchQuery string `form:"query"`
}

func (srv *Server) SearchPosts(ctx *gin.Context) {
	var req searchPostsReq
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*auth.Token)
	if authPayload.UID == "" {
		ctx.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}

	arg := db.SearchPostsParams{
		Limit:       req.Limit,
		Offset:      (req.Offset - 1) * req.Limit,
		Searchquery: req.SearchQuery,
	}

	posts, err := srv.Store.SearchPosts(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, posts)
}
