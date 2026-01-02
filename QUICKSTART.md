# ⚡ Quick Start - Rate Limiter

## 30 segundos para começar

### 1️⃣ Clone e Entre no Diretório
```bash
cd desafio-rate-limiter
```

### 2️⃣ Suba os Containers
```bash
docker compose up --build -d
```

### 3️⃣ Teste!
```bash
# Health check
curl http://localhost:8080/health

# Requisição aceita
curl http://localhost:8080/api/

# Com token
curl -H "API_KEY: token123" http://localhost:8080/api/test
```

---

## 📋 Configuração Rápida

### Limites Padrão (em `.env`)
- **IP**: 5 req/s, bloqueio de 300s (5 min)
- **Token**: 10 req/s, bloqueio de 600s (10 min)

### Customizar
Edite `.env` e reinicie:
```bash
docker compose restart app
```

---

## 🧪 Testar Limite Excedido

Envie 6 requisições rápidas (limite IP é 5):
```bash
for i in {1..6}; do
  curl http://localhost:8080/api/
done
```

Verá 200 OK para as 5 primeiras e 429 na sexta!

---

## 📊 Monitorar

```bash
# Logs em tempo real
docker compose logs -f app

# Status dos containers
docker compose ps

# Logs do Redis
docker compose logs redis
```

---

## 🛑 Parar

```bash
docker compose down
```

---

## 🚀 Próximos Passos

1. Leia [README.md](README.md) para documentação completa
2. Veja [api/requests.http](api/requests.http) para mais testes
3. Revise [internal/limiter/limiter.go](internal/limiter/limiter.go) para entender a lógica
4. Rode testes: `go test -v ./internal/limiter`

---

**Dúvidas?** Veja [README.md#-troubleshooting](README.md#-troubleshooting)
