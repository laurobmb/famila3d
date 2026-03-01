# Familia3D - Calculadora de Orçamentos 3D

Aplicação web desenvolvida em Go para automatizar o cálculo de custos e precificação de peças impressas em 3D. O sistema considera gastos com material, energia elétrica, margem de segurança operacional e markup de venda, fornecendo uma interface web responsiva e logs estruturados em formato JSON para fácil observabilidade.

## 🎯 Objetivo

Fornecer uma ferramenta de uso rápido e com baixo consumo de recursos (estaticamente compilada) para gerar orçamentos precisos de impressão 3D, pronta para rodar localmente ou ser deployada como um container em clusters (como OpenShift/Kubernetes).

## 🧮 Como Funciona o Cálculo

A lógica matemática segue a seguinte estrutura de composição de custos:

1. **Valor Filamento**: Calcula o custo proporcional ao peso da peça.
   `Custo Filamento = (Preço por Kg / 1000) * Peso em gramas`

2. **Valor Luz de Horas**: Estima o gasto energético com base no tempo de máquina.
   `Custo Energia = Tempo de Impressão (h) * Preço da Energia (R$/h)`

3. **Valor de Custo 1**: Soma direta do material e da energia.
   `Custo 1 = Custo Filamento + Custo Energia`

4. **Valor de Custo 2 (falhas + imp)**: Adiciona uma margem de segurança de 20% sobre o custo operacional (10% para possíveis falhas de impressão + 10% para depreciação/manutenção da impressora).
   `Custo 2 = Custo 1 * 1.20`

5. **Valor de Venda**: Preço final sugerido ao cliente com base no multiplicador de lucro.
   `Valor de Venda = Custo 2 * Markup`

---

## 🚀 Iniciando o Projeto em Go Puro (Desenvolvimento)

Para rodar o projeto localmente durante o desenvolvimento, sem compilar o binário final:

1. Certifique-se de ter o Go instalado (`go version`).
2. Clone o repositório e navegue até a pasta do projeto.
3. Execute o comando:
   ```bash
   go run main.go

```

4. Acesse no navegador: `http://localhost:8080`

---

## 🛠️ Executando o Script de Build e os Binários

O projeto possui um script Bash (`build.sh`) que realiza a compilação cruzada para Linux e Windows, gerando executáveis autônomos.

### Rodando o script:

1. Dê permissão de execução (apenas no primeiro uso no Linux/macOS ou WSL):
```bash
chmod +x build.sh

```


2. Execute o script:
```bash
./build.sh

```



Isso criará a pasta `dist/` com os arquivos `familia3d.bin` (Linux) e `familia3d.exe` (Windows).

### Executando os binários:

Como o sistema lê as pastas de *assets* estáticos diretamente do disco, **você deve executar o binário a partir da raiz do projeto**, onde as pastas `templates/` e `static/` estão localizadas.

* **No Linux:**
```bash
./dist/familia3d.bin

```


* **No Windows:**
```powershell
.\dist\familia3d.exe

```



---

## 🪟 Dicas para Executar no Windows 11

* **Uso do Binário Nativo:** Basta dar um duplo clique em `dist\familia3d.exe` pelo Windows Explorer ou executá-lo via PowerShell. Um terminal abrirá exibindo os logs estruturados e o servidor estará disponível no `localhost:8080`. Se o Firewall do Windows exibir um alerta, permita o acesso para conexões locais.
* **Executando o `build.sh` no Windows:** O prompt de comando padrão (CMD/PowerShell) não executa scripts `.sh` nativamente. Você tem duas alternativas:
1. Usar o terminal do **Git Bash** (que emula um ambiente Unix).
2. Usar o **WSL** (Windows Subsystem for Linux).


* **Containers via Podman:** Como você utiliza o Podman Desktop no Windows 11, você pode seguir as mesmas instruções da seção de Docker abaixo substituindo a palavra `docker` por `podman` no seu terminal.

---

## 🐳 Trabalhando com Docker / Podman

A aplicação inclui um `Dockerfile` otimizado, baseado no Alpine Linux, que copia o binário pré-compilado e as pastas essenciais para um ambiente leve e seguro.

### 1. Build da Imagem

Primeiro, garanta que o binário Linux foi gerado rodando o `./build.sh`. Depois, execute:

```bash
docker build -t familia3d:latest .

```

*(Se estiver usando Podman, use `podman build -t familia3d:latest .`)*

### 2. Executando o Container

Inicie o servidor mapeando a porta 8080 da sua máquina para a porta 8080 do container, rodando em background (`-d`):

```bash
docker run -d --name calc3d -p 8080:8080 familia3d:latest

```

### 3. Acompanhando os Logs

Para visualizar as requisições e cálculos processados em tempo real pelo sistema de logs nativo do Go:

```bash
docker logs -f calc3d

```

Para parar a aplicação:

```bash
docker stop calc3d

```
