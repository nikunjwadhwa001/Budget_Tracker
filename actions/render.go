package actions

import (
	"os"

	"github.com/gobuffalo/buffalo/render"
)

var r *render.Engine

func init() {
	r = render.New(render.Options{
		HTMLLayout:         "application.plush.html",
		TemplatesFS:        os.DirFS("templates"),
		DefaultContentType: "application/json",
	})
}
