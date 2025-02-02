package api

import (
	"errors"
	"firebase.google.com/go/v4/auth"
	"fmt"
	"github.com/gin-gonic/gin"
	db "github.com/imrishuroy/legal-referral/db/sqlc"
	"github.com/rs/zerolog/log"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

type feedPost struct {
	OwnerID           string      `json:"owner_id"`
	OwnerFirstName    string      `json:"owner_first_name"`
	OwnerLastName     string      `json:"owner_last_name"`
	OwnerAvatarUrl    *string     `json:"owner_avatar_url"`
	OwnerPracticeArea *string     `json:"owner_practice_area"`
	PostID            int32       `json:"post_id"`
	Content           *string     `json:"content"`
	Media             []string    `json:"media"`
	PostType          db.PostType `json:"post_type"`
	PollID            *int32      `json:"poll_id"`
	CreatedAt         time.Time   `json:"created_at"`
	LikesCount        int64       `json:"likes_count"`
	CommentsCount     int64       `json:"comments_count"`
	IsLiked           bool        `json:"is_liked"`
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

type post struct {
	PostID    int32     `json:"post_id"`
	OwnerID   string    `json:"owner_id"`
	Content   *string   `json:"content"`
	Media     []string  `json:"media"`
	PostType  PostType  `json:"post_type"`
	PollID    *int32    `json:"poll_id"`
	CreatedAt time.Time `json:"created_at"`
}

type postMetaData struct {
	PostID            int32   `json:"post_id"`
	OwnerFirstName    string  `json:"owner_first_name"`
	OwnerLastName     string  `json:"owner_last_name"`
	OwnerAvatarUrl    *string `json:"owner_avatar_url"`
	OwnerPracticeArea *string `json:"owner_practice_area"`
	LikesCount        int64   `json:"likes_count"`
	CommentsCount     int64   `json:"comments_count"`
	IsLiked           bool    `json:"is_liked"`
}

//func (server *Server) listNewsFeedV3(ctx *gin.Context) {
//	var req listNewsFeedReq
//	if err := ctx.ShouldBindQuery(&req); err != nil {
//		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request body"})
//		return
//	}
//
//	userID := ctx.Param("user_id")
//	authPayload := ctx.MustGet(authorizationPayloadKey).(*auth.Token)
//	if authPayload.UID != userID {
//		ctx.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
//		return
//	}
//
//	arg := db.ListNewsFeedV3Params{
//		UserID: userID,
//		Limit:  req.Limit,
//		Offset: maxOffset(0, (req.Offset-1)*req.Limit),
//	}
//
//	newsFeed, err := server.Store.ListNewsFeedV3(ctx, arg)
//	if err != nil {
//		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching feed posts"})
//		return
//	}
//
//	var postsToFetchFromDB []int32
//	var posts []post
//	feedPostMap := make(map[int32]int32)
//
//	for _, feed := range newsFeed {
//		feedPostMap[feed.PostID] = feed.FeedID
//		redisKey := fmt.Sprintf("post:%d", feed.PostID)
//		post, err := server.getCachedPost(ctx, redisKey)
//		if err != nil {
//			postsToFetchFromDB = append(postsToFetchFromDB, feed.PostID)
//		} else {
//			posts = append(posts, *post)
//		}
//	}
//
//	if len(postsToFetchFromDB) > 0 {
//		postsFromDB, err := server.Store.ListPosts(ctx, postsToFetchFromDB)
//		if err != nil {
//			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching missing posts from database"})
//			return
//		}
//		postsMap := make(map[string]*post)
//		for _, p := range postsFromDB {
//			redisKey := fmt.Sprintf("post:%d", p.PostID)
//			post := post{
//				PostID:    p.PostID,
//				OwnerID:   p.OwnerID,
//				Content:   p.Content,
//				Media:     p.Media,
//				PostType:  PostType(p.PostType),
//				PollID:    p.PollID,
//				CreatedAt: p.CreatedAt,
//			}
//			postsMap[redisKey] = &post
//			posts = append(posts, post)
//		}
//		if err := server.cachePosts(ctx, postsMap, 12*time.Hour); err != nil {
//			log.Printf("Error caching posts: %v", err)
//		}
//	}
//
//	postsIDs := extractPostIDs(posts)
//	postMetaArg := db.PostsMetaDataParams{
//		UserID:  userID,
//		PostIds: postsIDs,
//	}
//
//	postMetaDataMap := make(map[int32]postMetaData)
//
//	metaData, err := server.Store.PostsMetaData(ctx, postMetaArg)
//	for _, md := range metaData {
//		postMetaDataMap[md.PostID] = postMetaData{
//			PostID:            md.PostID,
//			OwnerFirstName:    md.OwnerFirstName,
//			OwnerLastName:     md.OwnerLastName,
//			OwnerAvatarUrl:    md.OwnerAvatarUrl,
//			OwnerPracticeArea: md.OwnerPracticeArea,
//			LikesCount:        md.LikesCount,
//			CommentsCount:     md.CommentsCount,
//			IsLiked:           md.IsLiked,
//		}
//	}
//
//	feedLists := createFeedList(posts, postMetaDataMap, feedPostMap)
//	ctx.JSON(http.StatusOK, feedLists)
//}

func createFeedList(posts []post, postMetaDataMap map[int32]postMetaData, feedPostMap map[int32]int32) []feed {
	feedLists := make([]feed, 0, len(posts))

	for _, post := range posts {
		metaData := postMetaDataMap[post.PostID]
		feed := feed{
			FeedID:   feedPostMap[post.PostID],
			FeedType: "post",
			FeedPost: &feedPost{
				OwnerID:           post.OwnerID,
				OwnerFirstName:    metaData.OwnerFirstName,
				OwnerLastName:     metaData.OwnerLastName,
				OwnerAvatarUrl:    metaData.OwnerAvatarUrl,
				OwnerPracticeArea: metaData.OwnerPracticeArea,
				PostID:            post.PostID,
				Content:           post.Content,
				Media:             post.Media,
				PostType:          db.PostType(post.PostType),
				PollID:            post.PollID,
				CreatedAt:         post.CreatedAt,
				LikesCount:        metaData.LikesCount,
				CommentsCount:     metaData.CommentsCount,
				IsLiked:           metaData.IsLiked,
			},
		}
		feedLists = append(feedLists, feed)
	}
	return feedLists
}

func extractPostIDs(posts []post) []int32 {
	postsIDs := make([]int32, 0, len(posts))
	for _, post := range posts {
		postsIDs = append(postsIDs, post.PostID)
	}
	return postsIDs
}

//
//func (server *Server) getCachedPost(ctx context.Context, key string) (*post, error) {
//	cachedData, err := server.rdb.Get(ctx, key).Result()
//	if err != nil {
//		return nil, err
//	}
//
//	var post post
//	if err := json.Unmarshal([]byte(cachedData), &post); err != nil {
//		return nil, fmt.Errorf("error deserializing post data: %v", err)
//	}
//
//	return &post, nil
//}
//
//func (server *Server) cachePosts(ctx context.Context, posts map[string]*post, expiration time.Duration) error {
//	pipe := server.rdb.Pipeline() // Use a Redis pipeline for batch operations
//
//	for key, post := range posts {
//		data, err := json.Marshal(post)
//		if err != nil {
//			return fmt.Errorf("error serializing post data for key %s: %v", key, err)
//		}
//		pipe.Set(ctx, key, data, expiration)
//	}
//
//	// Execute all commands in the pipeline
//	_, err := pipe.Exec(ctx)
//	if err != nil {
//		return fmt.Errorf("error executing Redis pipeline: %v", err)
//	}
//
//	return nil
//}

//func (server *Server) listNewsFeedV2(ctx *gin.Context) {
//	var req listNewsFeedReq
//	if err := ctx.ShouldBindQuery(&req); err != nil {
//		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request body"})
//		return
//	}
//
//	userID := ctx.Param("user_id")
//	authPayload := ctx.MustGet(authorizationPayloadKey).(*auth.Token)
//	if authPayload.UID != userID {
//		ctx.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
//		return
//	}
//
//	redisKey := server.buildFeedCacheKey(userID, req.Limit, req.Offset)
//
//	// Try to get the feed from the cache
//	if feedList, err := server.getCachedFeed(ctx, redisKey); err == nil {
//		log.Printf("Feed list from cache: %v", feedList)
//		ctx.JSON(http.StatusOK, feedList)
//		return
//	} else if !errors.Is(err, redis.Nil) {
//		// Log non-cache-miss errors without returning 500 to client
//		log.Printf("Redis error: %v", err)
//	}
//
//	// Cache miss or error: Fetch from DB
//	arg := db.ListNewsFeedParams{
//		UserID: userID,
//		Limit:  req.Limit,
//		Offset: maxOffset(0, (req.Offset-1)*req.Limit),
//	}
//
//	feedPosts, err := server.Store.ListNewsFeed(ctx, arg)
//	if err != nil {
//		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching feed posts"})
//		return
//	}
//
//	feedList := server.buildFeedList(feedPosts)
//
//	// Attempt to cache asynchronously to avoid delays
//	go func() {
//		if err := server.cacheFeed(ctx, redisKey, feedList, 10*time.Minute); err != nil {
//			log.Printf("Error caching feed: %v", err)
//		}
//	}()
//
//	ctx.JSON(http.StatusOK, feedList)
//}

// Helper to build Redis key for feed
func (srv *Server) buildFeedCacheKey(userID string, limit, offset int32) string {
	return fmt.Sprintf("user:%s:feed:limit:%d:offset:%d", userID, limit, offset)
}

// Helper to get cached feed from Redis
//func (server *Server) getCachedFeed(ctx context.Context, key string) ([]feed, error) {
//	cachedData, err := server.rdb.Get(ctx, key).Result()
//	if err != nil {
//		return nil, err
//	}
//
//	var feedList []feed
//	if err := json.Unmarshal([]byte(cachedData), &feedList); err != nil {
//		return nil, fmt.Errorf("error deserializing feed data: %v", err)
//	}
//
//	return feedList, nil
//}

// Helper to cache feed in Redis
//func (server *Server) cacheFeed(ctx context.Context, key string, feedList []feed, expiration time.Duration) error {
//	data, err := json.Marshal(feedList)
//	if err != nil {
//		return fmt.Errorf("error serializing feed data: %v", err)
//	}
//
//	return server.rdb.Set(ctx, key, data, expiration).Err()
//}

// Helper to build feed list from database results
func (srv *Server) buildFeedList(feedPosts []db.ListNewsFeedRow) []feed {
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

//func (server *Server) CacheUserFeed(ctx *gin.Context, userID string, feed []feed) error {
//	redisKey := "user:" + userID + ":feed"
//
//	// Serialize feed data to JSON
//	jsonData, err := json.Marshal(feed)
//	if err != nil {
//		return fmt.Errorf("failed to marshal feed data: %v", err)
//	}
//
//	// Store data in Redis with a 10-minute expiration
//	err = server.rdb.Set(ctx, redisKey, jsonData, 10*time.Minute).Err()
//	if err != nil {
//		return fmt.Errorf("failed to cache feed in Redis: %v", err)
//	}
//
//	return nil
//}

func (srv *Server) ListNewsFeed(ctx *gin.Context) {
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
	feedPosts, err := srv.Store.ListNewsFeed(ctx, arg)
	if err != nil {
		log.Printf("Error fetching feed posts: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching feed posts"})
		return
	}

	feedList := srv.buildFeedList(feedPosts)

	if len(feedList) == 0 {
		ctx.JSON(http.StatusOK, feedList)
		return
	}

	// Fetch a random ad
	randomAd, err := srv.Store.GetRandomAd(ctx)
	if err != nil && !errors.Is(err, db.ErrRecordNotFound) {
		log.Printf("Error fetching ads: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching ads"})
		return
	}

	// Add the ad to the feed if available
	feedList = srv.insertAdAtRandomPosition(feedList, randomAd)

	ctx.JSON(http.StatusOK, feedList)
}

// insertAdAtRandomPosition inserts an ad at a random position within the feed list
func (srv *Server) insertAdAtRandomPosition(feedList []feed, ad db.Ad) []feed {
	randSource := rand.NewSource(time.Now().UnixNano())
	random := rand.New(randSource)

	// Generate a random index within the feed list bounds
	randomIndex := random.Intn(len(feedList) + 1)

	// Insert the ad at the random position
	newFeedList := append(feedList[:randomIndex], append([]feed{{Ad: &ad, FeedType: "ad"}}, feedList[randomIndex:]...)...)
	return newFeedList
}

func (srv *Server) IgnoreFeed(ctx *gin.Context) {
	feedIdStr := ctx.Param("feed_id")
	feedID, err := strconv.Atoi(feedIdStr)
	if err != nil {
		ctx.JSON(400, errorResponse(err))
		return
	}
	authPayload := ctx.MustGet(authorizationPayloadKey).(*auth.Token)
	if authPayload.UID == "" {
		ctx.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}

	arg := db.IgnoreFeedParams{
		FeedID: int32(feedID),
		UserID: authPayload.UID,
	}

	err = srv.Store.IgnoreFeed(ctx, arg)
	if err != nil {
		log.Printf("Error ignoring post: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error ignoring post"})
		return
	}

	ctx.Status(http.StatusOK)

}
