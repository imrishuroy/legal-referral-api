package api

import (
	"context"
	"encoding/json"
	"errors"
	"firebase.google.com/go/v4/auth"
	"fmt"
	"github.com/gin-gonic/gin"
	db "github.com/imrishuroy/legal-referral/db/sqlc"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
	"math/rand"
	"net/http"
	"time"
)

type feedPost struct {
	OwnerID        string      `json:"owner_id"`
	OwnerFirstName string      `json:"owner_first_name"`
	OwnerLastName  string      `json:"owner_last_name"`
	OwnerAvatarUrl *string     `json:"owner_avatar_url"`
	PostID         int32       `json:"post_id"`
	Content        *string     `json:"content"`
	Media          []string    `json:"media"`
	PostType       db.PostType `json:"post_type"`
	PollID         *int32      `json:"poll_id"`
	CreatedAt      time.Time   `json:"created_at"`
	LikesCount     int64       `json:"likes_count"`
	CommentsCount  int64       `json:"comments_count"`
	IsLiked        bool        `json:"is_liked"`
}

type feed struct {
	FeedID   int32     `json:"feed_id"`
	FeedType string    `json:"feed_type"`
	FeedPost *feedPost `json:"feed_post"`
	Ad       *db.Ad    `json:"ad"`
}

type listNewsFeedReq struct {
	Limit  int32 `form:"limit" binding:"required"`
	Offset int32 `form:"offset" binding:"required"`
}

func (server *Server) listNewsFeedV2(ctx *gin.Context) {
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

	redisKey := server.buildFeedCacheKey(userID, req.Limit, req.Offset)

	// Try to get the feed from the cache
	if feedList, err := server.getCachedFeed(ctx, redisKey); err == nil {
		log.Printf("Feed list from cache: %v", feedList)
		ctx.JSON(http.StatusOK, feedList)
		return
	} else if !errors.Is(err, redis.Nil) {
		// Log non-cache-miss errors without returning 500 to client
		log.Printf("Redis error: %v", err)
	}

	// Cache miss or error: Fetch from DB
	arg := db.ListNewsFeedParams{
		UserID: userID,
		Limit:  req.Limit,
		Offset: maxOffset(0, (req.Offset-1)*req.Limit),
	}

	feedPosts, err := server.store.ListNewsFeed(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching feed posts"})
		return
	}

	feedList := server.buildFeedList(feedPosts)

	// Attempt to cache asynchronously to avoid delays
	go func() {
		if err := server.cacheFeed(ctx, redisKey, feedList, 10*time.Minute); err != nil {
			log.Printf("Error caching feed: %v", err)
		}
	}()

	ctx.JSON(http.StatusOK, feedList)
}

// Helper to build Redis key for feed
func (server *Server) buildFeedCacheKey(userID string, limit, offset int32) string {
	return fmt.Sprintf("user:%s:feed:limit:%d:offset:%d", userID, limit, offset)
}

// Helper to get cached feed from Redis
func (server *Server) getCachedFeed(ctx context.Context, key string) ([]feed, error) {
	cachedData, err := server.rdb.Get(ctx, key).Result()
	if err != nil {
		return nil, err
	}

	var feedList []feed
	if err := json.Unmarshal([]byte(cachedData), &feedList); err != nil {
		return nil, fmt.Errorf("error deserializing feed data: %v", err)
	}

	return feedList, nil
}

// Helper to cache feed in Redis
func (server *Server) cacheFeed(ctx context.Context, key string, feedList []feed, expiration time.Duration) error {
	data, err := json.Marshal(feedList)
	if err != nil {
		return fmt.Errorf("error serializing feed data: %v", err)
	}

	return server.rdb.Set(ctx, key, data, expiration).Err()
}

// Helper to build feed list from database results
func (server *Server) buildFeedList(feedPosts []db.ListNewsFeedRow) []feed {
	feedList := make([]feed, len(feedPosts))
	for i, post := range feedPosts {
		feedList[i] = feed{
			FeedID:   post.FeedID,
			FeedType: "post",
			FeedPost: &feedPost{
				OwnerID:        post.OwnerID,
				OwnerFirstName: post.OwnerFirstName,
				OwnerLastName:  post.OwnerLastName,
				OwnerAvatarUrl: post.OwnerAvatarUrl,
				PostID:         post.PostID,
				Content:        post.Content,
				Media:          post.Media,
				PostType:       post.PostType,
				PollID:         post.PollID,
				CreatedAt:      post.CreatedAt,
				LikesCount:     post.LikesCount,
				CommentsCount:  post.CommentsCount,
				IsLiked:        post.IsLiked,
			},
		}
	}
	return feedList
}

// max helper function to ensure non-negative offset
func maxOffset(a, b int32) int32 {
	if a > b {
		return a
	}
	return b
}

func (server *Server) CacheUserFeed(ctx *gin.Context, userID string, feed []feed) error {
	redisKey := "user:" + userID + ":feed"

	// Serialize feed data to JSON
	jsonData, err := json.Marshal(feed)
	if err != nil {
		return fmt.Errorf("failed to marshal feed data: %v", err)
	}

	// Store data in Redis with a 10-minute expiration
	err = server.rdb.Set(ctx, redisKey, jsonData, 10*time.Minute).Err()
	if err != nil {
		return fmt.Errorf("failed to cache feed in Redis: %v", err)
	}

	return nil
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
		Offset: maxOffset(0, (req.Offset-1)*req.Limit),
	}

	// Fetch the feed posts from the database
	feedPosts, err := server.store.ListNewsFeed(ctx, arg)
	if err != nil {
		log.Printf("Error fetching feed posts: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching feed posts"})
		return
	}

	feedList := server.buildFeedList(feedPosts)

	if len(feedList) == 0 {
		ctx.JSON(http.StatusOK, feedList)
		return
	}

	// Fetch a random ad
	randomAd, err := server.store.GetRandomAd(ctx)
	if err != nil && !errors.Is(err, db.ErrRecordNotFound) {
		log.Printf("Error fetching ads: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching ads"})
		return
	}

	// Add the ad to the feed if available
	feedList = server.insertAdAtRandomPosition(feedList, randomAd)

	ctx.JSON(http.StatusOK, feedList)
}

// insertAdAtRandomPosition inserts an ad at a random position within the feed list
func (server *Server) insertAdAtRandomPosition(feedList []feed, ad db.Ad) []feed {
	randSource := rand.NewSource(time.Now().UnixNano())
	random := rand.New(randSource)

	// Generate a random index within the feed list bounds
	randomIndex := random.Intn(len(feedList) + 1)

	// Insert the ad at the random position
	newFeedList := append(feedList[:randomIndex], append([]feed{{Ad: &ad, FeedType: "ad"}}, feedList[randomIndex:]...)...)
	return newFeedList
}
