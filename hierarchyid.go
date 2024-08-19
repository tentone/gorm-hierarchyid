package hierarchyid

import (
	"database/sql/driver"
	"encoding/json"
	"errors"

	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

// HierarchyId is a structure to represent database hierarchy ids.
type HierarchyId struct {
	// Path of the hierarchy (e.g "/1/2/3/4/")
	Data HierarchyIdData
}

// GormDataTypeInterface to specify the nema of data type.
func (HierarchyId) GormDataType() string {
	return "hierarchyid"
}

// GormDBDataTypeInterface defines the data type to apply in the database.
func (HierarchyId) GormDBDataType(db *gorm.DB, field *schema.Field) string {
	if db.Dialector.Name() != "sqlserver" {
		panic("hierarchyid is only supported on SQL Server")
	}

	return "hierarchyid"
}

// Get the tree level where this hierarchyid is located.
//
// '/1/2/3/4/' is at level 4, '/1/2/3/' is at level 3, etc.
func (j *HierarchyId) GetLevel() int {
	return len(j.Data)
}

// Get the root of the tree '\'.
//
// The root is the hierarchyid with an empty path.
func GetRoot() HierarchyId {
	return HierarchyId{Data: []int64{}}
}

// Create a new hierarchyid from a string.
func (j *HierarchyId) FromString(data string) error {
	var err error
	j.Data, err = FromString(data)
	return err
}

// Check if a hierarchyid is a descendant of another hierarchyid
func (j *HierarchyId) IsDescendantOf(parent HierarchyId) bool {
	return IsDescendantOf(j.Data, parent.Data)
}

// Calculate a new  hierarchyid when moving from a parent to another parent in the tree.
//
// The position will be calculated based on the old and new parents.
//
// E.g. if the element is on position '/1/2/57/8/' old parents is '/1/2/' and new parent is '/1/3/' the new position will be '/1/3/57/8/'
func (j *HierarchyId) GetReparentedValue(oldAncestor HierarchyId, newAncestor HierarchyId) HierarchyId {
	if !j.IsDescendantOf(oldAncestor) {
		return HierarchyId{}
	}

	path := j.Data
	path = append(newAncestor.Data, path[len(oldAncestor.Data):]...)

	return HierarchyId{Data: path}
}

// Get all ancestors of a hierarchyid.
//
// E.g. '/1/2/3/4/' will return ['/1/', '/1/2/', '/1/2/3/']
func (j *HierarchyId) GetAncestors() []HierarchyId {
	p := []HierarchyId{}
	pd := GetAncestors(j.Data)

	for _, d := range pd {
		p = append(p, HierarchyId{Data: d})
	}

	return p
}

// Create a string representation of the hierarchyid data type
func (j *HierarchyId) ToString() string {
	return ToString(j.Data)
}

// Get the direct parent of a hierarchyid.
func (j *HierarchyId) GetAncestor() HierarchyId {
	return HierarchyId{Data: GetAncestor(j.Data)}
}

// When marshaling to JSON, we want the field formatted as a string.
func (j HierarchyId) MarshalJSON() ([]byte, error) {
	return json.Marshal(ToString(j.Data))
}

// When unmarshaling from JSON, we want to parse the string into the field.
func (j *HierarchyId) UnmarshalJSON(data []byte) error {
	if len(data) == 0 {
		return nil
	}

	str := ""

	err := json.Unmarshal(data, &str)
	if err != nil {
		return err
	}

	j.Data, err = FromString(str)
	if err != nil {
		return err
	}

	return nil
}

// Value implements the driver.Valuer interface.
//
// Used to provide a value to the SQL server for storage.
func (j HierarchyId) Value() (driver.Value, error) {
	if j.Data == nil {
		return nil, nil
	}

	data, err := Encode(j.Data)
	if err != nil {
		return nil, err
	}

	return data, nil
}

// Scan implements the sql.Scanner interface.
//
// Used to read the value provided by the SQL server.
func (j *HierarchyId) Scan(src any) error {
	if src == nil {
		j.Data = nil
		return nil
	}

	switch src := src.(type) {
	case []byte:
		var err error
		j.Data, err = Decode(src)
		if err != nil {
			return err
		}
	default:
		return errors.New("incompatible type to scan")
	}

	return nil
}
