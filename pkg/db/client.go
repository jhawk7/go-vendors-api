package db

import (
	"errors"
	"fmt"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Vendor struct {
	gorm.Model
	Name  string
	Phone string
	Email string
	Cost  string
	Desc  string
}

type Vendors []Vendor

type UpdateRequest struct {
	Name   string
	Fields struct {
		Phone string `json:",omitempty"`
		Email string `json:",omitempty"`
		Cost  string `json:",omitempty"`
		Desc  string `json:",omitempty"`
	}
}

type DeleteRequest struct {
	Name string
}

type DBClient struct {
	svc *gorm.DB
}

func InitDB() (client DBClient, err error) {
	db, dbErr := gorm.Open(sqlite.Open("vendorsDB.db"), &gorm.Config{})
	if dbErr != nil {
		err = fmt.Errorf("failed to initialize db; [error: %v]", dbErr)
		return
	}

	db.AutoMigrate(&Vendor{})
	client.svc = db
	return
}

func (client *DBClient) GetActiveVendors() (vendors Vendors, err error) {
	if result := client.svc.Find(&vendors); result.Error != nil {
		err = fmt.Errorf("failed to retrieve all vendors from db; [error: %v]", result.Error)
		return
	}
	return
}

func (client *DBClient) GetVendorByName(name string) (vendor Vendor, err error, notFound bool) {
	if result := client.svc.Find(&vendor, "name = ?", name); result.Error != nil {
		notFound = errors.Is(result.Error, gorm.ErrRecordNotFound)
		err = fmt.Errorf("failed to retrieve vendor %v by name; [error: %v]", name, result.Error)
	}
	return
}

func (client *DBClient) CreateVendor(vendor *Vendor) (err error) {
	if result := client.svc.Create(&vendor); result.Error != nil {
		err = fmt.Errorf("failed to create vendor; [error: %v]", result.Error)
	}
	return
}

func (client *DBClient) UpdateVendor(update UpdateRequest) (err error, notFound bool) {
	var vendor Vendor
	if result := client.svc.Where("name = ?", update.Name).First(&vendor); result.Error != nil {
		notFound = errors.Is(result.Error, gorm.ErrRecordNotFound)
		err = fmt.Errorf("failed to retrieve vendor for update; [error: %v]", result.Error)
		return
	}

	if update.Fields.Cost != "" {
		vendor.Cost = update.Fields.Cost
	}

	if update.Fields.Email != "" {
		vendor.Email = update.Fields.Email
	}

	if update.Fields.Phone != "" {
		vendor.Phone = update.Fields.Phone
	}

	if update.Fields.Desc != "" {
		vendor.Desc = update.Fields.Desc
	}

	client.svc.Save(&vendor)
	return
}

func (client *DBClient) DeleteVendor(name string) {
	var vendor Vendor
	client.svc.Where("name = ?", name).Delete(&vendor) //soft delete
	return
}
