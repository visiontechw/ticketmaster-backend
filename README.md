# TicketMaster Backend

Sistema de venda de ingressos online de alta concorrência desenvolvido em Go com arquitetura limpa (Clean Architecture).

## 📋 Visão Geral

TicketMaster é um backend robusto e escalável para gerenciar vendas de ingressos para eventos com alta concorrência. O sistema foi projetado para lidar com picos de tráfego massivos, garantindo consistência de dados e experiência do usuário otimizada.

**Características principais:**
- ✅ Alta concorrência (suporta milhares de requisições simultâneas)
- ✅ Transações ACID garantidas via PostgreSQL
- ✅ API RESTful rápida com Gin
- ✅ ORM moderno com GORM
- ✅ Arquitetura limpa e testável
- ✅ Validação e tratamento robusto de erros
- ✅ Documentação de API com Swagger

## 🏗️ Arquitetura

O projeto segue os princípios de **Clean Architecture**, separando responsabilidades em camadas:

```
internal/
├── domain/           # Entidades, interfaces e regras de negócio
├── app/              # Use cases, serviços de aplicação
│   ├── services/     # Lógica de negócio
│   ├── dto/          # Data Transfer Objects
│   └── handlers/     # HTTP handlers
└── infra/            # Implementações de repositórios, banco de dados, cache
    ├── repositories/ # Camada de persistência
    ├── database/     # Conexão e migração de BD
    ├── cache/        # Cache (Redis)
    └── config/       # Configurações
```

### Fluxo de Requisição

```
HTTP Request → Handler (HTTP) → Service (Use Case) → Repository → Database
                    ↓
              Response
```

## 🛠️ Stack Tecnológico

