package class

// only used by jvm.go
func NewBootstrapMethod(cl *ClassLoader) *Method {
	method := &Method{}
	method.class = &Class{name: "~jvmgo", classLoader: cl}
	method.name = "<bootstrap>"
	method.accessFlags = ACC_STATIC
	method.maxStack = 8
	method.maxLocals = 8
	method.code = []byte{0xff, 0xb1} // bootstrap, return
	return method
}

// todo
func hackClass(class *Class) {
	if class.name == "java/lang/ClassLoader" {
		loadLibrary := class.GetStaticMethod("loadLibrary", "(Ljava/lang/Class;Ljava/lang/String;Z)V")
		loadLibrary.code = []byte{0xb1} // return void
	}
}
