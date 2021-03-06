package class

func (self *Class) IsArray() bool {
	return self.name[0] == '['
}

func (self *Class) IsPrimitiveArray() bool {
	return self.IsArray() && len(self.name) == 2
}

func (self *Class) ComponentClass() *Class {
	componentClassName := getComponentClassName(self.name)
	return self.classLoader.LoadClass(componentClassName)
}

func (self *Class) ArrayClass() *Class {
	arrayClassName := getArrayClassName(self.name)
	return self.classLoader.LoadClass(arrayClassName)
}
