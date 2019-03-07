package container

import (
	"fmt"
	"reflect"
	"sync"
)

type Closable interface {
	Close()
}

type Registry struct {
	sync.Mutex
	registry        map[string]interface{}
	registerTypes   map[string]reflect.Type
	registerRefVals map[string]reflect.Value
	typeSpecMap     map[string]*TypeSpec
	started         bool
	closed          bool
}

func NewRegistry() *Registry {
	return &Registry{
		registry:        make(map[string]interface{}),
		registerTypes:   make(map[string]reflect.Type),
		registerRefVals: make(map[string]reflect.Value),
		typeSpecMap:     make(map[string]*TypeSpec),
		started:         false,
		closed:          false,
	}
}

func newStruct(refType reflect.Type) interface{} {
	return reflect.New(refType).Elem().Interface()
}

func fullTypeName(refType reflect.Type) string {
	return fmt.Sprintf("%s/%s", refType.PkgPath(), refType.Name())
}

//
func (self *Registry) RegisterType(elem interface{}) error {
	refType := reflect.TypeOf(elem)
	if refType.Kind() != reflect.Struct {
		return fmt.Errorf("register target must be struct")
	}

	self.doRegister(refType)
	return nil
}

func (self *Registry) doRegister(refType reflect.Type) {
	elemType := refType.Elem()
	fullTypeName := fullTypeName(elemType)

	self.Lock()
	defer self.Unlock()
	if _, ok := self.registerTypes[fullTypeName]; !ok {
		//self.registerTypes[fullTypeName] = refType
		self.registerTypes[fullTypeName] = elemType
		self.typeSpecMap[fullTypeName] = NewTypeSpec(refType)
	}
}

func (self *Registry) GetObject(refType reflect.Type) interface{} {
	return self.registry[fullTypeName(refType)]
}

func (self *Registry) Start() error {
	self.Lock()
	defer self.Unlock()
	if !self.started {
		self.started = true
		self.initRegisterType()
		err := self.inject()

		if err != nil {
			return err
		}
	}
	//for name, define := range self.registerTypes {
	//	return self.createObject(define, name)
	//}

	return nil
}

//func (self *Registry) createObject(d *TypeSpec, name string) error {
//	if !d.inited {
//		val := newStruct(d.refType)
//		refVal := reflect.ValueOf(val).Elem()
//		for _, field := range d.fields {
//			if !field.Inject() {
//				continue
//			}
//
//			fd := self.registerTypes[field.fullTypeName]
//			if fd == nil {
//				return fmt.Errorf("%s's field %s is not be registed",
//					d.refType.Name(), field.fullTypeName)
//			}
//
//			refField := refVal.FieldByName(field.name)
//			if refField.CanSet() {
//				if !fd.inited {
//					self.createObject(fd, field.fullTypeName)
//				}
//				refField.Set(reflect.ValueOf(self.registry[field.fullTypeName]))
//			}
//
//		}
//
//		d.inited = true
//		self.registry[name] = val
//	}
//
//	return nil
//}

func (self *Registry) initRegisterType() {
	for k, v := range self.registerTypes {
		refVal := reflect.New(v)
		self.registerRefVals[k] = refVal
		self.registry[k] = refVal.Interface()
	}
}

func (self *Registry) inject() error {
	for k, refVal := range self.registerRefVals {
		spec := self.typeSpecMap[k]
		if spec.numField > 0 {
			for _, f := range spec.fields {
				if !f.Inject() {
					continue
				}
				field := refVal.Elem().FieldByName(f.name)
				if !field.CanSet() {
					fmt.Errorf("private field cannot inject")
				} else {
					field.Set(reflect.ValueOf(self.registry[f.fullTypeName]))
				}
			}
		}
	}
	return nil
}

func (self *Registry) Close() {
	self.Lock()
	defer self.Unlock()

	for _, refVal := range self.registerRefVals {
		if refVal.Type().Implements(reflect.TypeOf((*Closable)(nil)).Elem()) {
			refVal.MethodByName("Close").Call(nil)
		}
	}

	self.closed = true
}

type TypeSpec struct {
	refType      reflect.Type
	fullTypeName string
	numField     int
	fields       []*FieldSpec
	numMethod    int
	methods      []string
	inited       bool
}

func NewTypeSpec(refType reflect.Type) *TypeSpec {
	num := refType.Elem().NumField()
	fields := make([]*FieldSpec, num)
	for i := 0; i < num; i++ {
		fields[i] = NewFieldSpec(refType.Elem().Field(i))
	}

	numMethod := refType.NumMethod()
	methods := make([]string, numMethod)
	for i := 0; i < num; i++ {
		methods[i] = refType.Method(i).Name
	}

	return &TypeSpec{
		refType:      refType.Elem(),
		fullTypeName: fullTypeName(refType),
		numField:     num,
		fields:       fields,
		numMethod:    numMethod,
		methods:      methods,
		inited:       false,
	}
}

func (self *TypeSpec) Name() string {
	return self.refType.Name()
}

type FieldSpec struct {
	refType      reflect.Type
	fullTypeName string
	name         string
	tag          reflect.StructTag
}

func NewFieldSpec(sField reflect.StructField) *FieldSpec {
	refType := sField.Type
	return &FieldSpec{
		refType:      refType,
		fullTypeName: fullTypeName(refType),
		name:         sField.Name,
		tag:          sField.Tag,
	}
}

func (self *FieldSpec) Inject() bool {
	return self.tag.Get("inject") == "true"
}
