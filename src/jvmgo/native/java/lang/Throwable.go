package lang

import (
	. "jvmgo/any"
	"jvmgo/jvm/rtda"
	rtc "jvmgo/jvm/rtda/class"
)

func init() {
	_throwable(fillInStackTrace, "fillInStackTrace", "(I)Ljava/lang/Throwable;")
	_throwable(getStackTraceElement, "getStackTraceElement", "(I)Ljava/lang/StackTraceElement;")
	_throwable(getStackTraceDepth, "getStackTraceDepth", "()I")
}

func _throwable(method Any, name, desc string) {
	rtc.RegisterNativeMethod("java/lang/Throwable", name, desc, method)
}

type StackTraceElement struct {
	declaringClass string
	methodName     string
	fileName       string
	lineNumber     int
}

// private native Throwable fillInStackTrace(int dummy);
// (I)Ljava/lang/Throwable;
func fillInStackTrace(frame *rtda.Frame) {
	vars := frame.LocalVars()
	this := vars.GetRef(0)

	stack := frame.OperandStack()
	stack.PushRef(this)

	stes := createStackTraceElements(this, frame)
	this.SetExtra(stes)
}

func createStackTraceElements(tObj *rtc.Obj, frame *rtda.Frame) []*StackTraceElement {
	thread := frame.Thread()
	depth := thread.StackDepth()

	// skip unrelated frames
	i := uint(1)
	for k := tObj.Class(); k != nil; k = k.SuperClass() {
		i++
	}
	if thread.TopFrameN(i).Method().Name() == "<athrow>" {
		i++
	}

	stes := make([]*StackTraceElement, 0, depth)
	for ; i < depth; i++ {
		frameN := thread.TopFrameN(i)
		methodN := frameN.Method()
		classN := methodN.Class()
		lineNumber := methodN.GetLineNumber(frameN.NextPC() - 1)

		ste := &StackTraceElement{
			declaringClass: classN.NameJlsFormat(),
			methodName:     methodN.Name(),
			fileName:       classN.Attributes().SourceFile(),
			lineNumber:     lineNumber,
		}
		stes = append(stes, ste)
	}

	return stes
}

// native int getStackTraceDepth();
// ()I
func getStackTraceDepth(frame *rtda.Frame) {
	vars := frame.LocalVars()
	this := vars.GetRef(0)

	stes := this.Extra().([]*StackTraceElement)
	depth := int32(len(stes))

	stack := frame.OperandStack()
	stack.PushInt(depth)
}

// native StackTraceElement getStackTraceElement(int index);
// (I)Ljava/lang/StackTraceElement;
func getStackTraceElement(frame *rtda.Frame) {
	vars := frame.LocalVars()
	this := vars.GetRef(0)
	index := vars.GetInt(1)

	stes := this.Extra().([]*StackTraceElement)
	ste := stes[index]

	steObj := createStackTraceElementObj(ste, frame)
	stack := frame.OperandStack()
	stack.PushRef(steObj)
}

func createStackTraceElementObj(ste *StackTraceElement, frame *rtda.Frame) *rtc.Obj {
	declaringClass := rtda.NewJString(ste.declaringClass, frame)
	methodName := rtda.NewJString(ste.methodName, frame)
	fileName := rtda.NewJString(ste.fileName, frame)
	lineNumber := int32(ste.lineNumber)

	/*
	   public StackTraceElement(String declaringClass, String methodName,
	           String fileName, int lineNumber)
	*/
	steClass := frame.ClassLoader().LoadClass("java/lang/StackTraceElement")
	steObj := steClass.NewObj()
	// todo: call <init>
	steObj.SetFieldValue("declaringClass", "Ljava/lang/String;", declaringClass)
	steObj.SetFieldValue("methodName", "Ljava/lang/String;", methodName)
	steObj.SetFieldValue("fileName", "Ljava/lang/String;", fileName)
	steObj.SetFieldValue("lineNumber", "I", lineNumber)

	return steObj
}
