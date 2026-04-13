package rtype

import isa "riscv-instruction-encoder/pkg/isa"

type XOR struct {
	Type
}

func newXOR(t Type) *XOR {
	inst := &XOR{Type: t}

	inst.InstructionMeta = isa.InstructionMeta{
		Name:           "XOR",
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
