package main

import (
	"slices"
	"testing"

	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
)

func TestCreateRead(t *testing.T) {
	dsn := "sqlserver://sa:12345678@localhost:1433?database=test"
	db, err := gorm.Open(sqlserver.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatal("Failed to connect database", err)
	}

	type TestCreateReadTable struct {
		gorm.Model
		Path HierarchyId
	}

	_ = db.Migrator().DropTable(&TestCreateReadTable{})

	err = db.AutoMigrate(&TestCreateReadTable{})
	if err != nil {
		t.Fatal("Failed to migrate table", err)
	}

	new := &TestCreateReadTable{Path: HierarchyId{Data: []int64{1, 2, 3, 4}}}
	conn := db.Create(new)
	if conn.Error != nil {
		t.Fatal("Failed to create entry", new.Path.Data, err)
	}

	hid := TestCreateReadTable{}
	conn = db.First(&TestCreateReadTable{})
	if conn.Error != nil {
		t.Fatal("Failed to query database", conn.Error)
	}

	if slices.Compare([]int64{1, 2, 3, 4}, hid.Path.Data) == 0 {
		t.Fatal("Values read are not correct", hid.Path.Data)
	}
}

func TestUniqueNil(t *testing.T) {
	dsn := "sqlserver://sa:12345678@localhost:1433?database=test"
	db, err := gorm.Open(sqlserver.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatal("Failed to connect database", err)
	}

	type TestUniqueTable struct {
		gorm.Model
		Path HierarchyId `gorm:"unique;not null;"`
	}

	_ = db.Migrator().DropTable(&TestUniqueTable{})

	err = db.AutoMigrate(&TestUniqueTable{})
	if err != nil {
		t.Fatal("Failed to migrate table", err)
	}

	new := &TestUniqueTable{Path: HierarchyId{Data: []int64{1, 2, 3, 4}}}
	conn := db.Create(new)
	if conn.Error != nil {
		t.Fatal("Failed to create entry", new.Path.Data, err)
	}

	conn = db.Create(new)
	if conn.Error == nil {
		t.Fatal("Should not be able to create duplicated entry", new.Path.Data, err)
	}

	conn = db.Create(&TestUniqueTable{Path: HierarchyId{Data: nil}})
	if conn.Error == nil {
		t.Fatal("Should not be able to create null entry", new.Path.Data, err)
	}
}

func TestParents(t *testing.T) {
	dsn := "sqlserver://sa:12345678@localhost:1433?database=test"
	db, err := gorm.Open(sqlserver.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatal("Failed to connect database", err)
	}

	type TestParentsTable struct {
		gorm.Model
		Path HierarchyId `gorm:"unique;not null;type:hierarchyid;"`
	}

	_ = db.Migrator().DropTable(&TestParentsTable{})

	err = db.AutoMigrate(&TestParentsTable{})
	if err != nil {
		t.Fatal("Failed to migrate table", err)
	}

	child := &TestParentsTable{Path: HierarchyId{Data: []int64{1, 2, 3, 4}}}

	_ = db.Create(child)
	_ = db.Create(&TestParentsTable{Path: HierarchyId{Data: []int64{1, 2, 3}}})
	_ = db.Create(&TestParentsTable{Path: HierarchyId{Data: []int64{1, 2}}})
	_ = db.Create(&TestParentsTable{Path: HierarchyId{Data: []int64{1}}})
	_ = db.Create(&TestParentsTable{Path: HierarchyId{Data: []int64{2}}})
	_ = db.Create(&TestParentsTable{Path: HierarchyId{Data: []int64{2, 1}}})
	_ = db.Create(&TestParentsTable{Path: HierarchyId{Data: []int64{3}}})

	var count int64 = 0
	_ = db.Model(&TestParentsTable{}).Where("[path] IN (?)", child.Path.GetAncestors()).Count(&count)

	if count != 3 {
		t.Fatal("Expected 3 parents, got", count)
	}
}

func TestGetTreeLevelSearch(t *testing.T) {
	dsn := "sqlserver://sa:12345678@localhost:1433?database=test"
	db, err := gorm.Open(sqlserver.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatal("Failed to connect database", err)
	}

	type TestParentsTable struct {
		gorm.Model

		Path HierarchyId `gorm:"unique;not null;"`

		Name string `gorm:"not null;"`
	}

	_ = db.Migrator().DropTable(&TestParentsTable{})

	err = db.AutoMigrate(&TestParentsTable{})
	if err != nil {
		t.Fatal("Failed to migrate table", err)
	}

	_ = db.Create(&TestParentsTable{Name: "ade", Path: HierarchyId{Data: []int64{1, 2, 3, 4}}})
	_ = db.Create(&TestParentsTable{Name: "a", Path: HierarchyId{Data: []int64{1, 2, 3}}})
	_ = db.Create(&TestParentsTable{Name: "a", Path: HierarchyId{Data: []int64{1, 2}}})
	_ = db.Create(&TestParentsTable{Name: "a", Path: HierarchyId{Data: []int64{1}}})
	_ = db.Create(&TestParentsTable{Name: "b", Path: HierarchyId{Data: []int64{2}}})
	_ = db.Create(&TestParentsTable{Name: "b", Path: HierarchyId{Data: []int64{2, 1}}})
	_ = db.Create(&TestParentsTable{Name: "c", Path: HierarchyId{Data: []int64{3}}})
	_ = db.Create(&TestParentsTable{Name: "c", Path: HierarchyId{Data: []int64{3, 1}}})
	_ = db.Create(&TestParentsTable{Name: "cde", Path: HierarchyId{Data: []int64{3, 1, 1}}})

	elements := []TestParentsTable{}

	// Get all elements that are descendants of '/1/2/'
	conn := db.Where("[path].IsDescendantOf(?)=1", HierarchyId{Data: []int64{1, 2}}).Find(&elements)
	if conn.Error != nil {
		t.Fatal("Failed to query database", conn.Error)
	}

	if len(elements) != 3 {
		t.Fatal("Expected 3 elements, got", len(elements))
	}

	// Get all elements on the root that have a child with name that contains 'de'
	root := GetRoot()

	sub := db.Table("test_parents_tables AS b").Select("COUNT(*)").Where("[b].[path].IsDescendantOf([a].[path])=1 AND [b].[name] LIKE '%de%'")

	conn = db.Table("test_parents_tables AS a").
		Where("[a].[path].GetLevel()=? AND [a].[path].IsDescendantOf(?)=1 AND (?)>0", root.GetLevel()+1, root, sub).
		Find(&elements)
	if conn.Error != nil {
		t.Fatal("Failed to query database", conn.Error)
	}

	if len(elements) != 2 {
		t.Fatal("Expected 2 elements, got", len(elements))
	}
}
