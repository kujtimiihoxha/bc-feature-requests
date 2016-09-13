package models

import "time"
/*
 Client model.
*/
type FeatureRequestLog struct {
	BaseModel
	FeatureRequestID    string    `gorethink:"feature_request_id,omitempty" json:"feature_request_id"`
	UserId    string    `gorethink:"user_id,omitempty" json:"user_id"`
	Type    string    `gorethink:"type,omitempty" json:"type"`
	Icon    string    `gorethink:"icon,omitempty" json:"icon"`
	Description    string    `gorethink:"description,omitempty" json:"description"`
}
// The feature requests table name.
const feature_request_log_table = "feature_request_log"

func NewFeatureRequestLog(userId string,featureRid string,tp string, icon string, description string) *FeatureRequestLog {
	t:=time.Now().UTC()
	return &FeatureRequestLog{
		UserId:userId,
		FeatureRequestID:featureRid,
		Type:tp,
		Icon:icon,
		Description:description,
		BaseModel: BaseModel{
			CreatedAt:&t,
		},
	}
}
const(
	TITLE_UPDATE = "title_update"
	DESCRIPTION_UPDATE = "description_update"
	PRODUCT_ARE_UPDATE= "product_are_update"
	TICKET_URL_UPDATE= "ticket_url_update"
	TARGET_DATE= "target_date"
	CHANGED_CLIENTS= "changed_clients"
	CHANGED_PRIORITY= "changed_priority"
	STATE_CLOSE= "closed"
	STATE_OPEN= "reopen"
)
var LOG_MESSAGES map[string]string= map[string]string{
	TITLE_UPDATE : "<b><i>%s</i></b> updated the title of the feature request",
	DESCRIPTION_UPDATE : "<b><i>%s</i></b> changed the description of the feature request",
	PRODUCT_ARE_UPDATE : "<b><i>%s</i></b> modified the product area",
	TICKET_URL_UPDATE : "<b><i>%s</i></b> updated the ticket url of the feature request",
	TARGET_DATE : "<b><i>%s</i></b> changed the target date",
	CHANGED_CLIENTS : "<b><i>%s</i></b> updated the clients of the feature request",
	STATE_CLOSE : "<b><i>%s</i></b> closed feature request",
	STATE_OPEN : "<b><i>%s</i></b> reopened feature request",
	CHANGED_PRIORITY : "<b><i>%s</i></b> changed the priority of a feature request",
}
var ICONS map[string]string= map[string]string{
	TITLE_UPDATE : "title",
	DESCRIPTION_UPDATE : "subtitles",
	PRODUCT_ARE_UPDATE : "aspect_ratio",
	TICKET_URL_UPDATE : "link",
	TARGET_DATE : "date_range",
	CHANGED_CLIENTS : "folder",
	STATE_CLOSE : "highlight_off",
	STATE_OPEN : "info_outline",
	CHANGED_PRIORITY : "star_border",
}