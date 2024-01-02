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

type Vendors []Vendor

type UpdateRequest struct {
	Name         string
	UpdateFields struct {
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
	handlers.LogInfo("db initialized")
	client.svc = db
	return
}

func (client *DBClient) GetActiveVendors() (vendors Vendors, err error) {
	if result := client.svc.Find(&vendors); result.Error != nil {
		err = fmt.Errorf("failed to retrieve active vendors from db; [error: %v]", result.Error)
		return
	}

	handlers.LogInfo("retreived all active vendors")
	return
}

func (client *DBClient) GetAllVendors() (vendors Vendors, err error) {
	if result := client.svc.Where("deleted_At <> ?", "null").Find(&vendors); result.Error != nil {
		err = fmt.Errorf("failed to retrieve all vendors from db; [error: %v]", result.Error)
		return
	}

	handlers.LogInfo("retrieved all vendors")
	return
}

func (client *DBClient) GetVendorByName(name string) (vendor Vendor, err error, notFound bool) {
	if result := client.svc.Find(&vendor, "name = ?", name); result.Error != nil {
		notFound = errors.Is(result.Error, gorm.ErrRecordNotFound)
		err = fmt.Errorf("failed to retrieve vendor %v by name; [error: %v]", name, result.Error)
	}

	handlers.LogInfo(fmt.Sprintf("retrieved vendor %v", name))
	return
}

func (client *DBClient) CreateVendor(vendor *Vendor) (err error) {
	if result := client.svc.Create(&vendor); result.Error != nil {
		err = fmt.Errorf("failed to create vendor; [error: %v]", result.Error)
	}

	handlers.LogInfo(fmt.Sprintf("created vendor %v", vendor.Name))
	return
}

func (client *DBClient) UpdateVendor(update *UpdateRequest) (vendor Vendor, err error, notFound bool) {
	if result := client.svc.Where("name = ?", update.Name).First(&vendor); result.Error != nil {
		notFound = errors.Is(result.Error, gorm.ErrRecordNotFound)
		err = fmt.Errorf("failed to retrieve vendor for update; [error: %v]", result.Error)
		return
	}

	if update.UpdateFields.Cost != "" {
		vendor.Cost = update.UpdateFields.Cost
	}

	if update.UpdateFields.Email != "" {
		vendor.Email = update.UpdateFields.Email
	}

	if update.UpdateFields.Phone != "" {
		vendor.Phone = update.UpdateFields.Phone
	}

	if update.UpdateFields.Desc != "" {
		vendor.Desc = update.UpdateFields.Desc
	}

	handlers.LogInfo(fmt.Sprintf("updated vendor: %v", vendor))
	client.svc.Save(&vendor)
	return
}

func (client *DBClient) DeleteVendor(name string) {
	vendor := new(Vendor)
	client.svc.Where("name = ?", name).Delete(&vendor) //soft delete
	handlers.LogInfo(fmt.Sprintf("deleted vendor: %v", vendor))
}
