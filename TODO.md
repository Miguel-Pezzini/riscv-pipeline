# TODO — RISC-V Pipeline Simulator

Itens organizados por prioridade. Bugs silenciosos (decodificação errada sem erro) primeiro.

---

## Bugs menores

- [ ] **Parâmetro `forwarding` ignorado no control hazard detector**
  - `pkg/hazard/control_detector.go:5`
  - O parâmetro é aceito mas nunca usado dentro da função
  - Remover o parâmetro da assinatura ou documentar por que está lá

---

## Instruções faltando — R-Type

- [ ] **AND** — funct3=0x7, funct7=0x00
- [ ] **OR** — funct3=0x6, funct7=0x00
- [ ] **XOR** — funct3=0x4, funct7=0x00
- [ ] **SLL** (shift left logical) — funct3=0x1, funct7=0x00
- [ ] **SRL** (shift right logical) — funct3=0x5, funct7=0x00
- [ ] **SRA** (shift right arithmetic) — funct3=0x5, funct7=0x20
- [ ] **SLT** (set less than) — funct3=0x2, funct7=0x00
- [ ] **SLTU** — funct3=0x3, funct7=0x00

---

## Instruções faltando — I-Type

- [ ] **XORI** — funct3=0x4
- [ ] **SLTI** — funct3=0x2
- [ ] **SLTIU** — funct3=0x3
- [ ] **SLLI** (shift imediato) — funct3=0x1
- [ ] **SRLI** — funct3=0x5, imm[11:5]=0x00
- [ ] **SRAI** — funct3=0x5, imm[11:5]=0x20

---

## Melhorias no pipeline

- [ ] **Reduzir overhead de control hazard com forwarding ativo**
  - Hoje: insere NOP para qualquer branch/jump ainda no pipeline
  - Melhoria: se forwarding=true, branch resolve no estágio EX — reduzir stalls de 3 para 1
  - `pkg/hazard/control_detector.go`

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
