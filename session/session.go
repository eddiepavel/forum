package session

import (
	"forum-app/app"
	"forum-app/helpers"
	"forum-app/models"
	"time"
)

func SessionInit(app *app.Application, userId int) (models.Session, error) {
	session, exists, err := app.DB.SessionExistsDB(userId)
	if err != nil {
		return models.Session{}, err
	}

	if exists && !helpers.CompareDatesLess(session.ExpiresAt, time.Now().Format("2006-01-02 15:04:05")) {
		return session, nil
	} else if exists {
		app.DB.DeleteSession(session.ID)
	}

	newTokenId, err := app.DB.CreateSession(userId)
	if err != nil {
		return models.Session{}, err
	}

	session, err = app.DB.GetSessionById(newTokenId)
	if err != nil {
		return models.Session{}, err
	}

	return session, nil
}
