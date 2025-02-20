// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0

package db

import (
	"database/sql/driver"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

type AdType string

const (
	AdTypeImage AdType = "image"
	AdTypeVideo AdType = "video"
	AdTypeOther AdType = "other"
)

func (e *AdType) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = AdType(s)
	case string:
		*e = AdType(s)
	default:
		return fmt.Errorf("unsupported scan type for AdType: %T", src)
	}
	return nil
}

type NullAdType struct {
	AdType AdType `json:"ad_type"`
	Valid  bool   `json:"valid"` // Valid is true if AdType is not NULL
}

// Scan implements the Scanner interface.
func (ns *NullAdType) Scan(value interface{}) error {
	if value == nil {
		ns.AdType, ns.Valid = "", false
		return nil
	}
	ns.Valid = true
	return ns.AdType.Scan(value)
}

// Value implements the driver Valuer interface.
func (ns NullAdType) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return string(ns.AdType), nil
}

type DiscussionInviteStatus string

const (
	DiscussionInviteStatusPending  DiscussionInviteStatus = "pending"
	DiscussionInviteStatusAccepted DiscussionInviteStatus = "accepted"
	DiscussionInviteStatusRejected DiscussionInviteStatus = "rejected"
)

func (e *DiscussionInviteStatus) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = DiscussionInviteStatus(s)
	case string:
		*e = DiscussionInviteStatus(s)
	default:
		return fmt.Errorf("unsupported scan type for DiscussionInviteStatus: %T", src)
	}
	return nil
}

type NullDiscussionInviteStatus struct {
	DiscussionInviteStatus DiscussionInviteStatus `json:"discussion_invite_status"`
	Valid                  bool                   `json:"valid"` // Valid is true if DiscussionInviteStatus is not NULL
}

// Scan implements the Scanner interface.
func (ns *NullDiscussionInviteStatus) Scan(value interface{}) error {
	if value == nil {
		ns.DiscussionInviteStatus, ns.Valid = "", false
		return nil
	}
	ns.Valid = true
	return ns.DiscussionInviteStatus.Scan(value)
}

// Value implements the driver Valuer interface.
func (ns NullDiscussionInviteStatus) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return string(ns.DiscussionInviteStatus), nil
}

type InvitationStatus string

const (
	InvitationStatusPending   InvitationStatus = "pending"
	InvitationStatusAccepted  InvitationStatus = "accepted"
	InvitationStatusRejected  InvitationStatus = "rejected"
	InvitationStatusCancelled InvitationStatus = "cancelled"
	InvitationStatusNone      InvitationStatus = "none"
)

func (e *InvitationStatus) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = InvitationStatus(s)
	case string:
		*e = InvitationStatus(s)
	default:
		return fmt.Errorf("unsupported scan type for InvitationStatus: %T", src)
	}
	return nil
}

type NullInvitationStatus struct {
	InvitationStatus InvitationStatus `json:"invitation_status"`
	Valid            bool             `json:"valid"` // Valid is true if InvitationStatus is not NULL
}

// Scan implements the Scanner interface.
func (ns *NullInvitationStatus) Scan(value interface{}) error {
	if value == nil {
		ns.InvitationStatus, ns.Valid = "", false
		return nil
	}
	ns.Valid = true
	return ns.InvitationStatus.Scan(value)
}

// Value implements the driver Valuer interface.
func (ns NullInvitationStatus) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return string(ns.InvitationStatus), nil
}

type PaymentCycle string

const (
	PaymentCycleWeekly  PaymentCycle = "weekly"
	PaymentCycleMonthly PaymentCycle = "monthly"
	PaymentCycleYearly  PaymentCycle = "yearly"
)

func (e *PaymentCycle) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = PaymentCycle(s)
	case string:
		*e = PaymentCycle(s)
	default:
		return fmt.Errorf("unsupported scan type for PaymentCycle: %T", src)
	}
	return nil
}

type NullPaymentCycle struct {
	PaymentCycle PaymentCycle `json:"payment_cycle"`
	Valid        bool         `json:"valid"` // Valid is true if PaymentCycle is not NULL
}

