# 🚀 Rate Limiter em Go

[![Go Version](https://img.shields.io/badge/Go-1.21+-blue.svg)](https://golang.org)
[![Docker](https://img.shields.io/badge/Docker-Required-blue.svg)](https://docker.com)
[![Redis](https://img.shields.io/badge/Redis-7-red.svg)](https://redis.io)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)
[![Status](https://img.shields.io/badge/Status-Active-brightgreen.svg)](https://github.com)

## 📋 Sobre o Projeto

Um **rate limiter robusto e escalável** desenvolvido em Go, capaz de controlar o tráfego de requisições para aplicações web com suporte a limitação por **IP** ou **Token de Acesso**. Implementado com padrão **Strategy** para flexibilidade de armazenamento e **Middleware** para integração simples com servidores HTTP.

### 🎯 Objetivo do Desafio

Criar um rate limiter em Go que limite o número máximo de requisições por segundo com base em:
- 🔒 **Endereço IP** - Restringir requisições por IP
- 🔑 **Token de Acesso** - Limites customizados por token
- ⚡ **Priorização** - Token tem precedência sobre IP

### 🏆 Funcionalidades Implementadas

- ✅ **Rate Limiting por IP** - Limite padrão configurável
- ✅ **Rate Limiting por Token** - Limites customizados por token
- ✅ **Middleware HTTP** - Injeção simples em qualquer servidor
- ✅ **Strategy Pattern** - Fácil troca de persistência
- ✅ **Redis Storage** - Persistência distribuída
- ✅ **Configuração via .env** - Fácil customização
- ✅ **Block Duration** - Tempo customizável de bloqueio
- ✅ **Response HTTP 429** - Resposta padrão quando limite é excedido
- ✅ **Testes Automatizados** - Cobertura de testes completa
- ✅ **Docker Compose** - Deploy simplificado

---

## 🚀 Como Executar

### Pré-requisitos

- [Go 1.21+](https://golang.org/dl/)
- [Docker](https://www.docker.com/)
- [Docker Compose](https://docs.docker.com/compose/)

### ✨ Execução Automática

```bash
# Subir Redis + Aplicação com Docker Compose
docker compose up --build -d
```

> ⚠️ **Primeira execução**: Pode demorar alguns segundos para Redis e aplicação ficarem prontos

### ✅ Confirmação dos Serviços

```bash
# Verificar status dos containers
docker compose ps

# Ver logs em tempo real
docker compose logs -f app

# Testar endpoint de health
curl http://localhost:8080/health
```

### 🔄 Comandos Úteis

```bash
# Parar todos os serviços
docker compose down

# Rebuild completo (limpar volumes)
docker compose down -v
docker compose up --build -d

# Ver logs específicos
docker compose logs redis
docker compose logs app
```

---

## 🧪 Como Testar

### 🌐 **REST API** - Porta 8080

#### Health Check
```http
GET http://localhost:8080/health
```

#### Requisição Simples (sem token)
```http
GET http://localhost:8080/api/
```

#### Requisição com Token (limite maior)
```http
GET http://localhost:8080/api/test
API_KEY: token123
```

#### Teste com Premium Token
```http
GET http://localhost:8080/api/test
API_KEY: premium-token
```

> 📁 **Arquivo de testes:** [api/requests.http](api/requests.http)

#### Respostas Esperadas

**Sucesso (200):**
```json
{
  "message": "Request accepted! Rate limiter working correctly.",
  "timestamp": "2024-01-02T10:30:00Z"
}
```

**Limite Excedido (429):**
```json
{
  "error": "you have reached the maximum number of requests or actions allowed within a certain time frame"
}
```

### 🔧 Rodando Testes Unitários

```bash
# Local (com Redis rodando)
go test -v ./internal/limiter

# Com coverage
go test -cover ./internal/limiter
```

---

## ⚙️ Configuração

### Via Arquivo `.env`

```env
# Rate Limiting por IP
RATE_LIMIT_IP=5                # requisições por segundo
IP_BLOCK_DURATION=300          # segundos

# Rate Limiting por Token
RATE_LIMIT_TOKEN=10            # requisições por segundo
TOKEN_BLOCK_DURATION=600       # segundos

# Redis
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_DB=0

# Servidor
SERVER_PORT=8080
```

### Via Variáveis de Ambiente

```bash
export RATE_LIMIT_IP=5
export IP_BLOCK_DURATION=300
export RATE_LIMIT_TOKEN=10
export TOKEN_BLOCK_DURATION=600
docker compose up -d
```

### Exemplos de Configuração

| Cenário | RATE_LIMIT_IP | IP_BLOCK_DURATION | Comportamento |
|---------|---|---|---|
| Normal | 5 | 300 | Bloqueia após 5 req/s por 5 min |
| Restritivo | 1 | 600 | 1 req/s, 10 min bloqueado |
| Permissivo | 100 | 60 | 100 req/s, 1 min bloqueado |

---

## 📊 Arquitetura

```
🐳 DOCKER ARCHITECTURE
┌─────────────────────────────────────┐
│      Container Services             │
│   Redis       │       App           │
│   :6379       │ :8080               │
└─────────────────────────────────────┘
                 ⬇️
┌─────────────────────────────────────┐
│      APPLICATION LAYERS             │
│   HTTP Middleware                   │
├─────────────────────────────────────┤
│   Rate Limiter Logic                │
│   (IP / Token)                      │
├─────────────────────────────────────┤
│   Storage Strategy                  │
│   (Redis / Custom)                  │
└─────────────────────────────────────┘
```

### Fluxo de Requisição

```
┌─ HTTP Request ─┐
│  Header: API_KEY (opcional)
│  IP: 192.168.1.1
└────────────────┘
         ⬇️
┌──────────────────────────┐
│ RateLimiterMiddleware    │
│ Extrai IP e Token        │
└────────────────┘
         ⬇️
┌──────────────────────────┐
│ RateLimiter.Allow()      │
│ Token? Usa token config  │
│ Não? Usa IP config       │
└────────────────┘
         ⬇️
┌──────────────────────────┐
│ Redis Storage            │
│ Incrementa contador      │
│ Retorna valor atual      │
└────────────────┘
         ⬇️
┌──────────────────────────┐
│ Validação                │
│ contador < limite? ✅    │
│ contador >= limite? ❌   │
└────────────────┘
         ⬇️
┌─ HTTP Response ─┐
│ 200 OK ou 429   │
└────────────────┘
```

---

## 📁 Estrutura do Projeto

```
.
├── cmd/
│   └── main.go                    # Aplicação principal
│
├── internal/
│   ├── config/
│   │   └── config.go              # Carregamento de configuração
│   ├── limiter/
│   │   ├── limiter.go             # Lógica de rate limiting
│   │   └── limiter_test.go        # Testes unitários
│   ├── middleware/
│   │   └── middleware.go          # Middleware HTTP
│   └── strategy/
│       ├── strategy.go            # Interface de strategy
│       └── redis.go               # Implementação Redis
│
├── api/
│   └── requests.http              # Testes HTTP
│
├── docker-compose.yml             # Orquestração de containers
├── Dockerfile                     # Imagem Docker
├── go.mod                         # Dependências Go
├── .env                           # Configurações locais
├── .env.example                   # Exemplo de configuração
├── .gitignore                     # Arquivos ignorados
├── LICENSE                        # Licença MIT
└── README.md                      # Este arquivo
```

---

## 🛠️ Tecnologias Utilizadas

### Core
- **Go 1.21** - Linguagem principal
- **Redis 7** - Armazenamento em memória distribuído

### Padrões de Design
- **Strategy Pattern** - Interface de persistência
- **Middleware Pattern** - Integração HTTP
- **Factory Pattern** - Criação de componentes

### DevOps
- **Docker & Docker Compose** - Containerização
- **Health Checks** - Monitoramento de containers

---

## 🧠 Como Funciona

### Rate Limiting por IP

```
Requisição do IP 192.168.1.1
↓
Incrementa contador Redis: "limiter:ip:192.168.1.1"
↓
Primeiro acesso → Set TTL (Time To Live)
↓
Próximos acessos → Verifica contador
  ├─ contador < limite → Permite ✅
  └─ contador >= limite → Bloqueia (429) ❌
↓
TTL expira → Contador reseta
```

### Rate Limiting por Token

```
Requisição com header API_KEY: token123
↓
Busca configuração do token
↓
Incrementa contador Redis: "limiter:token:token123"
↓
Usa limite customizado do token (default maior)
↓
Token >= limite configurado → Bloqueia (429) ❌
↓
TTL expira → Contador reseta
```

### Exemplo Prático

**Configuração:**
- IP Limit: 5 req/s
- Token Limit: 100 req/s
- IP Block Duration: 300s (5 min)
- Token Block Duration: 600s (10 min)

**Cenário 1 - Limitação por IP:**
```
Tempo 0s:   IP 192.168.1.1 faz 5 requisições ✅
Tempo 0.5s: IP 192.168.1.1 tenta 6ª requisição ❌ → 429
Tempo 5min: Contador expira → IP pode requisitar novamente ✅
```

**Cenário 2 - Limitação por Token:**
```
Tempo 0s:   token123 faz 50 requisições ✅
Tempo 1s:   token123 faz mais 50 requisições ✅
Tempo 2s:   token123 tenta 101ª requisição ❌ → 429
Tempo 10min: Contador expira → token pode requisitar novamente ✅
```

---

## 🔒 Segurança

- 🔐 **Headers customizados** - Suporta X-Forwarded-For, X-Real-IP
- 🛡️ **TTL automático** - Contadores expiram automaticamente
- 🚫 **Bloqueio temporário** - Não banimento permanente
- 📊 **Sem vazamento de dados** - Contadores isolados por chave
- 🔑 **Tokens separados** - Cada token tem contador próprio

---

## 📈 Performance

- ⚡ **Sub-milissegundo** - Redis fornece latência ultra-baixa
- 🔄 **Escalável horizontalmente** - Redis é distribuído
- 💾 **Eficiente em memória** - Apenas 1-2 chaves por cliente
- 🚀 **Zero-downtime** - Novas configurações via .env

---

## 🐛 Troubleshooting

### Redis não conecta

```bash
# Verificar se Redis está rodando
docker compose ps

# Verificar logs do Redis
docker compose logs redis

# Reconectar
docker compose down
docker compose up -d
```

### Porta 8080 já em uso

```bash
# Usar porta diferente
SERVER_PORT=8081 docker compose up -d
```

---

## 📚 Exemplos de Uso

### Integração com seu servidor HTTP

```go
import (
    "github.com/SaraPMC/GO-desafio-rate-limiter/internal/limiter"
    "github.com/SaraPMC/GO-desafio-rate-limiter/internal/middleware"
    "github.com/SaraPMC/GO-desafio-rate-limiter/internal/strategy"
)

// Criar storage
storage, _ := strategy.NewRedisStorage("localhost", 6379, 0)

// Criar rate limiter
rl := limiter.NewRateLimiter(storage, 5, 300)

// Aplicar middleware
mux := http.NewServeMux()
handler := middleware.RateLimiterMiddleware(rl)(mux)
```

### Customizar tokens

```go
// Token com limite maior
rl.ConfigureToken("api-key-premium", 1000, 60)

// Token com expiração diferente
rl.ConfigureToken("api-key-hourly", 10000, 3600)
```

---

## 📝 Licença

Este projeto está licenciado sob a Licença MIT - veja o arquivo [LICENSE](LICENSE) para detalhes.
