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
