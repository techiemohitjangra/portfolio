package sample

import (
	"time"

	"github.com/techiemohitjangra/portfolio/model"
)

var SampleUser = model.User{
	FirstName:             "Mohit",
	LastName:              "Jangra",
	UserName:              "mohitjangra",
	EmailAddress:          "mohitjangra12@gmail.com",
	Password:              "Hidden@mj123",
	DOB:                   time.Now(),
	About:                 "I am mohit jangra.",
	ProfilePicturePath:    ".",
	City:                  "Chandigarh",
	LastUpdated:           time.Now(),
	LastLoggedIn:          time.Now(),
	LastLoggedInLocation:  "Chandigarh",
	LastLoggedInIPAddress: "192.168.1.47",
}
