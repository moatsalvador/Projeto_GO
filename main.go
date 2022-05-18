package main

import (
	"fmt"

	//classe com o dominio
	compra "github.com/moatsalvador/Projeto_GO/domain"

	// driver postgres
	_ "github.com/lib/pq"
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
