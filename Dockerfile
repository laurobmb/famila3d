# Utiliza uma imagem base super leve
FROM alpine:latest

# Define o diretório de trabalho dentro do container
WORKDIR /app

# Copia o binário Linux gerado previamente pelo script bash
COPY dist/familia3d.bin ./familia3d.bin

# Copia as pastas de assets vitais para o servidor funcionar (já que não há mais go:embed)
COPY templates/ ./templates/
COPY static/ ./static/

# Garante a permissão de execução no binário (boa prática em containers)
RUN chmod +x ./familia3d.bin

# Expõe a porta que a aplicação está escutando
EXPOSE 8080

# Define o comando de inicialização
CMD ["./familia3d.bin"]