// Scan implements the Scanner interface.
func (ns *NullPaymentCycle) Scan(value interface{}) error {
	if value == nil {
		ns.PaymentCycle, ns.Valid = "", false
		return nil
	}
	ns.Valid = true
	return ns.PaymentCycle.Scan(value)
}

// Value implements the driver Valuer interface.
func (ns NullPaymentCycle) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return string(ns.PaymentCycle), nil
}

type PostType string

const (
	PostTypeText     PostType = "text"
	PostTypeImage    PostType = "image"
	PostTypeVideo    PostType = "video"
	PostTypeAudio    PostType = "audio"
	PostTypeLink     PostType = "link"
	PostTypeDocument PostType = "document"
	PostTypePoll     PostType = "poll"
	PostTypeOther    PostType = "other"
)

func (e *PostType) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = PostType(s)
	case string:
		*e = PostType(s)
	default:
		return fmt.Errorf("unsupported scan type for PostType: %T", src)
	}
	return nil
}

type NullPostType struct {
	PostType PostType `json:"post_type"`
	Valid    bool     `json:"valid"` // Valid is true if PostType is not NULL
}

// Scan implements the Scanner interface.
func (ns *NullPostType) Scan(value interface{}) error {
	if value == nil {
		ns.PostType, ns.Valid = "", false
		return nil
	}
	ns.Valid = true
	return ns.PostType.Scan(value)
}

// Value implements the driver Valuer interface.
func (ns NullPostType) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return string(ns.PostType), nil
}

type ProjectStatus string

const (
	ProjectStatusActive            ProjectStatus = "active"
	ProjectStatusAwarded           ProjectStatus = "awarded"
	ProjectStatusAccepted          ProjectStatus = "accepted"
	ProjectStatusRejected          ProjectStatus = "rejected"
	ProjectStatusStarted           ProjectStatus = "started"
	ProjectStatusCompleteInitiated ProjectStatus = "complete_initiated"
	ProjectStatusCompleted         ProjectStatus = "completed"
	ProjectStatusCancelled         ProjectStatus = "cancelled"
)

func (e *ProjectStatus) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = ProjectStatus(s)
	case string:
		*e = ProjectStatus(s)
	default:
		return fmt.Errorf("unsupported scan type for ProjectStatus: %T", src)
	}
	return nil
}

type NullProjectStatus struct {
	ProjectStatus ProjectStatus `json:"project_status"`
	Valid         bool          `json:"valid"` // Valid is true if ProjectStatus is not NULL
}

// Scan implements the Scanner interface.
func (ns *NullProjectStatus) Scan(value interface{}) error {
	if value == nil {
		ns.ProjectStatus, ns.Valid = "", false
		return nil
	}
	ns.Valid = true
	return ns.ProjectStatus.Scan(value)
}

// Value implements the driver Valuer interface.
func (ns NullProjectStatus) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return string(ns.ProjectStatus), nil
}

type ProposalStatus string

const (
	ProposalStatusActive    ProposalStatus = "active"
	ProposalStatusAccepted  ProposalStatus = "accepted"
	ProposalStatusRejected  ProposalStatus = "rejected"
	ProposalStatusCancelled ProposalStatus = "cancelled"
)

func (e *ProposalStatus) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = ProposalStatus(s)
	case string:
		*e = ProposalStatus(s)
	default:
		return fmt.Errorf("unsupported scan type for ProposalStatus: %T", src)
	}
	return nil
}

type NullProposalStatus struct {
	ProposalStatus ProposalStatus `json:"proposal_status"`
	Valid          bool           `json:"valid"` // Valid is true if ProposalStatus is not NULL
}

// Scan implements the Scanner interface.
func (ns *NullProposalStatus) Scan(value interface{}) error {
	if value == nil {
		ns.ProposalStatus, ns.Valid = "", false
		return nil
	}
	ns.Valid = true
	return ns.ProposalStatus.Scan(value)
}

// Value implements the driver Valuer interface.
func (ns NullProposalStatus) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return string(ns.ProposalStatus), nil
}

