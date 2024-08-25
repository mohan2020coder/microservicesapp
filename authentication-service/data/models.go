package data

import (
	"errors"
	"log"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

//const dbTimeout = time.Second * 3

var db *gorm.DB

// New is the function used to create an instance of the data package. It returns the type
// Model, which embeds all the types we want to be available to our application.
func New(dbPool *gorm.DB) Models {
	db = dbPool

	return Models{
		User: User{},
	}
}

// Models is the type for this package. Note that any model that is included as a member
// in this type is available to us throughout the application, anywhere that the
// app variable is used, provided that the model is also added in the New function.
type Models struct {
	User User
}

// User is the structure which holds one user from the database.
type User struct {
	ID        int       `gorm:"primaryKey" json:"id"`
	Email     string    `gorm:"unique;not null" json:"email"`
	FirstName string    `json:"first_name,omitempty"`
	LastName  string    `json:"last_name,omitempty"`
	Password  string    `json:"-"`
	Active    bool      `json:"active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// GetAll returns a slice of all users, sorted by last name
func (u *User) GetAll() ([]*User, error) {
	var users []*User
	if err := db.Order("last_name").Find(&users).Error; err != nil {
		log.Println("Error querying users:", err)
		return nil, err
	}
	return users, nil
}

// GetByEmail returns one user by email
func (u *User) GetByEmail(email string) (*User, error) {
	var user User
	if err := db.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// GetOne returns one user by id
func (u *User) GetOne(id int) (*User, error) {
	var user User
	if err := db.First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// Update updates one user in the database, using the information stored in the receiver u
func (u *User) Update() error {
	if err := db.Save(u).Error; err != nil {
		return err
	}
	return nil
}

// Delete deletes one user from the database, by User.ID
func (u *User) Delete() error {
	if err := db.Delete(u).Error; err != nil {
		return err
	}
	return nil
}

// Insert inserts a new user into the database, and returns the ID of the newly inserted row
func (u *User) Insert(user User) (int, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), 12)
	if err != nil {
		return 0, err
	}
	user.Password = string(hashedPassword)
	if err := db.Create(&user).Error; err != nil {
		return 0, err
	}
	return user.ID, nil
}

// ResetPassword is the method we will use to change a user's password.
func (u *User) ResetPassword(password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return err
	}
	if err := db.Model(u).Update("password", string(hashedPassword)).Error; err != nil {
		return err
	}
	return nil
}

// PasswordMatches uses Go's bcrypt package to compare a user supplied password
// with the hash we have stored for a given user in the database. If the password
// and hash match, we return true; otherwise, we return false.
func (u *User) PasswordMatches(plainText string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(plainText))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}
