package models

import (
	r "github.com/dancannon/gorethink"
	"time"
	"github.com/kujtimiihoxha/bc-feature-requests/db"
)

/*
 Client model.
*/

// FeatureRequest model structure.
// ID: UUID of the record.
// Description: Request description.
// TargetDate: Date of completion.
// TicketUrl: Ticket URL.
// ProductAreaId: UUID of the product area.
// EmployID: UUID of the employ submitting the feature request.
// Clients: List of clients and priorities of clients for this request.
// CreatedAt: The date the record is created
// UpdatedAt: The date of the last update
type FeatureRequest struct {
	BaseModel
	Title         string    `gorethink:"title,omitempty" json:"title"`
	Description   string    `gorethink:"description,omitempty" json:"description"`
	TargetDate    string   `gorethink:"target_date,omitempty" json:"target_date"`
	TicketUrl     string    `gorethink:"ticket_url,omitempty" json:"ticket_url"`
	ProductAreaId string    `gorethink:"product_area_id,omitempty" json:"product_area_id"`
	EmployID      string    `gorethink:"employ_id,omitempty" json:"employ_id"`
	Clients     []ClientFeatureRequest  `gorethink:"-" json:"-"`
}
// The feature requests table name.
const feature_requests_table = "feature_requests"

// Create new feature request from FeatureRequestCreate data.
// fr: data.
// t: time to set CreatedAt.
// employID: the id of the employ submitting the request.
// Returns:
//	- The Feature request created.
func NewFeatureRequest(fr *FeatureRequestCreate, t time.Time, employID string) *FeatureRequest {
	frc := &FeatureRequest{
		Title: fr.Title,
		Description: fr.Description,
		EmployID: employID,
		ProductAreaId: fr.ProductAreaId,
		TicketUrl: fr.TicketUrl,
		Clients:[]ClientFeatureRequest{},
		TargetDate: fr.TargetDate,
		BaseModel: BaseModel{
			CreatedAt: &t,
		},
	}
	for _,v:= range fr.Clients {
		frc.Clients = append(frc.Clients,ClientFeatureRequest{
			ClientId: v.ClientId,
			Priority: v.Priority,
		})
	}
	return frc
}

// Insert new feature request.
// Feature request should have data before calling this method.
// Error :
// 	- Returns CodeInfo with the error informatiogn.
// Success :
//     - Sets the ID of the model calling the method on Success
//     - Returns CodeInfo with Code = 0 (No error)

func (c *FeatureRequest) Insert() *CodeInfo {
	for _, v := range c.Clients {
		CheckPriority(v)
	}
	id, result := c.insert(feature_requests_table, c)
	c.ID = id
	for _, v := range c.Clients {
		cfr := ClientFeatureRequest{
			FeatureRequestId: c.ID,
			ClientId: v.ClientId,
			Priority: v.Priority,
		}
		_, result = c.insert(client_feature_request_table, cfr)
		if result.Code != 0 {
			return result
		}
	}
	return result
}

func CheckPriority(cp  ClientFeatureRequest)  *CodeInfo{
	c_fr := ClientFeatureRequest{}
	c_frRes, err := r.Table(client_feature_request_table).Filter(r.And(r.Row.Field("client_id").Eq(
		cp.ClientId), r.Row.Field("priority").Eq(cp.Priority),r.Row.Field("id").Ne(cp.ID))).Run(db.GetSession().(r.QueryExecutor))
	if err != nil {
		return ErrorInfo(ErrSystem, err.Error())
	}
	if c_frRes.IsNil() {
		return OkInfo("")
	}
	c_frRes.One(&c_fr)
	c_fr.Priority++
	wr, err := r.Table(client_feature_request_table).Get(c_fr.ID).Update(c_fr).RunWrite(
		db.GetSession().(r.QueryExecutor))
	if err != nil {
		return ErrorInfo(ErrSystem, err.Error())
	}
	if wr.Errors > 0 {
		return ErrorInfo(ErrDatabase, wr.FirstError)
	}
	return CheckPriority(c_fr)
}