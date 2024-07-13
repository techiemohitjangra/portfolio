package model

import (
	"database/sql"
	"log"
)

type Reference struct {
	ID     int64  `param:"id" query:"id" form:"id" json:"id"`
	Title  string `param:"title" query:"title" form:"title" json:"title"`
	URL    string `param:"url" query:"url" form:"url" json:"url"`
	BlogID int64  `param:"blog_id" query:"blog_id" form:"blog_id" json:"blog_id"`
}

func GetReferences(txn *sql.Tx, blogID int64) ([]Reference, error) {
	// Query to get the references associated with the blog
	getReferences, err := txn.Prepare(`
		SELECT r.ID, r.Title, r.URL, r.BlogID
		FROM Reference r
		WHERE r.BlogID = ?;
	`)
	if err != nil {
		log.Println("failed to prepare statement: ", err)
		return nil, err
	}
	defer getReferences.Close()

	rows, err := getReferences.Query(blogID)
	if err != nil {
		log.Println("failed to query references: ", err)
	}

	var references []Reference
	for rows.Next() {
		var reference Reference
		err := rows.Scan(&reference.ID, &reference.Title, &reference.URL, &reference.BlogID)
		if err != nil {
			return nil, err
		}
		references = append(references, reference)
	}
	return references, nil
}

// TODO: add authentication
func AddReference(txn *sql.Tx, ref Reference) error {
	addReference, err := txn.Prepare(`
		INSERT OR IGNORE INTO Reference (Title, URL, BlogID)
		VALUES (?,?,?)
	`)
	if err != nil {
		log.Println("failed to prepare statement to insert Reference: ", err)
		return err
	}
	defer addReference.Close()

	_, err = addReference.Exec(ref.Title, ref.URL, ref.BlogID)
	if err != nil {
		log.Println("execution failed: ", err)
		return err
	}

	return nil
}

// TODO: add authentication
func UpdateReference(txn *sql.Tx, ref Reference) error {
	updateReference, err := txn.Prepare(`
		UPDATE Reference
		SET Title = ?, URL = ?
		WHERE BlogID = ? AND ID = ?;
	`)
	if err != nil {
		log.Println("failed to prepare statement to update Reference: ", err)
		return err
	}
	defer updateReference.Close()

	_, err = updateReference.Exec(ref.Title, ref.URL, ref.BlogID, ref.ID)
	if err != nil {
		log.Println("execution failed: ", err)
		return err
	}

	return nil
}

// TODO: add authentication
func DeleteReference(txn *sql.Tx, blogID int64, referenceID int64) error {
	deleteReference, err := txn.Prepare(`
		DELETE
		FROM Reference
		WHERE BlogID = ? AND ID = ?;
	`)
	if err != nil {
		log.Println("failed to prepare statement to delete Reference: ", err)
		return err
	}
	defer deleteReference.Close()

	_, err = deleteReference.Exec(blogID, referenceID)
	if err != nil {
		log.Println("execution failed: ", err)
		return err
	}
	return nil
}

// TODO: add authentication
func DeleteReferences(txn *sql.Tx, blogID int64) error {
	deleteReference, err := txn.Prepare(`
		DELETE
		FROM Reference
		WHERE BlogID = ?;
	`)
	if err != nil {
		log.Println("failed to prepare statement to delete Reference: ", err)
		return err
	}
	defer deleteReference.Close()

	_, err = deleteReference.Exec(blogID)
	if err != nil {
		log.Println("execution failed: ", err)
		return err
	}
	return nil
}
