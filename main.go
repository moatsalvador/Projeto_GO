package main

import (
	"fmt"
	"os"

	//classe com o dominio
	compra "github.com/moatsalvador/Projeto_GO/domain"
)

func checkErr(err error) {
	if err != nil {
		panic(err.Error())
	}
}

func main() {
	fmt.Println("Este é o programa para leitura de dados e armazenamento de dados")

	compra.ConectarBanco()
	compra.CriartabelaRegistro()

	//realiza a leitura do arquivo
	dados := compra.LeArquivo()
	//realiza o processamento do arquivo inserindo os dados no banco
	dadosCompras := compra.ProcessarDados(dados)
	fmt.Println("Foram processados: ", len(dadosCompras), "dados")
	fmt.Println("Dados inseridos")

	compra.ValidarDadosBanco()

	//desconecta no banco
	compra.DesconectarBanco()
}

func registraArquivo(linha string) {
	//para abrir um arquivo e caso não exista cria-lo usa-se a função
	//OpenFile que vc passa como segundo parametro o que deve ser feito, como so ler, escrever, ou caso não exista criado
	//o ultimo parametro é o da permissão do arquivo
	arquivo, err := os.OpenFile("inserts.json", os.O_RDONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		fmt.Println("Ocorreu um erro:", err)
	}
	//Escreve a linha do arquivo e após pula a linha
	arquivo.WriteString(linha + "\n")
	//	fmt.Println(arquivo)
}
