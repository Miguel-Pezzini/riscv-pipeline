package btype

import isa "riscv-instruction-encoder/pkg/isa"

type BGE struct {
	Type
}

func newBGE(t Type) *BGE {
	inst := &BGE{Type: t}

	inst.InstructionMeta = isa.InstructionMeta{
		Name:           "BGE",
		OpCode:         uint32(t.OpCode),
		IsLoad:         false,
		IsStore:        false,
		IsBranch:       true,
		IsJump:         false,
		WritesRegister: false,
		ReadsRegister:  true,
		Rs:             []int{int(t.Rs1), int(t.Rs2)},
		Rd:             nil,
		ProduceStage:   isa.EX,
		ConsumeStage:   isa.ID,
	}

	return inst
}
