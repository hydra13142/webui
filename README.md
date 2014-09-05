webui
=====

用网页模拟的go语言gui

本库是采用网页模拟GUI，以便进行交互。

类型Window，用来盛放控件（个控件都对应某个具体的HTML控件）。

支持如下控件

    Button    按钮
    Radio     单选框
    Check     复选框
    Select    选择列表
    Text      文本框
    Image     图像框
    Container 容器（可以盛放控件，或用于嵌入innerHTML）

以上控件都实现了Object接口：

	// 所有控件都实现的接口
	type Object interface {
		// 返回该控件的HTML排布
		Format(l, t int) string

		// 返回控件的ID
		ID() string

		// 返回控件的操作函数
		DO() func(*Context)
	}

Window和Container都具有[]Object类型的Sub字段，可以放入其他控件。

用于交互的数据类型

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

所有控件都包含有Common字段，Common字段如下：

    type Common struct {
        Id, Value     string
        Left, Top     int
        Width, Height int
        Do            func(*Context)
    }
	
控件的操作函数func(*Context)可以处理输入生成输出。输入可以从Call和Para字段获取
；输出应写入Ans和Err字段。只需要需要修改的控件和属性写入Ans字段，如果出错可以
将错误信息写入Err字段。Hold字段用于保管一些持久性信息。

因此，大小、位置、ID、值和操作函数是所有控件都有的属性。

Param类型即map[string]string，用来供函数获取参数信息，里面是控件的id到值的映射。

    注意，复选框的值是被选择的各项的值用"|"连接起来的结果。

使用时，会用到NewHandler函数：

    func NewHandler(*Window, string, func()interface{}) *http.ServeMux
  
该函数会返回一个handler，用它即可实现与页面的交互。

第一个参数表示窗口对象，其内包括多个控件。

第二个参数表示页面文件，如果不存在，函数会自动根据第一个参数生成一个该名字的文件；
如果已经存在，则会使用该文件（不做修改）。这是为了让使用者可以自行设计页面细节，
只要不修改控件的id和值，就不会对交互造成影响。

最后一个参数用于生成连接专一的内部参数，该参数会保存如Context.Hold字段，可以用类
型推断取出。建议该函数参数返回指针。

用法举例：

    http.ListenAndServe(":9999", webui.NewHandler(w, "index.htm", nil))
