package domain

type Compra struct {
	CPF           string
	Private       int
	Incompleto    int
	DtUltCompra   string
	TicketMedio   float64
	TicketUltComp float64
	LojaMaisFreq  string
	LojaUltComp   string
}

func NewCompra(cpf string, private int, incompleto int, dataCompra string, tiketmedio float64, ticketUltcomp float64, lojmaisfreq string, lojaultcomp string) *Compra {
	return &Compra{CPF: cpf, Private: private, Incompleto: incompleto, DtUltCompra: dataCompra, TicketMedio: tiketmedio, TicketUltComp: ticketUltcomp, LojaMaisFreq: lojmaisfreq, LojaUltComp: lojaultcomp}
}
