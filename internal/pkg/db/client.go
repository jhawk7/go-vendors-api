package db

import (
	"errors"
	"fmt"

	"github.com/jhawk7/go-vendors-api/internal/pkg/handlers"
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

//type Vendors []Vendor //sqlite needs pointer to struct or slice

type UpdateRequest struct {
	Name    string
	Changes struct {
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

func InitDB() (client *DBClient, err error) {
	client = new(DBClient)
	db, dbErr := gorm.Open(sqlite.Open("vendorsDB.db"), &gorm.Config{})
	if dbErr != nil {
		err = fmt.Errorf("failed to initialize db; [error: %v]", dbErr)
		return
	}

	db.AutoMigrate(&Vendor{})
	handlers.LogInfo("db initialized")
	client.svc = db
	return
}

func (client *DBClient) GetActiveVendors() (vendors *[]Vendor, err error) {
	if result := client.svc.Find(&vendors); result.Error != nil {
		err = fmt.Errorf("failed to retrieve active vendors from db; [error: %v]", result.Error)
		return
	}

	handlers.LogInfo("retreived all active vendors")
	return
}

func (client *DBClient) GetAllVendors() (vendors *[]Vendor, err error) {
	if result := client.svc.Where("deleted_At <> ?", "null").Find(&vendors); result.Error != nil {
		err = fmt.Errorf("failed to retrieve all vendors from db; [error: %v]", result.Error)
		return
	}

	handlers.LogInfo("retrieved all vendors")
	return
}

func (client *DBClient) GetVendorByName(name string) (vendor *Vendor, notFound bool, err error) {
	result := client.svc.Find(&vendor, "name = ?", name)
	if result.Error != nil {
		notFound = errors.Is(result.Error, gorm.ErrRecordNotFound)
		err = fmt.Errorf("failed to retrieve vendor %v by name; [error: %v]", name, result.Error)
		return
	}

	if vendor.Name == "" {
		notFound = true
		err = fmt.Errorf("vendor %v not found", name)
		return
	}

	handlers.LogInfo(fmt.Sprintf("retrieved vendor %v; %v", vendor.Name, result))
	return
}

func (client *DBClient) CreateVendor(vendor *Vendor) (err error) {
	if result := client.svc.Create(&vendor); result.Error != nil {
		err = fmt.Errorf("failed to create vendor; [error: %v]", result.Error)
	}

	handlers.LogInfo(fmt.Sprintf("created vendor %v", vendor.Name))
	return
}

func (client *DBClient) UpdateVendor(update *UpdateRequest) (vendor *Vendor, notFound bool, err error) {
	if result := client.svc.Where("name = ?", update.Name).First(&vendor); result.Error != nil {
		notFound = errors.Is(result.Error, gorm.ErrRecordNotFound)
		err = fmt.Errorf("failed to retrieve vendor for update; [error: %v]", result.Error)
		return
	}

	if update.Changes.Cost != "" {
		vendor.Cost = update.Changes.Cost
	}

	if update.Changes.Email != "" {
		vendor.Email = update.Changes.Email
	}

	if update.Changes.Phone != "" {
		vendor.Phone = update.Changes.Phone
	}

	if update.Changes.Desc != "" {
		vendor.Desc = update.Changes.Desc
	}

	handlers.LogInfo(fmt.Sprintf("updated vendor: %v", vendor))
	client.svc.Save(&vendor)
	return
}

func (client *DBClient) DeleteVendor(name string) (err error) {
	vendor := new(Vendor)
	if result := client.svc.Where("name = ?", name).Delete(&vendor); result.Error != nil {
		err = fmt.Errorf("failed to delete vendor %v; %v", name, result.Error)
		return
	} //soft delete
	handlers.LogInfo(fmt.Sprintf("deleted vendor: %v", vendor))
	return
}
