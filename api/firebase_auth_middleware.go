package api

const (
	authorizationHeaderKey  = "authorization"
	authorizationTypeBearer = "bearer"
	authorizationPayloadKey = "authorization_payload"
)

// func authMiddleware(auth *auth.Client) gin.HandlerFunc {
// 	return func(ctx *gin.Context) {
// 		authrorizationHeader := ctx.GetHeader(authorizationHeaderKey)
// 		if len(authrorizationHeader) == 0 {
// 			err := errors.New("authorization header is not provided")
// 			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
// 			return
// 		}
// 		fields := strings.Fields(authrorizationHeader)
// 		if len(fields) < 2 {
// 			err := errors.New("invalid authorization header format")
// 			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
// 			return
// 		}

// 		auhorizationType := strings.ToLower(fields[0])
// 		if auhorizationType != authorizationTypeBearer {
// 			err := fmt.Errorf("unsupported authrorization type %s", auhorizationType)
// 			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
// 			return
// 		}
// 		accessToken := fields[1]
// 		auth, err := auth.VerifyIDToken(context.Background(), accessToken)
// 		if err != nil {
// 			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
// 			return
// 		}

// 		fmt.Printf("firebase auth %v ", auth)

// 		ctx.Set(authorizationPayloadKey, auth)
// 		ctx.Next()

// 	}
// }
