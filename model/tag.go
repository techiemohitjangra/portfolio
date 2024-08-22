package model

import (
	"database/sql"
	"log"
)

type Tag struct {
	ID       int64  `param:"id" query:"id" form:"id" json:"id"`
	Name     string `param:"name" query:"name" form:"name" json:"name"`
	IconPath string `param:"icon_path" query:"icon_path" form:"icon_path" json:"icon_path"`
}

// TODO: add authentication
func AddTag(txn *sql.Tx, tag Tag) (int64, error) {
	insertTag, err := txn.Prepare(`
			INSERT OR IGNORE INTO Tag (Name, IconPath)
			VALUES (?,?)
		`)
	if err != nil {
		log.Println("failed to prepare statement to insert new Tag: ", err)
		return -1, err
	}
	defer insertTag.Close()

	result, err := insertTag.Exec(tag.Name, tag.IconPath)
	if err != nil {
		log.Println("execution failed: ", err)
		return -1, err
	}

	tagID, err := result.LastInsertId()
	if err != nil {
		log.Println("failed to get id of last inserted Tag: ", err)
		return -1, err
	}
	return tagID, nil
}

func GetBlogTags(txn *sql.Tx, blogID int64) ([]Tag, error) {
	getBlogTags, err := txn.Prepare(`
		SELECT t.ID, t.Name, t.IconPath
		FROM BlogTag bt
		JOIN Tag t ON bt.TagID = t.ID
		WHERE bt.BlogID = ?;
	`)
	if err != nil {
		log.Println("failed to prepare statement to get BlogTags: ", err)
		return nil, err
	}
	defer getBlogTags.Close()

	rows, err := getBlogTags.Query(blogID)
	if err != nil {
		log.Println("failed to query BlogTags: ", err)
		return nil, err
	}
	defer rows.Close()

	var tags []Tag
	for rows.Next() {
		var tag Tag
		err = rows.Scan(&tag.ID, &tag.Name, &tag.IconPath)
		if err != nil {
			log.Println("failed to parse tag object from query: ", err)
			return nil, err
		}

		tags = append(tags, tag)
	}
	if err = rows.Err(); err != nil {
		log.Println("error with rows while getting BlogTags: ", err)
		return nil, err
	}
	return tags, nil
}

// TODO: add authentication
func AddBlogTag(txn *sql.Tx, blogID int64, tag Tag) error {
	var tagID int64
	getTagIdFromName, err := txn.Prepare(`
		SELECT ID
		FROM Tag
		WHERE Name = ?;
	`)
	if err != nil {
		log.Println("failed to prepare statement to get Tag.ID from Tag.Name: ", err)
		return err
	}
	defer getTagIdFromName.Close()

	err = getTagIdFromName.QueryRow(tag.Name).Scan(&tagID)
	if err == sql.ErrNoRows {
		tagID, err = AddTag(txn, tag)
		if err != nil {
			log.Println("failed to insertTag: ", err)
			return err
		}
	} else if err != nil {
		log.Println("failed to query and parse tagID: ", err)
		return err
	}

	insertBlogTag, err := txn.Prepare(`
		INSERT OR IGNORE INTO BlogTag (BlogID, TagID)
		VALUES (?, ?);
	`)
	if err != nil {
		log.Println("failed to prepare statement to insert BlogTag: ", err)
		return err
	}
	defer insertBlogTag.Close()

	_, err = insertBlogTag.Exec(blogID, tagID)
	if err != nil {
		log.Println("execution failed: could not insert BlogTag: ", err)
		return err
	}

	return nil
}

// TODO: add authentication
func DeleteBlogTag(txn *sql.Tx, blogID int64, tagID int64) error {
	deleteBlogTag, err := txn.Prepare(`
		DELETE
		FROM BlogTag
		WHERE BlogID = ? AND TagID = ?;
	`)
	if err != nil {
		log.Println("failed to prepare statement to delete BlogTag: ", err)
		return err
	}
	defer deleteBlogTag.Close()

	_, err = deleteBlogTag.Exec(blogID, tagID)
	if err != nil {
		log.Println("execution failed: ", err)
		return err
	}

	return nil
}

