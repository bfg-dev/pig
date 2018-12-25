package db

import (
	"database/sql"
	"fmt"

	"github.com/lib/pq"
)

// RecManager - record manager
type RecManager struct {
	db               *sql.DB
	tableName        string
	historyTableName string
}

// NewRecManager - creates new DBRecManager
func NewRecManager(db *sql.DB, tableName, historyTableName string) *RecManager {
	return &RecManager{
		db:               db,
		tableName:        tableName,
		historyTableName: historyTableName,
	}
}

// CreateTable - create migration table
func (m *RecManager) CreateTable() error {
	tx, err := m.db.Begin()
	if err != nil {
		return err
	}

	query := fmt.Sprintf(sqlCreateTable, m.tableName)

	if _, err := tx.Exec(query); err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}

// CreateHistoryTable - create migration table
func (m *RecManager) CreateHistoryTable() error {
	tx, err := m.db.Begin()
	if err != nil {
		return err
	}

	query := fmt.Sprintf(sqlCreateHistoryTable, m.historyTableName)

	if _, err := tx.Exec(query); err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}

func (m *RecManager) getShort(where string) ([]*RecShort, error) {
	var result []*RecShort

	query := fmt.Sprintf(sqlSelectShort, m.tableName, where)

	rows, err := m.db.Query(query)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		row := &RecShort{}
		err := rows.Scan(
			&row.ID,
			&row.Name,
			&row.Applied,
			&row.Note,
			&row.GITinfo,
			&row.Filename,
			&row.TStamp,
			(*pq.StringArray)(&row.Requirements),
		)
		if err != nil {
			return nil, err
		}
		result = append(result, row)
	}

	return result, nil
}

func (m *RecManager) getFull(where string) ([]*RecFull, error) {
	var result []*RecFull

	query := fmt.Sprintf(sqlSelectFull, m.tableName, where)

	rows, err := m.db.Query(query)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		row := &RecFull{}
		err := rows.Scan(
			&row.ID,
			&row.Name,
			&row.Applied,
			&row.Note,
			&row.GITinfo,
			&row.Filename,
			&row.TStamp,
			(*pq.StringArray)(&row.Requirements),
			&row.SQLData,
		)
		if err != nil {
			return nil, err
		}
		result = append(result, row)
	}

	return result, nil
}

// GetAllShort get all short db records
func (m *RecManager) GetAllShort() ([]*RecShort, error) {
	return m.getShort("")
}

// GetUnappliedShort get unapplied short db records
func (m *RecManager) GetUnappliedShort() ([]*RecShort, error) {
	return m.getShort("WHERE applied = false")
}

// GetAllFull get all short db records
func (m *RecManager) GetAllFull() ([]*RecFull, error) {
	return m.getFull("")
}

// GetFullFromShort get full from short db records
func (m *RecManager) GetFullFromShort(short *RecShort) (*RecFull, error) {
	result, err := m.getFull(fmt.Sprintf("where id = %v", short.ID))
	if err != nil {
		return nil, err
	}
	if len(result) != 1 {
		return nil, fmt.Errorf("Can not find sql info for id %v", short.ID)
	}
	return result[0], nil
}

// GetUnappliedFull get unapplied short db records
func (m *RecManager) GetUnappliedFull() ([]*RecFull, error) {
	return m.getFull("WHERE applied = false")
}

func (m *RecManager) getHistory(where string) ([]*RecHistory, error) {
	var result []*RecHistory

	query := fmt.Sprintf(sqlSelectHistory, m.tableName, m.historyTableName, where)

	rows, err := m.db.Query(query)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		row := &RecHistory{}
		err := rows.Scan(
			&row.ID,
			&row.When,
			&row.Name,
			&row.Applied,
			&row.Note,
			&row.GITinfo,
			&row.Filename,
		)
		if err != nil {
			return nil, err
		}
		result = append(result, row)
	}

	return result, nil
}

// GetHistory - get history
func (m *RecManager) GetHistory() ([]*RecHistory, error) {
	return m.getHistory("")
}

// GetHistoryForName - get history
func (m *RecManager) GetHistoryForName(name string) ([]*RecHistory, error) {
	return m.getHistory(fmt.Sprintf("WHERE name = '%v'", name))
}

// GetHistoryForGITinfo - get history
func (m *RecManager) GetHistoryForGITinfo(gitinfo string) ([]*RecHistory, error) {
	return m.getHistory(fmt.Sprintf("WHERE gitinfo = '%v'", gitinfo))
}

// GetHistoryForNote - get history
func (m *RecManager) GetHistoryForNote(note string) ([]*RecHistory, error) {
	return m.getHistory(fmt.Sprintf("WHERE note = '%v'", note))
}

// FindByFilenameShort - find by filename in array
func (m *RecManager) FindByFilenameShort(records []*RecShort, filename string) *RecShort {
	for _, rec := range records {
		if rec.Filename == filename {
			return rec
		}
	}
	return nil
}

// InsertRecords - insert new record (only full)
func (m *RecManager) InsertRecords(records []*RecFull, tx *sql.Tx) error {
	var (
		id uint64
	)
	query := fmt.Sprintf(sqlInsert, m.tableName)
	historyQuery := fmt.Sprintf(sqlInsertHistory, m.historyTableName)

	for _, record := range records {

		err := tx.QueryRow(
			query,
			record.Name,
			record.Applied,
			record.Note,
			record.GITinfo,
			record.Filename,
			record.TStamp,
			(pq.StringArray)(record.Requirements),
			record.SQLData,
		).Scan(&id)

		if err != nil {
			return err
		}

		_, err = tx.Exec(
			historyQuery,
			id,
			record.Applied,
		)

		if err != nil {
			return err
		}
	}

	return nil
}

// UpdateRecords - insert new record (only full)
func (m *RecManager) UpdateRecords(records []*RecFull, tx *sql.Tx) error {
	historyQuery := fmt.Sprintf(sqlInsertHistory, m.historyTableName)

	for _, record := range records {
		query := fmt.Sprintf(sqlUpdate, m.tableName, record.ID)
		_, err := tx.Exec(
			query,
			record.Name,
			record.Applied,
			record.Note,
			record.GITinfo,
			record.Filename,
			record.TStamp,
			(pq.StringArray)(record.Requirements),
			record.SQLData,
		)
		if err != nil {
			return err
		}
		_, err = tx.Exec(
			historyQuery,
			record.ID,
			record.Applied,
		)
		if err != nil {
			return err
		}
	}

	return nil
}
