package model

import (
	"database/sql"
	"log"
)

type Link struct {
	ID     int64  `param:"id" query:"id" form:"id" json:"id"`
	Text   string `param:"text" query:"text" form:"text" json:"text"`
	URL    string `param:"url" query:"url" form:"url" json:"url"`
	BlogID int64  `param:"blog_id" query:"blog_id" form:"blog_id" json:"blog_id"`
}

// TODO: add authentication
func GetLinks(txn *sql.Tx, blogID int64) ([]Link, error) {
	getLinks, err := txn.Prepare(`
		SELECT ID, Text, URL, BlogID
		FROM Link
		WHERE BlogID = ?;
	`)
	if err != nil {
		log.Println("failed to prepare statement: ", err)
		return nil, err
	}
	defer getLinks.Close()

	rows, err := getLinks.Query(blogID)
	if err != nil {
		log.Println("execution failed: ", err)
		return nil, err
	}
	defer rows.Close()

	var links []Link
	for rows.Next() {
		var link Link
		err := rows.Scan(&link.ID, &link.Text, &link.URL, &link.BlogID)
		if err != nil {
			return nil, err
		}
		links = append(links, link)
	}

	return links, nil
}

// TODO: add authentication
func AddLink(txn *sql.Tx, link Link) error {
	addLink, err := txn.Prepare(`
		INSERT OR IGNORE INTO Link (Text, URL, BlogID)
		VALUES (?,?,?);
	`)
	if err != nil {
		log.Println("failed to prepare statement to add Link: ", err)
		return err
	}
	defer addLink.Close()

	_, err = addLink.Exec(link.Text, link.URL, link.BlogID)
	if err != nil {
		log.Println("execution failed: ", err)
		return err
	}

	return nil
}

// TODO: add authentication
func UpdateLink(txn *sql.Tx, link Link) error {
	updateLink, err := txn.Prepare(`
		UPDATE Link
		SET Text = ?, URL = ?
		WHERE BlogID = ? AND ID = ?;
	`)
	if err != nil {
		log.Println("failed to prepare statement to update Link: ", err)
		return err
	}
	defer updateLink.Close()

	_, err = updateLink.Exec(link.Text, link.URL, link.BlogID, link.ID)
	if err != nil {
		log.Println("execution failed: ", err)
		return err
	}

	return nil
}

// TODO: add authentication
func DeleteLink(txn *sql.Tx, blogID int64, linkID int64) error {
	deleteLink, err := txn.Prepare(`
		DELETE
		FROM Link
		WHERE BlogID = ? AND ID = ?;
	`)
	if err != nil {
		log.Println("failed to prepare statement to delete Link: ", err)
		return err
	}
	defer deleteLink.Close()

	_, err = deleteLink.Exec(blogID, linkID)
	if err != nil {
		log.Println("execution failed: ", err)
		return err
	}

	return nil
}

func DeleteLinks(txn *sql.Tx, blogID int64) error {
	deleteLink, err := txn.Prepare(`
		DELETE
		FROM Link
		WHERE BlogID = ?;
	`)
	if err != nil {
		log.Println("failed to prepare statement to delete Link: ", err)
		return err
	}
	defer deleteLink.Close()

	_, err = deleteLink.Exec(blogID)
	if err != nil {
		log.Println("execution failed: ", err)
		return err
	}

	return nil
}
