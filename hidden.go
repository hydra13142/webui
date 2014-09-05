package webui

import "fmt"

// 隐藏的对象，用于保管一些持久信息
type Hidden struct {
	Common
}

func (a Hidden) Format(_, _ int) string {
	return fmt.Sprintf(`<textarea id="%s" value="%s" style="visible:hidden"></textarea>`+"\n", a.Id, a.Value)
}

// 用于进行定时同步
type Timer struct {
	Common
	Ms int
}

func (a Timer) Format(_, _ int) string {
	if a.Ms == 0 {
		return ""
	}
	return fmt.Sprintf(`<div id="%s" style="visible:hidden"><script type="text/javascript">`+"\n"+
		`window.setInterval(function(){myfunc(document.getElementById("%s"));},%d);`+"\n"+
		"</script></div>\n", a.Id, a.Id, a.Ms)
}
