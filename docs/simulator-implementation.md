# Plano de Implementação — Simulação Real

Passos ordenados por dependência. Cada passo é buildável e testável independentemente.

---

## Passo 1 — `pkg/memory` (sem dependências novas)

Criar `pkg/memory/memory.go`:

```go
package memory

type Memory struct {
    data map[uint32]byte
}

func New() *Memory
func (m *Memory) LoadByte(addr uint32) (int8, error)
func (m *Memory) LoadHalf(addr uint32) (int16, error)
func (m *Memory) LoadWord(addr uint32) (int32, error)
func (m *Memory) LoadByteU(addr uint32) (uint8, error)
func (m *Memory) LoadHalfU(addr uint32) (uint16, error)
func (m *Memory) StoreByte(addr uint32, v int8)
func (m *Memory) StoreHalf(addr uint32, v int16)
func (m *Memory) StoreWord(addr uint32, v int32)
func (m *Memory) Dump(from, to uint32) string  // para debug
```

Regras de implementação:
- `map[uint32]byte` — sparse, sem limite de endereço
- Leitura de endereço não escrito retorna 0 (sem erro)
- `LoadHalf`/`LoadWord` verificam alinhamento (addr % 2 == 0 e addr % 4 == 0)
- Little-endian (RISC-V padrão)

**Teste mínimo:** store word em 0x10000000, load word no mesmo endereço → mesmo valor.

---

## Passo 2 — `pkg/cpu` (depende de `pkg/memory`)

Criar `pkg/cpu/state.go`:

```go
package cpu

import "riscv-instruction-encoder/pkg/memory"

const (
    TextBase  = uint32(0x00400000)
    DataBase  = uint32(0x10000000)
    StackTop  = uint32(0x7FFFFFFC)
)

type State struct {
    Regs [32]int32
    PC   uint32
    Mem  *memory.Memory
}

func NewState() *State   // PC = TextBase, sp = StackTop, Mem = memory.New()
func (s *State) ReadReg(n int) int32   // n=0 sempre retorna 0
func (s *State) WriteReg(n int, v int32)  // n=0 é no-op
func (s *State) DumpRegs() string     // para debug/output
```

**Teste mínimo:** WriteReg(0, 42) → ReadReg(0) == 0. WriteReg(5, 7) → ReadReg(5) == 7.

---

## Passo 3 — Adicionar `Execute` à interface `Instruction`

Em `pkg/isa/instruction.go`, adicionar à interface e ao `BaseInstruction`:

```go
// na interface
Execute(state *cpu.State) error

// no BaseInstruction
func (b *BaseInstruction) Execute(state *cpu.State) error {
    return fmt.Errorf("execute not implemented: %s", b.InstructionMeta.Name)
}
```

Isso causa import cycle se `cpu` importar `isa` e `isa` importar `cpu`. Solução: criar interface intermediária em `pkg/isa`:

```go
// pkg/isa/cpu_state.go
package isa

// CPUState é a interface que o executor fornece às instruções.
// Evita import cycle entre isa e cpu.
type CPUState interface {
    ReadReg(n int) int32
    WriteReg(n int, v int32)
    GetPC() uint32
    SetPC(pc uint32)
    LoadWord(addr uint32) (int32, error)
    LoadHalf(addr uint32) (int16, error)
    LoadByte(addr uint32) (int8, error)
    LoadHalfU(addr uint32) (uint16, error)
    LoadByteU(addr uint32) (uint8, error)
    StoreWord(addr uint32, v int32)
    StoreHalf(addr uint32, v int16)
    StoreByte(addr uint32, v int8)
}

// Na interface Instruction:
Execute(state CPUState) error
```

`cpu.State` implementa `isa.CPUState` automaticamente (duck typing do Go).

**Resultado:** build não quebra — `BaseInstruction.Execute` retorna erro "not implemented", o pipeline simulator continua funcionando igual.

---

## Passo 4 — `pkg/loader`

Criar `pkg/loader/loader.go`:

```go
package loader

import (
    "riscv-instruction-encoder/pkg/cpu"
    "riscv-instruction-encoder/pkg/decoder"
)

// Lê arquivo de instruções (bin ou hex), carrega no segmento .text
func LoadFile(path, format string) (*cpu.State, error)
```

Reutiliza `decoder.DecodeFromFile` para ler as palavras brutas, depois:
```go
for i, raw := range rawInstructions {
    addr := cpu.TextBase + uint32(i)*4
    state.Mem.StoreWord(addr, int32(raw.Value))
}
state.PC = cpu.TextBase
state.WriteReg(2, int32(cpu.StackTop)) // sp
```

---

## Passo 5 — `pkg/executor`

Criar `pkg/executor/executor.go`:

```go
package executor

type Config struct {
    MaxSteps int
    Trace    bool
}

type StepResult struct {
    PC    uint32
    Instr isa.Instruction
    // campos opcionais para trace
}

type Executor struct {
    State  isa.CPUState
    Config Config
}

func New(state isa.CPUState, cfg Config) *Executor
func (e *Executor) Run() error
func (e *Executor) Step() (*StepResult, error)  // um passo — útil para debug/trace
```

