package itype

import isa "riscv-instruction-encoder/pkg/isa"

type LW struct {
	Type
}

// LW – I-type, load word
func newLW(t Type) *LW {
	inst := &LW{Type: t}
	inst.InstructionMeta = isa.InstructionMeta{
		Name:           "LW",
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

func (l *LW) Execute(state isa.CPUState) error {
	base := state.ReadReg(int(l.Rs1))
	offset := isa.SignExtend12(l.Imm)
	addr := uint32(base + offset)
	val, err := state.LoadWord(addr)
	if err != nil {
		return err
	}
	state.WriteReg(int(l.Rd), val)
	return nil
}
