package rtype

import isa "riscv-instruction-encoder/pkg/isa"

type OR struct {
	Type
}

func newOR(t Type) *OR {
	inst := &OR{Type: t}

	inst.InstructionMeta = isa.InstructionMeta{
		Name:           "OR",
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

func (o *OR) Execute(state isa.CPUState) error {
	rs1 := state.ReadReg(int(o.Rs1))
	rs2 := state.ReadReg(int(o.Rs2))
	state.WriteReg(int(o.Rd), rs1|rs2)
	return nil
}
