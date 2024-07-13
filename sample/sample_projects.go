package sample

import (
	"time"

	"github.com/techiemohitjangra/portfolio/model"
)

var SampleProject1 = model.Project{
	ID:            1,
	Title:         "Project",
	Subtitle:      "subtitle",
	RepositoryURL: "github.com/techiemohitjangra",
	DeploymentURL: "www.mohitjangra.com",
	LastUpdated:   time.Now(),
	PublishedOn:   time.Now(),
	Status:        model.Complete,
	Tags: []model.Tag{
		{
			ID:       1,
			Name:     "tag_name_1",
			IconPath: ".",
		},
		{
			ID:       2,
			Name:     "tag_name_2",
			IconPath: ".",
		},
	},
	Summary:    "summary",
	Content:    "content",
	Visibility: true,
}

var SampleProject2 = model.Project{
	ID:            1,
	Title:         "Example Project",
	Subtitle:      "An Example Subtitle",
	RepositoryURL: "http://example.com/repo",
	DeploymentURL: "http://example.com/deploy",
	LastUpdated:   time.Now(),
	PublishedOn:   time.Now(),
	Status:        model.Complete,
	Tags: []model.Tag{
		{
			ID:       1,
			Name:     "Go",
			IconPath: "path/to/icon",
		},
	},
	Summary:    "Example summary",
	Content:    "Example content",
	Visibility: true,
}
