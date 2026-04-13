package stype

import isa "riscv-instruction-encoder/pkg/isa"

type SH struct {
	Type
}

func newSH(t Type) *SH {
	inst := &SH{Type: t}
	inst.InstructionMeta = isa.InstructionMeta{
		Name:           "SH",
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

func (s *SH) Execute(state isa.CPUState) error {
	base := state.ReadReg(int(s.Rs1))
	offset := isa.SignExtend12(s.Imm)
	addr := uint32(base + offset)
	val := state.ReadReg(int(s.Rs2))
	return state.StoreHalf(addr, int16(val))
}
