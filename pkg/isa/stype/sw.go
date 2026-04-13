package stype

import isa "riscv-instruction-encoder/pkg/isa"

type SW struct {
	Type
}

func newSW(t Type) *SW {
	inst := &SW{Type: t}
	inst.InstructionMeta = isa.InstructionMeta{
		Name:           "SW",
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

func (s *SW) Execute(state isa.CPUState) error {
	base := state.ReadReg(int(s.Rs1))
	offset := isa.SignExtend12(s.Imm)
	addr := uint32(base + offset)
	val := state.ReadReg(int(s.Rs2))
	return state.StoreWord(addr, val)
}
