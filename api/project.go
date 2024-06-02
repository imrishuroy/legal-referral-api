package api

import (
	"firebase.google.com/go/v4/auth"
	"github.com/gin-gonic/gin"
	db "github.com/imrishuroy/legal-referral/db/sqlc"
	"net/http"
	"strconv"
	"time"
)

type createReferralReq struct {
	ReferredUserIDs           []string `json:"referred_user_ids"`
	ReferrerUserID            string   `json:"referrer_user_id"`
	Title                     string   `json:"title"`
	PreferredPracticeArea     string   `json:"preferred_practice_area"`
	PreferredPracticeLocation string   `json:"preferred_practice_location"`
	CaseDescription           string   `json:"case_description"`
}

func (server *Server) createReferral(ctx *gin.Context) {
	var req createReferralReq
	if err := ctx.BindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	authPayload := ctx.MustGet(authorizationPayloadKey).(*auth.Token)
	if authPayload.UID != req.ReferrerUserID {
		ctx.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}

	// add project to project table
	arg := db.CreateReferralParams{
		ReferrerUserID:            req.ReferrerUserID,
		Title:                     req.Title,
		PreferredPracticeArea:     req.PreferredPracticeArea,
		PreferredPracticeLocation: req.PreferredPracticeLocation,
		CaseDescription:           req.CaseDescription,
	}

	project, err := server.store.CreateReferral(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// add referral to referral table
	for _, referredUserID := range req.ReferredUserIDs {
		arg := db.AddReferredUserToProjectParams{
			ProjectID:      project.ProjectID,
			ReferredUserID: referredUserID,
		}
		_, err := server.store.AddReferredUserToProject(ctx, arg)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}
	}

	ctx.JSON(http.StatusOK, "Referral created")
}

func (server *Server) listActiveReferrals(ctx *gin.Context) {
	userID := ctx.Param("user_id")
	authPayload := ctx.MustGet(authorizationPayloadKey).(*auth.Token)
	if authPayload.UID != userID {
		ctx.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}
	projects, err := server.store.ListActiveReferrals(ctx, userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, projects)
}

func (server *Server) listReferredUsers(ctx *gin.Context) {
	projectIDStr := ctx.Param("project_id")

	authPayload := ctx.MustGet(authorizationPayloadKey).(*auth.Token)
	if authPayload.UID == "" {
		ctx.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}
	projectID, err := strconv.Atoi(projectIDStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	users, err := server.store.ListReferredUsers2(ctx, int32(projectID))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, users)
}

type projectUser struct {
	UserID       string  `json:"user_id"`
	FirstName    string  `json:"first_name"`
	LastName     string  `json:"last_name"`
	AvatarUrl    *string `json:"avatar_url"`
	PracticeArea *string `json:"practice_area"`
}

type project struct {
	ProjectID        int32            `json:"project_id"`
	Status           db.ProjectStatus `json:"status"`
	CreatedAt        *time.Time       `json:"created_at"`
	StartedAt        *time.Time       `json:"started_at"`
	CompletedAt      *time.Time       `json:"completed_at"`
	Title            string           `json:"title"`
	CaseDescription  string           `json:"case_description"`
	PracticeLocation *string          `json:"preferred_practice_location"`
	PracticeArea     *string          `json:"preferred_practice_area"`
	projectUser      `json:"user"`
}

