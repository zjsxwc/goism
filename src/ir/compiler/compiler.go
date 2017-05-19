package compiler

import (
	"bytes"
	"dt"
	"ir"
	"tu"
)

type Compiler struct {
	buf         bytes.Buffer
	cvec        *dt.ConstPool
	st          *dt.ExecutionStack
	lastLabelID uint16
}

func New() *Compiler {
	return &Compiler{
		cvec: &dt.ConstPool{},
	}
}

func (cl *Compiler) CompileFunc(f *tu.Func) *ir.Func {
	cl.reset(f.Params)

	compileStmtList(cl, f.Body.Forms)

	return &ir.Func{
		Name:       f.Name,
		ArgsDesc:   argsDescriptor(len(f.Params), f.Variadic),
		StackUsage: cl.st.MaxLen(),
		Body:       cl.buf.Bytes(),
		DocString:  docString(f),
		ConstVec:   cl.cvec,
	}
}

// Prepare compiler for re-use.
func (cl *Compiler) reset(bindings []string) {
	cl.buf.Truncate(0)
	cl.cvec.Clear()
	cl.st = dt.NewExecutionStack(bindings)
	cl.lastLabelID = 0
}