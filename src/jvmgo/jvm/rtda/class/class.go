package class

import (
	. "jvmgo/any"
	"jvmgo/util"
)

const (
	_initializing = 1
	_initialized  = 2
)

// name, superClassName and interfaceNames are all binary names(jvms8-4.2.1)
type Class struct {
	AccessFlags
	constantPool       *ConstantPool
	name               string // thisClassName
	superClassName     string
	interfaceNames     []string
	fields             []*Field
	methods            []*Method
	attributes         *Attributes
	staticFieldCount   uint
	instanceFieldCount uint
	staticFieldValues  []Any
	vtable             []*Method // virtual method table
	jClass             *Obj      // java.lang.Class instance
	superClass         *Class
	interfaces         []*Class
	classLoader        *ClassLoader // defining class loader
	state              int
}

func (self *Class) String() string {
	return "{Class name:" + self.name + "}"
}

// getters
func (self *Class) ConstantPool() *ConstantPool {
	return self.constantPool
}
func (self *Class) Name() string {
	return self.name
}
func (self *Class) Attributes() *Attributes {
	return self.attributes
}
func (self *Class) JClass() *Obj {
	return self.jClass
}
func (self *Class) SuperClass() *Class {
	return self.superClass
}
func (self *Class) Interfaces() []*Class {
	return self.interfaces
}
func (self *Class) ClassLoader() *ClassLoader {
	return self.classLoader
}

// todo
func (self *Class) NameJlsFormat() string {
	return util.ReplaceAll(self.name, "/", ".")
}
func NameCfFormat(jlsName string) string {
	return util.ReplaceAll(jlsName, ".", "/")
}

func (self *Class) InitializationNotStarted() bool {
	return self.state < _initializing // todo
}
func (self *Class) MarkInitializing() {
	self.state = _initializing
}
func (self *Class) MarkInitialized() {
	self.state = _initialized
}

func (self *Class) getField(name, descriptor string, isStatic bool) *Field {
	for k := self; k != nil; k = k.superClass {
		for _, field := range k.fields {
			if field.IsStatic() == isStatic &&
				field.name == name &&
				field.descriptor == descriptor {

				return field
			}
		}
	}
	// todo
	return nil
}
func (self *Class) getMethod(name, descriptor string, isStatic bool) *Method {
	for k := self; k != nil; k = k.superClass {
		for _, method := range k.methods {
			if method.IsStatic() == isStatic &&
				method.name == name &&
				method.descriptor == descriptor {

				return method
			}
		}
	}
	// todo
	return nil
}

// todo
func (self *Class) _getMethod(name, descriptor string, isStatic bool) *Method {
	for _, method := range self.methods {
		if method.IsStatic() == isStatic &&
			method.name == name &&
			method.descriptor == descriptor {

			return method
		}
	}
	return nil
}

func (self *Class) GetStaticField(name, descriptor string) *Field {
	return self.getField(name, descriptor, true)
}
func (self *Class) GetInstanceField(name, descriptor string) *Field {
	return self.getField(name, descriptor, false)
}

func (self *Class) GetStaticMethod(name, descriptor string) *Method {
	return self.getMethod(name, descriptor, true)
}
func (self *Class) GetInstanceMethod(name, descriptor string) *Method {
	return self.getMethod(name, descriptor, false)
}

func (self *Class) GetMainMethod() *Method {
	return self.GetStaticMethod(mainMethodName, mainMethodDesc)
}
func (self *Class) GetClinitMethod() *Method {
	return self._getMethod(clinitMethodName, clinitMethodDesc, true)
}

func (self *Class) getArrayClass() *Class {
	return self.classLoader.getRefArrayClass(self)
}

func (self *Class) NewObjWithExtra(extra Any) *Obj {
	obj := self.NewObj()
	obj.extra = extra
	return obj
}
func (self *Class) NewObj() *Obj {
	if self.instanceFieldCount > 0 {
		fields := make([]Any, self.instanceFieldCount)
		obj := newObj(self, fields, nil)
		obj.initFields()
		return obj
	} else {
		return newObj(self, nil, nil)
	}
}
func (self *Class) NewArray(count uint) *Obj {
	return NewRefArray(self, count)
}

func (self *Class) isJlObject() bool {
	return self == _jlObjectClass
}
func (self *Class) isJlCloneable() bool {
	return self == _jlCloneableClass
}
func (self *Class) isJioSerializable() bool {
	return self == _ioSerializableClass
}

// reflection
func (self *Class) GetStaticValue(fieldName, fieldDescriptor string) Any {
	field := self.GetStaticField(fieldName, fieldDescriptor)
	return field.GetStaticValue()
}
func (self *Class) SetStaticValue(fieldName, fieldDescriptor string, value Any) {
	field := self.GetStaticField(fieldName, fieldDescriptor)
	field.PutStaticValue(value)
}

func (self *Class) AsObj() *Obj {
	return &Obj{fields: self.staticFieldValues}
}
