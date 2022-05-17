package main

import (
	"bufio"
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"regexp"
	"strconv"
	"strings"
	"unicode"

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

	criartabela()
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

	sqlSelect()

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

//cria a tabela que sera usada pela aplicação
func criartabela() {

	insertStmt := `CREATE TABLE IF NOT EXISTS public.registrocompra ( id SERIAL, cpf text COLLATE pg_catalog."default" NOT NULL, private integer NOT NULL DEFAULT 0, incompleto integer NOT NULL DEFAULT 0,dtultcompra text COLLATE pg_catalog."default",ticketmedio numeric DEFAULT 0,tickectultcomp numeric DEFAULT 0,lojamaisfreq text COLLATE pg_catalog."default",lojaultcompra text COLLATE pg_catalog."default",CONSTRAINT registrocompra_pkey PRIMARY KEY (id))`
	_, err := db.Exec(insertStmt)
	checkErr(err)
	if err != nil {
		fmt.Println("Ocorreu um erro:", err)
		os.Exit(1)
	}
}

func sqlSelect() {

	sqlStatement, err := db.Query("SELECT cpf,dtultcompra,lojamaisfreq FROM " + dbconfig.TableName)
	checkErr(err)

	for sqlStatement.Next() {

		var compra compra.Compra

		err = sqlStatement.Scan(&compra.CPF, &compra.DtUltCompra, &compra.LojaMaisFreq)
		checkErr(err)

		if IsCPF(compra.CPF) {
			fmt.Print(compra.CPF, " Valido -- ")
		} else {
			fmt.Print(compra.CPF, " Invalido -- ")
		}

		fmt.Printf("\t%s\t%s\t%s \n", compra.CPF, compra.DtUltCompra, compra.LojaMaisFreq)
	}
}

//validação CPF - pacote BR DOC -  https://pkg.go.dev/github.com/paemuri/brdoc/v2

// Regexp pattern for CPF and CNPJ.
var (
	CPFRegexp  = regexp.MustCompile(`^\d{3}\.?\d{3}\.?\d{3}-?\d{2}$`)
	CNPJRegexp = regexp.MustCompile(`^\d{2}\.?\d{3}\.?\d{3}\/?(:?\d{3}[1-9]|\d{2}[1-9]\d|\d[1-9]\d{2}|[1-9]\d{3})-?\d{2}$`)
)

// IsCPF verifies if the given string is a valid CPF document.
func IsCPF(doc string) bool {

	const (
		size = 9
		pos  = 10
	)

	return isCPFOrCNPJ(doc, CPFRegexp, size, pos)
}

// IsCNPJ verifies if the given string is a valid CNPJ document.
func IsCNPJ(doc string) bool {

	const (
		size = 12
		pos  = 5
	)

	return isCPFOrCNPJ(doc, CNPJRegexp, size, pos)
}

// isCPFOrCNPJ generates the digits for a given CPF or CNPJ and compares it with the original digits.
func isCPFOrCNPJ(doc string, pattern *regexp.Regexp, size int, position int) bool {

	if !pattern.MatchString(doc) {
		return false
	}

	cleanNonDigits(&doc)

	// Invalidates documents with all digits equal.
	if allEq(doc) {
		return false
	}

	d := doc[:size]
	digit := calculateDigit(d, position)

	d = d + digit
	digit = calculateDigit(d, position+1)

	return doc == d+digit
}

// cleanNonDigits removes every rune that is not a digit.
func cleanNonDigits(doc *string) {

	//buf := bytes.NewBufferString("")
	buf := bytes.NewBufferString("")
	for _, r := range *doc {
		if unicode.IsDigit(r) {
			buf.WriteRune(r)
		}
	}

	*doc = buf.String()
}

// allEq checks if every rune in a given string is equal.
func allEq(doc string) bool {

	base := doc[0]
	for i := 1; i < len(doc); i++ {
		if base != doc[i] {
			return false
		}
	}

	return true
}

// calculateDigit calculates the next digit for the given document.
func calculateDigit(doc string, position int) string {

	var sum int
	for _, r := range doc {

		sum += toInt(r) * position
		position--

		if position < 2 {
			position = 9
		}
	}

	sum %= 11
	if sum < 2 {
		return "0"
	}

	return strconv.Itoa(11 - sum)
}

// toInt converts a rune to an int.
func toInt(r rune) int {
	return int(r - '0')
}
