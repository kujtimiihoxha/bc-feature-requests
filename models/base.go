package models

import (
	r "github.com/dancannon/gorethink"
	"github.com/kujtimiihoxha/bc-feature-requests/db"
	"time"
)

// BaseModel.
// All other models "inherit" this model.
// ID: UUID of the record.
// CreatedAt: The date the record is created
// UpdatedAt: The date of the last update
type BaseModel struct {
	ID        string     `gorethink:"id,omitempty" json:"id"`
	CreatedAt *time.Time `gorethink:"created_at,omitempty" json:"created_at"`
	UpdatedAt *time.Time `gorethink:"updated_at,omitempty" json:"updated_at,omitempty"`
}

// Get all records from the table.
// Returns:
// 	- Array of records (or empty if there are no records in the table).
// 	- CodeInfo with the error information.
func getAll(table string, arr interface{}) *CodeInfo {
	res, err := r.Table(table).Run(db.GetSession().(r.QueryExecutor))
	if err != nil {
		return ErrorInfo(ErrSystem, err.Error())
	}
	err = res.All(arr)
	if err != nil {
		return ErrorInfo(ErrSystem, err.Error())
	}
	return OkInfo("")
}

// Get record by id.
// Error :
// 	- Returns CodeInfo with the error information.
// Success :
//     - Fills the data of the model calling the method.
//     - Returns CodeInfo with Code = 0 (No error)
func (b *BaseModel) getById(table string, id string, model interface{}) *CodeInfo {
	res, err := r.Table(table).Get(id).Run(b.Session())
	if err != nil {
		return ErrorInfo(ErrSystem, err.Error())
	}
	if res.IsNil() {
		return ErrorInfo(ErrNotFound, "Not Found")
	}
	err = res.One(model)
	if err != nil {
		return ErrorInfo(ErrSystem, err.Error())
	}
	return OkInfo("")
}

// Insert new record.
// Model should have data before calling this method.
// Error :
// 	- Returns Empty Id and CodeInfo with the error information.
// Success :
//     - Returns Generated ID and CodeInfo with Code = 0 (No error)

func (b *BaseModel) insert(table string, model interface{}) (string, *CodeInfo) {
	wr, err := r.Table(table).Insert(model).RunWrite(b.Session())
	if wr.Errors > 0 {
		return "", ErrorInfo(ErrDatabase, wr.FirstError)
	}
	if err != nil {
		return "", ErrorInfo(ErrSystem, err.Error())
	}
	return wr.GeneratedKeys[0], OkInfo("Data inserted succesfully")
}

// Update row by id.
// Error :
// 	- Returns CodeInfo with the error information.
// Success :
//     - Fills the updated data to the model.
//     - Returns CodeInfo with Code = 0 (No error)
func (b *BaseModel) update(table string, id string, model interface{}) *CodeInfo {
	wr, err := r.Table(table).Get(id).Update(model).RunWrite(b.Session())
	if wr.Skipped > 0 {
		return ErrorInfo(ErrNotFound, "Not Found")
	}
	if wr.Errors > 0 {
		return ErrorInfo(ErrDatabase, wr.FirstError)
	}
	if err != nil {
		return ErrorInfo(ErrSystem, err.Error())
	}
	rs, err := r.Table(table).Get(id).Run(b.Session())
	if err != nil {
		return ErrorInfo(ErrSystem, err.Error())
	}
	if rs.IsNil() {
		return ErrorInfo(ErrNotFound, "Not Found")
	}
	rs.One(model)
	return OkInfo("Data updated succesfully")
}

// Delete record by id.
// Error :
// 	- Returns CodeInfo with the error information.
// Success :
//     - Fills the deleted data to the model.
//     - Returns CodeInfo with Code = 0 (No error)
func (b *BaseModel) delete(table string, id string, model interface{}) *CodeInfo {
	res, err := r.Table(table).Get(id).Run(b.Session())
	if err != nil {
		return ErrorInfo(ErrSystem, err.Error())
	}
	if res.IsNil() {
		return ErrorInfo(ErrNotFound, "Not Found")
	}
	err = res.One(model)
	if err != nil {
		return ErrorInfo(ErrSystem, err.Error())
	}
	wr, err := r.Table(table).Get(id).Delete().RunWrite(b.Session())

	if wr.Errors > 0 {
		return ErrorInfo(ErrDatabase, wr.FirstError)
	}
	if err != nil {
		return ErrorInfo(ErrSystem, err.Error())
	}
	return OkInfo("Data deleted succesfully")
}
func (_ *BaseModel) Session() r.QueryExecutor {
	return db.GetSession().(r.QueryExecutor)
}
