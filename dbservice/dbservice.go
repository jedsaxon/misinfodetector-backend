package dbservice

import (
	"database/sql"
	"errors"
	"fmt"
	"misinfodetector-backend/models"
	"sync"

	"github.com/google/uuid"
)

type (
	DbService struct {
		db    *sql.DB
		dbmut *sync.Mutex
	}
)

func NewDbService(db *sql.DB) *DbService {
	initDb(db)
	return &DbService{
		db: db,
		dbmut: &sync.Mutex{},
	}
}

// GetPostCount attempts to get the amount of posts in the database. Will
// return -1, and the error if the operation failed. Otherwise, nil error
func (dbservice *DbService) GetPostCount() (int64, error) {
	dbservice.dbmut.Lock()
	defer dbservice.dbmut.Unlock()

	row := dbservice.db.QueryRow("select count(*) from posts")
	var count int64
	if err := row.Scan(&count); err != nil {
		return -1, fmt.Errorf("unable to query post count: %v", err)
	}
	return count, nil
}

func (dbservice *DbService) GetPosts(pageNumber int64, resultAmount int64) ([]models.PostModelId, error) {
	dbservice.dbmut.Lock()
	defer dbservice.dbmut.Unlock()

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

// FindPost finds the post with the given id in the database. If none were found, it returns
// a nil pointer to the post and error. If an error occurred, the post will also be nil, but the
// error will not be.
func (dbs *DbService) FindPost(id string) (*models.PostModelId, error) {
	dbs.dbmut.Lock()
	defer dbs.dbmut.Unlock()

	stmt, err := dbs.db.Prepare("select * from posts where id = ?")
	if err != nil {
		return nil, fmt.Errorf("error preparing statement: %v", err)
	}
	defer stmt.Close()

	row := stmt.QueryRow(id)
	if err := row.Err(); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("error querying database for post: %v", err)
	}

	var current models.PostModelId
	var idBytes []byte
	row.Scan(&idBytes, &current.Message, &current.Username, &current.SubmittedDate, &current.ContainsMisinformation)
	idUuid, err := uuid.ParseBytes(idBytes)
	if err != nil {
		return nil, fmt.Errorf("unable to create uuid: %v", err)
	}
	current.Id = idUuid

	return &current, nil
}

func (service *DbService) InsertPost(p *models.PostModel) (*models.PostModelId, error) {
	service.dbmut.Lock()
	defer service.dbmut.Unlock()

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

func initDb(db *sql.DB) error {
	_, err := db.Exec("create table if not exists posts(id varchar(36), message text, username text, date_submitted text, is_misinformation int);")
	if err != nil {
		return err
	}
	return nil
}
