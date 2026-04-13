# RISC-V Pipeline Simulator

Simulador educacional de pipeline RISC-V de 5 estágios com detecção de hazards e forwarding. O programa decodifica instruções binárias/hexadecimais e simula a execução no pipeline, inserindo NOPs onde necessário e reportando o overhead gerado.

## Funcionalidades

- Decodificação de instruções RISC-V a partir de binário ou hexadecimal
- Simulação de pipeline clássico de 5 estágios: **IF → ID → EX → MEM → WB**
- Detecção de **data hazards** (RAW e WAR) com e sem forwarding
- Detecção de **control hazards** (branches e jumps não resolvidos)
- Inserção automática de NOPs para resolver hazards
- 6 cenários de simulação independentes
- Relatório de overhead (quantidade de NOPs inseridos vs instruções originais)

## Instruções suportadas

| Tipo | Instruções |
|------|-----------|
| R-Type | `ADD`, `SUB` |
| I-Type | `ADDI`, `ANDI`, `ORI`, `LW`, `LB`, `JALR` |
| S-Type | `SW` |
| B-Type | `BEQ`, `BNE`, `BLT` |
| J-Type | `JAL` |
| Utility | `NOP` (inserido automaticamente) |

## Pré-requisitos

- Go 1.21+

## Como usar

```bash
# Clone e entre no diretório
git clone <repo>
cd riscv-instruction-decoder

# Execute
go run cmd/resolver/main.go
```

O programa pergunta o formato de entrada:
```
1 - Binary
2 - Hexadecimal
```

Em seguida lê o arquivo correspondente em `testdata/` e gera 6 arquivos de saída em `pkg/files/`.

### Formato de entrada

**Binário** (`testdata/bin.txt`) — uma instrução de 32 bits por linha:
```
00000000000100000000000010010011
```

**Hexadecimal** (`testdata/hex.txt`) — uma instrução por linha:
```
00100093
```

### Formato de saída

Cada arquivo em `pkg/files/` contém as instruções com NOPs inseridos:
```
0x0     ADDI x1, x0, 1
0x4     NOP
0x8     ADD x3, x1, x2
```

## Cenários de simulação

| Arquivo gerado | Data Hazard | Control Hazard | Forwarding |
|----------------|:-----------:|:--------------:|:----------:|
| `output_data_no_forwarding.txt` | ✓ | ✗ | ✗ |
| `output_data_forwarding.txt` | ✓ | ✗ | ✓ |
| `output_control_no_forwarding.txt` | ✗ | ✓ | ✗ |
| `output_control_forwarding.txt` | ✗ | ✓ | ✓ |
| `output_integrated_no_forwarding.txt` | ✓ | ✓ | ✗ |
| `output_integrated_forwarding.txt` | ✓ | ✓ | ✓ |

## Arquitetura

```
cmd/
└── resolver/main.go        # Ponto de entrada

pkg/
├── decoder/
│   └── decode.go           # Parsing de uint32 → instrução tipada (roteamento por opcode)
├── isa/
│   ├── instruction.go      # Interface Instruction, Stage, InstructionMeta, PipelineInstruction
│   ├── nop.go              # Instrução NOP
│   ├── rtype/              # R-Type: ADD, SUB
│   ├── itype/              # I-Type: ADDI, ANDI, ORI, LW, LB, JALR
│   ├── stype/              # S-Type: SW
│   ├── btype/              # B-Type: BEQ, BNE, BLT
│   ├── utype/              # U-Type (genérico)
│   └── jtype/              # J-Type: JAL
├── hazard/
│   ├── data_detector.go    # Detecção RAW/WAR (com/sem forwarding)
│   └── control_detector.go # Detecção de branch/jump não resolvido
├── runner/
│   ├── run.go              # Simulação do pipeline (Step, Run, inserção de NOP)
│   └── output.go           # Relatório de resultados e escrita de arquivos
└── files/                  # Saídas geradas (ignoradas pelo git)

testdata/
├── bin.txt                 # Programa de teste em binário (fibonacci recursivo)
└── hex.txt                 # Mesmo programa em hexadecimal
```

### Fluxo de dados

```
Arquivo de entrada (bin.txt / hex.txt)
    ↓ decoder.DecodeFromFile()
[]uint32 (instruções brutas)
    ↓ decoder.DecodeInstructionFromUInt32()
[]Instruction (objetos tipados com metadata)
    ↓ runner.Run()
Pipeline simulation (por ciclo: Step → hazard check → NOP ou fetch)
    ↓ output.printResult() + output.writeFile()
Console stats + pkg/files/output_*.txt
```

### Detecção de hazards

**Data hazard (RAW — Read After Write):**
O hazard ocorre quando uma instrução precisa ler um registrador antes que a instrução anterior termine de escrever nele.

- **Sem forwarding:** espera até o estágio WB (5) da instrução anterior
- **Com forwarding:** usa o resultado assim que ele fica disponível no `ProduceStage` da instrução anterior (EX=3 para operações aritméticas, MEM=4 para loads)

**Control hazard:**
Ocorre quando um branch ou jump ainda está no pipeline e o PC da próxima instrução é desconhecido. O simulador insere NOPs até que a instrução de desvio alcance o estágio WB.

### Metadata por instrução

Cada instrução declara em sua `InstructionMeta`:
- `ReadRegs` / `WriteRegs` — registradores consumidos e produzidos
- `ProduceStage` — estágio em que o resultado fica disponível (EX para maioria, MEM para loads)
- `ConsumeStage` — estágio em que os operandos são necessários (sempre ID)
- Flags: `IsBranch`, `IsJump`, `IsLoad`, `IsStore`

Essa metadata permite que o detector de hazards seja genérico e independente do tipo de instrução.
