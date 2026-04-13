package executor

import "fmt"

// FormatTrace formats a step result for trace output.
// Format: [PC]  instruction_string
func FormatTrace(result *StepResult) string {
	return fmt.Sprintf("[0x%08X]  %s\n", result.PC, result.Instr.String())
}
