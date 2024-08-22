package model

import (
	"database/sql"
	"log"
	"time"
)

type Blog struct {
	ID          int64       `param:"id" query:"id" form:"id" json:"id"`
	Title       string      `param:"title" query:"title" form:"title" json:"title"`
	Subtitle    string      `param:"subtitle" query:"subtitle" form:"subtitle" json:"subtitle"`
	PublishedOn time.Time   `param:"published_on" query:"published_on" form:"published_on" json:"published_on"`
	LastUpdated time.Time   `param:"last_updated" query:"last_updated" form:"last_updated" json:"last_updated"`
	References  []Reference `param:"references" query:"references" form:"references" json:"references"`
	Links       []Link      `param:"link" query:"link" form:"link" json:"link"` // link to repo/code
	Content     string      `param:"content" query:"content" form:"content" json:"content"`
	Tags        []Tag       `param:"tags" query:"tags" form:"tags" json:"tags"` // list of tags
}

type BlogItem struct {
	ID          int64     `param:"id" query:"id" form:"id" json:"id"`
	Title       string    `param:"title" query:"title" form:"title" json:"title"`
	Subtitle    string    `param:"subtitle" query:"subtitle" form:"subtitle" json:"subtitle"`
	PublishedOn time.Time `param:"published_on" query:"published_on" form:"published_on" json:"published_on"`
	LastUpdated time.Time `param:"last_updated" query:"last_updated" form:"last_updated" json:"last_updated"`
	Tags        []Tag     `param:"tags" query:"tags" form:"tags" json:"tags"` // list of tags
}

// TODO: add authentication
// if authenticated present options to edit/update blog
func GetBlog(db *sql.DB, id int64) (*Blog, error) {
	txn, err := db.Begin()
	if err != nil {
		log.Println("failed to start transaction: ", err)
		return nil, err
	}
	defer txn.Rollback()

	getBlog, err := txn.Prepare(`
		SELECT ID, Title, Subtitle, PublishedOn, LastUpdated, Content
		FROM Blog
		WHERE ID = ?;
	`)
	if err != nil {
		log.Println("failed to prepare statement to get blog by id: ", err)
		return nil, err
	}
	defer getBlog.Close()

	var blog Blog
	var publishDate int64
	var lastUpdateTime int64

	err = getBlog.QueryRow(id).Scan(
		&blog.ID,
		&blog.Title,
		&blog.Subtitle,
		&publishDate,
		&lastUpdateTime,
		&blog.Content,
	)
	if err != nil {
		log.Println("failed to query blog: ", err)
		return nil, err
	}

	blog.PublishedOn = time.Unix(publishDate, 0)
	blog.LastUpdated = time.Unix(lastUpdateTime, 0)

	references, err := GetReferences(txn, blog.ID)
	if err != nil {
		log.Println("failed to fetch references: ", err)
		return nil, err
	}
	links, err := GetLinks(txn, blog.ID)
	if err != nil {
		log.Println("failed to fetch links: ", err)
		return nil, err
	}
	tags, err := GetBlogTags(txn, blog.ID)
	if err != nil {
		log.Println("failed to fetch BlogTags: ", err)
		return nil, err
	}

	blog.References = references
	blog.Links = links
	blog.Tags = tags

	err = txn.Commit()
	if err != nil {
		log.Println("failed to commit transaction: ", err)
		return nil, err
	}

	return &blog, nil
}

