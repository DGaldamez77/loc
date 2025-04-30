package dto

import "time"

type Publisher struct {
	ID        int        `json:"publisher_id"`
	Publisher string     `json:"publisher"`
	Active    bool       `json:"active"`
	CreatedBy string     `json:"created_by"`
	UpdatedBy string     `json:"updated_by,omitempty"`
	Created   time.Time  `json:"created"`
	Updated   *time.Time `json:"updated,omitempty"`
}
