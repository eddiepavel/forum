package database

import (
	"database/sql"
	"fmt"
	"forum-app/helpers"
	"forum-app/models"
	"time"
)

func (db *Connection) SessionExistsDB(userId int) (*models.Session, bool, error) {
	query := `SELECT * FROM session WHERE userId = ? LIMIT 1;`
	var session models.Session

	err := db.DB.QueryRow(query, userId).Scan(&session.ID, &session.Token, &session.ExpiresAt, &session.UserId)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, false, nil
		}

		fmt.Printf("Error checking session existence: %v\n", err)
		return nil, false, err
	}

	return &session, true, err
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

func (db *Connection) GetSession(column string, param any) (*models.Session, error){
	
	query := fmt.Sprintf("SELECT * FROM session WHERE %s = ? LIMIT 1;", column)

	var session models.Session

	err := db.DB.QueryRow(query, param).Scan(&session.ID, &session.Token, &session.ExpiresAt, &session.UserId)

	if err != nil {
		return nil, err
	}

	return &session, nil
}

func (db *Connection) DeleteSession(sessionId int) error {
	query := `DELETE FROM session WHERE id = ?;`

	_, err := db.DB.Exec(query, sessionId)

	if err != nil {
		fmt.Println("Error deleting session")
	}

	return err
}

func (db *Connection) SessionInit(userId int) (*models.Session, error) {
	session, exists, err := db.SessionExistsDB(userId)
	if err != nil {
		return nil, err
	}

	if exists && !helpers.CompareDatesLess(session.ExpiresAt, time.Now().Format("2006-01-02 15:04:05")) {
		return session, nil
	} else if exists {
		db.DeleteSession(session.ID)
	}

	newTokenId, err := db.CreateSession(userId)
	if err != nil {
		return nil, err
	}

	session, err = db.GetSession("id", newTokenId)

	if err != nil {
		return nil, err
	}

	return session, nil
}
