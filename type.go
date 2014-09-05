package webui

// 保管从客户端获取的信息
type Import struct {
	Call string            `json:"call"`
	Para map[string]string `json:"param"`
}

// 保管将发给客户端的信息
type Export struct {
	Err string            `json:"error,omitempty"`
	Ans map[string]string `json:"answer"`
}

// 函数运行的环境
type Context struct {
	Hold interface{} // 本地信息
	Import
	Export
}

// 所有控件都实现的接口
type Object interface {
	// 返回该控件的HTML排布
	Format(l, t int) string

	// 返回控件的ID
	ID() string

	// 返回控件的操作函数
	DO() func(*Context)
}
