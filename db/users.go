package db

import "github.com/go-ozzo/ozzo-dbx"

type User struct {
	ID          uint64 `db:"id"`
	Email       string `db:"email"`
	Password    string `db:"password"`
	Name        string `db:"name"`
	Phone       string `db:"phone"`
	DateOfBirth string `db:"date_of_birth"`
}

func (u User) TableName() string {
	return "users"
}

func (d *DB) GetUser(email string) (*User, error) {
	user := &User{}
	err := d.db.Select().Where(dbx.HashExp{"email": email}).One(user)
	return user, err
}

func (d *DB) GetUserByID(id uint64) (*User, error) {
	var user User
	err := d.db.Select().Model(id, &user)
	return &user, err
}

func (d *DB) CreateUser(user *User) error {
	return d.db.Model(user).Insert()
}

func (d *DB) SetUserNewPassword(user *User) error {
	params := dbx.Params{"password": user.Password}
	expression := dbx.HashExp{"id": user.ID}
	_, err := d.db.Update("users", params, expression).Execute()
	return err
}
