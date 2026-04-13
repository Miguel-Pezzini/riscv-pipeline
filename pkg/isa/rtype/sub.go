package rtype

import isa "riscv-instruction-encoder/pkg/isa"

type SUB struct {
	Type
}

func newSUB(t Type) *SUB {
	inst := &SUB{Type: t}

	inst.InstructionMeta = isa.InstructionMeta{
		Name:           "SUB",
		OpCode:         uint32(t.Opcode),
		IsLoad:         false,
		IsStore:        false,
		IsBranch:       false,
		IsJump:         false,
		WritesRegister: true,
		ReadsRegister:  true,
		Rs:             []int{(int(t.Rs1)), (int(t.Rs2))},
		Rd:             isa.IntPtr(int(t.Rd)),
		ProduceStage:   isa.EX,
		ConsumeStage:   isa.ID,
	}

	return inst
}

func (s *SUB) Execute(state isa.CPUState) error {
	rs1 := state.ReadReg(int(s.Rs1))
	rs2 := state.ReadReg(int(s.Rs2))
	state.WriteReg(int(s.Rd), rs1-rs2)
	return nil
}
