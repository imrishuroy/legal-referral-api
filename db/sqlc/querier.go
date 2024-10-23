// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0

package db

import (
	"context"
)

type Querier interface {
	AcceptConnection(ctx context.Context, id int32) (ConnectionInvitation, error)
	AcceptProject(ctx context.Context, arg AcceptProjectParams) (Project, error)
	AddConnection(ctx context.Context, arg AddConnectionParams) error
	AddEducation(ctx context.Context, arg AddEducationParams) (Education, error)
	AddExperience(ctx context.Context, arg AddExperienceParams) (Experience, error)
	AddFirm(ctx context.Context, arg AddFirmParams) (Firm, error)
	AddPrice(ctx context.Context, arg AddPriceParams) (Pricing, error)
	AddReferredUserToProject(ctx context.Context, arg AddReferredUserToProjectParams) (Project, error)
	AddReview(ctx context.Context, arg AddReviewParams) (Review, error)
	AddSocial(ctx context.Context, arg AddSocialParams) (Social, error)
	ApproveLicense(ctx context.Context, userID string) error
	AwardProject(ctx context.Context, arg AwardProjectParams) (Project, error)
	CancelCompleteProjectInitiation(ctx context.Context, arg CancelCompleteProjectInitiationParams) (Project, error)
	CancelRecommendation(ctx context.Context, arg CancelRecommendationParams) error
	CheckConnection(ctx context.Context, arg CheckConnectionParams) (bool, error)
	CheckConnectionStatus(ctx context.Context, arg CheckConnectionStatusParams) (interface{}, error)
	CheckPostLike(ctx context.Context, arg CheckPostLikeParams) (bool, error)
	CommentPost(ctx context.Context, arg CommentPostParams) (CommentPostRow, error)
	CompleteProject(ctx context.Context, arg CompleteProjectParams) (Project, error)
	CreateAd(ctx context.Context, arg CreateAdParams) (Ad, error)
	CreateChatRoom(ctx context.Context, arg CreateChatRoomParams) (CreateChatRoomRow, error)
	CreateDiscussion(ctx context.Context, arg CreateDiscussionParams) (Discussion, error)
	CreateFAQ(ctx context.Context, arg CreateFAQParams) (Faq, error)
	CreateMessage(ctx context.Context, arg CreateMessageParams) (Message, error)
	CreateNotification(ctx context.Context, arg CreateNotificationParams) (Notification, error)
	CreatePoll(ctx context.Context, arg CreatePollParams) (Poll, error)
	CreatePost(ctx context.Context, arg CreatePostParams) (Post, error)
	CreateProjectReview(ctx context.Context, arg CreateProjectReviewParams) (ProjectReview, error)
	CreateProposal(ctx context.Context, arg CreateProposalParams) (Proposal, error)
	CreateReferral(ctx context.Context, arg CreateReferralParams) (Project, error)
	CreateUser(ctx context.Context, arg CreateUserParams) (User, error)
	DeleteEducation(ctx context.Context, educationID int64) error
	DeleteExperience(ctx context.Context, experienceID int64) error
	DeleteNotificationById(ctx context.Context, notificationID int32) (Notification, error)
	DeletePost(ctx context.Context, arg DeletePostParams) error
	DeleteSocial(ctx context.Context, socialID int64) error
	ExtendAdPeriod(ctx context.Context, arg ExtendAdPeriodParams) (Ad, error)
	FetchUserProfile(ctx context.Context, userID string) (FetchUserProfileRow, error)
	GetChatRoom(ctx context.Context, arg GetChatRoomParams) (GetChatRoomRow, error)
	GetFirm(ctx context.Context, firmID int64) (Firm, error)
	GetNotificationById(ctx context.Context, notificationID int32) (Notification, error)
	GetPosIsLikedByCurrentUser(ctx context.Context, arg GetPosIsLikedByCurrentUserParams) (bool, error)
	GetPost(ctx context.Context, postID int32) (Post, error)
	GetPostCommentsCount(ctx context.Context, postID int32) (int64, error)
	GetPostLikesAndCommentsCount(ctx context.Context, postID int32) (GetPostLikesAndCommentsCountRow, error)
	GetPostLikesCount(ctx context.Context, postID *int32) (int64, error)
	GetProjectReview(ctx context.Context, arg GetProjectReviewParams) (ProjectReview, error)
	GetProjectStatus(ctx context.Context, projectID int32) (ProjectStatus, error)
	//
	GetProposal(ctx context.Context, arg GetProposalParams) (Proposal, error)
	GetRandomAd(ctx context.Context) (Ad, error)
	GetUserByEmail(ctx context.Context, email string) (User, error)
	GetUserById(ctx context.Context, userID string) (User, error)
	GetUserWizardStep(ctx context.Context, userID string) (int32, error)
	InitiateCompleteProject(ctx context.Context, arg InitiateCompleteProjectParams) (Project, error)
	InviteUserToDiscussion(ctx context.Context, arg InviteUserToDiscussionParams) error
	JoinDiscussion(ctx context.Context, arg JoinDiscussionParams) error
	LikeComment(ctx context.Context, arg LikeCommentParams) error
	LikePost(ctx context.Context, arg LikePostParams) error
	ListActiveDiscussions(ctx context.Context, authorID string) ([]ListActiveDiscussionsRow, error)
	ListActiveProposals(ctx context.Context, userID string) ([]ListActiveProposalsRow, error)
	ListActiveReferralProjects(ctx context.Context, referrerUserID string) ([]Project, error)
	ListActiveReferrals(ctx context.Context, userID string) ([]Project, error)
	ListAllReferralProjects(ctx context.Context, referrerUserID string) ([]Project, error)
	ListAttorneys(ctx context.Context, arg ListAttorneysParams) ([]ListAttorneysRow, error)
	ListAwardedProjects(ctx context.Context, userID string) ([]ListAwardedProjectsRow, error)
	ListChatRooms(ctx context.Context, user1ID string) ([]ListChatRoomsRow, error)
	ListComments(ctx context.Context, postID int32) ([]ListCommentsRow, error)
	ListComments2(ctx context.Context, arg ListComments2Params) ([]ListComments2Row, error)
	ListCompletedReferralProjects(ctx context.Context, referrerUserID string) ([]Project, error)
	ListConnectedUserIDs(ctx context.Context, userID string) ([]interface{}, error)
	ListConnectedUsers(ctx context.Context, arg ListConnectedUsersParams) ([]ListConnectedUsersRow, error)
	ListConnectionInvitations(ctx context.Context, arg ListConnectionInvitationsParams) ([]ListConnectionInvitationsRow, error)
	ListConnections(ctx context.Context, arg ListConnectionsParams) ([]ListConnectionsRow, error)
	ListDiscussionInvites(ctx context.Context, invitedUserID string) ([]ListDiscussionInvitesRow, error)
	ListDiscussionMessages(ctx context.Context, arg ListDiscussionMessagesParams) ([]ListDiscussionMessagesRow, error)
	ListDiscussionParticipants(ctx context.Context, discussionID int32) ([]ListDiscussionParticipantsRow, error)
	ListEducations(ctx context.Context, userID string) ([]Education, error)
	ListExperiences(ctx context.Context, userID string) ([]ListExperiencesRow, error)
	ListExpiredAds(ctx context.Context) ([]Ad, error)
	ListFAQs(ctx context.Context) ([]Faq, error)
	ListFeaturePosts(ctx context.Context) ([]ListFeaturePostsRow, error)
	ListFirms(ctx context.Context, arg ListFirmsParams) ([]Firm, error)
	ListFirmsByOwner(ctx context.Context, ownerUserID string) ([]Firm, error)
	// lawyers
	ListLawyers(ctx context.Context) ([]ListLawyersRow, error)
	//     AND u.license_rejected = false
	ListLicenseUnVerifiedUsers(ctx context.Context, arg ListLicenseUnVerifiedUsersParams) ([]ListLicenseUnVerifiedUsersRow, error)
	ListLicenseVerifiedUsers(ctx context.Context, arg ListLicenseVerifiedUsersParams) ([]ListLicenseVerifiedUsersRow, error)
	ListMessages(ctx context.Context, arg ListMessagesParams) ([]ListMessagesRow, error)
	ListNewsFeed(ctx context.Context, arg ListNewsFeedParams) ([]ListNewsFeedRow, error)
	ListNotifications(ctx context.Context, arg ListNotificationsParams) ([]ListNotificationsRow, error)
	ListPlayingAds(ctx context.Context) ([]Ad, error)
	ListPostLikedUsers(ctx context.Context, postID *int32) ([]ListPostLikedUsersRow, error)
	ListPostLikedUsers2(ctx context.Context, arg ListPostLikedUsers2Params) ([]ListPostLikedUsers2Row, error)
	ListPostLikes(ctx context.Context, postID *int32) ([]string, error)
	ListRandomAds(ctx context.Context, limit int32) ([]Ad, error)
	ListRecommendations(ctx context.Context, arg ListRecommendationsParams) ([]ListRecommendationsRow, error)
	ListRecommendations2(ctx context.Context, arg ListRecommendations2Params) ([]ListRecommendations2Row, error)
	//
	ListReferredActiveProjects(ctx context.Context, userID string) ([]ListReferredActiveProjectsRow, error)
	ListReferredCompletedProjects(ctx context.Context, userID string) ([]ListReferredCompletedProjectsRow, error)
	ListReferredUsers(ctx context.Context, projectID int32) ([]ListReferredUsersRow, error)
	ListReferredUsers2(ctx context.Context, projectID int32) ([]ListReferredUsers2Row, error)
	ListReferrerActiveProjects(ctx context.Context, userID string) ([]ListReferrerActiveProjectsRow, error)
	ListReferrerCompletedProjects(ctx context.Context, userID string) ([]ListReferrerCompletedProjectsRow, error)
	ListSavedPosts(ctx context.Context, userID string) ([]ListSavedPostsRow, error)
	ListSocials(ctx context.Context, arg ListSocialsParams) ([]Social, error)
	ListUninvitedParticipants(ctx context.Context, discussionID int32) ([]ListUninvitedParticipantsRow, error)
	ListUsers(ctx context.Context, arg ListUsersParams) ([]ListUsersRow, error)
	MarkNotificationAsRead(ctx context.Context, notificationID int32) (Notification, error)
	MarkWizardCompleted(ctx context.Context, arg MarkWizardCompletedParams) (User, error)
	PostToNewsFeed(ctx context.Context, arg PostToNewsFeedParams) error
	RejectConnection(ctx context.Context, id int32) error
	RejectDiscussion(ctx context.Context, arg RejectDiscussionParams) error
	RejectLicense(ctx context.Context, userID string) error
	//
	RejectProject(ctx context.Context, arg RejectProjectParams) (Project, error)
	SaveAboutYou(ctx context.Context, arg SaveAboutYouParams) (User, error)
	SaveDevice(ctx context.Context, arg SaveDeviceParams) error
	SaveFeaturePost(ctx context.Context, arg SaveFeaturePostParams) error
	SaveLicense(ctx context.Context, arg SaveLicenseParams) (License, error)
	SavePost(ctx context.Context, arg SavePostParams) error
	Search1stDegreeConnections(ctx context.Context, arg Search1stDegreeConnectionsParams) ([]Search1stDegreeConnectionsRow, error)
	// Exclude the current user
	// Retrieve user information for the second-degree connections
	Search2ndDegreeConnections(ctx context.Context, arg Search2ndDegreeConnectionsParams) ([]Search2ndDegreeConnectionsRow, error)
	SearchAllUsers(ctx context.Context, query string) ([]SearchAllUsersRow, error)
	SendConnection(ctx context.Context, arg SendConnectionParams) (int32, error)
	SendMessageToDiscussion(ctx context.Context, arg SendMessageToDiscussionParams) (DiscussionMessage, error)
	StartProject(ctx context.Context, arg StartProjectParams) (Project, error)
	ToggleOpenToRefferal(ctx context.Context, arg ToggleOpenToRefferalParams) error
	UnSaveFeaturePost(ctx context.Context, arg UnSaveFeaturePostParams) error
	UnlikeComment(ctx context.Context, arg UnlikeCommentParams) error
	UnlikePost(ctx context.Context, arg UnlikePostParams) error
	UnsavePost(ctx context.Context, savedPostID int32) error
	UpdateDiscussionTopic(ctx context.Context, arg UpdateDiscussionTopicParams) error
	UpdateEducation(ctx context.Context, arg UpdateEducationParams) (Education, error)
	UpdateEmailVerificationStatus(ctx context.Context, arg UpdateEmailVerificationStatusParams) (User, error)
	UpdateExperience(ctx context.Context, arg UpdateExperienceParams) (Experience, error)
	UpdateMobileVerificationStatus(ctx context.Context, arg UpdateMobileVerificationStatusParams) (User, error)
	UpdatePrice(ctx context.Context, arg UpdatePriceParams) (Pricing, error)
	//
	UpdateProposal(ctx context.Context, arg UpdateProposalParams) (Proposal, error)
	UpdateSocial(ctx context.Context, arg UpdateSocialParams) (Social, error)
	UpdateUserAvatar(ctx context.Context, arg UpdateUserAvatarParams) error
	UpdateUserAvatarUrl(ctx context.Context, arg UpdateUserAvatarUrlParams) (User, error)
	UpdateUserBannerImage(ctx context.Context, arg UpdateUserBannerImageParams) error
	UpdateUserInfo(ctx context.Context, arg UpdateUserInfoParams) (User, error)
	UpdateUserWizardStep(ctx context.Context, arg UpdateUserWizardStepParams) (User, error)
	UploadLicense(ctx context.Context, arg UploadLicenseParams) (License, error)
}

var _ Querier = (*Queries)(nil)
