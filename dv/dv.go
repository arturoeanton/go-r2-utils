package dv

import (
	"errors"
	"fmt"
	"reflect"
)

func Fill(param interface{}) error {
	ps := reflect.ValueOf(param)
	s := ps.Elem()
	t := reflect.TypeOf(param).Elem()
	if s.Kind() != reflect.Struct {
		return errors.New("dv: Fill: param is not a struct")
	}
	for i := 0; i < s.NumField(); i++ {
		f := s.Field(i)
		if f.Kind() == reflect.Struct {
			Fill(f.Addr().Interface())
			continue
		}
		if f.IsValid() {
			if f.CanSet() {
				fType := t.Field(i)
				if f.Kind() == reflect.Ptr {
					v := reflect.New(fType.Type.Elem())
					if v.Elem().Kind() == reflect.String {
						v.Elem().SetString(fType.Tag.Get("dv"))
					} else {
						fmt.Sscan(fType.Tag.Get("dv"), v.Elem().Addr().Interface())
					}
					f.Set(v)
					continue
				}
				if f.Kind() == reflect.String {
					f.SetString(fType.Tag.Get("dv"))
					continue
				}
				fmt.Sscan(fType.Tag.Get("dv"), f.Addr().Interface())
			}
		}
	}
	return nil
}
