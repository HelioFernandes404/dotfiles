# SSH Configuration

## ⚠️ IMPORTANTE - SEGURANÇA

Este diretório contém **APENAS** o arquivo `config` do SSH.

**NUNCA** adicione ao Git:
- ❌ Chaves privadas (`id_rsa`, `id_ed25519`, etc.)
- ❌ Chaves públicas com informações sensíveis
- ❌ `known_hosts` (contém histórico de conexões)
- ❌ Certificados ou tokens

## O que está versionado

✅ `ssh/config` - Configuração de hosts e aliases

## Uso

O arquivo `config` será copiado para `~/.ssh/config` via symlink.

### Permissões importantes

Após instalação, garanta que as permissões estejam corretas:

```bash
chmod 700 ~/.ssh
chmod 600 ~/.ssh/config
chmod 600 ~/.ssh/id_* (suas chaves privadas)
chmod 644 ~/.ssh/id_*.pub (suas chaves públicas)
```

## Revisar antes de commitar

Sempre revise o `ssh/config` antes de fazer commit para garantir que não há:
- Senhas em texto plano
- Tokens de autenticação
- Informações corporativas sensíveis
- IPs internos que não devem ser públicos
