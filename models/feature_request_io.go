package models

import "time"

// Feature request create structure.
// This is the model that will be sent from the frontend to create a new feature request.
// Title: FeatureRequest title. REQUIRED
// Description: Request description. REQUIRED
// TargetDate: Date of completion. REQUIRED
// TicketUrl: Ticket URL. REQUIRED
// ProductAreaId: UUID of the product area. REQUIRED
// Clients: List of clients and priorities of clients for this request. REQUIRED
type FeatureRequestCreate struct {
	Title        string    `json:"title" valid:"ascii,required"`
	Description string    `json:"description" valid:"ascii,required"`
	TargetDate   *time.Time   `json:"target_date" valid:"required"`
	TicketUrl string    `json:"ticket_url"  valid:"url,required"`
	ProductAreaId string    `json:"product_area_id"   valid:"uuid,required"`
	Clients []struct{
		ClientId string  `json:"client_id"   valid:"uuid,required"`
		Priority int `json:"priority"   valid:"required"`
	}  `json:"clients"  valid:"required"`
}
type FeatureRequestFilterResponse struct {
	Data []FeatureRequest `json:"data"`
	Total int `json:"total"`
}
type FeatureRequestFilter struct {
	FeatureRequestSort
	FeatureRequestPagination
	Employ string `json:"employ"`
	Client string `json:"client"`
	Closed int `json:"closed"`
	ProductArea string `json:"product_area"`

}
type FeatureRequestSort struct {
	Field string `json:"field"`
	Dir string `json:"dir"`
}
type FeatureRequestPagination struct {
	Skip int `json:"skip"`
	Get int `json:"get"`
}