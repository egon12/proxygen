package proxygen

import (
	"strconv"
	"strings"
)

type (
	// Var is struct to store variable name and it's type
	// it can be used as Receiver, Params, Return or used in
	// Body. It doesn't need to have Name. For examples in
	// Return ussually we used Something like
	// Var{Name:"", Type:"entity.Result"}, Var{Name:"", Type:"error"}
	Var struct {
		Name string
		Type string
	}

	// MultiVar is list of Var. Ussually used in Params and Returns
	// And ussually used in function body
	MultiVar []Var

	// Func is struct that used to store data that can be used into
	// Generate function.
	Func struct {
		Name     string
		Receiver Var
		Params   MultiVar
		Return   MultiVar
		BaseType string
	}

	// FuncText is struct that will be used as input to be used in
	// func template and struct template
	FuncText struct {
		Func
		ReceiverText string
		ParamsText   string
		ReturnText   string
		ParamsNames  string
	}

	// Proxy is struct that store data that can be used to generate
	// Struct and its function
	Proxy struct {
		PackageName string
		BaseType    string
		Type        string
		Funcs       []Func
	}
)

// Text will return string that can be used in Params, Return and
// Receiver
func (v Var) Text() string {
	return v.Name + " " + v.Type
}

// Text will return string that can be used in Params, Return
func (m MultiVar) Text() string {
	vars := make([]string, len(m))
	for i, v := range m {
		vars[i] = v.Text()
	}

	return "(" + strings.Join(vars, ",") + ")"
}

// Names will return string that can be used in calling based function
func (m MultiVar) Names() string {
	vars := make([]string, len(m))
	for i, v := range m {
		vars[i] = v.Name
	}

	return "(" + strings.Join(vars, ",") + ")"
}

// FuncText generate FuncText that can be used in template
func (f Func) FuncText() FuncText {
	return FuncText{
		Func:         f,
		ReceiverText: f.Receiver.Text(),
		ParamsText:   f.Params.Text(),
		ReturnText:   f.Return.Text(),
		ParamsNames:  f.Params.Names(),
	}
}

// FixEmptyParams will set empty params names and set it with
// arg0 arg1 and so on..
func (f *Func) FixEmptyParams() {
	for i := range f.Params {
		if f.Params[i].Name == "" {
			f.Params[i].Name = "arg" + strconv.Itoa(i)
		}
	}
}

// SetRecieverTypeSuffix will set receiver for functions and struct
func (p *Proxy) SetRecieverTypeSuffix(suffix string) {
	p.Type = p.Type + suffix
	for i := range p.Funcs {
		p.Funcs[i].Receiver.Type = p.Funcs[i].Receiver.Type + suffix
	}
}

// SetBaseType will set it BaseType and its functions BaseType
func (p *Proxy) SetBaseType(baseType string) {
}
