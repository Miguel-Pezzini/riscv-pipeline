package rtype

import isa "riscv-instruction-encoder/pkg/isa"

type SLT struct {
	Type
}

func newSLT(t Type) *SLT {
	inst := &SLT{Type: t}

	inst.InstructionMeta = isa.InstructionMeta{
		Name:           "SLT",
		OpCode:         uint32(t.Opcode),
		IsLoad:         false,
		IsStore:        false,
		IsBranch:       false,
		IsJump:         false,
		WritesRegister: true,
		ReadsRegister:  true,
		Rs:             []int{int(t.Rs1), int(t.Rs2)},
		Rd:             isa.IntPtr(int(t.Rd)),
		ProduceStage:   isa.EX,
		ConsumeStage:   isa.ID,
	}

	return inst
}

func (s *SLT) Execute(state isa.CPUState) error {
	rs1 := state.ReadReg(int(s.Rs1))
	rs2 := state.ReadReg(int(s.Rs2))
	if rs1 < rs2 {
		state.WriteReg(int(s.Rd), 1)
	} else {
		state.WriteReg(int(s.Rd), 0)
	}
	return nil
}
