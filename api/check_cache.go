package api

import (
	"context"
	"encoding/json"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

const KEY = "test_cache_data"

type TestCacheData struct {
	OwnerID  string `json:"owner_id"`
	PostType string `json:"post_type"`
	Content  string `json:"content"`
}

func (srv *Server) AddTestCacheData(ctx *gin.Context) {
	testCacheData := TestCacheData{
		OwnerID:  "12345",
		PostType: "text",
		Content:  "This is a test post",
	}
	jsonData, err := json.Marshal(testCacheData)
	if err != nil {
		log.Error().Err(err).Msg("Failed to marshal test cache data")
		return
	}

	context := context.Background()
	err = srv.ValkeyClient.Do(context, srv.ValkeyClient.B().Set().Key(KEY).Value(string(jsonData)).Ex(1*time.Minute).Build()).Error()
	if err != nil {
		log.Error().Err(err).Msg("Failed to add test cache data")
		return
	}
	log.Info().Msg("Test cache data added successfully")

	ctx.JSON(200, gin.H{
		"message": "Test cache data added successfully",
	})

}

func (srv *Server) GetTestCacheData(ctx *gin.Context) {
	context := context.Background()
	data, err := srv.ValkeyClient.Do(context, srv.ValkeyClient.B().Get().Key(KEY).Build()).AsBytes()

	if err != nil {
		log.Error().Err(err).Msg("Failed to get test cache data")
		ctx.JSON(500, gin.H{
			"error": "Failed to get test cache data",
		})
		return
	}
	if data == nil {
		log.Info().Msg("No data found in cache")
		ctx.JSON(200, gin.H{
			"message": "No data found in cache",
		})
		return
	}
	var testCacheData TestCacheData
	err = json.Unmarshal(data, &testCacheData)
	if err != nil {
		log.Error().Err(err).Msg("Failed to unmarshal test cache data")
		ctx.JSON(500, gin.H{
			"error": "Failed to unmarshal test cache data",
		})
		return
	}
	log.Info().Msg("Test cache data retrieved successfully")
	ctx.JSON(200, gin.H{
		"data": testCacheData,
	})
	log.Info().Msgf("Test cache data: %v", testCacheData)

}

func (srv *Server) ListCacheKeys(ctx *gin.Context) {
	// context := context.Background()
	// keys, err := srv.ValkeyClient.Do(context, srv.ValkeyClient.B().Keys().Pattern("*").Build()).AsStrSlice()
	// if err != nil {
	// 	log.Error().Err(err).Msg("Failed to list cache keys")
	// 	ctx.JSON(500, gin.H{
	// 		"error": "Failed to list cache keys",
	// 	})
	// 	return
	// }
	// log.Info().Msg("Cache keys retrieved successfully")
	// ctx.JSON(200, gin.H{
	// 	"keys": keys,
	// })

	var cursor uint64 = 0
	for {
		resss := srv.ValkeyClient.Do(ctx, srv.ValkeyClient.B().Scan().Cursor(cursor).Match("*").Count(2).Build())
		se, err := resss.AsScanEntry()
		if err != nil {
			log.Error().Err(err).Msg("Failed to list cache keys")
			ctx.JSON(500, gin.H{
				"error": "Failed to list cache keys",
			})
			return
		}
		log.Info().Msg("Cache keys retrieved successfully")
		list := make([]string, 0)
		for _, sitem := range se.Elements {
			log.Info().Msgf("Key: %s", sitem)
			list = append(list, sitem)
		}
		ctx.JSON(200, gin.H{
			"keys": list,
		})

		cursor = se.Cursor
		if cursor == 0 {
			break
		}
	}
}