type Ad struct {
	AdID         int32        `json:"ad_id"`
	AdType       AdType       `json:"ad_type"`
	Title        string       `json:"title"`
	Description  string       `json:"description"`
	Link         string       `json:"link"`
	Media        []string     `json:"media"`
	PaymentCycle PaymentCycle `json:"payment_cycle"`
	AuthorID     string       `json:"author_id"`
	StartDate    time.Time    `json:"start_date"`
	EndDate      time.Time    `json:"end_date"`
	CreatedAt    time.Time    `json:"created_at"`
}

type Attachment struct {
	AttachmentID   int32  `json:"attachment_id"`
	MessageID      int32  `json:"message_id"`
	AttachmentUrl  string `json:"attachment_url"`
	AttachmentType string `json:"attachment_type"`
}

type CanceledRecommendation struct {
	ID                int32     `json:"id"`
	UserID            string    `json:"user_id"`
	RecommendedUserID string    `json:"recommended_user_id"`
	CanceledAt        time.Time `json:"canceled_at"`
}

type ChatRoom struct {
	RoomID        string             `json:"room_id"`
	User1ID       string             `json:"user1_id"`
	User2ID       string             `json:"user2_id"`
	LastMessageAt pgtype.Timestamptz `json:"last_message_at"`
	CreatedAt     time.Time          `json:"created_at"`
}

type Comment struct {
	CommentID       int32     `json:"comment_id"`
	UserID          string    `json:"user_id"`
	PostID          int32     `json:"post_id"`
	Content         string    `json:"content"`
	CreatedAt       time.Time `json:"created_at"`
	ParentCommentID *int32    `json:"parent_comment_id"`
}

type Connection struct {
	ID          int32     `json:"id"`
	SenderID    string    `json:"sender_id"`
	RecipientID string    `json:"recipient_id"`
	CreatedAt   time.Time `json:"created_at"`
}

type ConnectionInvitation struct {
	ID          int32            `json:"id"`
	SenderID    string           `json:"sender_id"`
	RecipientID string           `json:"recipient_id"`
	Status      InvitationStatus `json:"status"`
	CreatedAt   time.Time        `json:"created_at"`
}

type Device struct {
	DeviceID    string    `json:"device_id"`
	DeviceToken string    `json:"device_token"`
	UserID      string    `json:"user_id"`
	LastUsedAt  time.Time `json:"last_used_at"`
}

type Discussion struct {
	DiscussionID int32     `json:"discussion_id"`
	AuthorID     string    `json:"author_id"`
	Topic        string    `json:"topic"`
	CreatedAt    time.Time `json:"created_at"`
}

type DiscussionInvite struct {
	DiscussionInviteID int32                  `json:"discussion_invite_id"`
	DiscussionID       int32                  `json:"discussion_id"`
	InviteeUserID      string                 `json:"invitee_user_id"`
	InvitedUserID      string                 `json:"invited_user_id"`
	Status             DiscussionInviteStatus `json:"status"`
	CreatedAt          time.Time              `json:"created_at"`
}

type DiscussionMessage struct {
	MessageID       int32     `json:"message_id"`
	ParentMessageID *int32    `json:"parent_message_id"`
	DiscussionID    int32     `json:"discussion_id"`
	SenderID        string    `json:"sender_id"`
	Message         string    `json:"message"`
	SentAt          time.Time `json:"sent_at"`
}

type Education struct {
	EducationID  int64       `json:"education_id"`
	UserID       string      `json:"user_id"`
	School       string      `json:"school"`
	Degree       string      `json:"degree"`
	FieldOfStudy string      `json:"field_of_study"`
	StartDate    pgtype.Date `json:"start_date"`
	EndDate      pgtype.Date `json:"end_date"`
	Current      bool        `json:"current"`
	Grade        string      `json:"grade"`
	Achievements string      `json:"achievements"`
	Skills       []string    `json:"skills"`
}

type Experience struct {
	ExperienceID     int64       `json:"experience_id"`
	UserID           string      `json:"user_id"`
	Title            string      `json:"title"`
	PracticeArea     string      `json:"practice_area"`
	FirmID           int64       `json:"firm_id"`
	PracticeLocation string      `json:"practice_location"`
	StartDate        pgtype.Date `json:"start_date"`
	EndDate          pgtype.Date `json:"end_date"`
	Current          bool        `json:"current"`
	Description      string      `json:"description"`
	Skills           []string    `json:"skills"`
}

