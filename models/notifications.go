package models

import (
	"time"
)

type Notification struct {
	BaseModel
	UserID  string             `gorethink:"user_id,omitempty" json:"user_id"`
	Details *FeatureRequestLog `gorethink:"details,omitempty" json:"details"`
	Viewed  bool               `gorethink:"viewed" json:"viewed"`
	Link    string             `gorethink:"link,omitempty" json:"link"`
}

const notifications_table = "notifications"

func NewNotification(userId string, link string, details *FeatureRequestLog, t time.Time) *Notification {
	return &Notification{
		BaseModel: BaseModel{
			CreatedAt: &t,
		},
		UserID:  userId,
		Details: details,
		Viewed:  false,
		Link:    link,
	}
}
func (n *Notification) Insert() *CodeInfo {
	id, result := n.insert(notifications_table, n)
	n.ID = id
	return result
}
