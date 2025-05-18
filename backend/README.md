# AppMercado – Backend (Go)

![Go](https://img.shields.io/badge/Go-1.24.x-blue?logo=go)
![Gin](https://img.shields.io/badge/Gin-1.10-green?logo=go)
![PostgreSQL](https://img.shields.io/badge/PostgreSQL-14-blue?logo=postgresql)
![Docker](https://img.shields.io/badge/Docker-Compose-blue?logo=docker)

> Serviço **API RESTful** para gerenciamento de listas de compras, preços e histórico de gastos do projeto **AppMercado**.

---

## ✨ Tecnologias Principais

| Camada | Tecnologia / Biblioteca                                    | Função                                             |
| ------ | ---------------------------------------------------------- | -------------------------------------------------- |
| HTTP   | **Gin**                | Framework web minimalista, roteamento & middleware |
| ORM    | **GORM** + `gorm.io/driver/postgres`   | Mapear structs → tabelas PostgreSQL                |
| Config | **Viper**                | Carrega variáveis de ambiente / `.env`             |
| Auth   | **golang‑jwt/jwt/v4** | Geração e validação de tokens JWT                  |
| Build  | **Docker** / **Docker Compose**                            |  Containeriza API + banco PostgreSQL               |
| Dev    | **Air**                | Hot‑reloading em desenvolvimento                   |

---

## 📂 Estrutura de Pastas

```
appmercado/back-end/
│
├── cmd/server/               # ponto de entrada (main.go)
│
├── internal/                 # código privado (não importável fora do módulo)
│   ├── handlers/             # controllers – HTTP handlers (Auth etc.)
│   ├── services/             # regra de negócio
│   ├── repositories/         # persistência (PostgreSQL, GORM)
│   └── models/               # structs refletindo tabelas
│
├── pkg/config/               # utilitários exportáveis (carrega .env via Viper)
│
├── Dockerfile                # imagem otimizada p/ produção (distroless)
├── Dockerfile.dev            # imagem dev com Hot Reload (Air)
├── docker-compose.yml        # base (produção)
├── docker-compose.override.yml# override dev (Hot Reload)
├── .air.toml                 # ajustes do Air
├── .env.example              # modelo de variáveis de ambiente
├── go.mod / go.sum           # dependências & integridade
└── README.md                 # este arquivo
```

---

## ⚙️ Configuração de Ambiente

Crie um arquivo **`.env`** na raiz do repositório com o conteúdo abaixo e ajuste conforme a sua máquina ou VPS:

```env
# Servidor
SERVER_PORT=8080

# Banco de Dados
DB_HOST=appmercado-postgres
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=senha_segura
DB_NAME=appmercado

# JWT
JWT_SECRET=troque-por-uma-string-secreta
JWT_EXPIRATION_HOURS=72
```

> **Importante:** O `.env` nunca deve ser versionado. Ele já está no `.gitignore`.

## 🚀 Executando

### Desenvolvimento (Hot Reload)

```bash
# sobe API + PostgreSQL com Air
docker-compose -f docker-compose.yml -f docker-compose.override.yml up --build
```

* Modificou código? O `air` recompila e reinicia automaticamente.

### Produção / Teste sem Hot Reload

```bash
# build e sobe containers otimizados
docker-compose up --build -d
```

* API disponível em `http://localhost:8080`.
* PostgreSQL em `localhost:5432` usando as credenciais do `.env`.

---

## 🗂️ Endpoints Principais (MVP)

| Método | Rota             | Descrição                  |
| ------ | ---------------- | -------------------------- |
| POST   | `/auth/register` | Registro de usuário        |
| POST   | `/auth/login`    | Login e emissão de JWT     |
| CRUD   | `/products`      | Gerenciar produtos         |
| CRUD   | `/categories`    | Gerenciar categorias       |
| CRUD   | `/purchases`     | Registrar compras & preços |
| CRUD   | `/user-category-products` | Gerenciar UserCategoryProduct |

> **Nota:** Endpoints adicionais serão adicionados conforme evoluir o projeto.

---

## 🧪 Testes

```bash
go test ./...
```

---

## 💡 Contribuindo

1. Faça um fork / crie branch.
2. Siga o padrão de pastas (`internal/`, `pkg/`).
3. Execute `go vet` e `go test` antes de submeter PR.

---

## 📜 Licença

MIT © 2025 Parron01 – Projeto acadêmico/pessoal.
