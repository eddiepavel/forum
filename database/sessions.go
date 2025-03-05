package database

import (
	"database/sql"
	"fmt"
	"forum-app/helpers"
	"forum-app/models"
	"time"
)

func (db *Connection) SessionExistsDB(userId int) (models.Session, bool, error) {
	query := `SELECT * FROM session WHERE userId = ? LIMIT 1;`
	var session models.Session

	err := db.DB.QueryRow(query, userId).Scan(&session.ID, &session.Token, &session.ExpiresAt, &session.UserId)

	if err != nil {
		if err == sql.ErrNoRows {
			return session, false, nil
		}

		fmt.Printf("Error checking session existence: %v\n", err)
		return session, false, err
	}

	return session, true, err
}

func (db *Connection) CreateSession(userId int) (int, error) {

	token, _ := helpers.GenerateToken()

	insert := `INSERT INTO session (token, expiresAt, userId) VALUES (?, ?, ?);`

	query, err := db.DB.Exec(insert, token, time.Now().Add(time.Hour*1).Format("2006-01-02 15:04:05"), userId)

	if err != nil {
		fmt.Println(err)
		return 0, err
	}

	lastId, err := query.LastInsertId()

	return int(lastId), err
}

func (db *Connection) GetSessionById(sessionId int) (models.Session, error) {
	query := `SELECT * FROM session WHERE id = ? LIMIT 1;`
	var session models.Session

	err := db.DB.QueryRow(query, sessionId).Scan(&session.ID, &session.Token, &session.ExpiresAt, &session.UserId)

	if err != nil {
		fmt.Println("Error fetching session")
		return session, err
	}

	return session, err
}

func (db *Connection) GetSessionByToken(token string) (models.Session, error) {
	query := `SELECT * FROM session WHERE token = ? LIMIT 1;`
	var session models.Session

	err := db.DB.QueryRow(query, token).Scan(&session.ID, &session.Token, &session.ExpiresAt, &session.UserId)

	if err != nil {
		fmt.Println("Error fetching session")
		return session, err
	}

	return session, err
}

func (db *Connection) GetSessionByUserId(userId int) (models.Session, error) {
	query := `SELECT * FROM session WHERE userId = ? LIMIT 1;`
	var session models.Session

	err := db.DB.QueryRow(query, userId).Scan(&session.ID, &session.Token, &session.ExpiresAt, &session.UserId)

	if err != nil {
		fmt.Println("Error fetching session")
		return session, err
	}

	return session, err
}

func (db *Connection) DeleteSession(sessionId int) error {
	query := `DELETE FROM session WHERE id = ?;`

	_, err := db.DB.Exec(query, sessionId)

	if err != nil {
		fmt.Println("Error deleting session")
	}

	return err
}
