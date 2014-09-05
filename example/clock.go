package main

import (
	"github.com/hydra13142/webui"
	"net/http"
	"time"
)

func main() {
	w := &webui.Window{
		Width:  200,
		Height: 50,
		Sub: []webui.Object{
			&webui.Timer{Common: webui.Common{Id: "clock", Do: func(c *webui.Context) {
				c.Ans["text"] = time.Now().Format("2006/01/02 15:04:05.00 -0700")
			}}, Ms: 50},
			&webui.Text{Common: webui.Common{"text", "", 5, 5, 190, 40, nil}, Readonly: true},
		},
	}
	http.ListenAndServe(":9999", webui.NewHandler(w, "clock.htm", nil))
}
