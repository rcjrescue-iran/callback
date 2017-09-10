package main

import (
	"log"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/tucnak/telebot"
)

var db *gorm.DB

const (
	USER_REPORT = iota
	USER_RESUME
	USER_DEFAULT
)

func init() {
	var err error
	db, err = gorm.Open("sqlite3", "database.sqlite3")
	if err != nil {
		log.Fatal(err)
	}

	user := &User{}
	if !db.HasTable(user) {
		db.CreateTable(user)
	}

	db.AutoMigrate()
}

// Model for each models
type Model struct {
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `sql:"index"`
}

// User table
type User struct {
	telebot.User
	Model
	Status int
}

// Create , create new user
func (u *User) Create() {
	db.Model(u).Create(u)
}

// Save , save user
func (u *User) Save() {
	db.Model(u).Save(u)
}

// GetUser for fetching user info from database
func GetUser(user telebot.User) *User {
	u := &User{}
	db.Model(u).Find(u, "id = ?", user.ID)
	return u
}

// Exists , ...
func (u *User) Exists() bool {
	return u.ID > 0
}
