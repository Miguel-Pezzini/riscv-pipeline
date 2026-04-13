package itype

import isa "riscv-instruction-encoder/pkg/isa"

type LH struct {
	Type
}

func newLH(t Type) *LH {
	inst := &LH{Type: t}
	inst.InstructionMeta = isa.InstructionMeta{
		Name:           "LH",
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

func (l *LH) Execute(state isa.CPUState) error {
	base := state.ReadReg(int(l.Rs1))
	offset := isa.SignExtend12(l.Imm)
	addr := uint32(base + offset)
	val, err := state.LoadHalf(addr)
	if err != nil {
		return err
	}
	state.WriteReg(int(l.Rd), int32(val))
	return nil
}
