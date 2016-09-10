package models


import (
	r "github.com/dancannon/gorethink"
	"time"
)

/*
 Client model.
*/

// ProductArea model structure.
// ID: UUID of the record.
// Name: Client name.
// Description: Additional data for product_areas
// CreatedAt: The date the record is created
// UpdatedAt: The date of the last update
type ProductArea struct {
	BaseModel
	Name        string    `gorethink:"name,omitempty" json:"name"`
	Description string    `gorethink:"description,omitempty" json:"description"`
}

// The product area table name.
const product_area_table = "product_areas"

// Create new product_area from ClientCreateEdit data.
// product_area: data.
// t: time to set CreatedAt.
// Returns:
//	- The Client created.
func NewProductArea(product_area *ProductAreaCreate, t time.Time) *ProductArea {
	return &ProductArea{
		Name: product_area.Name,
		Description: product_area.Description,
		BaseModel: BaseModel{
			CreatedAt: &t,
		},
	}
}

// Get all product areas from the DB.
// Returns:
// 	- Array of product areas (or empty if there are no product_areas in the DB).
// 	- CodeInfo with the error information.
func GetAllProductAreas() ([]ProductArea,*CodeInfo ){
	product_areas := []ProductArea{}
	result := getAll(product_area_table,&product_areas)
	return product_areas, result
}

// Get product area by id.
// Error :
// 	- Returns CodeInfo with the error information.
// Success :
//     - Fills the data of the model calling the method.
//     - Returns CodeInfo with Code = 0 (No error)
func (c *ProductArea) GetById(id string) *CodeInfo {
	return  c.getById(product_area_table,id,c)
}

// Insert new product area.
// Product area should have data before calling this method.
// Error :
// 	- Returns CodeInfo with the error informatiogn.
// Success :
//     - Sets the ID of the model calling the method on Success
//     - Returns CodeInfo with Code = 0 (No error)

func (c *ProductArea) Insert() *CodeInfo {
	id, result := c.insert(product_area_table,c)
	c.ID =id
	return result
}

// Update ProductArea by id.
// Error :
// 	- Returns CodeInfo with the error information.
// Success :
//     - Fills the updated data to the model calling the method.
//     - Returns CodeInfo with Code = 0 (No error)
func (c *ProductArea) Update(id string, data * ProductAreaEdit, t time.Time) *CodeInfo {
	c.setFromProductAreaEdit(data)
	c.UpdatedAt = &t
	return c.update(product_area_table,id,c)
}

// Delete product area by id.
// Error :
// 	- Returns CodeInfo with the error information.
// Success :
//     - Fills the deleted data to the model calling the method.
//     - Returns CodeInfo with Code = 0 (No error)
func (c *ProductArea) Delete(id string) *CodeInfo {
	result,err:= r.Table(feature_requests_table).Filter(r.Row.Field("product_area_id").Eq(id)).Run(c.Session())
	if err != nil {
		return  ErrorInfo(ErrSystem, err.Error())
	}
	if !result.IsNil() {
		return ErrorInfo(ErrRecordHasConnections,"There are feature requests in this product area")
	}
	return c.delete(product_area_table,id,c)
}

// Set Client data from ClientEdit model.
func (c *ProductArea) setFromProductAreaEdit(data *ProductAreaEdit)  {
	c.Name = data.Name
	c.Description = data.Description
}