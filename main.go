package main

import (
	"bufio"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	//configuração do banco de dados
	"./dbconfig"
	//classe com o dominio
	compra "./domain"

	// driver postgres
	_ "github.com/lib/pq"
)

var db *sql.DB
var err error

func checkErr(err error) {
	if err != nil {
		panic(err.Error())
	}
}

func main() {
	fmt.Println("Este é o programa para leitura de dados e armazenamento de dados")
	fmt.Printf("Accessing %s ... ", dbconfig.DbName)

	db, err = sql.Open(dbconfig.PostgresDriver, dbconfig.DataSourceName)

	if err != nil {
		panic(err.Error())
	} else {
		fmt.Println("Connected!")
	}

	defer db.Close()

	//realiza a leitura do arquivo seguindo o layout exemplo
	dados := leSitesdoArquivo()
	dadosCompras := processaDados(dados)
	fmt.Println("Foram processados: ", len(dadosCompras), "dados")
	fmt.Println("Dados inseridos")

	//transforma os dados do map em um json
	arqJson, err := json.Marshal(dadosCompras)
	if err != nil {
		fmt.Println("Ocorreu um erro: ", err)
	} else {
		//fmt.Println(string(arqJson))
		registraArquivo(string(arqJson))
	}
}

//converte o valor recebido com virugula para um float64
func converterValor(valor string) float64 {
	valorconvert := strings.Replace(strings.TrimSpace(valor), ",", ".", -1)
	//fmt.Println("Valor original: ", valor, "Valor ajustado: ", valorconvert)
	if valorconvert == "NULL" {
		return 0
	} else {
		valorF, err := strconv.ParseFloat(valorconvert, 64)
		if err != nil {
			return 0
		} else {
			return valorF
		}
	}
}

//Função para leitura de arquivo e inserir em um array de string com as linhas
func leSitesdoArquivo() []string {
	var compras []string
	//para abrir um arquivo usa-se o metodo Open do pacote OS
	arquivo, err := os.Open("Teste_BaseMenor.txt")
	if err != nil {
		fmt.Println("Ocorreu um erro:", err)
		os.Exit(1)
	}
	//o pacote bufio tem a função newReader que vai retonar um leitor do arquivo
	leitor := bufio.NewReader(arquivo)

	for {
		//faz a leitura até a quebra de linha
		linha, err := leitor.ReadString('\n')
		//remover o espaço final
		linha = strings.TrimSpace(linha)
		//ignora a primeira linha de cabeçalho
		if linha[0] == 'C' {
			fmt.Println("Linha Cabeçalho")
		} else {
			compras = append(compras, linha)
		}
		//se encontrar o final do arquivo encerra o loop
		if err == io.EOF {
			break
		}
	}
	//fechar o arquivo
	arquivo.Close()
	//testa se retornou erro e retorna a mensagem
	if err != nil {
		fmt.Println("Ocorreu um erro:", err)
		os.Exit(1)
	}
	return compras
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

//função para processar os dados e inseri-los no banco
func processaDados(dados []string) map[int]compra.Compra {
	//cria um map de compras e poupla com os dados lido do arquivo
	dadosCompras := make(map[int]compra.Compra)
	fmt.Println("Iniciado a inserção dos dados")
	for i, dado := range dados {
		cpf := dado[0:14]
		private, _ := strconv.Atoi(dado[20:21])
		incompleto, _ := strconv.Atoi(dado[32:33])
		dataCompra := strings.TrimSpace(dado[43:54])
		tiketmedio := converterValor(dado[65:87])
		ticketUltcomp := converterValor(dado[87:111])
		lojaultcomp := strings.TrimSpace(dado[111:130])
		lojmaisfreq := strings.TrimSpace(dado[131:])
		//	compra := compra.Compra{cpf, private, incompleto, dataCompra, tiketmedio, ticketUltcomp, lojmaisfreq, lojaultcomp}
		//		dadosCompras[i+1] = compra
		dadosCompras[i+1] = compra.Compra{CPF: cpf, Private: private, Incompleto: incompleto, DtUltCompra: dataCompra, TicketMedio: tiketmedio, TicketUltComp: ticketUltcomp, LojaMaisFreq: lojmaisfreq, LojaUltComp: lojaultcomp}
		//fmt.Println("Compra é do tipo ", reflect.TypeOf(compra))
		inserirDadosBanco(cpf, private, incompleto, dataCompra, tiketmedio, ticketUltcomp, lojmaisfreq, lojaultcomp)
	}
	return dadosCompras
}

//insere os dados no banco
func inserirDadosBanco(cpf string, private int, incompleto int, dataCompra string, tiketmedio float64, ticketUltcomp float64, lojmaisfreq string, lojaultcomp string) {
	sqlStatement := fmt.Sprintf("INSERT INTO %s (cpf, private, incompleto, dtultcompra, ticketmedio, tickectultcomp, lojamaisfreq, lojaultcompra) VALUES ($1,$2, $3, $4, $5, $6, $7, $8)", dbconfig.TableName)

	//inserir no banco
	insert, err := db.Prepare(sqlStatement)
	checkErr(err)

	//result, err := insert.Exec(dados.CPF, dados.Private, dados.Incompleto, dados.DtUltCompra, dados.TicketMedio, dados.TicketUltComp, dados.LojaMaisFreq, dados.LojaUltComp)
	result, err := insert.Exec(cpf, private, incompleto, dataCompra, tiketmedio, ticketUltcomp, lojmaisfreq, lojaultcomp)
	checkErr(err)

	affect, err := result.RowsAffected()
	checkErr(err)

	if affect > 1 {
		println("Inseriu mais de um registro")
	}
}
