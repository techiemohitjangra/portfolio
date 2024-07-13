package sample

import (
	"time"

	"github.com/techiemohitjangra/portfolio/model"
)

var SampleBlog = model.Blog{
	ID:          1,
	Title:       "blogTitle",
	Subtitle:    "blogSubtitle",
	PublishedOn: time.Now(),
	LastUpdated: time.Now(),
	References: []model.Reference{
		{
			ID:     1,
			Title:  "referenceTitle1",
			URL:    "url1",
			BlogID: 1,
		},
		{
			ID:     2,
			Title:  "referenceTitle2",
			URL:    "url2",
			BlogID: 1,
		},
	},
	Links: []model.Link{
		{
			ID:     1,
			Text:   "linkText1",
			URL:    "linkURL1",
			BlogID: 1,
		},
		{
			ID:     2,
			Text:   "linkText2",
			URL:    "linkURL2",
			BlogID: 1,
		},
	},
	Content: "blogContent",
	Tags: []model.Tag{
		{
			ID:       1,
			Name:     "Python",
			IconPath: ".",
		},
		{
			ID:       2,
			Name:     "C++",
			IconPath: ".",
		},
	},
}
