package models

import "time"
/*
 User Comment model.
*/
type UserComment struct {
	BaseModel
	UserId    string    `gorethink:"user_id,omitempty" json:"user_id"`
	Comment    string    `gorethink:"comment,omitempty" json:"comment"`
	FeatureRequestId    string    `gorethink:"feature_request_id,omitempty" json:"feature_request_id"`
}
const user_comment_table = "user_comments"

func NewComment(frc *FeatureRequestAddComment, frId string, t time.Time)  *UserComment{
	return &UserComment{
		UserId:frc.UserId,
		Comment:frc.Comment,
		FeatureRequestId:frId,
		BaseModel:BaseModel{
			CreatedAt:&t,
		},
	}
}