// TODO: add authentication
// if authenticated present option to delete/update blog from the list
func GetBlogs(db *sql.DB) ([]BlogItem, error) {
	txn, err := db.Begin()
	if err != nil {
		log.Println("failed to start transaction: ", err)
		return nil, err
	}
	defer txn.Rollback()

	getBlogs, err := txn.Prepare(`
		SELECT ID, Title, Subtitle, PublishedOn, LastUpdated
		FROM Blog;
	`)
	if err != nil {
		log.Println("failed to prepare statement to get blogs: ", err)
		return nil, err
	}
	defer getBlogs.Close()

	rows, err := getBlogs.Query()
	if err != nil {
		log.Println("failed to query list of blogs: ", err)
		return nil, err
	}
	defer rows.Close()

	var blogItems []BlogItem

	for rows.Next() {
		var blogItem BlogItem
		var publishedOn int64
		var lastUpdated int64

		err = rows.Scan(
			&blogItem.ID,
			&blogItem.Title,
			&blogItem.Subtitle,
			&publishedOn,
			&lastUpdated,
		)
		if err != nil {
			log.Println("failed to parse data from query to blog object: ", err)
			return nil, err
		}

		blogItem.PublishedOn = time.Unix(publishedOn, 0)
		blogItem.LastUpdated = time.Unix(lastUpdated, 0)

		tags, err := GetBlogTags(txn, blogItem.ID)
		if err != nil {
			log.Println("failed to fetch blogTags: ", err)
			return nil, err
		}
		blogItem.Tags = tags

		blogItems = append(blogItems, blogItem)

	}
	if err = rows.Err(); err != nil {
		log.Println("failed to parse query result to BlogItems: ", err)
		return nil, err
	}

	err = txn.Commit()
	if err != nil {
		log.Println("failed to commit transaction: ", err)
		return nil, err
	}
	return blogItems, nil
}

// TODO: add authentication
func AddBlog(db *sql.DB, blog Blog) (int64, error) {
	txn, err := db.Begin()
	if err != nil {
		log.Println("failed to start transaction to add blog: ", err)
		return -1, err
	}
	defer txn.Rollback()

	// Create Blog
	createBlog, err := txn.Prepare(`
		INSERT INTO Blog (Title, Subtitle, PublishedOn, LastUpdated, Content)
		VALUES (?, ?, ?, ?, ?);
	`)
	if err != nil {
		log.Println("failed to prepare statement to insert Blog: ", err)
		return -1, err
	}
	defer createBlog.Close()

	result, err := createBlog.Exec(blog.Title, blog.Subtitle, blog.PublishedOn.Unix(), blog.LastUpdated.Unix(), blog.Content)
	if err != nil {
		log.Println("execution failed: ", err)
		return -1, err
	}

	blogID, err := result.LastInsertId()
	if err != nil {
		log.Println("failed to get last inserted Blog's id: ", err)
		return -1, err
	}

	// Insert Links for the Blog
	for _, link := range blog.Links {
		link.BlogID = blogID
		err = AddLink(txn, link)
		if err != nil {
			log.Println("failed to add link: ", err)
			return -1, err
		}
	}

	// Insert References for the Blog
	for _, reference := range blog.References {
		reference.BlogID = blogID
		err = AddReference(txn, reference)
		if err != nil {
			log.Println("failed to add reference: ", err)
			return -1, err
		}
	}

	// Insert Tags for the Blog
	for _, tag := range blog.Tags {
		err = AddBlogTag(txn, blogID, tag)
		if err != nil {
			log.Println("failed to add BlogTag: ", err)
			return -1, err
		}
	}

	err = txn.Commit()
	if err != nil {
		log.Println("failed to commit transaction: ", err)
		return -1, err
	}

	return blogID, nil
}

