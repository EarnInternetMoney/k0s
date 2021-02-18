package eval

import (
	"src.elv.sh/pkg/diag"
	"src.elv.sh/pkg/parse"
)

type errorpfer interface {
	errorpf(r diag.Ranger, fmt string, args ...interface{})
}

// argsWalker is used by builtin special forms to implement argument parsing.
type argsWalker struct {
	cp   errorpfer
	form *parse.Form
	idx  int
}

func (cp *compiler) walkArgs(f *parse.Form) *argsWalker {
	return &argsWalker{cp, f, 0}
}

func (aw *argsWalker) more() bool {
	return aw.idx < len(aw.form.Args)
}

func (aw *argsWalker) peek() *parse.Compound {
	if !aw.more() {
		aw.cp.errorpf(aw.form, "need more arguments")
	}
	return aw.form.Args[aw.idx]
}

func (aw *argsWalker) next() *parse.Compound {
	n := aw.peek()
	aw.idx++
	return n
}

// nextIs returns whether the next argument's source matches the given text. It
// also consumes the argument if it is.
func (aw *argsWalker) nextIs(text string) bool {
	if aw.more() && parse.SourceText(aw.form.Args[aw.idx]) == text {
		aw.idx++
		return true
	}
	return false
}

// nextMustLambda fetches the next argument, raising an error if it is not a
// lambda.
func (aw *argsWalker) nextMustLambda() *parse.Primary {
	n := aw.next()
	if len(n.Indexings) != 1 {
		aw.cp.errorpf(n, "must be lambda")
	}
	if len(n.Indexings[0].Indicies) != 0 {
		aw.cp.errorpf(n, "must be lambda")
	}
	pn := n.Indexings[0].Head
	if pn.Type != parse.Lambda {
		aw.cp.errorpf(n, "must be lambda")
	}
	return pn
}

func (aw *argsWalker) nextMustLambdaIfAfter(leader string) *parse.Primary {
	if aw.nextIs(leader) {
		return aw.nextMustLambda()
	}
	return nil
}

func (aw *argsWalker) mustEnd() {
	if aw.more() {
		aw.cp.errorpf(diag.Ranging{From: aw.form.Args[aw.idx].Range().From, To: aw.form.Range().To}, "too many arguments")
	}
}
