package models

/*
 Client IO models.
*/


// Client create structure.
// This is the model that will be sent from the frontend to create a new client.
// Name: Client name. REQUIRED
// Description: Additional data for clients REQUIRED
type ClientCreate struct {
	Name        string    `json:"name,omitempty" valid:"ascii,required"`
	Description string    `json:"description,omitempty"  valid:"ascii,optional"`
}

// Client edit structure.
// This is the model that will be sent from the frontend to edit a client.
// Name: Client name. Optional
// Description: Additional data for clients Optional
type ClientEdit struct {
	Name        string    	`json:"name,omitempty" valid:"ascii,optional"`
	Description string    	`json:"description,omitempty"  valid:"ascii,optional"`
}
