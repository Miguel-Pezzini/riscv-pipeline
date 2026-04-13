package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"riscv-instruction-encoder/pkg/cpu"
	"riscv-instruction-encoder/pkg/executor"
	"riscv-instruction-encoder/pkg/loader"
	"strconv"
	"strings"
)

func main() {
	format := flag.String("format", "hex", "Input format: bin or hex")
	trace := flag.Bool("trace", false, "Enable execution trace")
	dumpRegs := flag.Bool("dump-regs", false, "Dump registers after execution")
	dumpMem := flag.String("dump-mem", "", "Dump memory range as <start>:<length>, e.g. 0x10000000:64")
	maxSteps := flag.Int("max-steps", 100000, "Maximum execution steps")
	flag.Parse()

	args := flag.Args()
	if len(args) < 1 {
		fmt.Fprintf(os.Stderr, "Usage: simulator [flags] <input-file>\n")
		flag.PrintDefaults()
		os.Exit(1)
	}

	filePath := args[0]

	state, err := loader.LoadFile(filePath, *format)
	if err != nil {
		log.Fatalf("erro ao carregar arquivo: %v", err)
	}

	exec := executor.New(state, executor.Config{
		MaxSteps: *maxSteps,
		Trace:    *trace,
	})

	fmt.Printf("=== RISC-V Simulator ===\n")
	fmt.Printf("Arquivo: %s (formato: %s)\n", filePath, *format)
	fmt.Printf("PC inicial: 0x%08X\n", cpu.TextBase)
	fmt.Println()

	err = exec.Run()
	if err != nil {
		log.Fatalf("erro na execução: %v", err)
	}

	fmt.Printf("\n=== Execução concluída ===\n")
	fmt.Printf("Passos executados: %d\n", exec.Steps)
	fmt.Printf("PC final: 0x%08X\n", state.GetPC())

	if *dumpRegs {
		fmt.Printf("\n=== Registradores ===\n")
		fmt.Print(state.DumpRegs())
	}

	if *dumpMem != "" {
		from, to, err := parseDumpRange(*dumpMem)
		if err != nil {
			log.Fatalf("erro em --dump-mem: %v", err)
		}
		fmt.Printf("\n=== Memória [0x%08X .. 0x%08X) ===\n", from, to)
		fmt.Print(state.Mem.Dump(from, to))
	}
}

func parseDumpRange(raw string) (uint32, uint32, error) {
	parts := strings.Split(raw, ":")
	if len(parts) != 2 {
		return 0, 0, fmt.Errorf("formato inválido %q; use <start>:<length>", raw)
	}

	from, err := strconv.ParseUint(parts[0], 0, 32)
	if err != nil {
		return 0, 0, fmt.Errorf("endereço inicial inválido: %w", err)
	}

	length, err := strconv.ParseUint(parts[1], 0, 32)
	if err != nil {
		return 0, 0, fmt.Errorf("comprimento inválido: %w", err)
	}

	return uint32(from), uint32(from + length), nil
}
