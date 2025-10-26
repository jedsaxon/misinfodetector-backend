package dbservice

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type (
	DbService struct {
		db *sql.DB
	}

	PostModel struct {
		Id                     uuid.UUID
		Message                string
		Username               string
		SubmittedDate          time.Time
		ContainsMisinformation bool
	}
)

func NewDbService(db *sql.DB) *DbService {
	return &DbService{
		db: db,
	}
}

func (dbservice *DbService) GetPosts(pageNumber int64, resultAmount int64) ([]PostModel, error) {
	stmt, err := dbservice.db.Prepare("select * from posts order by date(date_submitted) limit ? offset ?")
	if err != nil {
		return nil, fmt.Errorf("unable to prepare statement: %v", err)
	}
	defer stmt.Close()

	rows, err := stmt.Query(pageNumber*resultAmount, resultAmount)
	if err != nil {
		return nil, fmt.Errorf("unable to execute prepared statement: %v", err)
	}

	response := make([]PostModel, 0)
	for rows.Next() {
		var current PostModel
		var idBytes []byte
		rows.Scan(&idBytes, &current.Message, &current.Username, &current.SubmittedDate, &current.ContainsMisinformation)
		idUuid, err := uuid.FromBytes(idBytes)
		if err != nil {
			return nil, fmt.Errorf("unable to create uuid: %v", err)
		}
		current.Id = idUuid
		response = append(response, current)
	}

	return response, nil
}
