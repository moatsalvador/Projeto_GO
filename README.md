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
├── Teste_Base.txt
├── Teste_BaseMenor.txt
└── README.md
```
O pacote `brdoc/` foram copiados os arquivos do pacote https://github.com/paemuri/brdoc com o validador do documento CPF e CNPJ, 
foram copiados apenas os arquivos cpfcnpj e util do pacote de origem, os quais seriam utilizados nesse projeto.

O pacote `dbconfig` tem a configuração do banco de dados utilizado.
No pacote `domain` estão estruturados os controles de execução da aplicação e regras de negocio.
No pacote domain existe os seguintes arquivos:
`compra.go` - com a struct de compra usado para o insert.
`dadosdb.go` - controles de interação com o banco de dados.
`processador.go` - Faz o controle do processamento do arquivo e validação dos dados.

O arquivo usado para o processamento é o Teste_Base.txt, conforme esta definido na funcao LeArquivo() do `processador.go`

Para execução do projeto. 
Necessário ter uma base de dados Postgres com nome compras.
A apliação pode ser executada com o comando go run main.go do projeto.

Pontos para melhoria e estudo futuro
[]Funcionamento do arquivo docker-compose e conexão da aplicação ao BD do docker
[] Melhoria da estrutura conforme Domain Driven Design (DDD)
[] Criação de testes
