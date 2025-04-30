package dto

import "time"

type Author struct {
	ID        int        `json:"author_id"`
	FirstName string     `json:"first_name"`
	LastName  string     `json:"last_name"`
	Active    bool       `json:"active"`
	CreatedBy string     `json:"created_by"`
	UpdatedBy string     `json:"updated_by,omitempty"`
	Created   time.Time  `json:"created"`
	Updated   *time.Time `json:"updated,omitempty"`
}
