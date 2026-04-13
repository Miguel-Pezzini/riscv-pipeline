package itype

import isa "riscv-instruction-encoder/pkg/isa"

type SRAI struct {
	Type
}

func newSRAI(t Type) *SRAI {
	inst := &SRAI{Type: t}

	inst.InstructionMeta = isa.InstructionMeta{
		Name:           "SRAI",
		OpCode:         uint32(t.OpCode),
		IsLoad:         false,
		IsStore:        false,
		IsBranch:       false,
		IsJump:         false,
		WritesRegister: true,
		ReadsRegister:  true,
		Rs:             []int{int(t.Rs1)},
		Rd:             isa.IntPtr(int(t.Rd)),
		ProduceStage:   isa.EX,
		ConsumeStage:   isa.ID,
	}

	return inst
}

func (s *SRAI) Execute(state isa.CPUState) error {
	rs1 := state.ReadReg(int(s.Rs1))
	shamt := s.Imm & 0x1F
	state.WriteReg(int(s.Rd), rs1>>shamt) // arithmetic right shift on int32
	return nil
}
