package dto

import "time"

type Genre struct {
	ID        int        `json:"publisher_id"`
	Genre     string     `json:"genre"`
	Active    bool       `json:"active"`
	CreatedBy string     `json:"created_by"`
	UpdatedBy string     `json:"updated_by,omitempty"`
	Created   time.Time  `json:"created"`
	Updated   *time.Time `json:"updated,omitempty"`
}
