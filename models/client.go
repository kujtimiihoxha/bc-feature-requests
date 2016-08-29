package models

import (
	"time"
)

/*
 Client model.
*/

// Client model structure.
// ID: UUID of the record.
// Name: Client name.
// Description: Additional data for clients
// CreatedAt: The date the record is created
// UpdatedAt: The date of the last update
type Client struct {
	BaseModel
	Name        string    `gorethink:"name,omitempty" json:"name"`
	Description string    `gorethink:"description,omitempty" json:"description"`
}

// The clients table name.
const client_table = "clients"

// Create new client from ClientCreateEdit data.
// client: data.
// t: time to set CreatedAt.
// Returns:
//	- The Client created.
func NewClient(client *ClientCreate, t time.Time) *Client {
	return &Client{
		Name: client.Name,
		Description: client.Description,
		BaseModel: BaseModel{
			CreatedAt: &t,
		},
	}
}

// Get all clients from the DB.
// Returns:
// 	- Array of clients (or empty if there are no clients in the DB).
// 	- CodeInfo with the error information.
func GetAllClients() ([]Client,*CodeInfo ){
	clients := []Client{}
	result := getAll(client_table,&clients)
	return clients, result
}

// Get client by id.
// Error :
// 	- Returns CodeInfo with the error information.
// Success :
//     - Fills the data of the model calling the method.
//     - Returns CodeInfo with Code = 0 (No error)
func (c *Client) GetById(id string) *CodeInfo {
	return  c.getById(client_table,id,c)
}

// Insert new client.
// Client should have data before calling this method.
// Error :
// 	- Returns CodeInfo with the error informatiogn.
// Success :
//     - Sets the ID of the model calling the method on Success
//     - Returns CodeInfo with Code = 0 (No error)

func (c *Client) Insert() *CodeInfo {
	id, result := c.insert(client_table,c)
	c.ID =id
	return result
}

// Update client by id.
// Error :
// 	- Returns CodeInfo with the error information.
// Success :
//     - Fills the updated data to the model calling the method.
//     - Returns CodeInfo with Code = 0 (No error)
func (c *Client) Update(id string, data * ClientEdit, t time.Time) *CodeInfo {
	c.setFromClientEdit(data)
	c.UpdatedAt = &t
	return c.update(client_table,id,c)
}

// Delete client by id.
// Error :
// 	- Returns CodeInfo with the error information.
// Success :
//     - Fills the deleted data to the model calling the method.
//     - Returns CodeInfo with Code = 0 (No error)
func (c *Client) Delete(id string) *CodeInfo {
	return c.delete(client_table,id,c)
}

// Set Client data from ClientEdit model.
func (c *Client) setFromClientEdit(data *ClientEdit)  {
	c.Name = data.Name
	c.Description = data.Description
}