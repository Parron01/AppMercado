# Dockerfile - Produção
FROM golang:1.24.3 AS builder
WORKDIR /app

# Copiar os arquivos necessários para o build
COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Compilar o binário do servidor
RUN go build -o app cmd/server/main.go

# Etapa 2: Imagem final
FROM gcr.io/distroless/base-debian11
WORKDIR /app

# Copiar o binário compilado
COPY --from=builder /app/app .

# Garantir que o binário seja executável
RUN chmod +x /app/app

# Porta exposta (definida nas variáveis de ambiente)
EXPOSE 8080

# Comando para iniciar o servidor
CMD ["./app"]
