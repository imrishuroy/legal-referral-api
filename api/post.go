package api

import (
	"errors"
	"firebase.google.com/go/v4/auth"
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

func (server *Server) createPost(ctx *gin.Context) {

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
		//urls, err := server.handleFilesUpload(req.Files, s3BucketName(req.PostType))
		urls, err := server.handleFilesUpload(req.Files)
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

		poll, err := server.createPoll(ctx, &createPollReq)
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

	post, err := server.store.CreatePost(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	server.publishToKafka("publish-feed", req.OwnerID, string(post.PostID))

	// this should also throw an error
	//server.createProducer(req.OwnerID, string(post.PostID))
	//server.publishToKafka(req.OwnerID, string(post.PostID))

	//if err := server.postToNewsFeed(ctx, req.OwnerID, post.PostID); err != nil {
	//	ctx.JSON(http.StatusInternalServerError, errorResponse(err))
	//	return
	//}

	ctx.JSON(http.StatusOK, "Post created successfully")
}

func (server *Server) createPoll(ctx *gin.Context, req *createPollReq) (*db.Poll, error) {
	if req == nil {
		return nil, errors.New("poll request is nil")
	}

	arg := db.CreatePollParams{
		OwnerID: req.OwnerID,
		Title:   req.PollTitle,
		Options: req.Options,
		//EndDate: pgtype.Timestamptz{Time: *req.EndDate, Valid: req.EndDate != nil},
	}

	poll, err := server.store.CreatePoll(ctx, arg)
	if err != nil {
		return nil, err

	}
	return &poll, nil
}

func (server *Server) postToNewsFeed(ctx *gin.Context, userID string, postID int32) error {
	userIDs, err := server.store.ListConnectedUserIDs(ctx, userID)
	if err != nil {
		return err
	}

	userIDs = append(userIDs, userID)
	for _, id := range userIDs {
		//arg := db.PostToNewsFeedParams{
		//	UserID: id.(string),
		//	PostID: postID,
		//}
		//server.createProducer(id.(string), string(postID))
		server.publishToKafka("publish-feed", id.(string), string(postID))
		//if err := server.store.PostToNewsFeed(ctx, arg); err != nil {
		//	return err
		//}
	}
	return nil
}

func (server *Server) postLikesAndCommentsCount(ctx *gin.Context) {
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

	likes, err := server.store.GetPostLikesCount(ctx, &int32PostID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	comments, err := server.store.GetPostCommentsCount(ctx, int32(postID))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"likes":    likes,
		"comments": comments,
	})
}

func (server *Server) isPostLiked(ctx *gin.Context) {
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

	liked, err := server.store.GetPosIsLikedByCurrentUser(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, liked)
}

func (server *Server) deletePost(ctx *gin.Context) {
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

	err = server.store.DeletePost(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Post deleted successfully"})
}

func (server *Server) getPost(ctx *gin.Context) {
	postIDStr := ctx.Param("post_id")
	postID, err := strconv.Atoi(postIDStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	int32PostID := int32(postID)

	post, err := server.store.GetPost(ctx, int32PostID)
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

func (server *Server) searchPosts(ctx *gin.Context) {
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

	posts, err := server.store.SearchPosts(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, posts)
}

func s3BucketName(postType PostType) string {
	switch postType {
	case PostTypeImage:
		return "post-images"
	case PostTypeVideo:
		return "post-videos"
	case PostTypeAudio:
		return "post-audios"
	case PostTypeDocument:
		return "post-documents"
	default:
		return "post-others"
	}
}
