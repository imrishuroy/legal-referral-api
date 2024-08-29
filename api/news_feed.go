package api

import (
	"errors"
	"firebase.google.com/go/v4/auth"
	"github.com/gin-gonic/gin"
	db "github.com/imrishuroy/legal-referral/db/sqlc"
	"github.com/rs/zerolog/log"
	"math/rand"
	"net/http"
	"time"
)

type feedPost struct {
	FeedID        int32     `json:"feed_id"`
	User          db.User   `json:"user"`
	Post          db.Post   `json:"post"`
	CreatedAt     time.Time `json:"created_at"`
	LikesCount    int64     `json:"likes_count"`
	CommentsCount int64     `json:"comments_count"`
	IsLiked       bool      `json:"is_liked"`
}

type feed struct {
	FeedType string    `json:"feed_type"`
	FeedPost *feedPost `json:"feed_post"`
	Ad       *db.Ad    `json:"ad"`
}

type listNewsFeedReq struct {
	Limit  int32 `form:"limit" binding:"required"`
	Offset int32 `form:"offset" binding:"required"`
}

func (server *Server) listNewsFeed(ctx *gin.Context) {
	var req listNewsFeedReq
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request body"})
		return
	}

	userID := ctx.Param("user_id")

	authPayload := ctx.MustGet(authorizationPayloadKey).(*auth.Token)
	if authPayload.UID != userID {
		ctx.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}

	arg := db.ListNewsFeedParams{
		UserID: userID,
		Limit:  req.Limit,
		Offset: (req.Offset - 1) * req.Limit,
	}

	feedList := make([]feed, 0)

	// Fetch the feed posts
	feedPosts, err := server.store.ListNewsFeed(ctx, arg)
	if err != nil {
		log.Printf("Error fetching feed posts: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching feed posts"})
		return
	}

	// Fetch the ads
	randomAd, err := server.store.GetRandomAd(ctx)

	// If no ads are found, return the feed posts
	if errors.Is(err, db.ErrRecordNotFound) {
		for _, f := range feedPosts {
			feedList = append(feedList, feed{
				FeedPost: (*feedPost)(&f),
				FeedType: "post",
			})
		}
		ctx.JSON(http.StatusOK, feedList)
		return
	}

	if err != nil {
		log.Printf("Error fetching ads: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching ads"})
		return
	}

	for _, f := range feedPosts {
		feedList = append(feedList, feed{
			FeedPost: (*feedPost)(&f),
			FeedType: "post",
		})
	}

	randSource := rand.NewSource(time.Now().UnixNano())
	random := rand.New(randSource)

	log.Printf("Feed list: %v", len(feedList))

	if len(feedList) == 0 {
		ctx.JSON(http.StatusOK, feedList)
		return
	}

	// Generate a random index between 1 and 10 (inclusive)
	randomIndex := random.Intn(len(feedList)) + 1

	log.Printf("Random index: %d", randomIndex)

	newFeedList := make([]feed, len(feedList)+1)
	copy(newFeedList, feedList[:randomIndex])
	newFeedList[randomIndex] = feed{
		Ad:       &randomAd,
		FeedType: "ad",
	}
	copy(newFeedList[randomIndex+1:], feedList[randomIndex:])

	ctx.JSON(http.StatusOK, newFeedList)
}
