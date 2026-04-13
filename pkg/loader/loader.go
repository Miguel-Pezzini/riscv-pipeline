package loader

import (
	"riscv-instruction-encoder/pkg/cpu"
	"riscv-instruction-encoder/pkg/decoder"
)

// LoadFile lê arquivo de instruções (bin ou hex), carrega no segmento .text
func LoadFile(path, format string) (*cpu.State, error) {
	rawInstructions, err := decoder.DecodeFromFile(path, format)
	if err != nil {
		return nil, err
	}

	state := cpu.NewState()
	for i, raw := range rawInstructions {
		addr := cpu.TextBase + uint32(i)*4
		if err := state.Mem.StoreWord(addr, int32(raw.Value)); err != nil {
			return nil, err
		}
	}
	state.PC = cpu.TextBase
	state.WriteReg(2, int32(cpu.StackTop)) // sp

	return state, nil
}
