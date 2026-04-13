package rtype

import isa "riscv-instruction-encoder/pkg/isa"

type AND struct {
	Type
}

func newAND(t Type) *AND {
	inst := &AND{Type: t}

	inst.InstructionMeta = isa.InstructionMeta{
		Name:           "AND",
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

func (a *AND) Execute(state isa.CPUState) error {
	rs1 := state.ReadReg(int(a.Rs1))
	rs2 := state.ReadReg(int(a.Rs2))
	state.WriteReg(int(a.Rd), rs1&rs2)
	return nil
}
