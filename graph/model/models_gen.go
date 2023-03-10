// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package model

type NewVendor struct {
	Name  string  `json:"name"`
	Phone *string `json:"phone"`
	Email *string `json:"email"`
	Cost  *string `json:"cost"`
	Desc  *string `json:"desc"`
}

type UpdateFields struct {
	Phone *string `json:"phone"`
	Email *string `json:"email"`
	Cost  *string `json:"cost"`
	Desc  *string `json:"desc"`
}

type UpdateVendor struct {
	Name         string        `json:"name"`
	UpdateFields *UpdateFields `json:"updateFields"`
}

type Vendor struct {
	ID        *int    `json:"id"`
	Name      string  `json:"name"`
	Phone     *string `json:"phone"`
	Email     *string `json:"email"`
	Cost      *string `json:"cost"`
	Desc      *string `json:"desc"`
	CreatedAt *string `json:"createdAt"`
	UpdatedAt *string `json:"updatedAt"`
	DeletedAt *string `json:"deletedAt"`
}