func (server *Server) listActiveProposals(ctx *gin.Context) {
	userID := ctx.Param("user_id")

	authPayload := ctx.MustGet(authorizationPayloadKey).(*auth.Token)
	if authPayload.UID != userID {
		ctx.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}

	projects, err := server.store.ListActiveProposals(ctx, userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	var projectList []project
	for _, p := range projects {
		var startedAt, completedAt *time.Time

		projectList = append(projectList, project{
			ProjectID:        p.ProjectID,
			Status:           p.Status,
			StartedAt:        startedAt,
			CreatedAt:        &p.CreatedAt,
			CompletedAt:      completedAt,
			Title:            p.Title,
			CaseDescription:  p.CaseDescription,
			PracticeLocation: &p.PreferredPracticeLocation,
			PracticeArea:     &p.PreferredPracticeArea,
			projectUser: projectUser{
				UserID:       p.UserID,
				FirstName:    p.FirstName,
				LastName:     p.LastName,
				AvatarUrl:    p.AvatarUrl,
				PracticeArea: p.PracticeArea,
			},
		})
	}

	if len(projectList) == 0 {
		projectList = []project{}
	}

	ctx.JSON(http.StatusOK, projectList)
}

func (server *Server) awardProject(ctx *gin.Context) {
	var req db.AwardProjectParams
	if err := ctx.BindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*auth.Token)
	if authPayload.UID == "" {
		ctx.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}

	status, err := server.store.GetProjectStatus(ctx, req.ProjectID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	if status == "awarded" {
		ctx.JSON(http.StatusConflict, gin.H{"message": "Project already awarded"})
		return
	}

	project, err := server.store.AwardProject(ctx, req)

	if err != nil {
		if db.ErrorCode(err) == db.UniqueViolation {
			ctx.JSON(http.StatusConflict, gin.H{"message": "Project already awarded"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, project)
}

func (server *Server) listAwardedProjects(ctx *gin.Context) {
	userID := ctx.Param("user_id")

	authPayload := ctx.MustGet(authorizationPayloadKey).(*auth.Token)
	if authPayload.UID != userID {
		ctx.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}

	projects, err := server.store.ListAwardedProjects(ctx, userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	projectList := make([]project, 0, len(projects)) // Preallocate slice with known length

	for _, p := range projects {
		var startedAt, completedAt *time.Time

		if p.StartedAt.Valid {
			startedAt = &p.StartedAt.Time
		}
		if p.CompletedAt.Valid {
			completedAt = &p.CompletedAt.Time
		}

		projectList = append(projectList, project{
			ProjectID:       p.ProjectID,
			Status:          p.Status,
			StartedAt:       startedAt,
			CreatedAt:       &p.CreatedAt,
			CompletedAt:     completedAt,
			Title:           p.Title,
			CaseDescription: p.CaseDescription,
			projectUser: projectUser{
				UserID:       p.UserID,
				FirstName:    p.FirstName,
				LastName:     p.LastName,
				AvatarUrl:    p.AvatarUrl,
				PracticeArea: p.PracticeArea,
			},
		})
	}

	ctx.JSON(http.StatusOK, projectList)
}

func (server *Server) acceptProject(ctx *gin.Context) {
	projectIdParam := ctx.Param("project_id")
	projectID, err := strconv.Atoi(projectIdParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*auth.Token)
	if authPayload.UID == "" {
		ctx.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}

	arg := db.AcceptProjectParams{
		ProjectID: int32(projectID),
		UserID:    authPayload.UID,
	}

	project, err := server.store.AcceptProject(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, project)
}

func (server *Server) rejectProject(ctx *gin.Context) {
	projectIdParam := ctx.Param("project_id")
	projectID, err := strconv.Atoi(projectIdParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*auth.Token)
	if authPayload.UID == "" {
		ctx.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}

	arg := db.RejectProjectParams{
		ProjectID: int32(projectID),
		UserID:    authPayload.UID,
	}

	project, err := server.store.RejectProject(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, project)
}

func (server *Server) listActiveProjects(ctx *gin.Context) {
	userID := ctx.Param("user_id")
	role := ctx.Query("role")

	authPayload := ctx.MustGet(authorizationPayloadKey).(*auth.Token)
	if authPayload.UID != userID {
		ctx.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}

	if role != "referrer" && role != "referred" {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Invalid role"})
		return
	}

	if role == "referrer" {
		projects, err := server.store.ListReferrerActiveProjects(ctx, userID)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}

		var projectList []project
		for _, p := range projects {
			var startedAt, completedAt *time.Time

			if p.StartedAt.Valid {
				startedAt = &p.StartedAt.Time
			}
			if p.CompletedAt.Valid {
				completedAt = &p.CompletedAt.Time
			}

			projectList = append(projectList, project{
				ProjectID:       p.ProjectID,
				Status:          p.Status,
				StartedAt:       startedAt,
				CreatedAt:       &p.CreatedAt,
				CompletedAt:     completedAt,
				Title:           p.Title,
				CaseDescription: p.CaseDescription,
				projectUser: projectUser{
					UserID:       p.UserID,
					FirstName:    p.FirstName,
					LastName:     p.LastName,
					AvatarUrl:    p.AvatarUrl,
					PracticeArea: p.PracticeArea,
				},
			})
		}

		if len(projectList) == 0 {
			projectList = []project{}

		}

		ctx.JSON(http.StatusOK, projectList)
	} else {
		projects, err := server.store.ListReferredActiveProjects(ctx, userID)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}

		var projectList []project
		for _, p := range projects {
			var startedAt, completedAt *time.Time

			if p.StartedAt.Valid {
				startedAt = &p.StartedAt.Time
			}
			if p.CompletedAt.Valid {
				completedAt = &p.CompletedAt.Time
			}

			projectList = append(projectList, project{
				ProjectID:       p.ProjectID,
				Status:          p.Status,
				StartedAt:       startedAt,
				CreatedAt:       &p.CreatedAt,
				CompletedAt:     completedAt,
				Title:           p.Title,
				CaseDescription: p.CaseDescription,
				projectUser: projectUser{
					UserID:       p.UserID,
					FirstName:    p.FirstName,
					LastName:     p.LastName,
					AvatarUrl:    p.AvatarUrl,
					PracticeArea: p.PracticeArea,
				},
			})
		}

		if len(projectList) == 0 {
			projectList = []project{}
		}

		ctx.JSON(http.StatusOK, projectList)
	}

}

func (server *Server) startProject(ctx *gin.Context) {
	projectIdParam := ctx.Param("project_id")
	projectID, err := strconv.Atoi(projectIdParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*auth.Token)
	if authPayload.UID == "" {
		ctx.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}

	arg := db.StartProjectParams{
		ProjectID: int32(projectID),
		UserID:    authPayload.UID,
	}

	project, err := server.store.StartProject(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, project)
}

func (server *Server) initiateCompleteProject(ctx *gin.Context) {
	projectIdParam := ctx.Param("project_id")
	projectID, err := strconv.Atoi(projectIdParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*auth.Token)
	if authPayload.UID == "" {
		ctx.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}

	arg := db.InitiateCompleteProjectParams{
		ProjectID: int32(projectID),
		UserID:    authPayload.UID,
	}

	project, err := server.store.InitiateCompleteProject(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, project)
}

func (server *Server) cancelInitiateCompleteProject(ctx *gin.Context) {
	projectIdParam := ctx.Param("project_id")
	projectID, err := strconv.Atoi(projectIdParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*auth.Token)
	if authPayload.UID == "" {
		ctx.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}

	arg := db.CancelCompleteProjectInitiationParams{
		ProjectID: int32(projectID),
		UserID:    authPayload.UID,
	}

	project, err := server.store.CancelCompleteProjectInitiation(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, project)
}

func (server *Server) completeProject(ctx *gin.Context) {
	projectIdParam := ctx.Param("project_id")
	projectID, err := strconv.Atoi(projectIdParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*auth.Token)
	if authPayload.UID == "" {
		ctx.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}

	arg := db.CompleteProjectParams{
		ProjectID: int32(projectID),
		UserID:    authPayload.UID,
	}

	project, err := server.store.CompleteProject(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, project)
}

func (server *Server) listCompletedProjects(ctx *gin.Context) {

	userID := ctx.Param("user_id")
	role := ctx.Query("role")

	authPayload := ctx.MustGet(authorizationPayloadKey).(*auth.Token)
	if authPayload.UID != userID {
		ctx.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}

	if role != "referrer" && role != "referred" {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Invalid role"})
		return
	}

	if role == "referrer" {
		projects, err := server.store.ListReferrerCompletedProjects(ctx, userID)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}

		var projectList []project
		for _, p := range projects {
			var startedAt, completedAt *time.Time

			if p.StartedAt.Valid {
				startedAt = &p.StartedAt.Time
			}
			if p.CompletedAt.Valid {
				completedAt = &p.CompletedAt.Time
			}

			projectList = append(projectList, project{
				ProjectID:        p.ProjectID,
				Status:           p.Status,
				StartedAt:        startedAt,
				CreatedAt:        &p.CreatedAt,
				CompletedAt:      completedAt,
				Title:            p.Title,
				CaseDescription:  p.CaseDescription,
				PracticeLocation: &p.PreferredPracticeLocation,
				PracticeArea:     &p.PreferredPracticeArea,
				projectUser: projectUser{
					UserID:       p.UserID,
					FirstName:    p.FirstName,
					LastName:     p.LastName,
					AvatarUrl:    p.AvatarUrl,
					PracticeArea: p.PracticeArea,
				},
			})
		}

		if len(projectList) == 0 {
			projectList = []project{}

		}

		ctx.JSON(http.StatusOK, projectList)

	} else {
		projects, err := server.store.ListReferredCompletedProjects(ctx, userID)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}

		var projectList []project
		for _, p := range projects {
			var startedAt, completedAt *time.Time

			if p.StartedAt.Valid {
				startedAt = &p.StartedAt.Time
			}
			if p.CompletedAt.Valid {
				completedAt = &p.CompletedAt.Time
			}

			projectList = append(projectList, project{
				ProjectID:        p.ProjectID,
				Status:           p.Status,
				StartedAt:        startedAt,
				CreatedAt:        &p.CreatedAt,
				CompletedAt:      completedAt,
				Title:            p.Title,
				CaseDescription:  p.CaseDescription,
				PracticeLocation: &p.PreferredPracticeLocation,
				PracticeArea:     &p.PreferredPracticeArea,
				projectUser: projectUser{
					UserID:       p.UserID,
					FirstName:    p.FirstName,
					LastName:     p.LastName,
					AvatarUrl:    p.AvatarUrl,
					PracticeArea: p.PracticeArea,
				},
			})
		}

		if len(projectList) == 0 {
			projectList = []project{}
		}

		ctx.JSON(http.StatusOK, projectList)
	}

}
