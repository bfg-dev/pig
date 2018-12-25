package db

const (
	sqlCreateTable = `
		CREATE TABLE "%v" ( 
			id serial NOT NULL,
			name text NOT NULL UNIQUE,
			applied boolean NOT NULL,
			note text NULL default NULL,
			gitinfo text NULL default NULL,
			filename text NOT NULL UNIQUE,
			tstamp timestamp NULL default now(),
			requirements text[],
			sqldata text NULL default NULL,
			PRIMARY KEY(id)
		);
	`
	sqlCreateHistoryTable = `
		CREATE TABLE "%v" ( 
			id serial NOT NULL,
			"when" timestamp NULL default now(),
			migrationId INTEGER NOT null,
			applied boolean NOT NULL,
			PRIMARY KEY(id)
		);
	`
	sqlSelectShort = `
		SELECT 
			id, 
			name, 
			applied, 
			note, 
			gitinfo, 
			filename, 
			tstamp,
			requirements
		FROM %v
		%v
		ORDER BY id;
	`
	sqlSelectFull = `
		SELECT 
			id, 
			name, 
			applied, 
			note, 
			gitinfo, 
			filename, 
			tstamp,
			requirements,
			sqldata
		FROM %v
		%v
		ORDER BY id;
	`
	sqlSelectHistory = `
		SELECT
			%[2]v.id as id,
			"when",
			name,
			%[2]v.applied, 
			note, 
			gitinfo, 
			filename
		FROM %[2]v JOIN %[1]v ON %[2]v.migrationId = %[1]v.id
		%[3]v
		ORDER BY id;
	`
	sqlInsert = `
		INSERT INTO "%v" ( 
			name,
			applied,
			note,
			gitinfo,
			filename,
			tstamp,
			requirements,
			sqldata
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id;
	`
	sqlUpdate = `
		UPDATE "%v"
		SET
			name = $1,
			applied = $2,
			note = $3,
			gitinfo = $4,
			filename = $5,
			tstamp = $6,
			requirements = $7,
			sqldata = $8
		WHERE id = %v;
	`
	sqlInsertHistory = `
		INSERT INTO "%v" ( 
			migrationId,
			applied
		) VALUES ($1, $2);
	`
)
