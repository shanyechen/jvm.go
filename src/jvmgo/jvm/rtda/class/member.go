package class

type ClassMember struct {
	AccessFlags
	name       string
	descriptor string
	class      *Class
}

func (self *ClassMember) Name() string {
	return self.name
}
func (self *ClassMember) Descriptor() string {
	return self.descriptor
}
func (self *ClassMember) Class() *Class {
	return self.class
}

func (self *ClassMember) ClassLoader() *ClassLoader {
	return self.class.classLoader
}
func (self *ClassMember) ConstantPool() *ConstantPool {
	return self.class.constantPool
}
