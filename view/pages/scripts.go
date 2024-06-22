package pages

import (
	component "github.com/techiemohitjangra/portfolio/view/component"
)

var TailwindCSS = component.Script{
	Src:   "https://cdn.tailwindcss.com",
	Defer: true,
	Async: true,
	Type:  "text/javascript",
}
