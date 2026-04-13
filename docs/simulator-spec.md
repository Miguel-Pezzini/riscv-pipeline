# Simulação Real de Execução RISC-V

## Visão Geral

O simulador atual apenas **modela o pipeline** — conta ciclos, detecta hazards, insere NOPs. As instruções não executam de verdade: registradores não têm valores, memória não existe.

O objetivo desta feature é adicionar uma segunda camada: **execução sequencial real**, onde cada instrução lê e escreve estado concreto (registradores, memória). Isso possibilita:

- Validar se um programa produz o resultado correto
- Inspecionar o conteúdo da memória e dos registradores a cada passo
- Rastrear variáveis alocadas na stack/heap
- Futuramente: combinar execução real com o pipeline (pipeline-accurate simulation)

O pipeline simulator **permanece inalterado**. A execução real é uma nova camada ortogonal.

---

## Modelo de Execução

A execução é **sequencial e determinística**:

```
PC → fetch(mem[PC]) → decode → execute → atualiza estado → PC += 4 (ou branch target)
```

Não há conceito de ciclo aqui — cada instrução completa atomicamente. O foco é na **corretude funcional**, não no tempo.

---

## Estado da CPU (`pkg/cpu/state.go`)

```go
type State struct {
    Regs [32]int32      // registradores x0–x31 (x0 sempre zero)
    PC   uint32         // Program Counter
    Mem  *memory.Memory // memória virtual
}
```

### Registradores

32 registradores inteiros de 32 bits. Convenções ABI (informativas):

| Reg    | ABI Name | Uso convencional              |
|--------|----------|-------------------------------|
| x0     | zero     | Hardwired 0 (WriteReg ignora) |
| x1     | ra       | Return address                |
| x2     | sp       | Stack pointer                 |
| x3     | gp       | Global pointer                |
| x4     | tp       | Thread pointer                |
| x5–x7  | t0–t2    | Temporários caller-saved      |
| x8     | s0/fp    | Saved / Frame pointer         |
| x9     | s1       | Saved callee-saved            |
| x10–x11| a0–a1    | Args / return values          |
| x12–x17| a2–a7    | Args adicionais               |
| x18–x27| s2–s11   | Saved callee-saved            |
| x28–x31| t3–t6    | Temporários caller-saved      |

`WriteReg(0, v)` é silenciosamente ignorado — x0 é sempre 0.

---

## Memória Virtual (`pkg/memory/memory.go`)

### Layout de endereços (RV32I)

```
0x00000000 ┌──────────────────┐
           │   (reservado)    │
0x00400000 ├──────────────────┤
           │  Segmento .text  │  ← instruções do programa
           │  (cresce ↓)      │
0x10000000 ├──────────────────┤
           │  Segmento .data  │  ← variáveis globais / estáticas
           │  (tamanho fixo)  │
0x10010000 ├──────────────────┤
           │  Heap            │  ← cresce ↑ (malloc futuro)
           │                  │
           │        ...       │
0x7FFFFFFC ├──────────────────┤
           │  Stack           │  ← cresce ↓, sp inicia aqui
0x00000000 └──────────────────┘
```

**sp inicial:** `0x7FFFFFFC`

### Operações de memória

```go
type Memory struct {
    data map[uint32]byte  // sparse — só aloca o que é usado
}

func (m *Memory) LoadByte(addr uint32)  (int8,  error)
func (m *Memory) LoadHalf(addr uint32)  (int16, error)
func (m *Memory) LoadWord(addr uint32)  (int32, error)
func (m *Memory) LoadByteU(addr uint32) (uint8,  error)  // LBU
func (m *Memory) LoadHalfU(addr uint32) (uint16, error)  // LHU

func (m *Memory) StoreByte(addr uint32, v int8)
func (m *Memory) StoreHalf(addr uint32, v int16)
func (m *Memory) StoreWord(addr uint32, v int32)
```

Uso de `map[uint32]byte` (sparse map) em vez de slice contíguo: evita alocar GBs para um endereço de stack alto. Acesso a endereço não escrito retorna 0 (comportamento padrão RISC-V para memória inicializada a zero).

**Alinhamento:** loads/stores de half/word verificam alinhamento. Acesso desalinhado retorna erro.

### Segmento .text

O programa é carregado em `.text` a partir de `0x00400000`. Cada instrução de 32 bits ocupa 4 bytes consecutivos. O PC inicial é `0x00400000`.

---

## Interface de Execução

Adicionamos o método `Execute` à interface `Instruction` existente:

```go
// pkg/isa/instruction.go
type Instruction interface {
    // ... métodos existentes ...
    Execute(state *cpu.State) error
}
```

`BaseInstruction` fornece uma implementação padrão que retorna erro "not implemented", permitindo adicionar instruções gradualmente sem quebrar a build.

