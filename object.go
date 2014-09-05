// Web UI package
package webui

import "fmt"

// Param用来保管参数信息
type Param map[string]string

// 所有控件都实现的接口
type Object interface {

	// 返回该控件的HTML排布
	Format(l, t int) string

	// 返回控件的ID
	ID() string

	// 返回控件的操作函数
	DO() func(Param) (Param, error)
}

// 所有控件都有的属性
type Common struct {
	Id, Value     string
	Left, Top     int
	Width, Height int
	Do            func(Param) (Param, error)
}

func (a Common) Format(l, t int) string {
	return fmt.Sprintf(`id="%s" value="%s" style="position:absolute; left:%d; top:%d; width:%d; height:%d" `, a.Id, a.Value, a.Left+l, a.Top+t, a.Width, a.Height)
}
func (a Common) ID() string {
	return a.Id
}
func (a Common) DO() func(Param) (Param, error) {
	return a.Do
}

// 按钮
type Button struct {
	Common
}

func (a Button) Format(l, t int) string {
	s := `<input type="button" ` + a.Common.Format(l, t)
	if a.Do != nil {
		s += `onClick="javascript:myfunc(this);" `
	}
	s += "/>\n"
	return s
}

// 单选框
type Radio struct {
	Common
	Flag     []string
	Multirow bool
}

func (a Radio) Format(l, t int) string {
	s := `<form type="radio" ` + a.Common.Format(l, t)
	if a.Do != nil {
		s += `onChange="javascript:myfunc(this);" `
	}
	s += ">\n"
	if a.Multirow {
		for _, f := range a.Flag {
			s += fmt.Sprintf(`<label><input type="radio" value="%s" />%s</label><br/>`+"\n", f, f)
		}
	} else {
		for _, f := range a.Flag {
			s += fmt.Sprintf(`<label><input type="radio" value="%s" />%s</label>`+"\n", f, f)
		}
	}
	s += "</form>\n"
	return s
}

// 复选框
type Check struct {
	Common
	Flag     []string
	Multirow bool
}

func (a Check) Format(l, t int) string {
	s := `<form type="check" ` + a.Common.Format(l, t)
	if a.Do != nil {
		s += `onChange="javascript:myfunc(this);" `
	}
	s += ">\n"
	if a.Multirow {
		for _, f := range a.Flag {
			s += fmt.Sprintf(`<label><input type="checkbox" value="%s" />%s</label><br/>`+"\n", f, f)
		}
	} else {
		for _, f := range a.Flag {
			s += fmt.Sprintf(`<label><input type="checkbox" value="%s" />%s</label>`+"\n", f, f)
		}
	}
	s += "</form>\n"
	return s
}

// 选择列表
type Select struct {
	Common
	Flag []string
	Menu bool
}

func (a Select) Format(l, t int) string {
	var s string
	if a.Menu {
		s = `<select type="select" ` + a.Common.Format(l, t)
		if a.Do != nil {
			s += `onSelect="javascript:myfunc(this);" `
		}
		s += ">\n"
		for _, f := range a.Flag {
			s += fmt.Sprintf(`<option>%s</option>`+"\n", f)
		}
	} else {
		s = fmt.Sprintf(`<select size="%d" `, len(a.Flag)) + a.Common.Format(l, t) + `>`
		for _, f := range a.Flag {
			s += fmt.Sprintf(`<option>%s</option>`+"\n", f)
		}
	}
	s += "</select>\n"
	return s
}

// 文本框（单行、多行、密码）
type Text struct {
	Common
	Password bool
	Multirow bool
	Autofold bool
	Readonly bool
}

func (a Text) Format(l, t int) string {
	s := ""
	if a.Password {
		s = `<input type="password" ` + a.Common.Format(l, t)
	} else {
		if a.Multirow {
			if a.Autofold {
				s = `<textarea type="area" wrap="physical" ` + a.Common.Format(l, t)
			} else {
				s = `<textarea type="area" wrap="off" ` + a.Common.Format(l, t)
			}
		} else {
			s = `<input type="text" ` + a.Common.Format(l, t)
			if a.Readonly {
				s += `readonly="readonly" `
			}
		}
	}
	if a.Do != nil {
		s += `onChange="javascript:myfunc(this)" `
	}
	s += "/>\n"
	return s
}

// 标签
type Label struct {
	Common
}

func (a Label) Format(l, t int) string {
	s := `<label type="label" ` + a.Common.Format(l, t) + ">" + a.Value + "</label>\n"
	return s
}

// 图像框
type Image struct {
	Common
}

func (a Image) Format(l, t int) string {
	return `<input type="image" border="1" ` + a.Common.Format(l, t) + " />\n"
}

// 容器
type Container struct {
	Common
	Sub []Object
}

func (a Container) Format(l, t int) string {
	s := `<div type="container" ` + a.Common.Format(l, t) + ">\n"
	for _, f := range a.Sub {
		s += f.Format(l, t)
	}
	s += "</div>\n"
	return s
}

// 窗体
type Window struct {
	Width, Height int
	Sub           []Object
}

func (a Window) String() string {
	s := fmt.Sprintf("<form>\n"+`<div style="position:relative; margin:auto; width:%d; height:%d; border-style:solid; border-width:1px; border-color:#000">`+"\n", a.Width, a.Height)
	for _, f := range a.Sub {
		s += f.Format(0, 0)
	}
	s += "</div>\n</form>\n"
	return s
}

const (
	Head = `<html>
<head>
<script type="text/javascript">
var ws;
if(!("WebSocket" in window))
{
	alert("unsupport websocket!");
}
else
{
	ws = new WebSocket("ws://%s");
	ws.onopen = function()
	{
		// alert("ready to go!");
	};
	ws.onmessage = function(m)
	{
		var o, e = JSON.parse(m.data);
		if(e.call!=undefined)
		{
			alert(e.call);
			return;
		}
		e = e.param;
		for(var key in e)
		{
			o = document.getElementById(key);
			if(o.type==undefined)
			{
				continue;
			}
			if(o.type=="radio" || o.type=="check")
			{
				continue;
			}
			if(o.type=="container")
			{
				o.innerHTML=e[key];
			}
			else
			{
				o.value=e[key];
			}
		}
	};
	ws.onclose = function()
	{
		// alert("connection is closed");
	};
}
function findset(e)
{
	var n = e.childNodes;
	var s = new Array();
	if(e.type != "radio" || e.type != "check")
	{
		return;
	}
	for(i=0; i<n.length; i++)
	{
		if(n[i].checked)
		{
			s.push(n[i].value)
		}
	}
	return s.join("|");
}
function myfunc(e)
{
	var o, m, s = {};
	o = document.getElementsByTagName("input");
	for(var i=0; i<o.length;i++)
	{
		if(o[i].type=="text" || o[i].type=="password")
		{
			s[o[i].id]=o[i].value;
		}
	}
	o = document.getElementsByTagName("textarea");
	for(var i=0; i<o.length;i++)
	{
		s[o[i].id]=o[i].value;
	}
	o = document.getElementsByTagName("select");
	for(var i=0; i<o.length;i++)
	{
		s[o[i].id]=o[i].value;
	}
	o = document.getElementsByTagName("form");
	for(var i=0; i<o.length;i++)
	{
		m = findset(o[i]);
		if(m!=undefined)
		{
			s[o[i].id]=m;
		}
	}
	ws.send(JSON.stringify({"call":e.id,"param":s}));
}
</script>
</head>
<body>
`
	Tail = `</body>
</html>`
)
