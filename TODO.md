# TODO — RISC-V Pipeline Simulator

Itens organizados por prioridade. Bugs silenciosos (decodificação errada sem erro) primeiro.

---

## Simulação Real de Execução

Design completo em [docs/simulator-spec.md](docs/simulator-spec.md) e [docs/simulator-implementation.md](docs/simulator-implementation.md).

### Infraestrutura
- [ ] **Passo 1** — `pkg/memory/memory.go`: sparse map, load/store byte/half/word, little-endian, alinhamento
- [ ] **Passo 2** — `pkg/cpu/state.go`: 32 registradores, PC, referência à memória, x0 hardwired
- [ ] **Passo 3** — `pkg/isa/cpu_state.go`: interface `CPUState` (evita import cycle), adicionar `Execute(CPUState) error` à interface `Instruction` e ao `BaseInstruction`
- [ ] **Passo 4** — `pkg/loader/loader.go`: lê arquivo bin/hex, carrega instruções em `.text` (0x00400000), inicializa sp

### Executor
- [ ] **Passo 5** — `pkg/executor/executor.go`: loop fetch-decode-execute, `Step()` e `Run()`
- [ ] **Passo 7** — `cmd/simulator/main.go`: entrypoint com flags `--format`, `--trace`, `--dump-regs`
- [ ] **Passo 8** — `pkg/executor/trace.go`: output `[PC] mnemônico │ efeito`

### Implementar `Execute` nas instruções
- [ ] **Passo 6a** — R-Type: ADD, SUB, AND, OR, XOR, SLL, SRL, SRA, SLT, SLTU
- [ ] **Passo 6b** — I-Type aritmético: ADDI, ANDI, ORI, XORI, SLLI, SRLI, SRAI, SLTI, SLTIU
- [ ] **Passo 6c** — I-Type loads: LW, LH, LHU, LB (LBU ainda falta no decoder)
- [ ] **Passo 6d** — S-Type stores: SW, SH, SB
- [ ] **Passo 6e** — B-Type branches: BEQ, BNE, BLT, BGE, BLTU, BGEU
- [ ] **Passo 6f** — J-Type / JALR: JAL, JALR
- [ ] **Passo 6g** — U-Type: LUI, AUIPC

### Testes
- [ ] **Passo 9** — Teste de integração: executar fibonacci de `testdata/hex.txt`, verificar resultado em `a0`

---

## Instruções ainda faltando no decoder

- [ ] LBU (I-Type, funct3=0x4 em OP_LOAD)
- [ ] U-Type: LUI (0x37) e AUIPC (0x17) — `utype/type.go` existe mas instâncias específicas?