```go
func (b *BaseInstruction) Execute(state *cpu.State) error {
    return fmt.Errorf("execute not implemented for %s", b.InstructionMeta.Name)
}
```

### Exemplo — ADD

```go
// pkg/isa/rtype/add.go
func (a *ADD) Execute(state *cpu.State) error {
    rs1 := state.ReadReg(int(a.Rs1))
    rs2 := state.ReadReg(int(a.Rs2))
    state.WriteReg(int(a.Rd), rs1 + rs2)
    return nil
}
```

### Exemplo — LW

```go
func (l *LW) Execute(state *cpu.State) error {
    base := state.ReadReg(int(l.Rs1))
    addr := uint32(int32(base) + signExtend12(l.Imm))
    val, err := state.Mem.LoadWord(addr)
    if err != nil { return err }
    state.WriteReg(int(l.Rd), val)
    return nil
}
```

### Exemplo — BEQ

```go
func (b *BEQ) Execute(state *cpu.State) error {
    rs1 := state.ReadReg(int(b.Rs1))
    rs2 := state.ReadReg(int(b.Rs2))
    if rs1 == rs2 {
        state.PC = uint32(int32(state.PC) + signExtend13(b.Imm) - 4)
        // -4 porque o loop principal faz PC += 4 após Execute
    }
    return nil
}
```

---

## Loop de Execução (`pkg/executor/executor.go`)

```go
type Executor struct {
    State     *cpu.State
    MaxSteps  int       // proteção contra loop infinito
    Trace     bool      // se true, imprime estado a cada passo
}

func (e *Executor) Run() error {
    for step := 0; step < e.MaxSteps; step++ {
        raw := e.State.Mem.LoadWord(e.State.PC)
        instr := decoder.DecodeInstruction(raw)
        if instr == nil {
            return fmt.Errorf("instrução inválida em PC=0x%08X", e.State.PC)
        }
        if e.Trace {
            e.printStep(instr)
        }
        if err := instr.Execute(e.State); err != nil {
            return err
        }
        e.State.PC += 4
        // ECALL / EBREAK: parada especial (ver abaixo)
    }
    return fmt.Errorf("max steps (%d) atingido", e.MaxSteps)
}
```

**Condição de parada:** instrução `ECALL` com `a7=10` (exit). Ou `EBREAK`. Definir convenção no spec.

---

## Carregamento do Programa (`pkg/loader/loader.go`)

```go
type Program struct {
    Instructions []uint32  // palavras do segmento .text
    DataSegment  []byte    // conteúdo inicial de .data (opcional)
}

func LoadFromFile(path, format string) (*Program, error)

func (p *Program) IntoState() *cpu.State {
    state := cpu.NewState()
    addr := uint32(0x00400000)
    for _, word := range p.Instructions {
        state.Mem.StoreWord(addr, int32(word))
        addr += 4
    }
    // carrega .data se presente
    return state
}
```

---

## Entrypoint (`cmd/simulator/main.go`)

Novo binário separado do pipeline resolver:

```
go run cmd/simulator/main.go [arquivo]
```

Flags previstas:
- `--format bin|hex` — formato do arquivo de entrada
- `--trace` — imprime estado após cada instrução
- `--max-steps N` — limite de passos (default: 1.000.000)
- `--dump-regs` — imprime registradores ao final
- `--dump-mem 0x10000000:64` — dump de região de memória ao final

---

## Inspeção de Estado

### Dump de registradores

```
Registers after execution:
  zero (x0)  = 0x00000000
  ra   (x1)  = 0x00400028
  sp   (x2)  = 0x7FFFFFF0
  ...
  a0   (x10) = 0x0000000A   ← return value = 10
```

### Dump de memória

```
Memory [0x10000000 .. 0x10000040]:
  0x10000000:  DE AD BE EF  01 02 03 04  ...
```

### Trace por instrução

```
[PC=0x00400000] ADD  x5, x6, x7   | x5 <- 0x0000000A
[PC=0x00400004] ADDI x5, x5, -1   | x5 <- 0x00000009
[PC=0x00400008] BNE  x5, x0, -8   | branch taken -> 0x00400000
```

---

## Tratamento de ECALL (futuro)

`ECALL` com `a7` (x17) definindo o serviço:

| a7  | Serviço         | Args         | Retorno |
|-----|-----------------|--------------|---------|
| 1   | print_int       | a0 = inteiro | —       |
| 4   | print_string    | a0 = addr    | —       |
| 10  | exit            | —            | para    |
| 11  | print_char      | a0 = char    | —       |
| 17  | exit2           | a0 = código  | para    |

---

## O que NÃO está no escopo (por enquanto)

- Pipelining na execução (ainda sequencial)
- MMU / paginação / TLB
- Exceções e trap handlers
- CSR registers (exceto básico para ECALL)
- Floating-point (F/D extensions)
- Multiplicação/divisão (M extension)
- Multithreading
