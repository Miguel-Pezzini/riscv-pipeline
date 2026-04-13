package itype

import isa "riscv-instruction-encoder/pkg/isa"

type JALR struct {
	Type
}

func newJALR(t Type) *JALR {
	inst := &JALR{Type: t}
	inst.InstructionMeta = isa.InstructionMeta{
		Name:           "JALR",
		OpCode:         uint32(t.OpCode),
		IsLoad:         false,
		IsStore:        false,
		IsBranch:       false,
		IsJump:         true,
		WritesRegister: true, // grava o PC+4 em Rd
		ReadsRegister:  true,
		Rs:             []int{int(t.Rs1)},
		Rd:             isa.IntPtr(int(t.Rd)),
		ProduceStage:   isa.EX, // calcula destino no EX
		ConsumeStage:   isa.ID,
	}
	return inst
}

func (j *JALR) Execute(state isa.CPUState) error {
	pc := state.GetPC()
	rs1 := state.ReadReg(int(j.Rs1))
	imm := isa.SignExtend12(j.Imm)
	target := uint32((rs1 + imm)) & ^uint32(1)
	state.WriteReg(int(j.Rd), int32(pc+4))
	state.SetPC(target - 4) // loop will add 4
	return nil
}
