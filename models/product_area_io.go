package models

/*
 Client IO models.
*/


// ProductArea create structure.
// This is the model that will be sent from the frontend to create a new client.
// Name: ProductArea name. REQUIRED
// Description: Additional data for clients REQUIRED
type ProductAreaCreate struct {
	Name        string    `json:"name,omitempty" valid:"ascii,required"`
	Description string    `json:"description,omitempty"  valid:"ascii,optional"`
}

// ProductArea edit structure.
// This is the model that will be sent from the frontend to edit a client.
// Name: ProductArea name. Optional
// Description: Additional data for product atrea Optional
type ProductAreaEdit struct {
	Name        string    	`json:"name,omitempty" valid:"ascii,optional"`
	Description string    	`json:"description,omitempty"  valid:"ascii,optional"`
}
