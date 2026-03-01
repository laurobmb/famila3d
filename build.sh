#!/bin/bash

# Interrompe o script caso algum comando falhe
set -e

# Define o nome dos arquivos de saída
LINUX_BIN="dist/familia3d.bin"
WINDOWS_BIN="dist/familia3d.exe"

echo "⚙️  Verificando/Criando diretório 'dist'..."
mkdir -p dist

echo "🐧 Compilando binário para Linux amd64..."
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ${LINUX_BIN} main.go
echo "✅ Criado: ${LINUX_BIN}"

echo "🪟 Compilando executável para Windows amd64..."
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o ${WINDOWS_BIN} main.go
echo "✅ Criado: ${WINDOWS_BIN}"

echo "🚀 Build finalizado com sucesso!"