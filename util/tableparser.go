package util

import (
	"fmt"
	"strconv"

	lua "github.com/yuin/gopher-lua"
)

type LuaTableParser struct {
	data *lua.LTable
}

func (p *LuaTableParser) String(key string) (s string, err error) {
	v := p.data.RawGetString(key)
	if v.Type() != lua.LTString {
		err = fmt.Errorf("incorrect lua type, expected 'LTString' but got '%s'", v.Type())
		return
	}
	s = v.String()
	return
}

func (p *LuaTableParser) Int(key string) (i int, err error) {
	v := p.data.RawGetString(key)
	if v.Type() != lua.LTNumber && v.Type() != lua.LTString {
		err = fmt.Errorf("incorrect lua type, expected 'LTString' or 'LTNumber' but got '%s'", v.Type())
		return
	}
	i, err = strconv.Atoi(v.String())
	return
}

func (p *LuaTableParser) Int64(key string) (i int64, err error) {
	v := p.data.RawGetString(key)
	if v.Type() != lua.LTNumber && v.Type() != lua.LTString {
		err = fmt.Errorf("incorrect lua type, expected 'LTString' or 'LTNumber' but got '%s'", v.Type())
		return
	}
	i, err = strconv.ParseInt(v.String(), 10, 64)
	return
}

func (p *LuaTableParser) OptionalInt(key string) (i *int, err error) {
	v := p.data.RawGetString(key)
	if v.Type() != lua.LTNumber && v.Type() != lua.LTString {
		err = fmt.Errorf("incorrect lua type, expected 'LTString' or 'LTNumber' but got '%s'", v.Type())
		return
	}
	var iVal int
	iVal, err = strconv.Atoi(v.String())
	i = &iVal
	return
}

func (p *LuaTableParser) Float64(key string) (f float64, err error) {
	v := p.data.RawGetString(key)
	if v.Type() != lua.LTNumber && v.Type() != lua.LTString {
		err = fmt.Errorf("incorrect lua type, expected 'LTString' or 'LTNumber' but got '%s'", v.Type())
		return
	}
	f, err = strconv.ParseFloat(v.String(), 64)
	return
}

func (p *LuaTableParser) Bool(key string) (b bool, err error) {
	v := p.data.RawGetString(key)
	if v.Type() != lua.LTBool {
		err = fmt.Errorf("incorrect lua type, expected 'LTBool' but got '%s'", v.Type())
		return
	}
	b, err = strconv.ParseBool(v.String())
	return
}

func (p *LuaTableParser) Table(key string) (parser *LuaTableParser, err error) {
	v := p.data.RawGetString(key)
	if v.Type() != lua.LTTable {
		err = fmt.Errorf("incorrect lua type, expected 'LTTable' but got '%s'", v.Type())
		return
	}
	parser = &LuaTableParser{v.(*lua.LTable)}
	return
}
