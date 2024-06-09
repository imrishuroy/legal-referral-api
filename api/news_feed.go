package api

import (
	"firebase.google.com/go/v4/auth"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

type User struct {
	UserID    string  `json:"user_id"`
	FirstName string  `json:"first_name"`
	LastName  string  `json:"last_name"`
	AvatarURL *string `json:"avatar_url"`
}

type Post struct {
	PostID    int32     `json:"post_id"`
	OwnerID   string    `json:"owner_id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	Media     []string  `json:"media"`
	PostType  PostType  `json:"post_type"`
	CreatedAt time.Time `json:"created_at"`
}

type NewsFeed struct {
	Post Post `json:"post"`
	User User `json:"user"`
}

func (server *Server) listNewsFeed(ctx *gin.Context) {
	userID := ctx.Param("user_id")

	authPayload := ctx.MustGet(authorizationPayloadKey).(*auth.Token)
	if authPayload.UID != userID {
		ctx.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}

	feed, err := server.store.ListNewsFeed2(ctx, userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	//newsFeeds := make([]NewsFeed, 0)
	//for _, f := range feed {
	//	user := User{
	//		UserID:    f.UserID,
	//		FirstName: f.FirstName,
	//		LastName:  f.LastName,
	//		AvatarURL: f.AvatarUrl,
	//	}
	//	post := Post{
	//		PostID:   f.PostID,
	//		OwnerID:  f.UserID,
	//		Title:    f.Title,
	//		Content:  f.Content,
	//		Media:    f.Media,
	//		PostType: PostType(f.PostType),
	//	}
	//
	//	newsFeeds = append(newsFeeds, NewsFeed{
	//		User: user,
	//		Post: post,
	//	})
	//}
	//
	//ctx.JSON(http.StatusOK, newsFeeds)

	ctx.JSON(http.StatusOK, feed)

}
