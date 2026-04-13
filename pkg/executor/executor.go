package executor

import (
	"fmt"
	"riscv-instruction-encoder/pkg/decoder"
	"riscv-instruction-encoder/pkg/isa"
)

type Config struct {
	MaxSteps int
	Trace    bool
}

type StepResult struct {
	PC    uint32
	Instr isa.Instruction
}

type Executor struct {
	State   isa.CPUState
	Config  Config
	Steps   int
	History []StepResult
	Halted  bool
}

func New(state isa.CPUState, cfg Config) *Executor {
	if cfg.MaxSteps == 0 {
		cfg.MaxSteps = 100000
	}
	return &Executor{
		State:  state,
		Config: cfg,
	}
}

func (e *Executor) Step() (*StepResult, error) {
	if e.Halted {
		return nil, fmt.Errorf("CPU halted")
	}

	pc := e.State.GetPC()

	raw, err := e.State.LoadWord(pc)
	if err != nil {
		return nil, fmt.Errorf("fetch error at PC=0x%08X: %w", pc, err)
	}

	if raw == 0 {
		e.Halted = true
		return nil, fmt.Errorf("hit zero instruction at PC=0x%08X (halt)", pc)
	}

	// ECALL: opcode 0x73, all other bits zero => inst == 0x00000073
	if uint32(raw) == 0x00000073 {
		switch e.State.ReadReg(17) {
		case 10, 17:
			e.Halted = true
			return nil, nil
		default:
			return nil, fmt.Errorf("unsupported ECALL service %d at PC=0x%08X", e.State.ReadReg(17), pc)
		}
	}

	// EBREAK
	if uint32(raw) == 0x00100073 {
		e.Halted = true
		return nil, nil
	}

	instr := decoder.DecodeInstruction(uint32(raw))
	if instr == nil {
		return nil, fmt.Errorf("unknown instruction 0x%08X at PC=0x%08X", uint32(raw), pc)
	}

	result := &StepResult{PC: pc, Instr: instr}

	if err := instr.Execute(e.State); err != nil {
		return result, fmt.Errorf("execute error at PC=0x%08X (%s): %w", pc, instr.String(), err)
	}

	e.State.SetPC(e.State.GetPC() + 4)
	e.Steps++
	e.History = append(e.History, *result)

	return result, nil
}

func (e *Executor) Run() error {
	for e.Steps < e.Config.MaxSteps {
		result, err := e.Step()
		if e.Halted {
			return nil
		}
		if err != nil {
			return err
		}
		if e.Config.Trace && result != nil {
			fmt.Print(FormatTrace(result))
		}
	}
	return fmt.Errorf("max steps (%d) exceeded", e.Config.MaxSteps)
}
