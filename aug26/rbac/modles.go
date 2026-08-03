package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID                 uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	IdentityProviderID uuid.UUID `gorm:"type:uuid;not null"`
	ExternalID         string    `gorm:"size:255;not null"`
	Username           string    `gorm:"size:100;not null"`
	Email              string    `gorm:"size:255;not null"`
	PasswordHash       *string   `gorm:"size:255"`
	Roles              []Role    `gorm:"many2many:user_roles;"`
	CreatedAt          time.Time
}

type Role struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	Name        string    `gorm:"size:50;not null;unique"`
	Description string
	CreatedAt   time.Time
}

type Entity struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	Name        string    `gorm:"size:100;not null;unique"`
	Description string
}

type EntityInstance struct {
	ID       uuid.UUID       `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	EntityID uuid.UUID       `gorm:"type:uuid;not null"`
	ParentID *uuid.UUID      `gorm:"type:uuid"`
	Parent   *EntityInstance `gorm:"foreignKey:ParentID"`
	Name     string          `gorm:"size:255;not null"`
}

type Permission struct {
	ID     uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	Action string    `gorm:"size:20;not null;unique"` // 'CREATE', 'READ', 'UPDATE', 'DELETE'
}

type RoleEntityPermission struct {
	ID               uuid.UUID  `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	RoleID           uuid.UUID  `gorm:"type:uuid;not null"`
	EntityID         uuid.UUID  `gorm:"type:uuid;not null"`
	PermissionID     uuid.UUID  `gorm:"type:uuid;not null"`
	EntityInstanceID *uuid.UUID `gorm:"type:uuid"`
}
