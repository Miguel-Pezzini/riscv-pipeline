package rtype

import isa "riscv-instruction-encoder/pkg/isa"

type SLTU struct {
	Type
}

func newSLTU(t Type) *SLTU {
	inst := &SLTU{Type: t}

	inst.InstructionMeta = isa.InstructionMeta{
		Name:           "SLTU",
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

func (s *SLTU) Execute(state isa.CPUState) error {
	rs1 := uint32(state.ReadReg(int(s.Rs1)))
	rs2 := uint32(state.ReadReg(int(s.Rs2)))
	if rs1 < rs2 {
		state.WriteReg(int(s.Rd), 1)
	} else {
		state.WriteReg(int(s.Rd), 0)
	}
	return nil
}
