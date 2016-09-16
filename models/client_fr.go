package models

import "time"

type ClientFeatureRequest struct {
	BaseModel
	FeatureRequestId string `gorethink:"feature_request_id,omitempty" json:"feature_request_id"`
	ClientId         string `gorethink:"client_id,omitempty" json:"client_id"`
	Priority         int    `gorethink:"priority,omitempty" json:"priority"`
}

// The client feature request table name.
const client_feature_request_table = "client_feature_request"

func NewClientFeatureRequest(frID string, clID string, priority int, t time.Time) *ClientFeatureRequest {
	return &ClientFeatureRequest{
		FeatureRequestId: frID,
		ClientId:         clID,
		Priority:         priority,
		BaseModel: BaseModel{
			CreatedAt: &t,
		},
	}
}
