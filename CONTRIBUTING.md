# 🤝 Guia de Contribuição

## Como Começar

1. **Fork** este repositório
2. **Clone** seu fork localmente
3. **Crie uma branch** para sua feature: `git checkout -b feature/sua-feature`
4. **Faça commits** descritivos das suas mudanças
5. **Push** para sua branch: `git push origin feature/sua-feature`
6. **Abra um Pull Request**

## Padrões de Código

### Go Style Guide

Seguimos o [Go Style Guide](https://google.github.io/styleguide/go/) do Google.

```bash
# Formatar código
go fmt ./...

# Verificar lint
golangci-lint run ./...
```

### Commits

Use mensagens descritivas:

```
feat: adicionar suporte a MongoDB
fix: corrigir race condition em contadores
docs: atualizar documentação de configuração
test: adicionar testes para edge cases
```

## Testes

Certifique-se de que todos os testes passam:

```bash
# Rodar todos os testes
go test -v ./...

# Com coverage
go test -cover ./...

# Testes específicos
go test -v -run TestIPRateLimit ./internal/limiter
```

## Criando uma Issue

Antes de criar uma Issue, verifique se não existe uma similar aberta.

**Descreva claramente:**
- O comportamento esperado
- O comportamento atual
- Como reproduzir o problema
- Seu ambiente (versão Go, SO, versão do Redis)

## Pull Request

1. Atualize a documentação relevante
2. Adicione ou atualize testes se necessário
3. Certifique-se de que os testes passam
4. Escreva uma descrição clara do que foi mudado

## Código de Conduta

Somos dedicados a fornecer um ambiente acolhedor e inspirador para todos. Leia nosso [Código de Conduta](CODE_OF_CONDUCT.md).

---

Obrigado por contribuir! 🎉
