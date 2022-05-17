# Projeto_GO
Projeto em GO para processamento de dados.

Este projeto foi criado para realizar o processamento de um arquivo TXT que possuí um layout específico e 
neste arquivo possui dados de compra de pessoas.
A aplicação realiza a extração de dados deste arquivo e insere os dados em uma tabela em um banco postgres.

O projeto esta estruturado da seguinte forma

```bash
├── projeto_go
│   ├──brdoc
│   ├── dbconfig/
│   ├── domain/
│   │  ├── **/*.go
├── docker-compose.yaml
├── go.mod
├── main.go
└── README.md
```
O pacote `brdoc/` foram copiados os arquivos do pacote https://github.com/paemuri/brdoc com o validador do documento CPF e CNPJ, 
foram copiados apenas os arquivos cpfcnpj e util do pacote de origem, os quais seriam utilizados nesse projeto.

O pacote `dbconfig` tem a configuração do banco de dados utilizado.
No pacote `domain` estão estruturados os controles de execução da aplicação e regras de negocio.

Para execução do projeto. 
Necessário ter uma base de dados Postgres com nome compras.
Executar go run main.go do projeto.

Pontos para melhoria
[]funcionamento do arquivo docker-compose e conexão da aplicação ao BD do docker
[]melhoria da estrutura conforme Domain Driven Design (DDD)
[]criação de testes
