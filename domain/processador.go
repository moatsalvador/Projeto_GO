package domain

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	valid "github.com/moatsalvador/Projeto_GO/brdoc"
)

//Função para leitura de arquivo e inserir em um array de string com as linhas
func LeArquivo() []string {
	var compras []string
	//para abrir um arquivo usa-se o metodo Open do pacote OS
	arquivo, err := os.Open("Teste_Base.txt")
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

//função para processar os dados e inseri-los no banco
func ProcessarDados(dados []string) map[int]Compra {
	//cria um map de compras e poupla com os dados lido do arquivo
	dadosCompras := make(map[int]Compra)
	fmt.Println("Iniciado a inserção dos dados")
	for i, dado := range dados {
		cpf := strings.TrimSpace(dado[0:19])
		private, _ := strconv.Atoi(dado[20:21])
		incompleto, _ := strconv.Atoi(dado[32:33])
		dataCompra := strings.TrimSpace(dado[43:54])
		tiketmedio := converterValor(dado[65:87])
		ticketUltcomp := converterValor(dado[87:111])
		lojaultcomp := strings.TrimSpace(dado[111:130])
		lojmaisfreq := strings.TrimSpace(dado[131:])
		//dadosCompras[i+1] = Compra{CPF: cpf, Private: private, Incompleto: incompleto, DtUltCompra: dataCompra, TicketMedio: tiketmedio, TicketUltComp: ticketUltcomp, LojaMaisFreq: lojmaisfreq, LojaUltComp: lojaultcomp}
		compra := Compra{CPF: cpf, Private: private, Incompleto: incompleto, DtUltCompra: dataCompra, TicketMedio: tiketmedio, TicketUltComp: ticketUltcomp, LojaMaisFreq: lojmaisfreq, LojaUltComp: lojaultcomp}
		dadosCompras[i+1] = compra
		//fmt.Println("Compra é do tipo ", reflect.TypeOf(compra))
		InserirDadosDeCompra(compra)
		//InserirDadosBancoCompra(cpf, private, incompleto, dataCompra, tiketmedio, ticketUltcomp, lojmaisfreq, lojaultcomp)
	}
	return dadosCompras
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

//Valida os dados de CPF e CNPJ do banco de dados e insere em um arquivo txt
func ValidarDadosBanco() {
	dadosCompras := make(map[int]Compra)
	dadosCompras = SqlSelect()
	invalido := ""
	for _, dado := range dadosCompras {

		//valida o CPF
		if !valid.IsCPF(dado.CPF) {
			invalido = "CPF: " + dado.CPF + ";"
		}

		//valida os CNPJs da loja da loja mais frequente
		if !valid.IsCNPJ(dado.LojaMaisFreq) && dado.LojaMaisFreq != "NULL" {
			invalido = "CNPJ: " + invalido + dado.LojaMaisFreq + ";"
		}

		//valida os CNPJs da loja da ultima compra
		if !valid.IsCNPJ(dado.LojaUltComp) && dado.LojaUltComp != "NULL" {
			invalido = "CNPJ: " + invalido + dado.LojaUltComp + ";"
		}

		if invalido != "" {
			registrarDocumentosInconsistentes(invalido + "\n")
			invalido = ""
		}
	}
}

func registrarDocumentosInconsistentes(linha string) {
	//para abrir um arquivo e caso não exista cria-lo usa-se a função
	//OpenFile que vc passa como segundo parametro o que deve ser feito, como so ler, escrever, ou caso não exista criado
	//o ultimo parametro é o da permissão do arquivo
	arquivo, err := os.OpenFile("docs_inconsistentes.txt", os.O_RDONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		fmt.Println("Ocorreu um erro:", err)
	}
	//Escreve a linha do arquivo e após pula a linha
	arquivo.WriteString(linha)
}
