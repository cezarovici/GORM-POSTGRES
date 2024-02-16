package models

import "gorm.io/gorm"

type Users struct {
	gorm.Model
	Name string `json:"name" gorm:"text;not null;default:null`
	Rank string `json:"rank" gorm:"text;not null;default:null`
}
