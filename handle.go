package webui

import (
	"code.google.com/p/go.net/websocket"
	"net/http"
	"os"
)

func find(act map[string]func(*Context), sub []Object) map[string]func(*Context) {
	if act == nil {
		act = map[string]func(*Context){}
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
// 本函数最后一个参数用于生成连接专一的内部参数，建议该函数参数返回指针
// 该参数会保存如Context.Hold字段，可以用类型推断取出
func NewHandler(win *Window, page string, loc func() interface{}) *http.ServeMux {
	act := find(nil, win.Sub)
	mux := http.NewServeMux()
	mux.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
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
	mux.Handle("/interact", websocket.Handler(func(ws *websocket.Conn) {
		env := &Context{}
		if loc != nil {
			env.Hold = loc()
		}
		for {
			env.Para = map[string]string{}
			env.Call = ""
			err := websocket.JSON.Receive(ws, &env.Import)
			if err != nil {
				break
			}
			env.Ans = map[string]string{}
			env.Err = ""
			act[env.Import.Call](env)
			err = websocket.JSON.Send(ws, env.Export)
			if err != nil {
				break
			}
		}
	}))
	return mux
}
