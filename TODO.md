# TODO — RISC-V Pipeline Simulator

Itens organizados por prioridade. Bugs silenciosos (decodificação errada sem erro) primeiro.

---

## Instruções faltando — R-Type

- [x] **AND** — funct3=0x7, funct7=0x00
- [x] **OR** — funct3=0x6, funct7=0x00
- [x] **XOR** — funct3=0x4, funct7=0x00
- [x] **SLL** (shift left logical) — funct3=0x1, funct7=0x00
- [x] **SRL** (shift right logical) — funct3=0x5, funct7=0x00
- [x] **SRA** (shift right arithmetic) — funct3=0x5, funct7=0x20
- [x] **SLT** (set less than) — funct3=0x2, funct7=0x00
- [x] **SLTU** — funct3=0x3, funct7=0x00

---

## Instruções faltando — I-Type

- [x] **XORI** — funct3=0x4
- [x] **SLTI** — funct3=0x2
- [x] **SLTIU** — funct3=0x3
- [x] **SLLI** (shift imediato) — funct3=0x1
- [x] **SRLI** — funct3=0x5, imm[11:5]=0x00
- [x] **SRAI** — funct3=0x5, imm[11:5]=0x20

---

## Melhorias no pipeline

- [ ] **Suporte a `LH` e `LHU` no I-Type (loads de halfword)**
  - funct3=0x1 (LH) e funct3=0x5 (LHU)
  - Mesma estrutura de LW/LB

---

## Qualidade de código

- [ ] **Inconsistência de naming: `NewANDI` vs `newXXX`**
  - `pkg/isa/itype/andi.go:9` — único construtor com `N` maiúsculo
  - Padronizar para `newANDI`

- [ ] **Caminhos de arquivo hardcoded no main**
  - `cmd/resolver/main.go:16-17` — `testdata/bin.txt` e `testdata/hex.txt` fixos
  - Aceitar caminho via argumento de linha de comando (`os.Args`)
