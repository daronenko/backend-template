package model

import (
	"strings"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID        uuid.UUID `json:"id" db:"id" redis:"id" validate:"omitempty"`
	Username  string    `json:"username" db:"username" redis:"username" validate:"required,lte=30"`
	Email     string    `json:"email,omitempty" db:"email" redis:"email" validate:"omitempty,lte=60,email"`
	Password  string    `json:"password,omitempty" db:"password" redis:"password" validate:"omitempty,required,gte=6"`
	Role      *string   `json:"role,omitempty" db:"role" redis:"role" validate:"omitempty,lte=10"`
	Avatar    *string   `json:"avatar,omitempty" db:"avatar" redis:"avatar" validate:"omitempty,lte=512,url"`
	CreatedAt time.Time `json:"created_at,omitempty" db:"created_at" redis:"created_at"`
	UpdatedAt time.Time `json:"updated_at,omitempty" db:"updated_at" redis:"updated_at"`
	LoginDate time.Time `json:"login_date" db:"login_date" redis:"login_date"`
}

type UserWithToken struct {
	User  *User  `json:"user"`
	Token string `json:"token"`
}

type UsersList struct {
	TotalCount int     `json:"total_count"`
	TotalPages int     `json:"total_pages"`
	Page       int     `json:"page"`
	Size       int     `json:"size"`
	HasMore    bool    `json:"has_more"`
	Users      []*User `json:"users"`
}

func (u *User) HashPassword() error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)
	return nil
}

func (u *User) ComparePasswords(password string) error {
	if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password)); err != nil {
		return err
	}
	return nil
}

func (u *User) SanitizePassword() {
	u.Password = ""
}

func (u *User) PrepareCreate() error {
	u.Email = strings.ToLower(strings.TrimSpace(u.Email))
	u.Password = strings.TrimSpace(u.Password)

	if err := u.HashPassword(); err != nil {
		return err
	}

	if u.Role != nil {
		*u.Role = strings.ToLower(strings.TrimSpace(*u.Role))
	}
	return nil
}

func (u *User) PrepareUpdate() error {
	u.Email = strings.ToLower(strings.TrimSpace(u.Email))

	if u.Role != nil {
		*u.Role = strings.ToLower(strings.TrimSpace(*u.Role))
	}
	return nil
}
