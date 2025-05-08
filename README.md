# ✅ Resumo Completo do Projeto de Gerenciamento de Compras - Definições Até Aqui

---

## 🌐 Visão Geral

**Projeto:** Sistema de gerenciamento de compras de mercado.  
**Objetivo:** Facilitar o gerenciamento de listas de compras, armazenamento de preços pagos, cálculo de preços médios e visualização de gráficos históricos.

**Arquitetura:**  
Aplicação full-stack com frontend em **Angular (PWA)** e backend em **Go (API RESTful)**, hospedada em um servidor VPS com **Docker** e **Nginx** como proxy reverso.

---

## ✅ Frontend (Angular + PWA)

- **SPA (Single Page Application):** Aplicação web com PWA, suportando uso offline.
- **Service Worker:** Controle de cache e atualização dinâmica.
- **Módulos:**  
  - Lista de Compras  
  - Histórico de Compras  
  - Gráficos e Relatórios  
- **Comunicação com o Backend:** Via HTTP com Angular HttpClient.
- **Design Responsivo:** Mobile-first, com suporte offline garantido pelo Service Worker.

---

## ✅ Backend (Go)

- **API RESTful:** Desenvolvida em Go, estruturada com controllers, services e repositories.
- **Autenticação JWT:** Gerenciamento de sessão seguro.
- **Estrutura Clean Architecture:**
  - **Controllers (Handlers):** Recebem e processam requisições HTTP.
  - **Services:** Contêm a lógica de negócios.
  - **Repositories:** Gerenciam o acesso ao banco de dados.
- **Principais Endpoints:**
  - Gerenciamento de produtos.
  - Gerenciamento de categorias de produtos.
  - Registro de compras e preços.
  - Controle de usuários e autenticação.

---

## ✅ Banco de Dados (PostgreSQL)

- **Relacional, normalizado até a 3ª forma normal.**
- **Tabelas:**
  - **User:** Controle de usuários (Admin, Standard, Guest).
  - **Category:** Categorias personalizadas de produtos.
  - **Product:** Produtos com preço médio calculado.
  - **Purchase:** Compras registradas com data e local.
  - **PriceHistory:** Histórico de preços pagos por produto.
  - **UserCategoryProduct:** Tabela intermediária que conecta produtos e categorias por usuário.

---

## ✅ Modelo de Banco de Dados (DER e Diagrama de Classes)

- **Diagrama ER:** Modela as entidades principais e seus relacionamentos:
  - Produtos são associados a categorias por meio da tabela UserCategoryProduct.
  - Cada usuário tem suas próprias categorias e histórico de preços.
- **Diagrama de Classes:** Define a estrutura do backend:
  - Classe **User** com métodos para criar categorias e registrar compras.
  - Classe **Product** com métodos para calcular preço médio e registrar preços.
  - Classe **Category** que organiza produtos.
  - Classe **Purchase** que gerencia as compras.
  - Classe **PriceHistory** que armazena o histórico de preços.
  - Classe **UserCategoryProduct** que associa produtos a categorias por usuário.

---

## ✅ Infraestrutura de Deploy

- **Servidor VPS (Ubuntu 22.04):**
  - **Nginx Proxy Reverso** (Portas 80 e 443) com SSL (HTTPS).
- **Containers Docker:**
  - Frontend (Angular + PWA).
  - Backend (Go API).
  - Banco de Dados (PostgreSQL).
- **Rede Interna Docker:** Comunicação segura entre os containers.
- **Controle de Acesso:** Apenas portas 22 (SSH), 80 e 443 estão expostas.

---

## ✅ Recursos Visuais e Diagramas

- **Diagrama de Entidade-Relacionamento (DER).**
- **Diagrama de Classes (Mermaid).**
- **Imagem da Arquitetura do Sistema (VPS, Nginx, Containers).**