O loop em `Run`:
1. `raw = state.LoadWord(state.GetPC())`
2. `instr = decoder.DecodeInstruction(uint32(raw))`
3. `instr.Execute(state)`
4. `state.SetPC(state.GetPC() + 4)` — salvo se ECALL/EBREAK ou se a instrução já modificou PC (branch/jump)

Convenção de PC para branches/jumps: a instrução **escreve o PC destino diretamente** via `state.SetPC(target)`, e o loop **não** soma 4 se `GetPC()` mudou. Mais simples: o loop sempre soma 4, e a instrução subtrai 4 do offset para compensar.

> Escolha: instrução escreve `PC = target - 4`, loop soma 4. Consistente e simples.

---

## Passo 6 — Implementar `Execute` nas instruções

Ordem sugerida (do mais simples ao mais complexo):

### R-Type (sem memória, sem PC)
- `ADD`: `rd = rs1 + rs2`
- `SUB`: `rd = rs1 - rs2`
- `AND`, `OR`, `XOR`: operações bit a bit
- `SLL`, `SRL`, `SRA`: shifts (usar apenas os 5 bits menos significativos de rs2)
- `SLT`: `rd = (rs1 < rs2) ? 1 : 0` (signed)
- `SLTU`: igual, unsigned

### I-Type aritmético
- `ADDI`: `rd = rs1 + signExtend(imm12)`
- `ANDI`, `ORI`, `XORI`: bit a bit com imediato
- `SLLI`, `SRLI`, `SRAI`: shifts com imediato (shamt = imm[4:0])
- `SLTI`, `SLTIU`

### I-Type loads
- `LW`: `rd = mem[rs1 + signExtend(imm)]` (word)
- `LH`: `rd = signExtend(mem16[rs1 + imm])`
- `LHU`: `rd = zeroExtend(mem16[rs1 + imm])`
- `LB`: `rd = signExtend(mem8[rs1 + imm])`
- `LBU` (ainda não implementado no decoder — adicionar depois)

### S-Type stores
- `SW`: `mem[rs1 + signExtend(imm)] = rs2` (word)
- `SH`: store half
- `SB`: store byte

### B-Type branches
- `BEQ`: `if rs1 == rs2: PC = PC + signExtend(imm13) - 4`
- `BNE`, `BLT`, `BGE`, `BLTU`, `BGEU`

### J-Type / I-Type jumps
- `JAL`: `rd = PC + 4; PC = PC + signExtend(imm21) - 4`
- `JALR`: `rd = PC + 4; PC = (rs1 + signExtend(imm12)) & ~1 - 4`

### U-Type
- `LUI`: `rd = imm20 << 12`
- `AUIPC`: `rd = PC + (imm20 << 12)`

---

## Passo 7 — `cmd/simulator/main.go`

```
cmd/
  resolver/main.go    ← existente (pipeline)
  simulator/main.go   ← novo (execução real)
```

Interface mínima:
```
go run cmd/simulator/main.go --format hex --trace --dump-regs testdata/hex.txt
```

---

## Passo 8 — Trace e output

`pkg/executor/trace.go`:

```
[0x00400000]  ADDI x2, x2, -32    │ x2  = 0x7FFFFFDC
[0x00400004]  SW   x1, 28(x2)     │ mem[0x7FFFFFDC+28] = 0x00400044
[0x00400008]  JAL  x1, 40         │ x1  = 0x0040000C  PC -> 0x00400030
```

Formato: `[PC]  mnemônico  operandos  │  efeito`

---

## Passo 9 — Testes de integração

Usar o fibonacci do `testdata/` como caso de teste end-to-end:
- Carregar `testdata/hex.txt`
- Executar com trace desligado
- Verificar `a0` ao final == resultado esperado do fibonacci

Criar `testdata/expected_result.txt` com o valor esperado.

---

## Estrutura final de pacotes

```
pkg/
  memory/
    memory.go          ← NOVO
  cpu/
    state.go           ← NOVO
  loader/
    loader.go          ← NOVO
  executor/
    executor.go        ← NOVO
    trace.go           ← NOVO
  isa/
    instruction.go     ← modificado (+ Execute, + CPUState interface)
    cpu_state.go       ← NOVO (interface CPUState)
    ...
  isa/rtype/
    add.go             ← modificado (+ Execute)
    ...
  decoder/             ← sem mudanças
  hazard/              ← sem mudanças
  runner/              ← sem mudanças
cmd/
  resolver/main.go     ← sem mudanças
  simulator/main.go    ← NOVO
```

---

## Decisões de design abertas

| Questão | Opção A | Opção B | Status |
|---------|---------|---------|--------|
| PC após branch/jump | instrução escreve `target - 4`, loop sempre soma 4 | loop verifica se PC mudou | Decidir no Passo 5 |
| Sign extension de imediatos | funções utilitárias em `pkg/isa/` | inline em cada instrução | Utilitários (evita duplicação) |
| Tratamento de misaligned access | retorna error | panic | Retorna error (mais testável) |
| Condição de parada | ECALL a7=10 | instrução especial HALT | ECALL (mais realista) |
