package stype

import isa "riscv-instruction-encoder/pkg/isa"

type SB struct {
	Type
}

func newSB(t Type) *SB {
	inst := &SB{Type: t}
	inst.InstructionMeta = isa.InstructionMeta{
		Name:           "SB",
		OpCode:         uint32(t.OpCode),
		IsLoad:         false,
		IsStore:        true,
		IsBranch:       false,
		IsJump:         false,
		WritesRegister: false,
		ReadsRegister:  true,
		Rs:             []int{int(t.Rs1), int(t.Rs2)},
		Rd:             nil,
		ProduceStage:   isa.WB,
		ConsumeStage:   isa.ID,
	}
	return inst
}
