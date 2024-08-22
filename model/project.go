package model

import (
	"database/sql"
	"log"
	"time"
)

type Status int64

const (
	Complete Status = iota
	InProgress
	Archived
	OnHold
)

func (status Status) toString() string {
	return [...]string{
		"Complete",
		"In Progress",
		"Archived",
		"On Hold",
	}[status]
}

func (status Status) EnumIndex() int64 {
	return int64(status)
}

func StatusIdentifier(statusStr string) Status {
	switch statusStr {
	case "Complete":
		return Complete
	case "In Progress":
		return InProgress
	case "Archived":
		return Archived
	case "On Hold":
		return OnHold
	}
	return -1
}

type Project struct {
	ID            int64     `param:"id" query:"id" form:"id" json:"id"`
	Title         string    `param:"title" query:"title" form:"title" json:"title"`
	Subtitle      string    `param:"subtitle" query:"subtitle" form:"subtitle" json:"subtitle"`
	RepositoryURL string    `param:"repository_url" query:"repository_url" form:"repository_url" json:"repository_url"`
	DeploymentURL string    `param:"deployment_url" query:"deployment_url" form:"deployment_url" json:"deployment_url"`
	LastUpdated   time.Time `param:"last_updated" query:"last_updated" form:"last_updated" json:"last_updated"`
	PublishedOn   time.Time `param:"published_on" query:"published_on" form:"published_on" json:"published_on"`
	Status        Status    `param:"status" query:"status" form:"status" json:"status"`
	Tags          []Tag     `param:"tags" query:"tags" form:"tags" json:"tags"`
	Summary       string    `param:"summary" query:"summary" form:"summary" json:"summary"`
	Content       string    `param:"content" query:"content" form:"content" json:"content"`
	Visibility    bool      `param:"visibility" query:"visibility" form:"visibility" json:"visibility"`
}

type ProjectItem struct {
	ID            int64     `param:"id" query:"id" form:"id" json:"id"`
	Title         string    `param:"title" query:"title" form:"title" json:"title"`
	Subtitle      string    `param:"subtitle" query:"subtitle" form:"subtitle" json:"subtitle"`
	RepositoryURL string    `param:"repository_url" query:"repository_url" form:"repository_url" json:"repository_url"`
	DeploymentURL string    `param:"deployment_url" query:"deployment_url" form:"deployment_url" json:"deployment_url"`
	LastUpdated   time.Time `param:"last_updated" query:"last_updated" form:"last_updated" json:"last_updated"`
	Status        Status    `param:"status" query:"status" form:"status" json:"status"`
	Tags          []Tag     `param:"tags" query:"tags" form:"tags" json:"tags"`
	Visibility    bool      `param:"visibility" query:"visibility" form:"visibility" json:"visibility"`
}

// TODO: add authentication
func GetProject(db *sql.DB, projectID int64) (*Project, error) {
	txn, err := db.Begin()
	if err != nil {
		log.Println("failed to start transaction: ", err)
		return nil, err
	}
	defer txn.Rollback()

	getProject, err := txn.Prepare(`
		SELECT ID, Title, Subtitle, RepositoryURL, DeploymentURL, LastUpdated, PublishedOn, Status, Summary, Content, Visibility
		FROM Project
		WHERE ID = ?;
	`)
	if err != nil {
		log.Println("failed to prepare statement to get project by ID: ", err)
		return nil, err
	}
	defer getProject.Close()

	var project Project
	var lastUpdated int64
	var publishedOn int64
	var status string

	err = getProject.QueryRow(projectID).Scan(
		&project.ID,
		&project.Title,
		&project.Subtitle,
		&project.RepositoryURL,
		&project.DeploymentURL,
		&lastUpdated,
		&publishedOn,
		&status,
		&project.Summary,
		&project.Content,
		&project.Visibility,
	)
	if err != nil {
		log.Println("failed to fetch and parse Project: ", err)
		return nil, err
	}

	project.LastUpdated = time.Unix(lastUpdated, 0)
	project.PublishedOn = time.Unix(publishedOn, 0)
	project.Status = StatusIdentifier(status)

	tags, err := GetProjectTags(txn, project.ID)
	if err != nil {
		log.Println("failed to get ProjectTags: ", err)
		return nil, err
	}
	project.Tags = tags

	err = txn.Commit()
	if err != nil {
		log.Println("failed to commit transaction: ", err)
		return nil, err
	}
	return &project, nil
}

