package ds

import (
	// "gorm.io/gorm"
	"time"
)

const DRAFT string = "черновик"
const FORMED string = "сформирован"
const COMPELTED string = "завершён"
const REJECTED string = "отклонён"
const DELETED string = "удалён"

type User struct {
	UUID      string `gorm:"type:uuid;primary_key;default:gen_random_uuid()"  json:"-"`
	Login     string `gorm:"size:30;not null"  json:"-"`
	Password  string `gorm:"size:30;not null"  json:"-"`
	Name      string `gorm:"size:50;not null"  json:"name"`
	Moderator bool   `gorm:"not null"  json:"-"`
}

type Medicine struct {
	UUID      		string  `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"uuid" binding:"-"`
	Name       		string  `gorm:"size:100;not null" form:"name" json:"name" binding:"required"`
	ImageURL  		*string `gorm:"size:100" json:"image_url" binding:"-"`
	Dosage     		string  `gorm:"size:75;not null" form:"dosage" json:"dosage" binding:"required"`
	Amount       	int     `gorm:"not null" json:"amount" form:"amount" binding:"required"`
	Manufacturer    string  `gorm:"size:100;not null" form:"manufacturer" json:"manufacturer" binding:"required"`
	IsDeleted 		bool    `gorm:"not null;default:false" json:"-" binding:"-"`
}

type Component struct {
	UUID             string     `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	Status           string     `gorm:"size:20;not null"`
	CreationDate     time.Time  `gorm:"not null;type:timestamp"`
	FormationDate    *time.Time `gorm:"type:timestamp"`
	CompletionDate   *time.Time `gorm:"type:timestamp"`
	ModeratorId      *string    `json:"-"`
	CustomerId       string     `gorm:"not null"`
	ComponentName string     	`gorm:"size:50;not null"`

	Moderator *User
	Customer  User
}

type MedicineProduction struct {
	MedicineId    string `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"medicine_id"`
	ComponentId   string `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"component_id"`

	Medicine    *Medicine    `gorm:"foreignKey:MedicineId" json:"medicine"`
	Component 	*Component	 `gorm:"foreignKey:ComponentId" json:"component"`
}
