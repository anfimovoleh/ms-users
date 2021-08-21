package db

import "time"

type Token struct {
	Token      string    `db:"pk,token"`
	UserID     uint64    `db:"user_id"`
	LastSentAt time.Time `db:"last_sent_at"`
}

func (t Token) TableName() string {
	return "tokens"
}

func (d *DB) CreateToken(token *Token) error {
	return d.db.Model(token).Insert()
}

func (d *DB) GetUserByToken(tokenID string) (*Token, error) {
	var token Token
	err := d.db.Select().Model(tokenID, &token)
	return &token, err
}

func (d *DB) DeleteToken(tokenID string) error {
	token := &Token{Token: tokenID}
	return d.db.Model(token).Delete()
}
