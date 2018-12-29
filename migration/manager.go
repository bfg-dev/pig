package migration

import (
	"database/sql"
	"log"

	"github.com/bfg-dev/pig/db"
	"github.com/bfg-dev/pig/file"
)

// Manager - operation manager
type Manager struct {
	dbconnection *sql.DB
	dbmanager    *db.RecManager
	filemanager  *file.RecManager
}

func stringPointerToString(str *string) string {
	if str == nil {
		return ""
	}
	return *str
}

// NewManager - returns new OpManager
func NewManager(dbconnection *sql.DB, tableName, historyTableName, directory string) *Manager {
	return &Manager{
		dbconnection: dbconnection,
		dbmanager:    db.NewRecManager(dbconnection, tableName, historyTableName),
		filemanager:  file.NewRecManager(directory),
	}
}

// InitDB - create migration tables
func (o *Manager) InitDB() error {
	if err := o.dbmanager.CreateTable(); err != nil {
		return err
	}
	if err := o.dbmanager.CreateHistoryTable(); err != nil {
		return err
	}
	return nil
}

func (o *Manager) loadRawMigrations() (*Migrations, error) {
	var (
		newShortFilerecords []*file.RecShort
		newFullFilerecords  []*file.RecFull
		result              Migrations
		found               *db.RecShort
	)

	// 1: get all short records from db
	dbrecords, err := o.dbmanager.GetAllShort()
	if err != nil {
		return nil, err
	}

	// 2: get all short records from directory
	filerecords, err := o.filemanager.GetAllShort()
	if err != nil {
		return nil, err
	}

	// 3: find new and unapplied file records
	for _, filerec := range filerecords {
		found = o.dbmanager.FindByFilenameShort(dbrecords, filerec.Filename)
		if found == nil || !found.Applied {
			newShortFilerecords = append(newShortFilerecords, filerec)
		}
	}

	// 4: load full file record for new files
	newFullFilerecords, err = o.filemanager.GetFullFromShort(newShortFilerecords)
	if err != nil {
		return nil, err
	}

	// 5: add db info to migration
	for _, dbrec := range dbrecords {
		metaItem := &Meta{
			ID:           dbrec.ID,
			Requirements: nil,
			Children:     nil,
			Name:         dbrec.Name,
			Filename:     dbrec.Filename,
			Note:         stringPointerToString(dbrec.Note),
			GITinfo:      stringPointerToString(dbrec.GITinfo),
			Applied:      dbrec.Applied,
			TStamp:       dbrec.TStamp,
			DBShortRec:   dbrec,
		}
		if !dbrec.Applied {
			metaItem.Pending = true

			metaItem.FileFullRec = o.filemanager.FindByFilenameFull(newFullFilerecords, dbrec.Filename)
			if metaItem.FileFullRec != nil {
				newFullFilerecords = o.filemanager.RemoveFromListFull(newFullFilerecords, metaItem.FileFullRec)
			} else {
				// No file - disabling pending
				metaItem.Pending = false
			}
		}
		result.Items = append(result.Items, metaItem)
	}

	// 6: add file info to migration
	for _, filerec := range newFullFilerecords {
		metaItem := &Meta{
			ID:           0,
			Requirements: nil,
			Children:     nil,
			Name:         filerec.Name,
			Filename:     filerec.Filename,
			Applied:      false,
			Pending:      true,
			TStamp:       filerec.TStamp,
			FileFullRec:  filerec,
		}
		result.Items = append(result.Items, metaItem)
	}

	return &result, nil
}

func (o *Manager) parseRequirements(migrations *Migrations) error {
	var rawRequirements []string
	for _, mig := range migrations.Items {

		if mig.DBShortRec != nil {
			rawRequirements = mig.DBShortRec.Requirements
		}

		if mig.FileFullRec != nil {
			rawRequirements = mig.FileFullRec.Requirements
		}

		if len(rawRequirements) == 1 && rawRequirements[0] == "" {
			continue
		}

		for _, rawReq := range rawRequirements {
			req := migrations.GetByName(rawReq)
			if req == nil {
				return &RequirementNotFoundError{Migration: mig, Requirement: rawReq}
			}

			if err := mig.AddRequirement(req); err != nil {
				return err
			}
		}
	}

	return nil
}

