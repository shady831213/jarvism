package buildIncmd

import (
	"github.com/shady831213/jarvisSim/parser"
)

type RepeatOption struct {
	parser.AstOption
}

func newRepeatOption() *parser.AstOption {
	inst := new(RepeatOption)
	inst.AstOption.Init("repeat")
	inst.WithValue = parser.NewAstOptionAction()
	if err := parser.AstParse(inst.WithValue, map[interface{}]interface{}{"compile_option": "-q +lint=none"}); err != nil {
		panic(err)
	}
	return &inst.AstOption
}

func init() {
	//parser.RegisterOption(newRepeatOption())
}