// TODO: add authentication
func UpdateBlog(db *sql.DB, blog Blog) error {
	txn, err := db.Begin()
	if err != nil {
		log.Println("failed to start transaction: ", err)
		return err
	}
	defer txn.Rollback()

	updateBlog, err := txn.Prepare(`
		UPDATE Blog
		SET Title = ?, Subtitle = ?, PublishedOn = ?, LastUpdated = ?, Content = ?
		WHERE ID = ?;
	`)
	if err != nil {
		log.Println("failed to prepare statement to updateBlog: ", err)
		return err
	}
	defer updateBlog.Close()

	_, err = updateBlog.Exec(
		blog.Title,
		blog.Subtitle,
		blog.PublishedOn.Unix(),
		blog.LastUpdated.Unix(),
		blog.Content,
	)
	if err != nil {
		log.Println("execution failed: ", err)
		return err
	}

	// Process references
	existingReferenceIDs := make(map[int64]bool)
	for _, ref := range blog.References {
		if ref.ID == 0 {
			err := AddReference(txn, ref)
			if err != nil {
				log.Println("failed to add reference: ", err)
				return nil
			}
		} else {
			err := UpdateReference(txn, ref)
			if err != nil {
				log.Println("failed to update reference: ", err)
				return nil
			}
		}
		existingReferenceIDs[ref.ID] = true
	}

	// get blog References
	blogReferences, err := txn.Prepare(`
		SELECT ID
		FROM Reference
		WHERE BlogID = ?;
	`)
	if err != nil {
		log.Println("failed to prepare statement to delete removed References: ", err)
		return err
	}
	defer blogReferences.Close()

	referenceRows, err := blogReferences.Query(blog.ID)
	if err != nil {
		log.Println("failed to query References:", err)
	}
	defer referenceRows.Close()

	// delete old reference not present in new data
	for referenceRows.Next() {
		var id int64
		err = referenceRows.Scan(&id)
		if err != nil {
			log.Println("failed to parse Reference ID: ", err)
			return err
		}
		if !existingReferenceIDs[id] {
			err = DeleteReference(txn, blog.ID, id)
		}
	}

	// Process links
	existingLinkIDs := make(map[int64]bool)
	for _, link := range blog.Links {
		if link.ID == 0 {
			err := AddLink(txn, link)
			if err != nil {
				log.Println("failed to add link: ", err)
				return err
			}
		} else {
			err := UpdateLink(txn, link)
			if err != nil {
				log.Println("failed to update link: ", err)
				return err
			}
		}
		existingLinkIDs[link.ID] = true
	}

	// fetch blog links
	blogLinks, err := txn.Prepare(`
		SELECT ID
		FROM Link
		WHERE BlogID = ?;
	`)
	if err != nil {
		log.Println("failed to prepare statement to get Links: ", err)
		return err
	}
	defer blogLinks.Close()

	linkRows, err := blogLinks.Query(blog.ID)
	if err != nil {
		log.Println("execution failed: ", err)
		return err
	}
	defer linkRows.Close()

	// delete old links not present in new data
	for linkRows.Next() {
		var id int64
		err := linkRows.Scan(&id)
		if err != nil {
			log.Println("failed to parse Link ID: ", err)
			return err
		}
		if !existingLinkIDs[id] {
			err := DeleteLink(txn, blog.ID, id)
			if err != nil {
				log.Println("failed to delete Link: ", err)
				return err
			}
		}
	}

	blogTags, err := txn.Prepare(`
		SELECT TagID
		FROM BlogTag
		WHERE BlogID = ?;
	`)
	if err != nil {
		log.Println("failed to prepare statement to get TagIDs for the Blog: ", err)
		return err
	}
	defer blogTags.Close()

	tagRows, err := blogTags.Query(blog.ID)
	if err != nil {
		log.Println()
		return err
	}
	defer tagRows.Close()

	for tagRows.Next() {
		var tagID int64
		err := tagRows.Scan(&tagID)
		if err != nil {
			log.Println("failed to parse tagID: ", err)
			return err
		}
		if !existingLinkIDs[tagID] {
			err := DeleteBlogTag(txn, blog.ID, tagID)
			if err != nil {
				log.Println("failed to delete blogTag: ", err)
				return err
			}
		}
	}

	err = txn.Commit()
	if err != nil {
		log.Println("failed to commit the transaction: ", err)
		return err
	}

	return nil
}

// TODO: add authentication
// deletes blogs and tags relations for that blog
func DeleteBlog(db *sql.DB, blogID int64) error {
	txn, err := db.Begin()
	if err != nil {
		log.Println("failed to begin transaction: ", err)
		return err
	}
	defer txn.Rollback()

	// delete tags from blog_tags table
	err = DeleteBlogTags(txn, blogID)
	if err != nil {
		log.Println("failed to delete all BlogTags for the blog: ", err)
		return err
	}

	err = DeleteLinks(txn, blogID)
	if err != nil {
		log.Println("failed to delete all the Links for the blog: ", err)
		return err
	}

	err = DeleteReferences(txn, blogID)
	if err != nil {
		log.Println("failed to delete all the References for the blog: ", err)
	}

	// delete blog from blogs table
	deleteBlog, err := txn.Prepare(`
		DELETE
		FROM Blog
		WHERE ID = ?;
	`)
	if err != nil {
		log.Println("failed to prepare statement: ", err)
		return err
	}
	defer deleteBlog.Close()

	_, err = deleteBlog.Exec(blogID)
	if err != nil {
		log.Println("execution failed: ", err)
		return err
	}

	err = txn.Commit()
	if err != nil {
		log.Println("failed to commit transaction: ", err)
		return err
	}
	return nil
}
