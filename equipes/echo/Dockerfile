# Estágio de Build - ATUALIZADO PARA 1.24
FROM golang:1.25-alpine AS builder

RUN apk add --no-cache git

WORKDIR /app

# Copia os arquivos de dependências
COPY backend/go.mod ./
COPY backend/go.sum* ./

# Agora com Go 1.25 ele vai conseguir ler o mod
RUN go mod download || true

# Copia o restante do código
COPY backend/ . 

# Compila o binário
RUN go build -o main .

# Estágio Final
FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/main .
# Garante que a pasta frontend existe para o servidor Go achar
COPY frontend/ ./frontend/
EXPOSE 8080
CMD ["./main"]