// TODO: add authentication
func GetProjects(db *sql.DB) ([]ProjectItem, error) {
	txn, err := db.Begin()
	if err != nil {
		log.Println("failed to start transaction: ", err)
		return nil, err
	}
	defer txn.Rollback()

	getProjects, err := txn.Prepare(`
		SELECT ID, Title, Subtitle, RepositoryURL, DeploymentURL, LastUpdated, Status, Visibility
		FROM Project;
	`)
	if err != nil {
		log.Println("failed to prepare statement to get Project: ", err)
		return nil, err
	}
	defer getProjects.Close()

	var projects []ProjectItem

	rows, err := getProjects.Query()
	if err != nil {
		log.Println("failed to query projects for ProjectItems: ", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var project ProjectItem
		var lastUpdated int64
		var status string

		err = rows.Scan(
			&project.ID,
			&project.Title,
			&project.Subtitle,
			&project.RepositoryURL,
			&project.DeploymentURL,
			&lastUpdated,
			&status,
			&project.Visibility,
		)
		if err != nil {
			log.Println("failed to parse ProjectItem: ", err)
			return nil, err
		}
		project.LastUpdated = time.Unix(lastUpdated, 0)
		project.Status = StatusIdentifier(status)
		tags, err := GetProjectTags(txn, project.ID)
		if err != nil {
			log.Println("failed to fetch and parse ProjectTags: ", err)
			return nil, err
		}
		project.Tags = tags
		projects = append(projects, project)
	}

	err = txn.Commit()
	if err != nil {
		log.Println("failed to commit transaction: ", err)
		return nil, err
	}
	return projects, nil
}

// TODO: add authentication
func AddProject(db *sql.DB, project Project) (int64, error) {
	txn, err := db.Begin()
	if err != nil {
		log.Println("failed to start transaction: ", err)
		return -1, err
	}
	defer txn.Rollback()

	createProject, err := txn.Prepare(`
		INSERT OR IGNORE INTO Project (Title, Subtitle, RepositoryURL, DeploymentURL, LastUpdated, PublishedOn, Status, Summary, Content, Visibility)
		VALUES (?,?,?,?,?,?,?,?,?,?);
	`)
	if err != nil {
		log.Println("failed to create statement to create Project: ", err)
		return -1, err
	}
	defer createProject.Close()

	result, err := createProject.Exec(
		project.Title,
		project.Subtitle,
		project.RepositoryURL,
		project.DeploymentURL,
		project.LastUpdated.Unix(),
		project.PublishedOn.Unix(),
		project.Status.toString(),
		project.Summary,
		project.Content,
		project.Visibility,
	)
	if err != nil {
		log.Println("execution failed: ", err)
		return -1, err
	}

	projectID, err := result.LastInsertId()
	if err != nil {
		log.Println("failed to get ID of the last inserted Project: ", err)
		return -1, err
	}

	for _, tag := range project.Tags {
		err := AddProjectTag(txn, projectID, tag)
		if err != nil {
			log.Println("failed to add AddProjectTag: ", err)
			return -1, err
		}
	}

	err = txn.Commit()
	if err != nil {
		log.Println("failed to commit transaction: ", err)
		return -1, err
	}
	return projectID, nil
}

// TODO: add authentication
func UpdateProject(db *sql.DB, project Project) (int64, error) {
	txn, err := db.Begin()
	if err != nil {
		log.Println("failed to start transaction: ", err)
		return -1, err
	}
	defer txn.Rollback()

	updateProject, err := txn.Prepare(`
		UPDATE Project
		SET Title = ?, Subtitle = ?, RepositoryURL = ?, DeploymentURL = ?, LastUpdated = ?, PublishedOn = ?, Status = ?, Summary = ?, Content = ?, Visibility = ?
		WHERE ID = ?;
	`)
	if err != nil {
		log.Println("failed to prepare statement to update Project: ", err)
		return -1, err
	}
	defer updateProject.Close()

	_, err = updateProject.Exec(
		project.Title,
		project.Subtitle,
		project.RepositoryURL,
		project.DeploymentURL,
		project.LastUpdated.Unix(),
		project.PublishedOn.Unix(),
		project.Status.toString(),
		project.Summary,
		project.Content,
		project.Visibility,
	)
	if err != nil {
		log.Println("execution failed: ", err)
		return -1, err
	}

	existingTagIDs := make(map[int64]bool)
	for _, tag := range project.Tags {
		err := AddProjectTag(txn, project.ID, tag)
		if err != nil {
			log.Println("failed to add ProjectTag: ", err)
			return -1, err
		}
		existingTagIDs[tag.ID] = true
	}

	// delete removed ProjectTags
	deleteProjectTags, err := txn.Prepare(`
		SELECT TagID
		FROM ProjectTag
		WHERE ProjectID = ?;
	`)
	defer deleteProjectTags.Close()

	projectTagRows, err := deleteProjectTags.Query(project.ID)
	if err != nil {
		log.Println("failed to get tagIDs of projectTags: ", err)
		return -1, err
	}
	defer projectTagRows.Close()

	for projectTagRows.Next() {
		var tagID int64
		err := projectTagRows.Scan(&tagID)
		if err != nil {
			log.Println("failed to parse tagID: ", err)
			return -1, err
		}
		if !existingTagIDs[tagID] {
			err := DeleteProjectTag(txn, project.ID, tagID)
			if err != nil {
				log.Println("failed to delete ProjectTag: ", err)
				return -1, err
			}
		}
	}

	err = txn.Commit()
	if err != nil {
		log.Println("failed to commit transaction: ", err)
		return -1, err
	}
	return project.ID, nil
}

// TODO: add authentication
func DeleteProject(db *sql.DB, projectID int64) error {
	txn, err := db.Begin()
	if err != nil {
		log.Println("failed to start transaction: ", err)
		return err
	}
	defer txn.Rollback()

	err = DeleteProjectTags(txn, projectID)
	if err != nil {
		log.Println("failed to delete all the ProjectTags for the Project:", err)
		return err
	}

	deleteProject, err := txn.Prepare(`
		DELETE
		FROM Projects
		WHERE ID = ?;
	`)
	if err != nil {
		log.Println("failed to prepare statement to delete Project: ", err)
		return err
	}
	defer deleteProject.Close()

	_, err = deleteProject.Exec()
	if err != nil {
		log.Println("execution failed: could not delete Project: ", err)
	}

	err = txn.Commit()
	if err != nil {
		log.Println("failed to commit transaction: ", err)
		return err
	}

	return nil
}
