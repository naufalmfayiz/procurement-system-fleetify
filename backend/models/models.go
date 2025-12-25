package models

import (
	"time"

	"gorm.io/gorm"
)

// User model
type User struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Username  string         `gorm:"uniqueIndex;not null;size:100" json:"username"`
	Password  string         `gorm:"not null" json:"-"`
	Role      string         `gorm:"default:user;size:50" json:"role"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// Supplier model
type Supplier struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Name      string         `gorm:"not null;size:200" json:"name"`
	Email     string         `gorm:"size:100" json:"email"`
	Address   string         `gorm:"type:text" json:"address"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// Item model
type Item struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Name      string         `gorm:"not null;size:200" json:"name"`
	Stock     int            `gorm:"not null;default:0" json:"stock"`
	Price     float64        `gorm:"not null;default:0" json:"price"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// Purchasing (Header) model
type Purchasing struct {
	ID                 uint                `gorm:"primaryKey" json:"id"`
	Date               time.Time           `gorm:"not null" json:"date"`
	SupplierID         uint                `gorm:"not null" json:"supplier_id"`
	Supplier           Supplier            `gorm:"foreignKey:SupplierID" json:"supplier,omitempty"`
	UserID             uint                `gorm:"not null" json:"user_id"`
	User               User                `gorm:"foreignKey:UserID" json:"user,omitempty"`
	GrandTotal         float64             `gorm:"not null;default:0" json:"grand_total"`
	PurchasingDetails  []PurchasingDetail  `gorm:"foreignKey:PurchasingID" json:"details,omitempty"`
	CreatedAt          time.Time           `json:"created_at"`
	UpdatedAt          time.Time           `json:"updated_at"`
	DeletedAt          gorm.DeletedAt      `gorm:"index" json:"-"`
}

// PurchasingDetail model
type PurchasingDetail struct {
	ID           uint           `gorm:"primaryKey" json:"id"`
	PurchasingID uint           `gorm:"not null" json:"purchasing_id"`
	ItemID       uint           `gorm:"not null" json:"item_id"`
	Item         Item           `gorm:"foreignKey:ItemID" json:"item,omitempty"`
	Qty          int            `gorm:"not null" json:"qty"`
	SubTotal     float64        `gorm:"not null;default:0" json:"sub_total"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}
