package btype

import isa "riscv-instruction-encoder/pkg/isa"

type BLT struct {
	Type
}

func newBLT(t Type) *BLT {
	inst := &BLT{Type: t}

	inst.InstructionMeta = isa.InstructionMeta{
		Name:           "BLT",
		OpCode:         uint32(t.OpCode),
		IsLoad:         false,
		IsStore:        false,
		IsBranch:       true,
		IsJump:         false,
		WritesRegister: false,
		ReadsRegister:  true,
		Rs:             []int{int(t.Rs1), int(t.Rs2)},
		Rd:             nil,
		ProduceStage:   isa.EX,
		ConsumeStage:   isa.ID,
	}

	return inst
}

func (b *BLT) Execute(state isa.CPUState) error {
	rs1 := state.ReadReg(int(b.Rs1))
	rs2 := state.ReadReg(int(b.Rs2))
	if rs1 < rs2 {
		pc := state.GetPC()
		offset := isa.SignExtend13(b.Imm)
		state.SetPC(uint32(int32(pc)+offset) - 4)
	}
	return nil
}
