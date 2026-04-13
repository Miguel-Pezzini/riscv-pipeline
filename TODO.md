# TODO — RISC-V Pipeline Simulator

Itens organizados por prioridade. Bugs silenciosos (decodificação errada sem erro) primeiro.

---

## Melhorias no pipeline

- [x] **Suporte a `LH` e `LHU` no I-Type (loads de halfword)**
  - funct3=0x1 (LH) e funct3=0x5 (LHU)
  - Mesma estrutura de LW/LB

---

## Qualidade de código

- [x] **Inconsistência de naming: `NewANDI` vs `newXXX`**
  - `pkg/isa/itype/andi.go:9` — único construtor com `N` maiúsculo
  - Padronizar para `newANDI`

- [x] **Caminhos de arquivo hardcoded no main**
  - `cmd/resolver/main.go:16-17` — `testdata/bin.txt` e `testdata/hex.txt` fixos
  - Aceitar caminho via argumento de linha de comando (`os.Args`)
