package models

import (
	r "github.com/dancannon/gorethink"
	"time"
	"github.com/kujtimiihoxha/bc-feature-requests/db"
	"github.com/kujtimiihoxha/bc-feature-requests/helpers"
	"fmt"
	"github.com/bradfitz/slice"
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
	TargetDate    *time.Time   `gorethink:"target_date,omitempty" json:"target_date"`
	TicketUrl     string    `gorethink:"ticket_url,omitempty" json:"ticket_url"`
	ProductAreaId string    `gorethink:"product_area_id,omitempty" json:"product_area_id"`
	EmployID      string    `gorethink:"employ_id,omitempty" json:"employ_id"`
	Closed        bool    `gorethink:"closed" json:"closed"`
	Clients       []ClientFeatureRequest  `gorethink:"-" json:"clients"`
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
		Closed: false,
		BaseModel: BaseModel{
			CreatedAt: &t,
		},
	}
	for _, v := range fr.Clients {
		frc.Clients = append(frc.Clients, ClientFeatureRequest{
			ClientId: v.ClientId,
			Priority: v.Priority,
		})
	}
	return frc
}
// Get all clients from the DB.
// Returns:
// 	- Array of clients (or empty if there are no clients in the DB).
// 	- CodeInfo with the error information.
func GetFeatureRequestByFilterSort(filter *FeatureRequestFilter) (FeatureRequestFilterResponse, *CodeInfo) {
	feature_requests := []FeatureRequest{}
	statements, errC := generateFeatureRequestQuery(filter)
	if errC.Code != 0 {
		return FeatureRequestFilterResponse{}, errC
	}
	total := 0
	cntRm, err := r.Table(feature_requests_table).Filter(statements).Count().Run(db.GetSession().(r.QueryExecutor))
	cntRm.One(&total)
	var result *r.Cursor;
	term := r.Table(feature_requests_table).Filter(statements)
	if filter.Skip != 0 {
		term = term.Skip(filter.Skip)
	}
	if filter.Get != 0 {
		term = term.Limit(filter.Get)
	}
	if filter.Field != "" {
		if filter.Dir != "" {
			if filter.Dir == "asc" {
				term = term.OrderBy(r.Asc(filter.Field))
			} else if filter.Dir == "desc" {
				term = term.OrderBy(r.Desc(filter.Field))
			}
		} else {
			term = term.OrderBy(filter.Field)
		}
	}
	result, err = term.Run(db.GetSession().(r.QueryExecutor))
	if err != nil {
		return FeatureRequestFilterResponse{}, ErrorInfo(ErrSystem, err.Error())
	}
	err = result.All(&feature_requests)
	if err != nil {
		return FeatureRequestFilterResponse{}, ErrorInfo(ErrSystem, err.Error())
	}
	for i,v := range feature_requests {
		feature_requests[i].Clients = []ClientFeatureRequest{}
		clientRewRes, err := r.Table(client_feature_request_table).Filter(
			r.Row.Field("feature_request_id").Eq(v.ID)).Run(db.GetSession().(r.QueryExecutor))
		if err != nil {
			return FeatureRequestFilterResponse{}, ErrorInfo(ErrSystem, err.Error())
		}
		if !clientRewRes.IsNil() {
			clientRewRes.All(&feature_requests[i].Clients)
		}
	}
	if filter.ClientPriorityDir != "" && filter.Client != ""{
		if filter.ClientPriorityDir == "asc" {
			slice.Sort(feature_requests[:], func(i, j int) bool {
				var iClient ClientFeatureRequest;
				var jClient ClientFeatureRequest;
				for _,v :=  range feature_requests[i].Clients {
					if v.ClientId == filter.Client {
						iClient = v
					}
				}
				for _,v :=  range feature_requests[j].Clients {
					if v.ClientId == filter.Client {
						jClient = v
					}
				}
				return iClient.Priority > jClient.Priority
			})
		}
		if filter.ClientPriorityDir == "desc" {
			slice.Sort(feature_requests[:], func(i, j int) bool {
				var iClient ClientFeatureRequest;
				var jClient ClientFeatureRequest;
				for _,v :=  range feature_requests[i].Clients {
					if v.ClientId == filter.Client {
						iClient = v
					}
				}
				for _,v :=  range feature_requests[j].Clients {
					if v.ClientId == filter.Client {
						jClient = v
					}
				}
				return iClient.Priority < jClient.Priority
			})
		}
	}
	return FeatureRequestFilterResponse{
		feature_requests,
		total,
	}, errC
}

func generateFeatureRequestQuery(filter *FeatureRequestFilter) (interface{}, *CodeInfo) {
	filterStatements := []interface{}{}
	if filter.Client != "" {
		fr_match := []ClientFeatureRequest{};
		query := r.Table(client_feature_request_table).Filter(
			r.Row.Field("client_id").Eq(filter.Client));
		c_frRes, err := query.Run(db.GetSession().(r.QueryExecutor))
		if err != nil {
			return filterStatements, ErrorInfo(ErrSystem, err.Error())
		}
		err = c_frRes.All(&fr_match)
		if err != nil {
			return filterStatements, ErrorInfo(ErrSystem, err.Error())
		}
		statements := []interface{}{}
		ids := []interface{}{}
		for _, v := range fr_match {
			ids = append(ids,v.FeatureRequestId)
		}
		fmt.Println(fr_match)
		helpers.RemoveDuplicates(&ids)
		for _,v := range ids {
			statements = append(statements, r.Row.Field("id").Eq(v))
		}
		filterStatements = append(filterStatements, r.Or(statements...))
	}
	if filter.Closed != 0 {
		if filter.Closed == 1 {
			filterStatements = append(filterStatements, r.Row.Field("closed").Eq(true))
		} else if filter.Closed == 2 {
			filterStatements = append(filterStatements, r.Row.Field("closed").Eq(false))
		}
	}
	if filter.Employ != "" {
		filterStatements = append(filterStatements, r.Row.Field("employ_id").Eq(filter.Employ))
	}
	if filter.ProductArea != "" {
		filterStatements = append(filterStatements, r.Row.Field("product_area_id").Eq(filter.ProductArea))
	}
	return r.And(filterStatements...), OkInfo("")
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

func CheckPriority(cp  ClientFeatureRequest) *CodeInfo {
	c_fr := ClientFeatureRequest{}
	c_frRes, err := r.Table(client_feature_request_table).Filter(r.And(r.Row.Field("client_id").Eq(
		cp.ClientId), r.Row.Field("priority").Eq(cp.Priority), r.Row.Field("id").Ne(cp.ID))).Run(db.GetSession().(r.QueryExecutor))
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