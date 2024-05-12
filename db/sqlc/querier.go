// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0

package db

import (
	"context"
)

type Querier interface {
	AcceptConnection(ctx context.Context, id int32) (ConnectionInvitation, error)
	AddConnection(ctx context.Context, arg AddConnectionParams) error
	AddEducation(ctx context.Context, arg AddEducationParams) (Education, error)
	AddExperience(ctx context.Context, arg AddExperienceParams) (Experience, error)
	AddFirm(ctx context.Context, arg AddFirmParams) (Firm, error)
	AddPrice(ctx context.Context, arg AddPriceParams) (Pricing, error)
	AddReview(ctx context.Context, arg AddReviewParams) (Review, error)
	AddSocial(ctx context.Context, arg AddSocialParams) (Social, error)
	CancelRecommendation(ctx context.Context, arg CancelRecommendationParams) error
	CreateChatRoom(ctx context.Context, arg CreateChatRoomParams) (CreateChatRoomRow, error)
	CreateMessage(ctx context.Context, arg CreateMessageParams) (Message, error)
	CreateUser(ctx context.Context, arg CreateUserParams) (User, error)
	DeleteEducation(ctx context.Context, educationID int64) error
	DeleteExperience(ctx context.Context, experienceID int64) error
	DeleteSocial(ctx context.Context, socialID int64) error
	FetchUserProfile(ctx context.Context, userID string) (FetchUserProfileRow, error)
	GetChatRoom(ctx context.Context, roomID string) (GetChatRoomRow, error)
	GetFirm(ctx context.Context, firmID int64) (Firm, error)
	GetUserByEmail(ctx context.Context, email string) (User, error)
	GetUserById(ctx context.Context, userID string) (User, error)
	GetUserWizardStep(ctx context.Context, userID string) (int32, error)
	ListChatRooms(ctx context.Context, user1ID string) ([]ListChatRoomsRow, error)
	ListConnectionInvitations(ctx context.Context, arg ListConnectionInvitationsParams) ([]ListConnectionInvitationsRow, error)
	ListConnections(ctx context.Context, arg ListConnectionsParams) ([]ListConnectionsRow, error)
	ListEducations(ctx context.Context, userID string) ([]Education, error)
	// -- name: ListExperiences :many
	// SELECT sqlc.embed(experiences), sqlc.embed(firms)
	// FROM experiences
	// JOIN firms ON experiences.firm_id = firms.firm_id
	// WHERE user_id = $1;
	ListExperiences(ctx context.Context, userID string) ([]ListExperiencesRow, error)
	ListFirms(ctx context.Context, arg ListFirmsParams) ([]Firm, error)
	ListMessages(ctx context.Context, arg ListMessagesParams) ([]ListMessagesRow, error)
	ListRecommendations(ctx context.Context, arg ListRecommendationsParams) ([]ListRecommendationsRow, error)
	ListRecommendations2(ctx context.Context, arg ListRecommendations2Params) ([]ListRecommendations2Row, error)
	ListSocials(ctx context.Context, arg ListSocialsParams) ([]Social, error)
	MarkWizardCompleted(ctx context.Context, arg MarkWizardCompletedParams) (User, error)
	RejectConnection(ctx context.Context, arg RejectConnectionParams) error
	SaveAboutYou(ctx context.Context, arg SaveAboutYouParams) (User, error)
	SaveLicense(ctx context.Context, arg SaveLicenseParams) (License, error)
	Search1stDegreeConnections(ctx context.Context, arg Search1stDegreeConnectionsParams) ([]Search1stDegreeConnectionsRow, error)
	// Exclude the current user
	// Retrieve user information for the second-degree connections
	Search2ndDegreeConnections(ctx context.Context, arg Search2ndDegreeConnectionsParams) ([]Search2ndDegreeConnectionsRow, error)
	SearchAllUsers(ctx context.Context, query string) ([]SearchAllUsersRow, error)
	SendConnection(ctx context.Context, arg SendConnectionParams) (int32, error)
	ToggleOpenToRefferal(ctx context.Context, arg ToggleOpenToRefferalParams) error
	UpdateEducation(ctx context.Context, arg UpdateEducationParams) (Education, error)
	UpdateEmailVerificationStatus(ctx context.Context, arg UpdateEmailVerificationStatusParams) (User, error)
	//- name: ListExperiences3 :many
	// SELECT
	//     experiences.experience_id,
	//     experiences.user_id,
	//     experiences.title,
	//     experiences.practice_area,
	//     experiences.firm_id,
	//     experiences.practice_location,
	//     experiences.start_date,
	//     experiences.end_date,
	//     experiences.current,
	//     experiences.description,
	//     experiences.skills,
	//     firms.firm_id,
	//     firms.name,
	//     firms.logo_url,
	//     firms.org_type,
	//     firms.website,
	//     firms.location
	// FROM experiences
	// JOIN firms ON experiences.firm_id = firms.firm_id
	// WHERE user_id = $1;
	UpdateExperience(ctx context.Context, arg UpdateExperienceParams) (Experience, error)
	UpdateMobileVerificationStatus(ctx context.Context, arg UpdateMobileVerificationStatusParams) (User, error)
	UpdatePrice(ctx context.Context, arg UpdatePriceParams) (Pricing, error)
	UpdateSocial(ctx context.Context, arg UpdateSocialParams) (Social, error)
	UpdateUserAvatar(ctx context.Context, arg UpdateUserAvatarParams) error
	UpdateUserAvatarUrl(ctx context.Context, arg UpdateUserAvatarUrlParams) (User, error)
	UpdateUserBannerImage(ctx context.Context, arg UpdateUserBannerImageParams) error
	UpdateUserInfo(ctx context.Context, arg UpdateUserInfoParams) (User, error)
	UpdateUserWizardStep(ctx context.Context, arg UpdateUserWizardStepParams) (User, error)
	UploadLicense(ctx context.Context, arg UploadLicenseParams) (License, error)
}

var _ Querier = (*Queries)(nil)