type Faq struct {
	FaqID     int32     `json:"faq_id"`
	Question  string    `json:"question"`
	Answer    string    `json:"answer"`
	CreatedAt time.Time `json:"created_at"`
}

type FeaturedPost struct {
	FeaturePostID int32     `json:"feature_post_id"`
	PostID        int32     `json:"post_id"`
	UserID        string    `json:"user_id"`
	CreatedAt     time.Time `json:"created_at"`
}

type Firm struct {
	FirmID      int64     `json:"firm_id"`
	OwnerUserID string    `json:"owner_user_id"`
	Name        string    `json:"name"`
	LogoUrl     string    `json:"logo_url"`
	OrgType     string    `json:"org_type"`
	Website     string    `json:"website"`
	Location    string    `json:"location"`
	About       string    `json:"about"`
	CreatedAt   time.Time `json:"created_at"`
}

type License struct {
	LicenseID     int64       `json:"license_id"`
	UserID        string      `json:"user_id"`
	Name          string      `json:"name"`
	LicenseNumber string      `json:"license_number"`
	IssueDate     pgtype.Date `json:"issue_date"`
	IssueState    string      `json:"issue_state"`
	LicenseUrl    *string     `json:"license_url"`
}

type Like struct {
	LikeID    int32     `json:"like_id"`
	UserID    string    `json:"user_id"`
	PostID    *int32    `json:"post_id"`
	CommentID *int32    `json:"comment_id"`
	Type      string    `json:"type"`
	CreatedAt time.Time `json:"created_at"`
}

type Message struct {
	MessageID       int32     `json:"message_id"`
	ParentMessageID *int32    `json:"parent_message_id"`
	SenderID        string    `json:"sender_id"`
	RecipientID     string    `json:"recipient_id"`
	Message         string    `json:"message"`
	HasAttachment   bool      `json:"has_attachment"`
	AttachmentID    *int32    `json:"attachment_id"`
	IsRead          bool      `json:"is_read"`
	RoomID          string    `json:"room_id"`
	SentAt          time.Time `json:"sent_at"`
}

type NewsFeed struct {
	FeedID    int32     `json:"feed_id"`
	UserID    string    `json:"user_id"`
	PostID    int32     `json:"post_id"`
	CreatedAt time.Time `json:"created_at"`
}

type Notification struct {
	NotificationID   int32     `json:"notification_id"`
	UserID           string    `json:"user_id"`
	SenderID         string    `json:"sender_id"`
	TargetID         int32     `json:"target_id"`
	TargetType       string    `json:"target_type"`
	NotificationType string    `json:"notification_type"`
	Message          string    `json:"message"`
	IsRead           bool      `json:"is_read"`
	CreatedAt        time.Time `json:"created_at"`
}

type Poll struct {
	PollID    int32              `json:"poll_id"`
	OwnerID   string             `json:"owner_id"`
	Title     string             `json:"title"`
	Options   []string           `json:"options"`
	CreatedAt time.Time          `json:"created_at"`
	EndTime   pgtype.Timestamptz `json:"end_time"`
}

type PollResult struct {
	PollResultID int32     `json:"poll_result_id"`
	PollID       int32     `json:"poll_id"`
	OptionIndex  int32     `json:"option_index"`
	UserID       string    `json:"user_id"`
	CreatedAt    time.Time `json:"created_at"`
}

type Post struct {
	PostID    int32     `json:"post_id"`
	OwnerID   string    `json:"owner_id"`
	Content   *string   `json:"content"`
	Media     []string  `json:"media"`
	PostType  PostType  `json:"post_type"`
	PollID    *int32    `json:"poll_id"`
	CreatedAt time.Time `json:"created_at"`
}

