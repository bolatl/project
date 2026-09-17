package models

import (
	"crypto/sha256"
	"database/sql"
	"encoding/base64"
	"fmt"

	"github.com/bolatl/lenslocked/rand"
)

const (
	MinBytesPerToken = 32
)

type Session struct {
	ID     int
	UserID int
	// here Token is only set when creating a new session. When looking up it is empty
	// since we only store TokenHash in our db and it can't be reversed to raw token
	Token     string
	TokenHash string
}

type SessionService struct {
	DB *sql.DB
	// if less than MinBytesPerToken --> it is automatically set to min
	BytesPerToken int
}

// for creating session
func (ss *SessionService) Create(userID int) (*Session, error) {
	bytesPerToken := ss.BytesPerToken
	if bytesPerToken < MinBytesPerToken {
		bytesPerToken = MinBytesPerToken
	}
	token, err := rand.String(bytesPerToken)
	if err != nil {
		return nil, fmt.Errorf("create: %w", err)
	}
	session := Session{
		// skipping ID so that it's set by DB itself
		UserID:    userID,
		Token:     token,
		TokenHash: ss.hash(token),
	}
	row := ss.DB.QueryRow(`INSERT INTO sessions (user_id, token_hash) 
	VALUES($1, $2) ON CONFLICT (user_id) 
	DO UPDATE 
	SET token_hash = $2 
	RETURNING id;`, session.UserID, session.TokenHash)
	err = row.Scan(&session.ID)

	if err != nil {
		return nil, fmt.Errorf("create func, error returning id of the created session: %w", err)
	}
	return &session, nil
}

// for looking up a User with given token in our db
func (ss *SessionService) User(token string) (*User, error) {
	tokenHash := ss.hash(token)
	var user User
	row := ss.DB.QueryRow(`
		SELECT users.id, 
			users.email, 
			users.password_hash
		FROM sessions 
			JOIN users on users.id = sessions.user_id
		WHERE sessions.token_hash = $1;`, tokenHash)
	err := row.Scan(&user.ID, &user.Email, &user.PasswordHash)
	if err != nil {
		return nil, fmt.Errorf("session/user: %w", err)
	}
	return &user, nil
}

func (ss *SessionService) Delete(token string) error {
	tokenHash := ss.hash(token)
	_, err := ss.DB.Exec(`DELETE FROM sessions WHERE token_hash=$1;`, tokenHash)
	if err != nil {
		return fmt.Errorf("deleting session: %w", err)
	}
	return nil
}

func (ss *SessionService) hash(token string) string {
	tokenHash := sha256.Sum256([]byte(token))
	return base64.URLEncoding.EncodeToString(tokenHash[:])
}