// GetAllMigrations - get migrations
func (o *Manager) GetAllMigrations() (*Migrations, error) {
	migs, err := o.loadRawMigrations()
	if err != nil {
		return nil, err
	}

	err = o.parseRequirements(migs)
	if err != nil {
		return nil, err
	}

	err = migs.Prepare()
	if err != nil {
		return nil, err
	}

	return migs, nil
}

// GetUpPlan - get UP plan for all pending migraions
func (o *Manager) GetUpPlan() (*Migrations, error) {
	var result Migrations

	migs, err := o.GetAllMigrations()
	if err != nil {
		return nil, err
	}

	for _, mig := range migs.FindBottoms().Items {
		result.Items = append(result.Items, mig.GetUpPlan().Items...)
	}

	return result.RemoveDuplicates(), nil
}

// GetDownPlan - get DOWN plan for all pending migraions
func (o *Manager) GetDownPlan() (*Migrations, error) {
	var result Migrations

	migs, err := o.GetAllMigrations()
	if err != nil {
		return nil, err
	}

	for _, mig := range migs.Items {
		if mig.Applied && len(mig.Requirements) == 0 {
			result.Items = append(result.Items, mig.GetDownPlan().Items...)
		}
	}

	return result.RemoveDuplicates(), nil
}

// GetUpPlanForName - get UP plan for one migration
func (o *Manager) GetUpPlanForName(name string) (*Migrations, error) {
	migs, err := o.GetAllMigrations()
	if err != nil {
		return nil, err
	}

	mig := migs.GetByName(name)
	if mig == nil {
		return nil, &NotFound{Filter: "Name", Name: name}
	}

	return mig.GetUpPlan().RemoveDuplicates(), nil
}

// GetDownPlanForName - get DOWN plan for one migration
func (o *Manager) GetDownPlanForName(name string) (*Migrations, error) {
	migs, err := o.GetAllMigrations()
	if err != nil {
		return nil, err
	}

	mig := migs.GetByName(name)
	if mig == nil {
		return nil, &NotFound{Filter: "Name", Name: name}
	}

	return mig.GetDownPlan().RemoveDuplicates(), nil
}

// GetUpPlanForGITinfo - get UP plan for one migration
func (o *Manager) GetUpPlanForGITinfo(gitinfo string) (*Migrations, error) {
	migs, err := o.GetAllMigrations()
	if err != nil {
		return nil, err
	}

	mig := migs.GetByGITinfo(gitinfo)
	if mig == nil {
		return nil, &NotFound{Filter: "GITinfo", Name: gitinfo}
	}

	return mig.GetUpPlan().RemoveDuplicates(), nil
}

// GetDownPlanForGITinfo - get DOWN plan for one migration
func (o *Manager) GetDownPlanForGITinfo(gitinfo string) (*Migrations, error) {
	migs, err := o.GetAllMigrations()
	if err != nil {
		return nil, err
	}

	mig := migs.GetByGITinfo(gitinfo)
	if mig == nil {
		return nil, &NotFound{Filter: "GITinfo", Name: gitinfo}
	}

	return mig.GetDownPlan().RemoveDuplicates(), nil
}

// GetUpPlanForNote - get UP plan for one migration
func (o *Manager) GetUpPlanForNote(note string) (*Migrations, error) {
	migs, err := o.GetAllMigrations()
	if err != nil {
		return nil, err
	}

	mig := migs.GetByNote(note)
	if mig == nil {
		return nil, &NotFound{Filter: "Note", Name: note}
	}

	return mig.GetUpPlan().RemoveDuplicates(), nil
}

// GetDownPlanForNote - get DOWN plan for one migration
func (o *Manager) GetDownPlanForNote(note string) (*Migrations, error) {
	migs, err := o.GetAllMigrations()
	if err != nil {
		return nil, err
	}

	mig := migs.GetByNote(note)
	if mig == nil {
		return nil, &NotFound{Filter: "Note", Name: note}
	}

	return mig.GetDownPlan().RemoveDuplicates(), nil
}

