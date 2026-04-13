# CLAUDE.md — RISC-V Pipeline Simulator

## O que é este projeto

Simulador educacional de pipeline RISC-V de 5 estágios. O programa lê instruções em binário ou hexadecimal, faz o decode para objetos tipados, simula a execução no pipeline e detecta/resolve hazards inserindo NOPs. Gera relatórios de overhead por cenário.

## Como rodar

```bash
go run cmd/resolver/main.go
# Escolhe formato: 1=binário, 2=hexadecimal
# Lê de testdata/bin.txt ou testdata/hex.txt
# Gera 6 arquivos em pkg/files/
```

## Estrutura de pacotes

```
cmd/resolver/main.go          # Entrada; decodifica e roda 6 cenários
pkg/decoder/decode.go         # Parsing de uint32 → Instruction pelo opcode
pkg/isa/instruction.go        # Interface Instruction, Stage enum, InstructionMeta, PipelineInstruction
pkg/isa/nop.go                # Instrução NOP (inserida automaticamente em hazards)
pkg/isa/rtype/                # ADD, SUB
pkg/isa/itype/                # ADDI, ANDI, ORI, LW, LB, JALR
pkg/isa/stype/                # SW
pkg/isa/btype/                # BEQ, BNE, BLT
pkg/isa/utype/                # U-Type genérico
pkg/isa/jtype/                # JAL
pkg/hazard/data_detector.go   # Detecção RAW e WAR (com/sem forwarding)
pkg/hazard/control_detector.go# Detecção de branch/jump não resolvido
pkg/runner/run.go             # Loop de simulação, inserção de NOP, Step()
pkg/runner/output.go          # Console stats e escrita de arquivos
testdata/                     # Entradas de teste (fibonacci recursivo)
pkg/files/                    # Saídas geradas (não versionar)
```

## Conceitos centrais

### InstructionMeta (pkg/isa/instruction.go)
Cada instrução declara:
- `ReadRegs []int` / `WriteRegs []int` — registradores lidos/escritos
- `ProduceStage Stage` — em qual estágio o resultado fica disponível (geralmente EX=3, ou MEM=4 para loads)
- `ConsumeStage Stage` — em qual estágio os operandos são necessários (sempre ID=2)
- `IsBranch`, `IsJump`, `IsLoad`, `IsStore` — flags de categoria

### Pipeline de 5 estágios
```
IF(1) → ID(2) → EX(3) → MEM(4) → WB(5)
```
Cada `Step()` avança todas as instruções em execução +1 estágio e tenta buscar a próxima.

### Detecção de hazard
**Data hazard (RAW):** instrução posterior lê registrador antes da anterior terminar de escrever.
- Sem forwarding: espera até WB (estágio 5)
- Com forwarding: usa resultado a partir de ProduceStage (EX=3, ou MEM=4 para loads)

**Control hazard:** instrução branch/jump ainda não saiu do pipeline — next-PC desconhecido.

Quando há hazard: um NOP é inserido antes da instrução problemática.

### 6 cenários de simulação
| # | Data Hazard | Control Hazard | Forwarding |
|---|-------------|----------------|------------|
| 1 | ✓ | ✗ | ✗ |
| 2 | ✓ | ✗ | ✓ |
| 3 | ✗ | ✓ | ✗ |
| 4 | ✗ | ✓ | ✓ |
| 5 | ✓ | ✓ | ✗ |
| 6 | ✓ | ✓ | ✓ |

## Como adicionar uma nova instrução

1. Escolha o pacote do tipo correto (`rtype`, `itype`, etc.)
2. Crie um arquivo `<nome>.go` com `InstructionMeta` preenchido (especialmente `ProduceStage`, `ConsumeStage`, `ReadRegs`, `WriteRegs`)
3. Registre o funct3/funct7 no `type.go` do pacote (switch de roteamento)
4. O decoder já repassa pelo opcode — nenhuma mudança em `decode.go` é necessária se o opcode já existe

## Convenções

- Instruções implementam a interface `Instruction` de `pkg/isa/instruction.go`
- `BaseInstruction` fornece implementações padrão (no-op) para todos os métodos — embedde e sobrescreva apenas o necessário
- PC é calculado como `index * 4` no início da simulação (`runner/run.go:InstructionsToPipeline`)
- Arquivos de saída ficam em `pkg/files/` — não commitar conteúdo gerado
- Módulo Go: `riscv-instruction-encoder` (go.mod)
