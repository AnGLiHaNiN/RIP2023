package ds

import (
	"R_I_P_labs/internal/app/role"
	"time"
)

const DRAFT string = "черновик"
const FORMED string = "сформирован"
const COMPELTED string = "завершён"
const REJECTED string = "отклонён"
const DELETED string = "удалён"

type User struct {
	UUID     string    `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"-"`
	Role     role.Role `sql:"type:string;"`
	Login    string    `gorm:"size:30;not null" json:"login"`
	Password string    `gorm:"size:40;not null" json:"-"`
	Name     *string   `gorm:"size:60" json:"name"`
	Email    *string   `gorm:"size:40" json:"email"`
	// The SHA-1 hash is 20 bytes. When encoded in hexadecimal, each byte is represented by two characters. Therefore, the resulting hash string will be 40 characters long
}

type Component struct {
	UUID       string  `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"uuid" binding:"-"`
	Name       string  `gorm:"size:100;not null" form:"name" json:"name" binding:"required"`
	ImageURL   *string `gorm:"size:100" json:"image_url" binding:"-"`
	WorldName  string  `gorm:"size:75;not null" form:"world_name" json:"world_name" binding:"required"`
	Amount     int     `gorm:"not null" json:"amount" form:"amount" binding:"required"`
	Properties string  `gorm:"size:200;not null" form:"properties" json:"properties" binding:"required"`
	IsDeleted  bool    `gorm:"not null;default:false" json:"-" binding:"-"`
}

type Medicine struct {
	UUID               string     `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	Status             string     `gorm:"size:20;not null"`
	CreationDate       time.Time  `gorm:"not null;type:timestamp"`
	FormationDate      *time.Time `gorm:"type:timestamp"`
	CompletionDate     *time.Time `gorm:"type:timestamp"`
	ModeratorId        *string    `json:"-"`
	CustomerId         string     `gorm:"not null"`
	Name               *string    `gorm:"size:50"`
	VerificationStatus *string    `gorm:"size:40"`

	Moderator *User
	Customer  User
}

const VerificationCompleted string = "пройдена"
const VerificationFailed string = "провалена"
const VerificationStarted string = "находится в проверке"

type MedicineProduction struct {
	MedicineId  string `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"medicine_id"`
	ComponentId string `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"component_id"`
	Count       int    `gorm:"not null;default:1"`

	Medicine  *Medicine  `gorm:"foreignKey:MedicineId" json:"medicine"`
	Component *Component `gorm:"foreignKey:ComponentId" json:"component"`
}
