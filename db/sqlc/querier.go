// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0

package db

import (
	"context"
)

type Querier interface {
	AddEducation(ctx context.Context, arg AddEducationParams) (Education, error)
	AddExperience(ctx context.Context, arg AddExperienceParams) (Experience, error)
	AddFirm(ctx context.Context, arg AddFirmParams) (Firm, error)
	AddPrice(ctx context.Context, arg AddPriceParams) (Pricing, error)
	AddReview(ctx context.Context, arg AddReviewParams) (Review, error)
	AddSocial(ctx context.Context, arg AddSocialParams) (Social, error)
	CreateUser(ctx context.Context, arg CreateUserParams) (User, error)
	DeleteEducation(ctx context.Context, educationID int64) error
	DeleteExperience(ctx context.Context, experienceID int64) error
	FetchUserProfile(ctx context.Context, userID string) (FetchUserProfileRow, error)
	GetFirm(ctx context.Context, firmID int64) (Firm, error)
	GetUserByEmail(ctx context.Context, email string) (User, error)
	GetUserById(ctx context.Context, userID string) (User, error)
	GetUserWizardStep(ctx context.Context, userID string) (int32, error)
	ListEducations(ctx context.Context, userID string) ([]Education, error)
	ListExperiences(ctx context.Context, userID string) ([]ListExperiencesRow, error)
	ListFirms(ctx context.Context, arg ListFirmsParams) ([]Firm, error)
	ListSocials(ctx context.Context, arg ListSocialsParams) ([]Social, error)
	MarkWizardCompleted(ctx context.Context, arg MarkWizardCompletedParams) (User, error)
	SaveAboutYou(ctx context.Context, arg SaveAboutYouParams) (User, error)
	SaveLicense(ctx context.Context, arg SaveLicenseParams) (License, error)
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
