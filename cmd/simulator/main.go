package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"riscv-instruction-encoder/pkg/cpu"
	"riscv-instruction-encoder/pkg/executor"
	"riscv-instruction-encoder/pkg/loader"
	"riscv-instruction-encoder/pkg/runner"
	"strconv"
	"strings"
)

func main() {
	format := flag.String("format", "hex", "Input format: bin or hex")
	trace := flag.Bool("trace", false, "Enable execution trace")
	pipeline := flag.Bool("pipeline", false, "Run pipeline analysis over the executed trace")
	pipelineOutDir := flag.String("pipeline-out-dir", "pkg/files", "Directory for generated pipeline reports")
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
	maxStepsReached := false
	if err != nil {
		var maxErr *executor.MaxStepsError
		if errors.As(err, &maxErr) {
			maxStepsReached = true
			fmt.Fprintf(os.Stderr, "aviso: limite de passos (%d) atingido; usando traço parcial da execução\n", maxErr.MaxSteps)
		} else {
			log.Fatalf("erro na execução: %v", err)
		}
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

	if *pipeline {
		if len(exec.History) == 0 {
			fmt.Println("\nNenhuma instrução executada; relatórios de pipeline não foram gerados.")
			return
		}

		if err := os.MkdirAll(*pipelineOutDir, 0o755); err != nil {
			log.Fatalf("erro ao preparar diretório de saída do pipeline: %v", err)
		}

		runPipelineReports(exec.History, *pipelineOutDir)
		if maxStepsReached {
			fmt.Println("\nRelatórios de pipeline gerados a partir de traço parcial.")
		}
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

func runPipelineReports(history []executor.StepResult, outDir string) {
	scenarios := []struct {
		forwarding    bool
		dataHazard    bool
		controlHazard bool
		fileName      string
	}{
		{false, true, false, "output_trace_data_no_forwarding.txt"},
		{true, true, false, "output_trace_data_forwarding.txt"},
		{false, false, true, "output_trace_control_no_forwarding.txt"},
		{true, false, true, "output_trace_control_forwarding.txt"},
		{false, true, true, "output_trace_integrated_no_forwarding.txt"},
		{true, true, true, "output_trace_integrated_forwarding.txt"},
	}

	fmt.Printf("\n=== Pipeline Sobre Traço Executado ===\n")
	fmt.Printf("Diretório de saída: %s\n", outDir)

	for _, scenario := range scenarios {
		runner.RunFromExecutionHistory(
			history,
			scenario.forwarding,
			scenario.dataHazard,
			scenario.controlHazard,
			filepath.Join(outDir, scenario.fileName),
		)
	}
}
