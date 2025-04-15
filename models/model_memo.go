package models

import (

	"github.com/google/uuid"
)

type Memo struct {
	ID          uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Date      	string 	  `gorm:"type:date" json:"date"`
	Note	 	string    `json:"note"`
	UserID      uuid.UUID `gorm:"type:uuid;not null"`
	User        User      `gorm:"foreignKey:UserID;references:ID"`
}