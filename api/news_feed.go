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

	// Check cache
	if feedList, err := server.getCachedFeed(ctx, redisKey); err == nil {
		// print
		log.Printf("Feed list from cache: %v", feedList)
		ctx.JSON(http.StatusOK, feedList)
		return
	} else if !errors.Is(err, redis.Nil) {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Error fetching cached feed"})
		return
	}

	// Cache miss: Generate the feed
	arg := db.ListNewsFeedParams{
		UserID: userID,
		Limit:  req.Limit,
		Offset: maxOffset(0, (req.Offset-1)*req.Limit),
	}

	// Fetch feed posts from the database
	feedPosts, err := server.store.ListNewsFeed(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching feed posts"})
		return
	}

	// Prepare the feed list
	feedList := server.buildFeedList(feedPosts)

	// Cache the feed list with a 10-minute expiration
	if err := server.cacheFeed(ctx, redisKey, feedList, 10*time.Minute); err != nil {
		log.Printf("Error caching feed: %v", err)
	}

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

//
//func (server *Server) listNewsFeed(ctx *gin.Context) {
//	var req listNewsFeedReq
//	if err := ctx.ShouldBindQuery(&req); err != nil {
//		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request body"})
//		return
//	}
//
//	userID := ctx.Param("user_id")
//
//	authPayload := ctx.MustGet(authorizationPayloadKey).(*auth.Token)
//	if authPayload.UID != userID {
//		ctx.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
//		return
//	}
//
//	arg := db.ListNewsFeedParams{
//		UserID: userID,
//		Limit:  req.Limit,
//		Offset: (req.Offset - 1) * req.Limit,
//	}
//
//	feedList := make([]feed, 0)
//
//	// Fetch the feed posts
//	feedPosts, err := server.store.ListNewsFeed(ctx, arg)
//	if err != nil {
//		log.Printf("Error fetching feed posts: %v", err)
//		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching feed posts"})
//		return
//	}
//
//	// Fetch the ads
//	randomAd, err := server.store.GetRandomAd(ctx)
//
//	// If no ads are found, return the feed posts
//	if errors.Is(err, db.ErrRecordNotFound) {
//		for _, f := range feedPosts {
//			feedList = append(feedList, feed{
//				FeedPost: (*feedPost)(&f),
//				FeedType: "post",
//			})
//		}
//
//		err = server.CacheUserFeed(ctx, userID, feedList)
//		if err != nil {
//			log.Error().Err(err).Msg("Error caching feed")
//		}
//
//		err = server.rdb.Set(ctx, "foo", "bar", 10*time.Minute).Err()
//		if err != nil {
//			log.Error().Err(err).Msg("Error setting foo")
//		}
//
//		ctx.JSON(http.StatusOK, feedList)
//		return
//	}
//
//	if err != nil {
//		log.Printf("Error fetching ads: %v", err)
//		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching ads"})
//		return
//	}
//
//	for _, f := range feedPosts {
//		feedList = append(feedList, feed{
//			FeedPost: (*feedPost)(&f),
//			FeedType: "post",
//		})
//	}
//
//	randSource := rand.NewSource(time.Now().UnixNano())
//	random := rand.New(randSource)
//
//	log.Printf("Feed list: %v", len(feedList))
//
//	if len(feedList) == 0 {
//		ctx.JSON(http.StatusOK, feedList)
//		return
//	}
//
//	// Generate a random index between 1 and 10 (inclusive)
//	randomIndex := random.Intn(len(feedList)) + 1
//
//	log.Printf("Random index: %d", randomIndex)
//
//	newFeedList := make([]feed, len(feedList)+1)
//	copy(newFeedList, feedList[:randomIndex])
//	newFeedList[randomIndex] = feed{
//		Ad:       &randomAd,
//		FeedType: "ad",
//	}
//	copy(newFeedList[randomIndex+1:], feedList[randomIndex:])
//
//	ctx.JSON(http.StatusOK, newFeedList)
//}
