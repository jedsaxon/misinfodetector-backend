package dbservice

import (
	"database/sql"
	"fmt"
	"misinfodetector-backend/models"

	"github.com/google/uuid"
)

type (
	DbService struct {
		db *sql.DB
	}
)

func NewDbService(db *sql.DB) *DbService {
	return &DbService{
		db: db,
	}
}

func (dbservice *DbService) GetPosts(pageNumber int64, resultAmount int64) ([]models.PostModelId, error) {
	stmt, err := dbservice.db.Prepare("select * from posts order by date(date_submitted) limit ? offset ?")
	if err != nil {
		return nil, fmt.Errorf("unable to prepare statement: %v", err)
	}
	defer stmt.Close()

	rows, err := stmt.Query(resultAmount, (pageNumber-1)*resultAmount)
	if err != nil {
		return nil, fmt.Errorf("unable to execute prepared statement: %v", err)
	}

	response := make([]models.PostModelId, 0)
	for rows.Next() {
		var current models.PostModelId
		var idBytes []byte
		rows.Scan(&idBytes, &current.Message, &current.Username, &current.SubmittedDate, &current.ContainsMisinformation)
		idUuid, err := uuid.ParseBytes(idBytes)
		if err != nil {
			return nil, fmt.Errorf("unable to create uuid: %v", err)
		}
		current.Id = idUuid
		response = append(response, current)
	}

	return response, nil
}

func (service *DbService) InsertPost(p *models.PostModel) (*models.PostModelId, error) {
	// _, err := db.Exec("create table if not exists posts(id varchar(36), message text, username text, date_submitted text, is_misinformation int);")
	id, err := uuid.NewRandom()
	if err != nil {
		return nil, fmt.Errorf("unable to generate new id: %v", err)
	}
	stmt, err := service.db.Prepare("insert into posts(id, message, username, date_submitted, is_misinformation) values (?, ?, ?, ?, ?)")
	if err != nil {
		return nil, fmt.Errorf("unable to prepare statement: %v", err)
	}

	_, err = stmt.Exec(id, p.Message, p.Username, p.SubmittedDate, p.ContainsMisinformation)
	if err != nil {
		return nil, fmt.Errorf("error while executing prepared statement: %v", err)
	}

	return p.WithId(id), nil
}
