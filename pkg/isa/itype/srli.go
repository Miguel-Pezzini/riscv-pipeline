package itype

import isa "riscv-instruction-encoder/pkg/isa"

type SRLI struct {
	Type
}

func newSRLI(t Type) *SRLI {
	inst := &SRLI{Type: t}

	inst.InstructionMeta = isa.InstructionMeta{
		Name:           "SRLI",
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
