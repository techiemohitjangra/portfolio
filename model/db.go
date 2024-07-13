package model

import (
	"database/sql"
	"log"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

// open a sqlite DB context
func OpenDB(dbFile string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", dbFile)
	if err != nil {
		log.Println("failed to open sqlite DB: ", err)
		return nil, err
	}
	if err := db.Ping(); err != nil {
		db.Close()
		log.Println("failed to ping DB, closing DB: ", err)
		return nil, err
	}
	db.SetConnMaxIdleTime(5 * time.Minute)
	db.SetConnMaxLifetime(1 * time.Hour)
	return db, nil
}

// close a DB context
func CloseDB(db *sql.DB) error {
	err := db.Close()
	if err != nil {
		log.Panic(err)
	}
	return nil
}

func InitTables(db *sql.DB) error {
	txn, err := db.Begin()
	if err != nil {
		log.Println("failed to start transaction: ", err)
		return err
	}
	defer txn.Rollback()

	// Create Blog Table
	createBlogTable, err := db.Prepare(`
		CREATE TABLE IF NOT EXISTS Blog (
			ID INTEGER PRIMARY KEY AUTOINCREMENT,
			Title TEXT NOT NULL,
			Subtitle TEXT,
			PublishedOn INTEGER NOT NULL,
			LastUpdated INTEGER NOT NULL,
			Content TEXT
		);
	`)
	if err != nil {
		log.Println("failed to prepare statement: ", err)
		return err
	}
	defer createBlogTable.Close()

	_, err = createBlogTable.Exec()
	if err != nil {
		log.Println("execution failed: ", err)
		return err
	}
	log.Println(`
		CREATE TABLE IF NOT EXISTS Blog (
			ID INTEGER PRIMARY KEY AUTOINCREMENT,
			Title TEXT NOT NULL,
			Subtitle TEXT,
			PublishedOn INTEGER NOT NULL,
			LastUpdatedOn INTEGER NOT NULL,
			Content TEXT
		);
	`)

	// Create Link Table
	createLinkTable, err := txn.Prepare(`
		CREATE TABLE IF NOT EXISTS Link (
			ID INTEGER PRIMARY KEY AUTOINCREMENT,
			Text TEXT NOT NULL,
			URL TEXT NOT NULL,
			BlogID INTEGER,
			FOREIGN KEY (BlogID) REFERENCES Blog(ID) ON DELETE CASCADE
		);
	`)
	if err != nil {
		log.Println("failed to prepare statement: ", err)
		return err
	}
	defer createLinkTable.Close()

	_, err = createLinkTable.Exec()
	if err != nil {
		log.Println("execution failed: ", err)
		return err
	}
	log.Println(`
		CREATE TABLE IF NOT EXISTS Link (
			ID INTEGER PRIMARY KEY AUTOINCREMENT,
			Text TEXT NOT NULL,
			URL TEXT NOT NULL,
			BlogID INTEGER,
			FOREIGN KEY (BlogID) REFERENCES Blog(ID) ON DELETE CASCADE
		);
	`)

	// Create Project Table
	createProjectTable, err := txn.Prepare(`
		CREATE TABLE IF NOT EXISTS Project (
			ID INTEGER PRIMARY KEY AUTOINCREMENT,
			Title TEXT NOT NULL,
			Subtitle TEXT,
			RepositoryURL TEXT,
			DeploymentURL TEXT,
			LastUpdated INTEGER NOT NULL,
			PublishedOn INTEGER NOT NULL,
			Status TEXT NOT NULL,
			Summary TEXT,
			Content TEXT,
			Visibility INTEGER
		);
	`)
	if err != nil {
		log.Println("failed to prepare statement to create Project table: ", err)
		return err
	}
	defer createProjectTable.Close()

	_, err = createProjectTable.Exec()
	if err != nil {
		log.Println("execution failed: ", err)
		return err
	}
	log.Println(`
		CREATE TABLE IF NOT EXISTS Project (
			ID INTEGER PRIMARY KEY AUTOINCREMENT,
			Title TEXT NOT NULL,
			Subtitle TEXT,
			RepositoryURL TEXT,
			DeploymentURL TEXT,
			LastUpdated INTEGER NOT NULL,
			PublishedOn INTEGER NOT NULL,
			Status TEXT NOT NULL,
			Summary TEXT,
			Content TEXT,
			Visibility INTEGER
		);
	`)

	// Create Reference Table
	referenceTable, err := txn.Prepare(`
		CREATE TABLE IF NOT EXISTS Reference (
			ID INTEGER PRIMARY KEY AUTOINCREMENT,
			Title TEXT NOT NULL,
			URL TEXT NOT NULL,
			BlogID INTEGER,
			FOREIGN KEY (BlogID) REFERENCES Blog(ID) ON DELETE CASCADE
		);
	`)
	if err != nil {
		log.Println("failed to prepare statement: ", err)
		return err
	}
	defer referenceTable.Close()

	_, err = referenceTable.Exec()
	if err != nil {
		log.Println("execution failed: ", err)
		return err
	}
	log.Println(`
		CREATE TABLE IF NOT EXISTS Reference (
			ID INTEGER PRIMARY KEY AUTOINCREMENT,
			Title TEXT NOT NULL,
			URL TEXT NOT NULL,
			BlogID INTEGER,
			FOREIGN KEY (BlogID) REFERENCES Blog(ID) ON DELETE CASCADE
		);
	`)

	// Create Tag Table
	createTagTable, err := txn.Prepare(`
		CREATE TABLE IF NOT EXISTS Tag (
			ID INTEGER PRIMARY KEY AUTOINCREMENT,
			Name TEXT NOT NULL UNIQUE,
			IconPath TEXT
		);
	`)
	if err != nil {
		log.Println("failed to prepare statement to create Tag table: ", err)
		return err
	}
	defer createTagTable.Close()

	_, err = createTagTable.Exec()
	if err != nil {
		log.Println("execution failed: ", err)
		return err
	}
	log.Println(`
		CREATE TABLE IF NOT EXISTS Tag (
			ID INTEGER PRIMARY KEY AUTOINCREMENT,
			Name TEXT NOT NULL UNIQUE,
			IconPath TEXT
		);
	`)

	// created BlogTags relation table
	createBlogTagsTable, err := txn.Prepare(`
		CREATE TABLE IF NOT EXISTS BlogTag (
			BlogID INTEGER,
			TagID INTEGER,
			FOREIGN KEY (BlogID) REFERENCES Blog(ID) ON DELETE CASCADE,
			FOREIGN KEY (TagID) REFERENCES Tag(ID) ON DELETE CASCADE,
			PRIMARY KEY (BlogID, TagID)
		);
	`)
	if err != nil {
		log.Println("failed to prepare statement to create BlogTag table: ", err)
		return err
	}
	defer createBlogTagsTable.Close()

	_, err = createBlogTagsTable.Exec()
	if err != nil {
		log.Println("execution failed: ", err)
		return err
	}
	log.Println(`
		CREATE TABLE IF NOT EXISTS BlogTag (
			BlogID INTEGER,
			TagID INTEGER,
			FOREIGN KEY (BlogID) REFERENCES Blog(ID) ON DELETE CASCADE,
			FOREIGN KEY (TagID) REFERENCES Tag(ID) ON DELETE CASCADE,
			PRIMARY KEY (BlogID, TagID)
		);
	`)

	// created ProjectTags relation table
	createProjectTagsTable, err := txn.Prepare(`
		CREATE TABLE IF NOT EXISTS ProjectTag (
			ProjectID INTEGER,
			TagID INTEGER,
			FOREIGN KEY (ProjectID) REFERENCES Project(ID) ON DELETE CASCADE,
			FOREIGN KEY (TagID) REFERENCES Tag(ID) ON DELETE CASCADE,
			PRIMARY KEY (ProjectID, TagID)
		);
	`)
	if err != nil {
		log.Println("failed to prepare statement to create ProjectTag table: ", err)
	}
	defer createProjectTagsTable.Close()

	_, err = createProjectTagsTable.Exec()
	if err != nil {
		log.Println("execution failed: ", err)
		return err
	}
	log.Println(`
		CREATE TABLE IF NOT EXISTS ProjectTag (
			ProjectID INTEGER,
			TagID INTEGER,
			FOREIGN KEY (ProjectID) REFERENCES Project(ID) ON DELETE CASCADE,
			FOREIGN KEY (TagID) REFERENCES Tag(ID) ON DELETE CASCADE,
			PRIMARY KEY (ProjectID, TagID)
		);
	`)

	// Create User Table
	createUserTable, err := txn.Prepare(`
		CREATE TABLE IF NOT EXISTS User (
			FirstName TEXT NOT NULL,
			LastName TEXT NOT NULL,
			UserName TEXT NOT NULL UNIQUE,
			EmailAddress TEXT NOT NULL UNIQUE,
			Password TEXT NOT NULL,
			DOB INTEGER NOT NULL,
			About TEXT,
			ProfilePicturePath TEXT,
			City TEXT,
			LastUpdated INTEGER NOT NULL,
			LastLoggedIn INTEGER NOT NULL,
			LastLoggedInLocation TEXT,
			LastLoggedInIPAddress TEXT
		);
	`)
	if err != nil {
		log.Println("failed to prepare statement: ", err)
		return err
	}
	defer createUserTable.Close()

	_, err = createUserTable.Exec()
	if err != nil {
		log.Println("execution failed: ", err)
		return err
	}
	log.Println(`
		CREATE TABLE IF NOT EXISTS User (
			FirstName TEXT NOT NULL,
			LastName TEXT NOT NULL,
			UserName TEXT NOT NULL UNIQUE,
			EmailAddress TEXT NOT NULL UNIQUE,
			Password TEXT NOT NULL,
			DOB INTEGER NOT NULL,
			About TEXT,
			ProfilePicturePath TEXT,
			City TEXT,
			LastUpdated INTEGER NOT NULL,
			LastLoggedIn INTEGER NOT NULL,
			LastLoggedInLocation TEXT,
			LastLoggedInIPAddress TEXT
		);
	`)

	err = txn.Commit()
	if err != nil {
		log.Println("failed to commit transaction: ", err)
		return err
	}

	return nil
}
