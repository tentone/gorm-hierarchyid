package main

import (
	"database/sql/driver"
	"encoding/json"
	"errors"

	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

// HierarchyId is a structure to represent database hierarchy ids.
type HierarchyIdDb struct {
	// Path of the hierarchy (e.g "/1/2/3/4/")
	Data HierarchyId
}

// GormDataType gorm common data type
func (HierarchyIdDb) GormDataType() string {
	return "hierarchyid"
}

// GormDBDataType gorm db data type
func (HierarchyIdDb) GormDBDataType(db *gorm.DB, field *schema.Field) string {
	return "hierarchyid"
}

// When marshaling to JSON, we want the field formatted as a string.
func (j HierarchyIdDb) MarshalJSON() ([]byte, error) {
	return json.Marshal(j.Data)
}

// When unmarshaling from JSON, we want to parse the string into the field.
func (j *HierarchyIdDb) UnmarshalJSON(data []byte) error {
	if len(data) == 0 || string(data) == "null" {
		return nil
	}

	var err = json.Unmarshal(data, &j.Data)
	if err != nil {
		return err
	}

	return nil
}

// Value implements the driver.Valuer interface.
//
// Used to provide a value to the SQL server for storage.
func (j HierarchyIdDb) Value() (driver.Value, error) {
	return j.Data, nil
}

// Scan implements the sql.Scanner interface.
//
// Used to read the value provided by the SQL server.
func (j *HierarchyIdDb) Scan(src any) error {
	if src == nil {
		j.Data = nil
		return nil
	}

	switch src := src.(type) {
	case []byte:
		var err error
		j.Data, err = Parse(src)
		if err != nil {
			return err
		}
	default:
		return errors.New("incompatible type to scan")
	}

	return nil
}
