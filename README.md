webui
=====

用网页模拟的go语言gui

本库是采用网页模拟GUI，以便进行交互。

类型Window，用来盛放控件（个控件都对应某个具体的HTML控件）。

支持如下控件

  Button 按钮
  Radio 单选框
  Check 复选框
  Select 选择列表
  Text 文本框
  Image 图像框
  Container 容器（可以盛放控件，或者用于嵌入innerHTML）

以上控件都实现了如下三个方法：

  // 生成排布html代码
  Format(l,t int) string
  // 返回控件id
  ID() string
  // 返回控件操作函数
  DO() func(Param) (Param, error)
  
实现了这三个方法，即实现了Object接口。

Window和Container都具有[]Object类型Sub字段，可以放入其他控件。

所有控件都包含有Common字段，Common字段如下：

  type Common struct {
    Id, Value     string
    Left, Top     int
    Width, Height int
    Do            func(Param) (Param, error)
  }
  
因此，大小、位置、ID、值和操作函数是所有控件都有的属性。

Param类型即map[string]string，用来供函数获取参数信息，里面是控件的id到值的映射。

  注意，复选框的值是被选择的各项的值用"|"连接起来的结果。

使用时，会用到NewHandler函数：

  func NewHandler(*Window, page, wsp string) (http.Handler, websocket.Handler)
  
该函数会返回两个handler，用它们即可实现与页面的交互。

page表示页面文件，如果不存在，函数会自动根据Window及其下控件的信息生成一个该名字的文件；
如果已经存在，则会使用该文件（不做修改）。这是为了让使用者可以自行设计页面细节，只要不
修改控件的id和值，就不会对交互造成影响。

wsp表示websocket的连接地址，因为某些原因发现websocket的Handler没法挂载到http.ServeMux上；
只能让它独立出来开两个服务端处理。

  h, s := webui.NewHandler(w, "index.htm", "127.0.0.1:8888")
  go http.ListenAndServe(":8888", s)
  http.ListenAndServe(":9999", h)
