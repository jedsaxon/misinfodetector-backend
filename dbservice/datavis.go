package dbservice

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"misinfodetector-backend/models"
	"strconv"
	"sync"
)

func (service *DbService) ImportTnseEmbeddings(f io.Reader) error {
	if err := service.deleteAllTnseEmbeddings(); err != nil {
		return fmt.Errorf("error deleting tnse embeddings: %v", err)
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
			service.importTnseEmbeddingRecord(record, i)
		})
		i++
	}

	wg.Wait()
	log.Printf("import -> done")
	return nil
}

func (service *DbService) deleteAllTnseEmbeddings() error {
	_, err := service.db.Exec("delete from tnse_embeddings")
	return err
}

func (service *DbService) importTnseEmbeddingRecord(record []string, idx int) {
	service.dbmut.Lock()
	defer service.dbmut.Unlock()

	parsedRecord, err := parseTnseRecord(record)
	if err != nil {
		log.Printf("unable to parse record on line %b: %v", idx, err)
		return
	}

	stmt, err := service.db.Prepare("insert into tnse_embeddings(record_id, label, pred_label, correct, tnse_x, tnse_y) values (?, ?, ?, ?, ?, ?)")
	if err != nil {
		log.Printf("error preparing sql statement on line %b: %v", idx, err)
		return
	}
	defer stmt.Close()

	if _, err := stmt.Exec(parsedRecord.RecordId, parsedRecord.Label, parsedRecord.PredictionLabel, parsedRecord.Correct, parsedRecord.TnseY, parsedRecord.TnseY); err != nil {
		log.Printf("error executing sql statement on line %b: %v", idx, err)
		return
	}
}

func (service *DbService) GetAllTnseEmbeddings() ([]*models.TnseEmbeddingRecord, error) {
	service.dbmut.Lock()
	defer service.dbmut.Unlock()

	rows, err := service.db.Query("select record_id, label, pred_label, correct, tnse_x, tnse_y from tnse_embeddings")
	if err != nil {
		return nil, fmt.Errorf("error querying tnse embeddings: %v", err)
	}

	var records []*models.TnseEmbeddingRecord
	for rows.Next() {
		var record models.TnseEmbeddingRecord

		if err := rows.Scan(&record.RecordId, &record.Label, &record.PredictionLabel, &record.Correct, &record.TnseX, &record.TnseY); err != nil {
			return nil, fmt.Errorf("error scanning tnse embedding records: %v", err)
		}

		records = append(records, &record)
	}

	return records, nil
}

func parseTnseRecord(record []string) (*models.TnseEmbeddingRecord, error) {
	id, err := strconv.ParseInt(record[0], 10, 64)
	if err != nil {
		return nil, fmt.Errorf("unable to parse id: %v", err)
	}

	label, err := strconv.ParseInt(record[1], 10, 64)
	if err != nil {
		return nil, fmt.Errorf("unable to parse label: %v", err)
	}

	predLabel, err := strconv.ParseInt(record[2], 10, 64)
	if err != nil {
		return nil, fmt.Errorf("unalbe to parse prediction label: %v", err)
	}

	correct := record[3]

	tnseX, err := strconv.ParseFloat(record[4], 64)
	if err != nil {
		return nil, fmt.Errorf("unable to parse tnse_x: %v", err)
	}

	tnseY, err := strconv.ParseFloat(record[5], 64)
	if err != nil {
		return nil, fmt.Errorf("unable to parse tnse_y: %v", err)
	}

	return &models.TnseEmbeddingRecord{
		RecordId:        id,
		Label:           label,
		PredictionLabel: predLabel,
		Correct:         correct,
		TnseX:           tnseX,
		TnseY:           tnseY,
	}, nil
}
