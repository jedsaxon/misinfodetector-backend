package dbservice

import (
	"database/sql"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"log"
	"misinfodetector-backend/models"
	"strconv"
	"sync"
	"time"

	"github.com/go-faker/faker/v4"
	"github.com/google/uuid"
	"github.com/huandu/go-sqlbuilder"
)

type (
	DbService struct {
		db    *sql.DB
		dbmut *sync.Mutex
	}
)

// NewDbService creates a new sqlite connection. If it was successful, it will return
// a DbService instance, with a function to close the database connection.
func NewDbService(sqliteDsn string) (*DbService, func() error, error) {
	db, err := sql.Open("sqlite3", sqliteDsn)
	if err != nil {
		return nil, nil, err
	}

	if err := initDb(db); err != nil {
		return nil, nil, err
	}

	return &DbService{
		db:    db,
		dbmut: &sync.Mutex{},
	}, db.Close, nil
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

		var submittedDate string
		err := rows.Scan(&idBytes, &current.Message, &current.Username, &submittedDate, &current.MisinfoState)
		if err != nil {
			return nil, err
		}

		time, err := time.Parse(time.RFC3339, submittedDate)
		if err != nil {
			return nil, err
		}
		current.SubmittedDateUTC = time

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

	var current models.PostModelId
	var submittedDate string
	var idBytes []byte

	err = stmt.QueryRow(id).Scan(&idBytes, &current.Message, &current.Username, &submittedDate, &current.MisinfoState)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("error querying database for post: %v", err)
	}

	t, err := time.Parse(time.RFC3339, submittedDate)
	if err != nil {
		return nil, err
	}
	current.SubmittedDateUTC = t.UTC()

	idUuid, err := uuid.ParseBytes(idBytes)
	if err != nil {
		return nil, fmt.Errorf("unable to create uuid: %v", err)
	}
	current.Id = idUuid

	return &current, nil
}

func (service *DbService) ImportPosts(f io.Reader) error {
	log.Printf("import -> start")
	var wg sync.WaitGroup

	i := 0
	r := csv.NewReader(f)
	for {
		defer func() { i++ }()

		record, err := r.Read()
		if err != nil {
			if err == io.EOF {
				break
			} else {
				return err
			}
		}

		wg.Go(func() {
			service.importPostRecord(record, i)
		})
	}

	wg.Wait()
	return nil
}

// importPostRecord inserts a single post from the python AI predictions
// It expects the following records, in order:
// id,text,label,pred_label,pred_prob,correct
func (service *DbService) importPostRecord(record []string, i int) {
	aiCorrect := record[6]
	if aiCorrect != "True" {
		log.Printf("import -> skipping record on line %b: correct != \"True\"", i)
		return
	}

	randUsername := faker.Name()

	message := record[1]
	rawDate := record[2]
	predictionLabel := record[3]
	_ = record[4]

	predictionFormatted, err := misinfoLabel(predictionLabel)
	if err != nil {
		log.Printf("import -> skipping record on line %b: %v", i, err)
		return
	}

	dateFormatted, err := time.Parse("2006-01-06", rawDate)
	if err != nil {
		log.Printf("import -> skipping record on line %b: bad date: %v", i, err)
		return
	}

	post := models.NewPost(message, randUsername, dateFormatted, predictionFormatted)
	service.InsertPost(post)
}

func misinfoLabel(lbl string) (models.MisinfoState, error) {
	if lbl == "0" {
		return models.MisinfoStateFake, nil
	} else if lbl == "1" {
		return models.MisinfoStateTrue, nil
	}
	return -1, fmt.Errorf("unknown/unsupported misinformation label: %s", lbl)
}

func (service *DbService) InsertPost(p *models.PostModel) (*models.PostModelId, error) {
	service.dbmut.Lock()
	defer service.dbmut.Unlock()

	id, err := uuid.NewRandom()
	if err != nil {
		return nil, fmt.Errorf("unable to generate new id: %v", err)
	}
	stmt, err := service.db.Prepare("insert into posts(id, message, username, date_submitted, misinfo_state_id) values (?, ?, ?, ?, ?)")
	if err != nil {
		return nil, fmt.Errorf("unable to prepare statement: %v", err)
	}

	_, err = stmt.Exec(id.String(), p.Message, p.Username, p.SubmittedDateUTC.Format(time.RFC3339), p.MisinfoState)
	if err != nil {
		return nil, fmt.Errorf("error while executing prepared statement: %v", err)
	}

	return p.WithId(id), nil
}

// UpdatePost will compare old to updated, and update the record in the database with only the changed fields
// This function will use `old.Id` for the update statement. This function also assumes that there is an actual
// change, and will perform the update regardless of whether old and updated are fully equal.
// Returns the amount of records affected, or an error. Will return -1 if an error occurred.
func (service *DbService) UpdatePost(old *models.PostModelId, updated *models.PostModel) (int64, error) {
	sql := sqlbuilder.Update("posts")

	if updated.Message != old.Message {
		sql.Set(sql.Assign("message", updated.Message))
	}
	if updated.Username != old.Username {
		sql.Set(sql.Assign("username", updated.Username))
	}
	if updated.SubmittedDateUTC.Format(time.RFC3339) != old.SubmittedDateUTC.Format(time.RFC3339) {
		sql.Set(sql.Assign("date_submitted", updated.SubmittedDateUTC.Format(time.RFC3339)))
	}
	if updated.MisinfoState != old.MisinfoState {
		sql.Set(sql.Assign("misinfo_state_id", strconv.FormatInt(int64(updated.MisinfoState), 10)))
	}

	sql.Where(sql.Equal("id", old.Id.String()))
	sqlstmt, args := sql.Build()

	service.dbmut.Lock()
	defer service.dbmut.Unlock()

	stmt, err := service.db.Prepare(sqlstmt)
	if err != nil {
		log.Printf("error: %v", err)
		return -1, err
	}
	defer stmt.Close()

	upd, err := stmt.Exec(args...)
	if err != nil {
		return -1, err
	}
	return upd.RowsAffected()
}

func initDb(db *sql.DB) error {
	_, err := db.Exec("create table if not exists misinfo_state(id int primary key, name varchar(64) not null);")
	if err != nil {
		return err
	}

	if err = insertMisinfoState(db, int64(models.MisinfoStateFake), "Fake"); err != nil {
		return err
	}

	if err = insertMisinfoState(db, int64(models.MisinfoStateTrue), "True"); err != nil {
		return err
	}

	if err = insertMisinfoState(db, int64(models.MisinfoStateNotChecked), "Not Checked"); err != nil {
		return err
	}

	_, err = db.Exec("create table if not exists posts(id varchar(36), message text, username text, date_submitted text, misinfo_state_id int references misinfo_state(id));")
	if err != nil {
		return err
	}
	return nil
}

// Inserts a new state if it does not yet exist
func insertMisinfoState(db *sql.DB, id int64, name string) error {
	existsStmt, err := db.Prepare("select count(*) from misinfo_state where id=?")
	if err != nil {
		return err
	}
	defer existsStmt.Close()
	res := existsStmt.QueryRow(id)
	if res.Err() != nil {
		return res.Err()
	}
	var count int64
	res.Scan(&count)
	if count > 0 {
		// record exists
		log.Printf("inserting misinfo state: misinfo state with id %b already exists, skipping", id)
		return nil
	}

	stmt, err := db.Prepare("insert into misinfo_state values(?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(id, name)
	return err
}