// GetHistory - get history
func (o *Manager) GetHistory() ([]*db.RecHistory, error) {
	return o.dbmanager.GetHistory()
}

// GetHistoryForName - get history for migration
func (o *Manager) GetHistoryForName(name string) ([]*db.RecHistory, error) {
	return o.dbmanager.GetHistoryForName(name)
}

// GetHistoryForGITinfo - get history for migration
func (o *Manager) GetHistoryForGITinfo(gitinfo string) ([]*db.RecHistory, error) {
	return o.dbmanager.GetHistoryForGITinfo(gitinfo)
}

// GetHistoryForNote - get history for migration
func (o *Manager) GetHistoryForNote(note string) ([]*db.RecHistory, error) {
	return o.dbmanager.GetHistoryForNote(note)
}

// ExecuteUp - execute UP migrations
func (o *Manager) ExecuteUp(migrationMeta *Meta) error {
	if migrationMeta.FileFullRec == nil {
		return &NoFileInformation{Migration: migrationMeta}
	}

	return o.executeMigration(migrationMeta, true)
}

// ExecuteDown - execute DOWN migrations
func (o *Manager) ExecuteDown(migrationMeta *Meta) error {
	if migrationMeta.DBShortRec == nil {
		return &NoDBInformation{Migration: migrationMeta}
	}

	fullRec, err := o.dbmanager.GetFullFromShort(migrationMeta.DBShortRec)
	if err != nil {
		return err
	}

	migrationMeta.DBFullRec = fullRec

	return o.executeMigration(migrationMeta, false)
}

func (o *Manager) updateMigrationTable(migrationMeta *Meta, up bool, tx *sql.Tx) error {
	if up {
		newDBFullRec := db.RecFull{
			RecShort: db.RecShort{
				Name:         migrationMeta.Name,
				Applied:      true,
				Note:         &migrationMeta.Note,
				GITinfo:      &migrationMeta.GITinfo,
				Filename:     migrationMeta.Filename,
				TStamp:       migrationMeta.FileFullRec.TStamp,
				Requirements: migrationMeta.FileFullRec.Requirements,
			},
			SQLData: migrationMeta.FileFullRec.SQLData,
		}

		if migrationMeta.DBShortRec == nil {
			if err := o.dbmanager.InsertRecords([]*db.RecFull{&newDBFullRec}, tx); err != nil {
				return err
			}
		} else {
			newDBFullRec.ID = migrationMeta.DBShortRec.ID
			if err := o.dbmanager.UpdateRecords([]*db.RecFull{&newDBFullRec}, tx); err != nil {
				return err
			}
		}

	} else {
		migrationMeta.DBFullRec.Applied = false
		if err := o.dbmanager.UpdateRecords([]*db.RecFull{migrationMeta.DBFullRec}, tx); err != nil {
			return err
		}
	}

	return nil
}

func (o *Manager) executeMigration(migrationMeta *Meta, up bool) error {

	statements, err := migrationMeta.GetSQLStatements(up)
	if err != nil {
		return err
	}

	if statements.Transactional {
		// TRANSACTION.

		tx, err := o.dbconnection.Begin()
		if err != nil {
			log.Fatal(err)
		}

		for _, query := range statements.Lines {
			if _, err = tx.Exec(query); err != nil {
				tx.Rollback()
				return err
			}
		}

		if err := o.updateMigrationTable(migrationMeta, up, tx); err != nil {
			tx.Rollback()
			return err
		}

		return tx.Commit()
	}

	// NO TRANSACTION.
	for _, query := range statements.Lines {
		if _, err := o.dbconnection.Exec(query); err != nil {
			return err
		}
	}

	tx, err := o.dbconnection.Begin()
	if err != nil {
		log.Fatal(err)
	}

	if err := o.updateMigrationTable(migrationMeta, up, tx); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}
