package dbservice

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"misinfodetector-backend/models"
	"strconv"
	"sync"
	"time"
)

func (service *DbService) ImportTopicActivities(f io.Reader) error {
	if err := service.deleteAllTopicActivities(); err != nil {
		return fmt.Errorf("error deleting topic activities: %v", err)
	}

	log.Printf("import -> start")
	var wg sync.WaitGroup

	i := 0
	r := csv.NewReader(f)
	for {
		record, err := r.Read()
		if err != nil {
			i++
			if err == io.EOF {
				break
			} else {
				return err
			}
		}

		wg.Go(func() {
			service.importTopicActivityRecord(record, i)
		})
		i++
	}

	wg.Wait()
	log.Printf("import -> done")
	return nil
}

func (service *DbService) deleteAllTopicActivities() error {
	_, err := service.db.Exec("delete from topic_activities")
	return err
}

func (service *DbService) importTopicActivityRecord(record []string, idx int) {
	service.dbmut.Lock()
	defer service.dbmut.Unlock()

	parsedRecord, err := parseTopicActivityRecord(record)
	if err != nil {
		log.Printf("unable to parse record on line %b: %v", idx, err)
		return
	}

	stmt, err := service.db.Prepare("insert into topic_activities(date, contents, topic_id, topic_name) values (?, ?, ?, ?)")
	if err != nil {
		log.Printf("error preparing sql statement on line %b: %v", idx, err)
		return
	}
	defer stmt.Close()

	if _, err := stmt.Exec(parsedRecord.DateUtc.Format(time.RFC3339), parsedRecord.Text, parsedRecord.TopicId, parsedRecord.TopicName); err != nil {
		log.Printf("error executing sql statement on line %b: %v", idx, err)
		return
	}
}

func (service *DbService) GetAllTopicActivityRecords() ([]*models.TopicActivityRecord, error) {
	service.dbmut.Lock()
	defer service.dbmut.Unlock()

	rows, err := service.db.Query("select record_id, date, contents, topic_id, topic_name from topic_activities")
	if err != nil {
		return nil, fmt.Errorf("error querying topic activity records: %v", err)
	}

	records := make([]*models.TopicActivityRecord, 0)
	for rows.Next() {
		var record models.TopicActivityRecord

		var recordDateRaw string
		if err := rows.Scan(&record.RecordId, &recordDateRaw, &record.Text, &record.TopicId, &record.TopicName); err != nil {
			return nil, fmt.Errorf("error scanning topic activity records: %v", err)
		}

		record.DateUtc, err = time.Parse(time.RFC3339, recordDateRaw)
		if err != nil {
			return nil, fmt.Errorf("error parsing date: %v", err)
		}

		records = append(records, &record)
	}

	return records, nil
}

func parseTopicActivityRecord(record []string) (*models.TopicActivityRecord, error) {
	date, err := time.Parse("2006-01-02", record[0])
	if err != nil {
		return nil, fmt.Errorf("error parsing record date: %v", err)
	}

	text := record[1]

	topicId, err := strconv.ParseInt(record[2], 10, 64)
	if err != nil {
		return nil, fmt.Errorf("error parsing topic id: %v", err)
	}

	topicName := record[3]

	return &models.TopicActivityRecord{
		DateUtc:   date.UTC(),
		Text:      text,
		TopicId:   topicId,
		TopicName: topicName,
	}, nil
}
