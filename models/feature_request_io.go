package models

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
	TargetDate  string   `json:"target_date" valid:"required"`
	TicketUrl string    `json:"ticket_url"  valid:"url,required"`
	ProductAreaId string    `json:"product_area_id"   valid:"uuid,required"`
	Clients []struct{
		ClientId string  `json:"client_id"   valid:"uuid,required"`
		Priority int `json:"priority"   valid:"required"`
	}  `json:"clients"  valid:"required"`
}
