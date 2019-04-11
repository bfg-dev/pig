package file

import (
	"fmt"
	"os"
	"text/template"
	"time"
)

const (
	fileTemplate = `
-- +pig Name: {{.Name}}
-- +pig Requirements:
-- +pig Up

-- CREATE TABLE EXAMPLE
CREATE SEQUENCE "<!!!TABLE NAME!!!>_id_seq" START 1;
 
CREATE TABLE "<!!!TABLE NAME!!!>" (
	id bigint NOT NULL DEFAULT nextval('<!!!TABLE NAME!!!>_id_seq'::regclass),
	-- fields here
);

-- ALTER TABLE EXAMPLE
ALTER TABLE "<!!!TABLE NAME!!!>" ADD COLUMN "<!!!COLUMN NAME!!!>" ...
 
-- +pig Down

-- DROP TABLE EXAMPLE
DROP TABLE "<!!!TABLE NAME!!!>";
DROP SEQUENCE "<!!!TABLE NAME!!!>_id_seq";

-- ALTER TABLE EXAMPLE
ALTER TABLE "<!!!TABLE NAME!!!>" DROP COLUMN "<!!!COLUMN NAME!!!>" ...
	`
)

type templateData struct {
	Name string
}

// GenerateNewSQLFile - generate new SQL file from template
func GenerateNewSQLFile(dir *string, name string) error {
	var (
		err error
	)

	filename := fmt.Sprintf("%v_%v.sql", time.Now().Unix(), name)
	if dir != nil {
		filename = fmt.Sprintf("%v/%v", *dir, filename)
	}

	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	t := template.Must(template.New("sql").Parse(fileTemplate))

	err = t.Execute(f, templateData{Name: name})
	if err != nil {
		return err
	}

	return nil
}