| Componente | Ferramenta | Versão |
|------------|-----------|--------|
| Linguagem | Go | 1.26.2+ |
| Framework Web | [Gin](https://gin-gonic.com/) | v1.9.x+ |
| Banco de Dados | PostgreSQL | 14+ |
| ORM | [GORM](https://gorm.io/) | v1.25.x+ |
| Validação | [go-playground/validator](https://github.com/go-playground/validator) | v10.x+ |
| Logging | [Logrus](https://github.com/sirupsen/logrus) | v1.9.x+ |
| Utilitários | [Viper](https://github.com/spf13/viper) | v1.17.x+ |

## 📦 Dependências

Principais dependências do projeto:

```bash
go get -u gorm.io/gorm
go get -u gorm.io/driver/postgres
go get -u github.com/gin-gonic/gin
go get -u github.com/go-playground/validator/v10
go get -u github.com/sirupsen/logrus
go get -u github.com/spf13/viper
go get -u github.com/google/uuid
```

## 🚀 Quick Start

### Pré-requisitos

- Go 1.26.2 ou superior
- PostgreSQL 14 ou superior
- Docker e Docker Compose (opcional, para ambiente local)

### Instalação

1. **Clone o repositório:**
```bash
git clone https://github.com/visiontechw/ticketmaster.git
cd ticketmaster-backend
```

2. **Instale as dependências:**
```bash
go mod download
go mod tidy
```

3. **Configure variáveis de ambiente:**
```bash
cp .env.example .env
# Edite o .env com suas configurações
```

4. **Execute as migrações do banco de dados:**
```bash
go run ./cmd/migrate/main.go
```

5. **Inicie o servidor:**
```bash
go run ./cmd/server/main.go
```

O servidor estará disponível em `http://localhost:3001`

### Com Docker Compose

```bash
docker-compose up -d
```

## 📝 Estrutura de Pastas

```
ticketmaster-backend/
├── api/                    # Definições de API (OpenAPI, Swagger)
├── cmd/                    # Pontos de entrada da aplicação
│   ├── server/             # Servidor HTTP
│   
├── internal/              # Código interno privado da aplicação
│   ├── domain/            # Modelos de domínio e interfaces
│   ├── app/               # Use cases e serviços de aplicação
│   │   ├── services/      # Logica de negócio central
│   │   ├── dto/           # Data Transfer Objects
│   │   ├── handlers/      # HTTP request handlers
│   │   └── middleware/    # Middleware HTTP
│   └── infra/             # Camada de infraestrutura
│       ├── database/      # Configuração do BD e modelos
│       ├── repositories/  # Implementação de repositórios
│       ├── cache/         # Cache e Redis
│       └── config/        # Carregamento de configurações
├── pkg/                   # Código reutilizável (utilitários, helpers)
├── scripts/               # Scripts de utilitário
├── tests/                 # Testes de integração e E2E
├── deployments/           # Arquivos de deployment (Kubernetes, Docker)
├── go.mod                 # Definição de módulo Go
├── go.sum                 # Hash das dependências
└── README.md              
```

## 🗄️ Modelos de Dados

### Usuários
```json
{
  "id": "uuid",
  "email": "string",
  "name": "string",
  "phone": "string",
  "password_hash": "string",
  "created_at": "timestamp",
  "updated_at": "timestamp"
}
```

### Eventos
```json
{
  "id": "uuid",
  "title": "string",
  "description": "string",
  "starts_at": "timestamp",
  "ends_at": "timestamp",
  "venue": "string",
  "total_seats": "integer",
  "available_seats": "integer",
  "created_at": "timestamp"
}
```

### Ingressos
```json
{
  "id": "uuid",
  "event_id": "uuid",
  "user_id": "uuid",
  "price": "decimal",
  "status": "available|sold|cancelled",
  "purchased_at": "timestamp",
  "created_at": "timestamp"
}
```

## 🔌 API Endpoints

### Autenticação
- `POST /auth/register` - Registrar novo usuário
- `POST /auth/login` - Login de usuário
- `POST /auth/logout` - Logout

### Eventos
- `GET /events` - Listar eventos
- `GET /events/:id` - Detalhes do evento
- `POST /events` - Criar evento (admin)
- `PUT /events/:id` - Atualizar evento (admin)
- `DELETE /events/:id` - Deletar evento (admin)

### Ingressos
- `POST /tickets/buy` - Comprar ingresso
- `GET /tickets/my-tickets` - Meus ingressos
- `POST /tickets/:id/cancel` - Cancelar ingresso
- `GET /tickets/:id` - Detalhes do ingresso

### Pagamentos
- `POST /payments` - Processar pagamento
- `GET /payments/:id` - Status do pagamento

## ⚡ Performance e Concorrência

O sistema foi otimizado para alta concorrência:

### Estratégias Implementadas:
1. **Connection Pooling**: GORM com pool de conexões configurável
2. **Indexação Eficiente**: Índices em colunas frequentemente consultadas
3. **Caching**: Redis para dados que mudam raramente
4. **Rate Limiting**: Proteção contra abuso
5. **Middleware de Concorrência**: Limita requisições simultâneas por usuário
6. **Transações**: Garante atomicidade em operações críticas

### Recomendações de Produção:
```env
DB_MAX_OPEN_CONNS=50
DB_MAX_IDLE_CONNS=10
DB_CONN_MAX_LIFETIME=5m
CACHE_TTL=5m
RATE_LIMIT=1000/minute
```

## 🧪 Testes

### Executar todos os testes:
```bash
go test ./...
```

### Com cobertura:
```bash
go test -cover ./...
```

### Testes específicos:
```bash
go test -run TestName ./internal/app/services
```

### Teste de carga (load test):
```bash
go run ./tests/load/main.go
```

## 📊 Monitoramento

O sistema possui endpoints para monitoramento:

- `GET /health` - Health check
- `GET /metrics` - Métricas Prometheus
- `GET /version` - Versão da aplicação

## 🔐 Segurança

- ✅ Senhas com hash bcrypt
- ✅ JWT para autenticação
- ✅ HTTPS em produção
- ✅ CORS configurável
- ✅ Rate limiting
- ✅ Validação de entrada
- ✅ SQL Injection prevention (GORM parametrizado)

## 🌍 Variáveis de Ambiente

```env
# Servidor
SERVER_PORT=8080
SERVER_ENV=development|production

# Banco de Dados
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=password
DB_NAME=ticketmaster
DB_SSLMODE=disable
DB_MAX_OPEN_CONNS=50
DB_MAX_IDLE_CONNS=10

# JWT
JWT_SECRET_KEY=sua-chave-secreta-aqui
JWT_EXPIRATION=24h

# Redis (Opcional)
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASS=

# Logging
LOG_LEVEL=info
LOG_FORMAT=json|text
```

## 🐳 Docker

### Build da imagem:
```bash
docker build -t ticketmaster:latest .
```

### Run do container:
```bash
docker run -p 8080:8080 --env-file .env ticketmaster:latest
```

## 📚 Documentação da API

Swagger está disponível em: `http://localhost:8080/swagger/index.html`

Para gerar documentação:
```bash
swag init -g ./cmd/server/main.go
```

## 🤝 Contribuindo

1. Fork o projeto
2. Crie uma branch para sua feature (`git checkout -b feature/AmazingFeature`)
3. Commit suas mudanças (`git commit -m 'Add some AmazingFeature'`)
4. Push para a branch (`git push origin feature/AmazingFeature`)
5. Abra um Pull Request

### Padrões de Código:
- Siga as convenções Go (use `go fmt`)
- Adicione testes para novas funcionalidades
- Mantenha a cobertura de testes acima de 80%
- Documente funções públicas

## 📝 Convenções de Commit

```
feat: Adiciona nova funcionalidade
fix: Corrige um bug
refactor: Refactora código sem mudar funcionalidade
docs: Altera documentação
test: Adiciona ou modifica testes
chore: Mudanças que não afetam o código (deps, config)
```

## 📖 Recursos Adicionais

- [Clean Architecture](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html)
- [Go Documentation](https://golang.org/doc/)
- [Gin Framework](https://gin-gonic.com/)
- [GORM Guide](https://gorm.io/docs/)
- [PostgreSQL Best Practices](https://www.postgresql.org/)

## 📱 Suporte

Para relatórios de bugs ou sugestões, abra uma issue no repositório.

## 📄 Licença

Distribuído sob a licença MIT. Veja `LICENSE` para mais informações.

## 👥 Autor

Desenvolvido por **VisioNTech Workspace**

---

**Última atualização:** Abril de 2026

migrate create -ext sql -dir db/migrations -seq nome_migrate

migrate -path db/migrations -database "postgres://dev:local%21%21dev@127.0.0.1:5432/ticketfacil?sslmode=disable" up




