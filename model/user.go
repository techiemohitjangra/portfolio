package model

import (
	"database/sql"
	"log"
	"time"
)

type User struct {
	FirstName             string    `param:"first_name" query:"first_name" form:"first_name" json:"first_name"`
	LastName              string    `param:"last_name" query:"last_name" form:"last_name" json:"last_name"`
	UserName              string    `param:"user_name" query:"user_name" form:"user_name" json:"user_name"`
	EmailAddress          string    `param:"email_address" query:"email_address" form:"email_address" json:"email_address"`
	Password              string    `param:"password" query:"password" form:"password" json:"password"`
	DOB                   time.Time `param:"dob" query:"dob" form:"dob" json:"dob"`
	About                 string    `param:"about" query:"about" form:"about" json:"about"`
	ProfilePicturePath    string    `param:"profile_picture_path" query:"profile_picture_path" form:"profile_picture_path" json:"profile_picture_path"`
	City                  string    `param:"city" query:"city" form:"city" json:"city"`
	LastUpdated           time.Time `param:"last_updated" query:"last_updated" form:"last_updated" json:"last_updated"`
	LastLoggedIn          time.Time `param:"last_logged_in" query:"last_logged_in" form:"last_logged_in" json:"last_logged_in"`
	LastLoggedInLocation  string    `param:"last_logged_in_location" query:"last_logged_in_location" form:"last_logged_in_location" json:"last_logged_in_location"`
	LastLoggedInIPAddress string    `param:"last_logged_in_ip_address" query:"last_logged_in_ip_address" form:"last_logged_in_ip_address" json:"last_logged_in_ip_address"`
}

type UserShow struct {
	FirstName          string `param:"first_name" query:"first_name" form:"first_name" json:"first_name"`
	LastName           string `param:"last_name" query:"last_name" form:"last_name" json:"last_name"`
	UserName           string `param:"user_name" query:"user_name" form:"user_name" json:"user_name"`
	EmailAddress       string `param:"email_address" query:"email_address" form:"email_address" json:"email_address"`
	About              string `param:"about" query:"about" form:"about" json:"about"`
	ProfilePicturePath string `param:"profile_picture_path" query:"profile_picture_path" form:"profile_picture_path" json:"profile_picture_path"`
}

// TODO: add authentication
func AddUser(db *sql.DB, user User) error {
	txn, err := db.Begin()
	if err != nil {
		log.Println("failed to start transaction: ", err)
		return err
	}
	defer txn.Rollback()

	createUser, err := txn.Prepare(`
		INSERT OR IGNORE INTO User (FirstName, LastName, UserName, EmailAddress, Password, DOB, About, ProfilePicturePath, City, LastUpdated, LastLoggedIn, LastLoggedInLocation, LastLoggedInIPAddress)
		VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?);
	`)
	if err != nil {
		log.Println("failed to prepare statement to create User: ", err)
		return err
	}
	defer createUser.Close()

	_, err = createUser.Exec(
		user.FirstName,
		user.LastName,
		user.UserName,
		user.EmailAddress,
		user.Password,
		user.DOB.Unix(),
		user.About,
		user.ProfilePicturePath,
		user.City,
		user.LastUpdated.Unix(),
		user.LastLoggedIn.Unix(),
		user.LastLoggedInLocation,
		user.LastLoggedInIPAddress,
	)
	if err != nil {
		log.Println("execution failed: failed to insert User: ", err)
		return err
	}

	err = txn.Commit()
	if err != nil {
		log.Println("failed to commit transaction: ", err)
		return err
	}
	return nil
}

// TODO: add authentication
func GetUser(db *sql.DB, userName string) (*User, error) {
	txn, err := db.Begin()
	if err != nil {
		log.Println("failed to start transaction: ", err)
		return nil, err
	}
	defer txn.Rollback()

	addUser, err := txn.Prepare(`
		SELECT FirstName, LastName, UserName, EmailAddress, Password, DOB, About, ProfilePicturePath, City, LastUpdate, LastLoggedIn, LastLoggedInLocation, LastLoggedInIPAddress
		FROM User;
	`)
	if err != nil {
		log.Println("failed to prepare statement to add user: ", err)
		return nil, err
	}
	defer addUser.Close()

	var user User
	var dob int64
	var lastUpdated int64
	var lastLoggedIn int64

	err = addUser.QueryRow(userName).Scan(
		&user.FirstName,
		&user.LastName,
		&user.UserName,
		&user.EmailAddress,
		&user.Password,
		dob,
		&user.About,
		&user.ProfilePicturePath,
		&user.City,
		lastUpdated,
		lastLoggedIn,
		&user.LastLoggedInLocation,
		&user.LastLoggedInIPAddress,
	)
	if err != nil {
		log.Println("failed to parse User object from query result: ", err)
		return nil, err
	}

	user.DOB = time.Unix(dob, 0)
	user.LastUpdated = time.Unix(lastUpdated, 0)
	user.LastLoggedIn = time.Unix(lastLoggedIn, 0)

	err = txn.Commit()
	if err != nil {
		log.Println("failed to commit transaction: ", err)
		return nil, err
	}
	return &user, nil
}

func GetUserShow(db *sql.DB, userName string) (*UserShow, error) {
	var user UserShow
	txn, err := db.Begin()
	if err != nil {
		log.Println("failed to start transaction: ", err)
		return nil, err
	}
	defer txn.Rollback()

	addUser, err := txn.Prepare(`
		SELECT FirstName, LastName, UserName, EmailAddress, About, ProfilePicturePath
		FROM User
		Where UserName = ?;
	`)
	if err != nil {
		log.Println("failed to prepare statement to add user: ", err)
		return nil, err
	}
	defer addUser.Close()

	err = addUser.QueryRow(userName).Scan(
		&user.FirstName,
		&user.LastName,
		&user.UserName,
		&user.EmailAddress,
		&user.About,
		&user.ProfilePicturePath,
	)
	if err != nil {
		log.Println("failed to parse UserShow object from query result: ", err)
		return nil, err
	}

	err = txn.Commit()
	if err != nil {
		log.Println("failed to commit transaction: ", err)
		return nil, err
	}

	return &user, nil
}
