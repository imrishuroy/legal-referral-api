package api

import (
	"errors"
	"firebase.google.com/go/v4/auth"
	"github.com/rs/zerolog/log"
	"mime/multipart"
	"net/http"
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
	Content   string                  `form:"content" binding:"required"`
	Files     []*multipart.FileHeader `form:"files"`
	PostType  PostType                `form:"post_type" binding:"required"`
	PollTitle string                  `form:"poll_title"`
	Options   []string                `form:"options"`
	EndTime   *time.Time              `form:"end_time"`
	//PollID   *int32                  `form:"poll_id"`
	//Poll     *createPollReq          `form:"poll"`
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

	if req.PostType == PostTypeImage {
		urls, err := server.handleFileUpload(ctx, req.Files)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}
		log.Info().Msgf("URLs: %+v", urls)
		imageUrls = append(imageUrls, urls...)
	}

	var pollID *int32

	createPollReq := createPollReq{
		OwnerID:   req.OwnerID,
		PollTitle: req.PollTitle,
		Options:   req.Options,
		EndTime:   req.EndTime,
	}

	if req.PostType == PostTypePoll {
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

	// this should aslo throw an error
	server.createProducer(req.OwnerID, string(post.PostID))
	//server.publishToKafka(req.OwnerID, string(post.PostID))

	//if err := server.postToNewsFeed(ctx, req.OwnerID, post.PostID); err != nil {
	//	ctx.JSON(http.StatusInternalServerError, errorResponse(err))
	//	return
	//}

	ctx.JSON(http.StatusOK, "Post created successfully")
}

func (server *Server) handleFileUpload(ctx *gin.Context, files []*multipart.FileHeader) ([]string, error) {
	if len(files) == 0 {
		return nil, errors.New("no file uploaded")
	}

	urls := make([]string, 0, len(files))
	for _, file := range files {
		url, err := server.uploadFileHandler(file)
		if err != nil {
			return nil, err
		}
		urls = append(urls, url)
	}
	return urls, nil
}

func (server *Server) uploadFileHandler(file *multipart.FileHeader) (string, error) {
	fileName := generateUniqueFilename() + getFileExtension(file)
	multiPartFile, err := file.Open()
	if err != nil {
		return "", err
	}
	defer multiPartFile.Close()

	return server.uploadFile(multiPartFile, fileName, file.Header.Get("Content-Type"), "post-images")
}

func (server *Server) createPoll(ctx *gin.Context, req *createPollReq) (*db.Poll, error) {
	if req == nil {
		return nil, errors.New("poll request is nil")
	}

	arg := db.CreatePollParams{
		OwnerID: req.OwnerID,
		Title:   req.PollTitle,
		Options: req.Options,
		//EndTime: pgtype.Timestamptz{Time: *req.EndTime, Valid: req.EndTime != nil},
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
		server.createProducer(id.(string), string(postID))
		//if err := server.store.PostToNewsFeed(ctx, arg); err != nil {
		//	return err
		//}
	}
	return nil
}
