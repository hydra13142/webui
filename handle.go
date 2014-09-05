package webui

import (
	"code.google.com/p/go.net/websocket"
	"net/http"
	"os"
)

// 保管页面和服务端交互的信息
type Action struct {
	Call  string `json:"call,omitempty"`
	Param Param  `json:"param"`
}

func find(act map[string]func(Param) (Param, error), sub []Object) map[string]func(Param) (Param, error) {
	if act == nil {
		act = map[string]func(Param) (Param, error){}
	}
	for _, s := range sub {
		if Do := s.DO(); Do != nil {
			act[s.ID()] = Do
		}
		if c, ok := s.(*Container); ok {
			act = find(act, c.Sub)
			continue
		}
		if c, ok := s.(Container); ok {
			act = find(act, c.Sub)
			continue
		}
	}
	return act
}

// 创建一个监听服务，并返回两个Handler，一个用来管理页面访问，一个用来处理websocket
func NewHandler(win *Window, page string) *http.ServeMux {
	act := find(nil, win.Sub)
	n := http.NewServeMux()
	n.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		if _, err := os.Stat(page); err != nil {
			if !os.IsNotExist(err) {
				http.NotFound(w, req)
				return
			}
			file, err := os.Create(page)
			if err != nil {
				http.NotFound(w, req)
				return
			}
			file.WriteString(Head)
			file.WriteString(win.String())
			file.WriteString(Tail)
			file.Close()
		}
		http.ServeFile(w, req, page)
	}))
	n.Handle("/interact", websocket.Handler(func(ws *websocket.Conn) {
		data := Action{"", Param{}}
		for {
			err := websocket.JSON.Receive(ws, &data)
			if err != nil {
				break
			}
			param, err := act[data.Call](data.Param)
			if err != nil {
				data.Call = err.Error()
			} else {
				data.Call = ""
			}
			data.Param = param
			err = websocket.JSON.Send(ws, data)
			if err != nil {
				break
			}
		}
	}))
	return n
}
