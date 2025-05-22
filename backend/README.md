# AppMercado – Backend (Go)

![Go](https://img.shields.io/badge/Go-1.24.x-blue?logo=go)
![Gin](https://img.shields.io/badge/Gin-1.10-green?logo=go)
![PostgreSQL](https://img.shields.io/badge/PostgreSQL-14-blue?logo=postgresql)
![Docker](https://img.shields.io/badge/Docker-Compose-blue?logo=docker)

> API RESTful para gerenciamento de listas de compras, produtos, categorias, histórico de preços e relacionamento entre usuários, categorias e produtos.

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
backend/
│
├── cmd/server/               # ponto de entrada (main.go)
│
├── internal/                 # código privado (não importável fora do módulo)
│   ├── handlers/             # controllers – HTTP handlers (Auth, User, Category, Product, Purchase, PriceHistory, UserCategoryProduct)
│   ├── services/             # regra de negócio
│   ├── repositories/         # persistência (PostgreSQL, GORM)
│   └── models/               # structs refletindo tabelas
│
├── pkg/config/               # utilitários exportáveis (carrega .env via Viper)
├── pkg/utils/                # funções utilitárias (formatação, etc)
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

---

## 🚀 Executando

### Desenvolvimento (Hot Reload)

```bash
docker-compose -f docker-compose.yml -f docker-compose.override.yml up --build
```

* Modificou código? O `air` recompila e reinicia automaticamente.

### Produção / Teste sem Hot Reload

```bash
docker-compose up --build -d
```

* API disponível em `http://localhost:8080`.
* PostgreSQL em `localhost:5432` usando as credenciais do `.env`.

---

## 🗂️ Endpoints Principais

| Método | Rota             | Descrição                                      |
| ------ | ---------------- | ---------------------------------------------- |
| POST   | `/auth/register` | Registro de usuário                            |
| POST   | `/auth/login`    | Login e emissão de JWT                         |
| GET    | `/users/all`     | Listar todos os usuários (admin)               |
| DELETE | `/users/delete/:id` | Deletar usuário (próprio ou admin)           |
| CRUD   | `/categories`    | Gerenciar categorias do usuário                |
| CRUD   | `/products`      | Gerenciar produtos (admin)                     |
| CRUD   | `/purchases`     | Registrar e consultar compras                  |
| CRUD   | `/price-history` | Consultar histórico de preços                  |
| CRUD   | `/user-category-products` | Relacionar produtos a categorias do usuário |

> **Nota:** Endpoints adicionais e detalhes de payloads podem ser consultados no código dos handlers.

---

## 🔒 Autenticação & Permissões

- JWT obrigatório para todas as rotas (exceto `/auth/register` e `/auth/login`).
- Papéis de usuário: `Admin`, `Standard`, `Guest`.
- Permissões de escrita em produtos são restritas a administradores.
- Categorias e compras são privadas por usuário.
- Admin pode listar e gerenciar todos os registros.

---

## 🗃️ Principais Modelos

- **User**: Usuário do sistema, com papel (role).
- **Category**: Categoria de produtos, associada a um usuário.
- **Product**: Produto global, gerenciado por admin.
- **Purchase**: Compra realizada por um usuário, com itens.
- **PurchaseItem**: Item de uma compra (produto, quantidade, preço).
- **PriceHistory**: Histórico de preços de produtos por compra.
- **UserCategoryProduct**: Relação entre usuário, categoria e produto.

---

## 💡 Contribuindo

1. Faça um fork / crie branch.
2. Siga o padrão de pastas (`internal/`, `pkg/`).
3. Execute `go vet` antes de submeter PR.

---

## 📜 Licença

MIT © 2025 Parron01 – Projeto acadêmico/pessoal.
