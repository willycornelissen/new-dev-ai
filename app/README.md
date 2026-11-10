# new-dev-ai CLI

Ferramenta de linha de comando desenvolvida em **Go (Go 1.23)** com **Cobra** para geração de scaffold automatizado de projetos baseados no ecossistema [new-dev-ai](https://github.com/willycornelissen/new-dev-ai).

## 🚀 O que a ferramenta faz

1. Clona o repositório base `https://github.com/willycornelissen/new-dev-ai`.
2. Remove o diretório `app/` e o arquivo `quickstart.md` do novo projeto gerado.
3. Remove os metadados de histórico original do Git (`.git/`) e inicializa um repositório limpo (`git init`).
4. Inicializa o arquivo `README.md` contendo exclusivamente o nome do projeto passado como parâmetro como título (`# <nome-do-projeto>\n`).
5. Preserva todas as instruções de agentes e configurações de desenvolvimento com IA (`.agents/`, `AGENTS.md`, etc.).

---

## 📦 Como Compilar

```bash
# Na pasta app/
go build -o new-dev-ai main.go
```

---

## 💻 Como Usar

### Uso Básico

```bash
./new-dev-ai <nome-do-projeto>
```

**Exemplo:**
```bash
./new-dev-ai hello-world
```

### Flags Opcionais

```bash
Flags:
      --dir string    Caminho do diretório de destino (padrão: nome do projeto)
      --repo string   URL do repositório template (padrão: "https://github.com/willycornelissen/new-dev-ai")
  -h, --help          Exibe ajuda
```

**Exemplo com diretório customizado:**
```bash
./new-dev-ai meu-projeto --dir /caminho/para/meu-projeto
```

---

## 🧪 Testes e Cobertura

O projeto foi construído seguindo rigorosamente a metodologia **TDD (Red-Green-Refactor)**, com cobertura de testes unitários e de integração superior a 89%:

```bash
# Executar todos os testes
go test -v ./...

# Executar testes com relatório de cobertura
go test -coverprofile=coverage.out ./...
go tool cover -func=coverage.out
```
