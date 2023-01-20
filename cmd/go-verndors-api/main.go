package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jhawk7/go-vendors-api/pkg/db"
	log "github.com/sirupsen/logrus"
)

var dbClient *db.DBClient

func main() {
	client, dbErr := db.InitDB()
	if dbErr != nil {
		panic(dbErr)
	}
	dbClient = &client

	r := gin.Default()
	r.GET("/vendors", GetAllVendors)
	r.GET("/vendors/:name", GetVendor)
	r.POST("/vendors", CreateVendor)
	r.PATCH("/vendors", UpdateVendor)
	r.DELETE("/vendors", DeleteVendor)
	r.Run() //runs on port 8080 by default
}

func GetAllVendors(c *gin.Context) {
	vendors, err := dbClient.GetActiveVendors()
	if err != nil {
		ErrorHandler(c, err, http.StatusInternalServerError, false)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": vendors,
	})
}

func GetVendor(c *gin.Context) {
	vendor, err, notFound := dbClient.GetVendorByName(c.Param("name"))
	if err != nil {
		var status int
		if notFound {
			status = http.StatusNotFound
		} else {
			status = http.StatusBadRequest
		}
		ErrorHandler(c, err, status, false)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": vendor,
	})
}

func CreateVendor(c *gin.Context) {
	var vendor db.Vendor
	if bindErr := c.Bind(&vendor); bindErr != nil {
		ErrorHandler(c, bindErr, 0, false)
		return
	}

	if createErr := dbClient.CreateVendor(&vendor); createErr != nil {
		ErrorHandler(c, createErr, http.StatusBadRequest, false)
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"data": vendor,
	})
}

func UpdateVendor(c *gin.Context) {
	var updateRequest db.UpdateRequest
	if bindErr := c.Bind(&updateRequest); bindErr != nil {
		ErrorHandler(c, bindErr, 0, false)
		return
	}

	if updateErr, notFound := dbClient.UpdateVendor(updateRequest); updateErr != nil {
		var status int
		if notFound {
			status = http.StatusNotFound
		} else {
			status = http.StatusBadRequest
		}

		ErrorHandler(c, updateErr, status, false)
		return
	}

	c.JSON(http.StatusAccepted, nil)
}

func DeleteVendor(c *gin.Context) {
	var req db.DeleteRequest
	if bindErr := c.Bind(&req); bindErr != nil {
		ErrorHandler(c, bindErr, 0, false)
		return
	}

	dbClient.DeleteVendor(req.Name)
	c.JSON(http.StatusNoContent, nil)
}

func ErrorHandler(c *gin.Context, err error, status int, fatal bool) {
	if err != nil {
		log.Error(fmt.Errorf("error: %v", err.Error()))

		if fatal {
			panic(err)
		}

		if status != 0 {
			c.JSON(status, gin.H{
				"error": err.Error(),
			})
		}
	}
}
