package api

import (
	"firebase.google.com/go/v4/auth"
	"github.com/gin-gonic/gin"
	db "github.com/imrishuroy/legal-referral/db/sqlc"
	"github.com/rs/zerolog/log"
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

func (s *Server) CreateReferral(ctx *gin.Context) {
	var req createReferralReq
	if err := ctx.BindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	log.Info().Msgf("Referred: %s", req.ReferredUserIDs)

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

	project, err := s.Store.CreateReferral(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	log.Info().Msgf("Project: %v", project.ProjectID)

	// add referral to referral table
	for _, referredUserID := range req.ReferredUserIDs {
		log.Info().Msgf("Referred -----: %s", referredUserID)
		arg := db.AddReferredUserToProjectParams{
			ProjectID:      project.ProjectID,
			ReferredUserID: &referredUserID,
		}
		_, err := s.Store.AddReferredUserToProject(ctx, arg)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}
	}

	ctx.JSON(http.StatusOK, "Referral created")
}

func (s *Server) ListActiveReferrals(ctx *gin.Context) {
	userID := ctx.Param("user_id")
	authPayload := ctx.MustGet(authorizationPayloadKey).(*auth.Token)
	if authPayload.UID != userID {
		ctx.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}
	projects, err := s.Store.ListActiveReferrals(ctx, userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, projects)
}

func (s *Server) ListReferredUsers(ctx *gin.Context) {
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

	users, err := s.Store.ListReferredUsers2(ctx, int32(projectID))
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

func (s *Server) ListActiveProposals(ctx *gin.Context) {
	userID := ctx.Param("user_id")

	authPayload := ctx.MustGet(authorizationPayloadKey).(*auth.Token)
	if authPayload.UID != userID {
		ctx.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}

	projects, err := s.Store.ListActiveProposals(ctx, userID)
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

func (s *Server) AwardProject(ctx *gin.Context) {
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

	status, err := s.Store.GetProjectStatus(ctx, req.ProjectID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	if status == "awarded" {
		ctx.JSON(http.StatusConflict, gin.H{"message": "Project already awarded"})
		return
	}

	project, err := s.Store.AwardProject(ctx, req)

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

func (s *Server) ListAwardedProjects(ctx *gin.Context) {
	userID := ctx.Param("user_id")

	authPayload := ctx.MustGet(authorizationPayloadKey).(*auth.Token)
	if authPayload.UID != userID {
		ctx.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}

	projects, err := s.Store.ListAwardedProjects(ctx, userID)
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

func (s *Server) AcceptProject(ctx *gin.Context) {
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

	project, err := s.Store.AcceptProject(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, project)
}

func (s *Server) RejectProject(ctx *gin.Context) {
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

	project, err := s.Store.RejectProject(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, project)
}

func (s *Server) ListActiveProjects(ctx *gin.Context) {
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
		projects, err := s.Store.ListReferrerActiveProjects(ctx, userID)
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
		projects, err := s.Store.ListReferredActiveProjects(ctx, userID)
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

func (s *Server) StartProject(ctx *gin.Context) {
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

	project, err := s.Store.StartProject(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, project)
}

func (s *Server) InitiateCompleteProject(ctx *gin.Context) {
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

	project, err := s.Store.InitiateCompleteProject(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, project)
}

func (s *Server) CancelInitiateCompleteProject(ctx *gin.Context) {
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

	project, err := s.Store.CancelCompleteProjectInitiation(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, project)
}

func (s *Server) CompleteProject(ctx *gin.Context) {
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

	project, err := s.Store.CompleteProject(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, project)
}

func (s *Server) ListCompletedProjects(ctx *gin.Context) {

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
		projects, err := s.Store.ListReferrerCompletedProjects(ctx, userID)
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
		projects, err := s.Store.ListReferredCompletedProjects(ctx, userID)
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
