package models

import (
	r "github.com/dancannon/gorethink"
	"time"
	"github.com/kujtimiihoxha/bc-feature-requests/db"
	"github.com/kujtimiihoxha/bc-feature-requests/helpers"
	"fmt"
	"github.com/bradfitz/slice"
	"strings"
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
	TitleNormalized         string    `gorethink:"title_normalized,omitempty" json:"-"`
	Description   string    `gorethink:"description,omitempty" json:"description"`
	TargetDate    *time.Time   `gorethink:"target_date,omitempty" json:"target_date"`
	TicketUrl     string    `gorethink:"ticket_url,omitempty" json:"ticket_url"`
	ProductAreaId string    `gorethink:"product_area_id,omitempty" json:"product_area_id"`
	EmployID      string    `gorethink:"employ_id,omitempty" json:"employ_id"`
	Closed        bool    `gorethink:"closed" json:"closed"`
	Clients       []ClientFeatureRequest  `gorethink:"-" json:"clients"`
	Modifications       []FeatureRequestLog  `gorethink:"-" json:"modifications"`
	Comments       []UserComment  `gorethink:"-" json:"comments"`
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
		TitleNormalized: strings.ToLower(fr.Title),
		Description: fr.Description,
		EmployID: employID,
		ProductAreaId: fr.ProductAreaId,
		TicketUrl: fr.TicketUrl,
		Clients:[]ClientFeatureRequest{},
		TargetDate: fr.TargetDate,
		Closed: false,
		BaseModel: BaseModel{
			CreatedAt: &t,
			UpdatedAt: &t,
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
	if filter.Skip != 0 {
		term = term.Skip(filter.Skip)
	}
	if filter.Get != 0 {
		term = term.Limit(filter.Get)
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

// Get client by id.
// Error :
// 	- Returns CodeInfo with the error information.
// Success :
//     - Fills the data of the model calling the method.
//     - Returns CodeInfo with Code = 0 (No error)
func (c *FeatureRequest) GetById(id string) *CodeInfo {
	ci:=c.getById(feature_requests_table, id, c)
	if ci.Code != 0 {
		return  ci
	}
	c.Clients = []ClientFeatureRequest{}
	clientRewRes, err := r.Table(client_feature_request_table).Filter(
		r.Row.Field("feature_request_id").Eq(c.ID)).Run(c.Session())
	if err != nil {
		return ErrorInfo(ErrSystem, err.Error())
	}
	if !clientRewRes.IsNil() {
		clientRewRes.All(&c.Clients)
	}
	c.getLogs()
	c.getComments()
	return ci
}


// Update feature request target date.
// Error :
// 	- Returns CodeInfo with the error information.
// Success :
//     - Fills the updated data to the model calling the method.
//     - Returns CodeInfo with Code = 0 (No error)
func (c *FeatureRequest) UpdateTargetDate(id string,userId string, username string, data *FeatureRequestEditTargetDate, t time.Time) *CodeInfo {
	c.setFromFeatureRequestEditTargetDate(data)
	fmt.Println(c)
	c.UpdatedAt = &t
	result := c.update(feature_requests_table, id, c)
	if result.Code != 0 {
		return result
	}
	wr, err := r.Table(feature_request_log_table).Insert(NewFeatureRequestLog(
		userId,
		id,
		TARGET_DATE,
		ICONS[TARGET_DATE],
		fmt.Sprintf(LOG_MESSAGES[TARGET_DATE],
		username))).RunWrite(c.Session())
	if wr.Errors > 0 {
		return  ErrorInfo(ErrDatabase, wr.FirstError)
	}
	if err != nil {
		return  ErrorInfo(ErrSystem, err.Error())
	}
	c.getLogs()
	return OkInfo("Data updated succesfully")
}
// Update feature request details.
// Error :
// 	- Returns CodeInfo with the error information.
// Success :
//     - Fills the updated data to the model calling the method.
//     - Returns CodeInfo with Code = 0 (No error)
func (c *FeatureRequest) UpdateDetails(id string,userId string, username string,data *FeatureRequestEditDetails,t time.Time) *CodeInfo {
	c.setFromFeatureRequestEditDetails(data)
	c.UpdatedAt = &t
	result := c.update(feature_requests_table, id, c)
	if result.Code != 0 {
		return  result
	}

	updates := []*FeatureRequestLog{}
	for _,v := range data.Modifications {
		updates = append(updates,NewFeatureRequestLog(userId,id,v,ICONS[v],fmt.Sprintf(LOG_MESSAGES[v],username)))
	}
	wr, err := r.Table(feature_request_log_table).Insert(updates).RunWrite(c.Session())
	if wr.Errors > 0 {
		return ErrorInfo(ErrDatabase, wr.FirstError)
	}
	if err != nil {
		return ErrorInfo(ErrSystem, err.Error())
	}
	c.getLogs()
	return OkInfo("Data updated succesfully")
}
func (c *FeatureRequest) getLogs() *CodeInfo {
	res,err := r.Table(feature_request_log_table).Filter(r.Row.Field("feature_request_id").Eq(c.ID)).OrderBy(r.Desc("created_at")).Run(db.GetSession().(r.QueryExecutor))
	if err != nil {
		return  ErrorInfo(ErrSystem, err.Error())
	}
	err = res.All(&c.Modifications)
	if err != nil {
		return  ErrorInfo(ErrSystem, err.Error())
	}
	return OkInfo("")
}
// Update feature request state.
// Error :
// 	- Returns CodeInfo with the error information.
// Success :
//     - Fills the updated data to the model calling the method.
//     - Returns CodeInfo with Code = 0 (No error)
func (c *FeatureRequest) UpdateState(id string,userId string, username string, state bool, t time.Time) *CodeInfo {
	c.Closed = state
	c.UpdatedAt = &t
	result := c.update(feature_requests_table, id, c)
	if result.Code != 0 {
		return result
	}
	stateChange := STATE_OPEN
	if state {
		stateChange = STATE_CLOSE
	}
	wr, err := r.Table(feature_request_log_table).Insert(NewFeatureRequestLog(
		userId,
		id,
		stateChange,
		ICONS[stateChange],
		fmt.Sprintf(LOG_MESSAGES[stateChange],
			username))).RunWrite(c.Session())
	if wr.Errors > 0 {
		return ErrorInfo(ErrDatabase, wr.FirstError)
	}
	if err != nil {
		return ErrorInfo(ErrSystem, err.Error())
	}
	c.getLogs()
	return OkInfo("Data updated succesfully")
}

// Add Remove Clients
// Error :
// 	- Returns CodeInfo with the error information.
// Success :
//     - Fills the updated data to the model calling the method.
//     - Returns CodeInfo with Code = 0 (No error)
func (c *FeatureRequest) AddRemoveClients(id string,userId string, username string, addRemove *FeatureRequestAddRemoveClients, t time.Time) *CodeInfo {
	if len(addRemove.ClientsToRemove) > 0 {
		statements := []interface{}{}
		for _,v := range addRemove.ClientsToRemove {
			statements = append(statements,r.Row.Field("client_id").Eq(v))
		}
		wr,err := r.Table(client_feature_request_table).Filter(r.And(r.Or(statements...),r.Row.Field("feature_request_id").Eq(id))).Delete().RunWrite(c.Session())

		if wr.Errors > 0 {
			return ErrorInfo(ErrDatabase, wr.FirstError)
		}
		if err != nil {
			return  ErrorInfo(ErrSystem, err.Error())
		}
	}
	clientsToAdd:=[]*ClientFeatureRequest{}
	for _,v:= range addRemove.ClientsToAdd {
		clientFR := NewClientFeatureRequest(id,v.Client_id,v.Priority,time.Now().UTC());
		CheckPriority(*clientFR)
		clientsToAdd = append(clientsToAdd,clientFR)
	}
	wr,err := r.Table(client_feature_request_table).Insert(clientsToAdd).RunWrite(c.Session())

	if wr.Errors > 0 {
		return ErrorInfo(ErrDatabase, wr.FirstError)
	}
	if err != nil {
		return  ErrorInfo(ErrSystem, err.Error())
	}
	c.UpdatedAt = &t
	result := c.update(feature_requests_table, id, c)
	if result.Code != 0 {
		return  result
	}
	clientRewRes, err := r.Table(client_feature_request_table).Filter(
		r.Row.Field("feature_request_id").Eq(c.ID)).Run(c.Session())
	if err != nil {
		return ErrorInfo(ErrSystem, err.Error())
	}
	if !clientRewRes.IsNil() {
		clientRewRes.All(&c.Clients)
	}
	wr, err = r.Table(feature_request_log_table).Insert(NewFeatureRequestLog(
		userId,
		id,
		CHANGED_CLIENTS,
		ICONS[CHANGED_CLIENTS],
		fmt.Sprintf(LOG_MESSAGES[CHANGED_CLIENTS],
			username))).RunWrite(c.Session())
	if wr.Errors > 0 {
		return ErrorInfo(ErrDatabase, wr.FirstError)
	}
	if err != nil {
		return ErrorInfo(ErrSystem, err.Error())
	}
	c.getLogs()
	return OkInfo("")
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
	if result.Code != 0 {
		return result
	}
	c.ID = id
	cfrs:=[]*ClientFeatureRequest{}
	for _,v:= range c.Clients{
		cfrs = append(cfrs,NewClientFeatureRequest(c.ID,v.ClientId, v.Priority,time.Now().UTC()))
	}
	_, result = c.insert(client_feature_request_table, cfrs)
	return result
}
func (c *FeatureRequest) AddComment(id string,data *FeatureRequestAddComment) *CodeInfo {
	comment := NewComment(data,id,time.Now().UTC())
	_, result := c.insert(user_comment_table, comment)
	fmt.Println(result)
	if result.Code != 0 {
		return result
	}
	c.ID = id
	c.getComments()
	return result
}
func (c *FeatureRequest) getComments()  *CodeInfo{
	res,err := r.Table(user_comment_table).Filter(r.Row.Field("feature_request_id").Eq(c.ID)).OrderBy(r.Asc("created_at")).Run(db.GetSession().(r.QueryExecutor))
	if err != nil {
		return  ErrorInfo(ErrSystem, err.Error())
	}
	err = res.All(&c.Comments)
	if err != nil {
		return  ErrorInfo(ErrSystem, err.Error())
	}
	return OkInfo("")
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

// Set FeatureRequest data from FeatureRequestEditTargetDate model.
func (c *FeatureRequest) setFromFeatureRequestEditTargetDate(data *FeatureRequestEditTargetDate) {
	c.TargetDate = data.TargetDate
}
// Set FeatureRequest data from FeatureRequestEditTargetDate model.
func (c *FeatureRequest) setFromFeatureRequestEditDetails(data *FeatureRequestEditDetails) {
	c.Title = data.Title
	c.Description = data.Description
	c.TicketUrl = data.TicketUrl
	c.ProductAreaId = data.ProductAreaId
}