type PostStatistic struct {
	PostID    int32     `json:"post_id"`
	Views     int64     `json:"views"`
	Likes     int64     `json:"likes"`
	Comments  int64     `json:"comments"`
	Shares    int64     `json:"shares"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Pricing struct {
	PriceID          int64          `json:"price_id"`
	UserID           string         `json:"user_id"`
	ServiceType      string         `json:"service_type"`
	PerHourPrice     pgtype.Numeric `json:"per_hour_price"`
	PerHearingPrice  pgtype.Numeric `json:"per_hearing_price"`
	ContingencyPrice *string        `json:"contingency_price"`
	HybridPrice      *string        `json:"hybrid_price"`
}

type Project struct {
	ProjectID                 int32              `json:"project_id"`
	Title                     string             `json:"title"`
	PreferredPracticeArea     string             `json:"preferred_practice_area"`
	PreferredPracticeLocation string             `json:"preferred_practice_location"`
	CaseDescription           string             `json:"case_description"`
	ReferrerUserID            string             `json:"referrer_user_id"`
	ReferredUserID            *string            `json:"referred_user_id"`
	Status                    ProjectStatus      `json:"status"`
	CreatedAt                 time.Time          `json:"created_at"`
	StartedAt                 pgtype.Timestamptz `json:"started_at"`
	CompletedAt               pgtype.Timestamptz `json:"completed_at"`
}

type ProjectReview struct {
	ReviewID  int32          `json:"review_id"`
	ProjectID int32          `json:"project_id"`
	UserID    string         `json:"user_id"`
	Review    string         `json:"review"`
	Rating    pgtype.Numeric `json:"rating"`
	CreatedAt time.Time      `json:"created_at"`
}

type Proposal struct {
	ProposalID int32          `json:"proposal_id"`
	ProjectID  int32          `json:"project_id"`
	UserID     string         `json:"user_id"`
	Title      string         `json:"title"`
	Proposal   string         `json:"proposal"`
	Status     ProposalStatus `json:"status"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
}

type ReferralUser struct {
	ReferralUserID int32  `json:"referral_user_id"`
	ProjectID      int32  `json:"project_id"`
	ReferredUserID string `json:"referred_user_id"`
}

type ReportReason struct {
	ReasonID int32  `json:"reason_id"`
	Reason   string `json:"reason"`
}

type ReportedPost struct {
	ReportID   int32     `json:"report_id"`
	PostID     int32     `json:"post_id"`
	ReportedBy string    `json:"reported_by"`
	ReasonID   int32     `json:"reason_id"`
	CreatedAt  time.Time `json:"created_at"`
}

type Review struct {
	ReviewID   int64     `json:"review_id"`
	UserID     string    `json:"user_id"`
	ReviewerID string    `json:"reviewer_id"`
	Review     string    `json:"review"`
	Rating     float64   `json:"rating"`
	Timestamp  time.Time `json:"timestamp"`
}

type SavedPost struct {
	SavedPostID int32     `json:"saved_post_id"`
	PostID      int32     `json:"post_id"`
	UserID      string    `json:"user_id"`
	CreatedAt   time.Time `json:"created_at"`
}

type Social struct {
	SocialID   int64  `json:"social_id"`
	EntityID   string `json:"entity_id"`
	EntityType string `json:"entity_type"`
	Platform   string `json:"platform"`
	Link       string `json:"link"`
}

type User struct {
	UserID                  string    `json:"user_id"`
	Email                   string    `json:"email"`
	FirstName               string    `json:"first_name"`
	LastName                string    `json:"last_name"`
	About                   *string   `json:"about"`
	Mobile                  *string   `json:"mobile"`
	Address                 *string   `json:"address"`
	AvatarUrl               *string   `json:"avatar_url"`
	BannerUrl               *string   `json:"banner_url"`
	EmailVerified           bool      `json:"email_verified"`
	MobileVerified          bool      `json:"mobile_verified"`
	WizardStep              int32     `json:"wizard_step"`
	WizardCompleted         bool      `json:"wizard_completed"`
	SignupMethod            int32     `json:"signup_method"`
	PracticeArea            *string   `json:"practice_area"`
	PracticeLocation        *string   `json:"practice_location"`
	Experience              *string   `json:"experience"`
	AverageBillingPerClient *int32    `json:"average_billing_per_client"`
	CaseResolutionRate      *int32    `json:"case_resolution_rate"`
	OpenToReferral          bool      `json:"open_to_referral"`
	LicenseVerified         bool      `json:"license_verified"`
	LicenseRejected         bool      `json:"license_rejected"`
	JoinDate                time.Time `json:"join_date"`
}
