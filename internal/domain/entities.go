package domain

import (
	"time"

	"github.com/google/uuid"
)

type AccessLevel string

const (
	AccessLevelFree    AccessLevel = "free"
	AccessLevelBasic   AccessLevel = "basic"
	AccessLevelPremium AccessLevel = "premium"
)

type SubscriptionStatus string

const (
	SubscriptionStatusActive    SubscriptionStatus = "active"
	SubscriptionStatusExpired   SubscriptionStatus = "expired"
	SubscriptionStatusCancelled SubscriptionStatus = "cancelled"
)

type WatchStatus string

const (
	WatchStatusStarted   WatchStatus = "started"
	WatchStatusPaused    WatchStatus = "paused"
	WatchStatusCompleted WatchStatus = "completed"
)

type User struct {
	ID           uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	Email        string    `gorm:"uniqueIndex;not null" json:"email"`
	PasswordHash string    `gorm:"not null" json:"-"`
	Name         string    `gorm:"not null" json:"name"`
	Bio          string    `json:"bio"`
	Picture      string    `json:"picture"`
	Phone        string    `json:"phone"`
	IsAdmin      bool      `gorm:"default:false" json:"is_admin"`
	CreatedAt    time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt    time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

func (User) TableName() string {
	return "users"
}

type Content struct {
	ID              uuid.UUID   `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	Title           string      `gorm:"not null;index" json:"title"`
	Description     string      `json:"description"`
	AccessLevel     AccessLevel `gorm:"type:varchar(20);not null;index" json:"access_level"`
	DurationSeconds int         `gorm:"not null" json:"duration_seconds"`
	ThumbnailURL    string      `json:"thumbnail_url"`
	VideoURL        string      `json:"video_url"`
	Published       bool        `gorm:"default:false;index" json:"published"`
	CreatedAt       time.Time   `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt       time.Time   `gorm:"autoUpdateTime" json:"updated_at"`
}

func (Content) TableName() string {
	return "contents"
}

type Plan struct {
	ID                uuid.UUID   `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	Name              string      `gorm:"not null;uniqueIndex" json:"name"`
	Price             int64       `gorm:"not null" json:"price"`
	ValidityDays      int         `gorm:"not null" json:"validity_days"`
	AccessLevel       AccessLevel `gorm:"type:varchar(20);not null" json:"access_level"`
	MaxDevicesAllowed int         `gorm:"not null;default:1" json:"max_devices_allowed"`
	Resolution        string      `json:"resolution"`
	Description       string      `json:"description"`
	IsActive          bool        `gorm:"default:true;index" json:"is_active"`
	CreatedAt         time.Time   `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt         time.Time   `gorm:"autoUpdateTime" json:"updated_at"`
}

func (Plan) TableName() string { return "plans" }

type Subscription struct {
	ID        uuid.UUID          `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	UserID    uuid.UUID          `gorm:"type:uuid;not null;index" json:"user_id"`
	PlanID    uuid.UUID          `gorm:"type:uuid;not null" json:"plan_id"`
	User      *User              `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Plan      *Plan              `gorm:"foreignKey:PlanID" json:"plan,omitempty"`
	StartDate time.Time          `gorm:"not null" json:"start_date"`
	EndDate   time.Time          `gorm:"not null;index" json:"end_date"`
	IsActive  bool               `gorm:"default:true;index" json:"is_active"`
	Status    SubscriptionStatus `gorm:"type:varchar(20);not null;index" json:"status"`
	CreatedAt time.Time          `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time          `gorm:"autoUpdateTime" json:"updated_at"`
}

func (Subscription) TableName() string  { return "subscriptions" }
func (s *Subscription) IsExpired() bool { return time.Now().After(s.EndDate) }

type WatchHistory struct {
	ID             uuid.UUID   `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	UserID         uuid.UUID   `gorm:"type:uuid;not null;index" json:"user_id"`
	ContentID      uuid.UUID   `gorm:"type:uuid;not null;index" json:"content_id"`
	User           *User       `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Content        *Content    `gorm:"foreignKey:ContentID" json:"content,omitempty"`
	WatchedSeconds int         `gorm:"not null;default:0" json:"watched_seconds"`
	TotalSeconds   int         `gorm:"not null" json:"total_seconds"`
	Status         WatchStatus `gorm:"type:varchar(20);not null;default:'started'" json:"status"`
	LastWatchedAt  time.Time   `gorm:"not null" json:"last_watched_at"`
	CreatedAt      time.Time   `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt      time.Time   `gorm:"autoUpdateTime" json:"updated_at"`
}

func (WatchHistory) TableName() string { return "watch_histories" }
func (w *WatchHistory) ProgressPercentage() float64 {
	if w.TotalSeconds == 0 {
		return 0
	}
	return (float64(w.WatchedSeconds) / float64(w.TotalSeconds)) * 100
}
func (w *WatchHistory) IsCompleted() bool {
	return w.Status == WatchStatusCompleted || w.ProgressPercentage() >= 90.0
}
