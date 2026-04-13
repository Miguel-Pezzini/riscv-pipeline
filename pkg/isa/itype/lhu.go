package itype

import isa "riscv-instruction-encoder/pkg/isa"

type LHU struct {
	Type
}

func newLHU(t Type) *LHU {
	inst := &LHU{Type: t}
	inst.InstructionMeta = isa.InstructionMeta{
		Name:           "LHU",
		OpCode:         uint32(t.OpCode),
		IsLoad:         true,
		IsStore:        false,
		IsBranch:       false,
		IsJump:         false,
		WritesRegister: true,
		ReadsRegister:  true,
		Rs:             []int{int(t.Rs1)},
		Rd:             isa.IntPtr(int(t.Rd)),
		ProduceStage:   isa.MEM,
		ConsumeStage:   isa.ID,
	}
	return inst
}
