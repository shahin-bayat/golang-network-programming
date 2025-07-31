package server

import (
	"fmt"
	"reflect"
)

type Server struct {
	methods map[string]reflect.Value
}

func (s *Server) Register(name string, service any) {
	if s.methods == nil {
		s.methods = make(map[string]reflect.Value)
	}

	serviceType := reflect.TypeOf(service)
	serviceValue := reflect.ValueOf(service)

	for i := 0; i < serviceType.NumMethod(); i++ {
		method := serviceType.Method(i)
		value := serviceValue.Method(i)
		methodName := fmt.Sprintf("%s.%s", name, method.Name)
		s.methods[methodName] = value
	}
}
