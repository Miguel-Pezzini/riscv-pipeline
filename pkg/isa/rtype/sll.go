package rtype

import isa "riscv-instruction-encoder/pkg/isa"

type SLL struct {
	Type
}

func newSLL(t Type) *SLL {
	inst := &SLL{Type: t}

	inst.InstructionMeta = isa.InstructionMeta{
		Name:           "SLL",
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
