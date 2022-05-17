package domain

/*Classe responsável pelos comandos no banco de dados*/

import (
	"database/sql"
	"fmt"
	"os"

	dbconfig "github.com/moatsalvador/Projeto_GO/dbconfig"

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

//conecta no banco de dados
func ConectarBanco() {
	fmt.Printf("Accessing %s ... ", dbconfig.DbName)

	db, err = sql.Open(dbconfig.PostgresDriver, dbconfig.DataSourceName)

	if err != nil {
		panic(err.Error())
	} else {
		fmt.Println("Connected!")
	}

	//defer db.Close()
}

//desconecta do banco de dados
func DesconectarBanco() {
	defer db.Close()
}

//cria a tabela que sera usada pela aplicação
func CriartabelaRegistro() {
	insertStmt := `CREATE TABLE IF NOT EXISTS public.registrocompra ( id SERIAL, cpf text COLLATE pg_catalog."default" NOT NULL, private integer NOT NULL DEFAULT 0, incompleto integer NOT NULL DEFAULT 0,dtultcompra text COLLATE pg_catalog."default",ticketmedio numeric DEFAULT 0,tickectultcomp numeric DEFAULT 0,lojamaisfreq text COLLATE pg_catalog."default",lojaultcompra text COLLATE pg_catalog."default",CONSTRAINT registrocompra_pkey PRIMARY KEY (id))`
	_, err := db.Exec(insertStmt)
	checkErr(err)
	if err != nil {
		fmt.Println("Ocorreu um erro:", err)
		os.Exit(1)
	}
}

//insere os dados no banco
func InserirDadosBancoCompra(cpf string, private int, incompleto int, dataCompra string, tiketmedio float64, ticketUltcomp float64, lojmaisfreq string, lojaultcomp string) {
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

//Inserir dados no banco com os dados do tipo Compra
func InserirDadosDeCompra(dados Compra) {
	sqlStatement := fmt.Sprintf("INSERT INTO %s (cpf, private, incompleto, dtultcompra, ticketmedio, tickectultcomp, lojamaisfreq, lojaultcompra) VALUES ($1,$2, $3, $4, $5, $6, $7, $8)", dbconfig.TableName)

	//inserir no banco
	insert, err := db.Prepare(sqlStatement)
	checkErr(err)

	result, err := insert.Exec(dados.CPF, dados.Private, dados.Incompleto, dados.DtUltCompra, dados.TicketMedio, dados.TicketUltComp, dados.LojaMaisFreq, dados.LojaUltComp)
	checkErr(err)

	affect, err := result.RowsAffected()
	checkErr(err)

	//	println("inserido :", dados.CPF, affect)

	if affect > 1 {
		println("Inseriu mais de um registro")
	}
}

//Seleciona os dados do banco para fazer a validação de documentos
func SqlSelect() map[int]Compra {
	dadosCompras := make(map[int]Compra)

	sqlStatement, err := db.Query("SELECT cpf,dtultcompra,lojamaisfreq, lojaultcompra FROM " + dbconfig.TableName)
	checkErr(err)

	indice := 0
	for sqlStatement.Next() {

		var compra Compra

		err = sqlStatement.Scan(&compra.CPF, &compra.DtUltCompra, &compra.LojaMaisFreq, &compra.LojaUltComp)
		checkErr(err)

		dadosCompras[indice+1] = compra
		indice += 1
		//	fmt.Printf("\t%s\t%s\t%s \n", compra.CPF, compra.DtUltCompra, compra.LojaMaisFreq)
	}
	return dadosCompras
}