// TODO: add authentication
func DeleteBlogTags(txn *sql.Tx, blogID int64) error {
	deleteBlogTag, err := txn.Prepare(`
		DELETE
		FROM BlogTag
		WHERE BlogID = ?;
	`)
	if err != nil {
		log.Println("failed to prepare statement to delete BlogTags: ", err)
		return err
	}
	defer deleteBlogTag.Close()

	_, err = deleteBlogTag.Exec(blogID)
	if err != nil {
		log.Println("execution failed: ", err)
		return err
	}

	return nil
}

func GetProjectTags(txn *sql.Tx, projectID int64) ([]Tag, error) {
	getProjectTags, err := txn.Prepare(`
		SELECT t.ID, t.Name, t.IconPath
		FROM ProjectTag pt
		JOIN Tag t ON pt.TagID = t.ID
		WHERE pt.ProjectID = ?;
	`)
	if err != nil {
		log.Println("failed to prepare statement to get ProjectTags: ", err)
		return nil, err
	}
	defer getProjectTags.Close()

	projectTagRows, err := getProjectTags.Query(projectID)
	if err != nil {
		log.Println("failed to query ProjectTags: ", err)
		return nil, err
	}
	defer projectTagRows.Close()

	var tags []Tag
	for projectTagRows.Next() {
		var tag Tag
		err = projectTagRows.Scan(&tag.ID, &tag.Name, &tag.IconPath)
		if err != nil {
			log.Println("failed to parse tag object from query: ", err)
			return nil, err
		}

		tags = append(tags, tag)
	}
	if err = projectTagRows.Err(); err != nil {
		log.Println("error with rows while getting ProjectTags: ", err)
		return nil, err
	}
	return tags, nil
}

// TODO: add authentication
func AddProjectTag(txn *sql.Tx, projectID int64, tag Tag) error {
	var tagID int64
	getTagIdFromName, err := txn.Prepare(`
		SELECT ID
		FROM Tag
		WHERE Name = ?;
	`)
	if err != nil {
		log.Println("failed to prepare statement to get Tag.ID from Tag.Name: ", err)
		return err
	}
	defer getTagIdFromName.Close()

	err = getTagIdFromName.QueryRow(tag.Name).Scan(&tagID)
	if err == sql.ErrNoRows {
		tagID, err = AddTag(txn, tag)
		if err != nil {
			log.Println("failed to insertTag: ", err)
			return err
		}
	} else if err != nil {
		log.Println("failed to query and parse tagID: ", err)
		return err
	}

	insertBlogTag, err := txn.Prepare(`
		INSERT OR IGNORE INTO ProjectTag (ProjectID, TagID)
		VALUES (?, ?);
	`)
	if err != nil {
		log.Println("failed to prepare statement to insert ProjectTag: ", err)
		return err
	}
	defer insertBlogTag.Close()

	_, err = insertBlogTag.Exec(projectID, tagID)
	if err != nil {
		log.Println("execution failed: could not insert ProjectTag: ", err)
		return err
	}

	return nil
}

// TODO: add authentication
func DeleteProjectTag(txn *sql.Tx, projectID int64, tagID int64) error {
	deleteProjectTag, err := txn.Prepare(`
		DELETE
		FROM ProjectTag
		WHERE ProjectID = ? AND TagID = ?;
	`)
	if err != nil {
		log.Println("failed to prepare statement to delete ProjectTag: ", err)
		return err
	}
	defer deleteProjectTag.Close()

	_, err = deleteProjectTag.Exec(projectID, tagID)
	if err != nil {
		log.Println("execution failed: ", err)
		return err
	}

	return nil
}

func DeleteProjectTags(txn *sql.Tx, projectID int64) error {
	deleteProjectTags, err := txn.Prepare(`
		DELETE
		FROM ProjectTag
		WHERE ProjectID = ?;
	`)
	if err != nil {
		log.Println("failed to prepare statement to delete ProjectTags for the project: ", err)
		return err
	}
	defer deleteProjectTags.Close()

	_, err = deleteProjectTags.Exec(projectID)
	if err != nil {
		log.Println("execution failed: ", err)
		return err
	}

	return nil
}
