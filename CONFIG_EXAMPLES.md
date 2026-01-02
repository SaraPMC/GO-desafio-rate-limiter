# Configuration Examples - Rate Limiter

## Default Configuration

```env
# Rate Limiting por IP (requisições por segundo)
RATE_LIMIT_IP=5

# Tempo de bloqueio para IP (segundos)
IP_BLOCK_DURATION=300

# Rate Limiting por Token (requisições por segundo)
RATE_LIMIT_TOKEN=10

# Tempo de bloqueio para Token (segundos)
TOKEN_BLOCK_DURATION=600

# Redis Configuration
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_DB=0

# Server Configuration
SERVER_PORT=8080
```

---

## Development Configuration

```env
RATE_LIMIT_IP=100
IP_BLOCK_DURATION=60
RATE_LIMIT_TOKEN=1000
TOKEN_BLOCK_DURATION=60
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_DB=0
SERVER_PORT=8080
```

---

## Strict Configuration (API Protection)

```env
RATE_LIMIT_IP=1
IP_BLOCK_DURATION=600
RATE_LIMIT_TOKEN=5
TOKEN_BLOCK_DURATION=1200
REDIS_HOST=redis.example.com
REDIS_PORT=6379
REDIS_DB=1
SERVER_PORT=8080
```

---

## High Throughput Configuration

```env
RATE_LIMIT_IP=1000
IP_BLOCK_DURATION=30
RATE_LIMIT_TOKEN=10000
TOKEN_BLOCK_DURATION=30
REDIS_HOST=redis-cluster.example.com
REDIS_PORT=6379
REDIS_DB=0
SERVER_PORT=8080
```

---

## Production Configuration (Recommended)

```env
# Conservative IP limits
RATE_LIMIT_IP=10
IP_BLOCK_DURATION=300

# Higher token limits (for trusted clients)
RATE_LIMIT_TOKEN=100
TOKEN_BLOCK_DURATION=600

# Production Redis
REDIS_HOST=redis-prod.example.com
REDIS_PORT=6379
REDIS_DB=2

# Standard port
SERVER_PORT=8080
```

---

## Environment Variables Override

You can override `.env` with system environment variables:

```bash
# Linux/macOS
export RATE_LIMIT_IP=50
export IP_BLOCK_DURATION=180
docker compose up -d

# Windows PowerShell
$env:RATE_LIMIT_IP = 50
$env:IP_BLOCK_DURATION = 180
docker compose up -d
```

---

## Docker Compose Override

Create `docker-compose.override.yml`:

```yaml
version: '3.8'
services:
  app:
    environment:
      - RATE_LIMIT_IP=20
      - IP_BLOCK_DURATION=180
      - RATE_LIMIT_TOKEN=50
      - TOKEN_BLOCK_DURATION=300
```

---

## Token Configuration Examples

### Basic Tokens (in main.go)

```go
// Standard tier
rl.ConfigureToken("token-standard", 10, 600)

// Premium tier
rl.ConfigureToken("token-premium", 100, 600)

// Enterprise tier
rl.ConfigureToken("token-enterprise", 1000, 600)
```

### Expiring Tokens

```go
// Hourly limit
rl.ConfigureToken("token-hourly", 10000, 3600)

// Daily limit
rl.ConfigureToken("token-daily", 100000, 86400)
```

---

## Testing Configurations

### Tight Testing
```env
RATE_LIMIT_IP=2
IP_BLOCK_DURATION=5
RATE_LIMIT_TOKEN=5
TOKEN_BLOCK_DURATION=5
```

### Load Testing
```env
RATE_LIMIT_IP=10000
IP_BLOCK_DURATION=1
RATE_LIMIT_TOKEN=100000
TOKEN_BLOCK_DURATION=1
```

---

**Escolha a configuração apropriada para seu caso de uso!**
