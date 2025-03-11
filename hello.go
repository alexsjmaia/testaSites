package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

const MONITORAMENTOS = 3
const DELAY = 5

func main() {
	exibeIntroducao()

	for {
		exibeMenu()
		comando := leComando()

		switch comando {
		case 1:
			iniciarMonitoramento()
		case 2:
			fmt.Println("Exibindo Logs...")
			imprimeLogs()
		case 0:
			fmt.Println("Saindo do Programa")
			os.Exit(0)
		default:
			fmt.Println("Não conheço este comando")
			os.Exit(-1)
		}
	}
}

func exibeIntroducao() {
	nome := "Alex"
	versao := 1.2

	fmt.Println("Olá Mundo, Sr.", nome)
	fmt.Println("Este programa está na versão", versao)
}

func exibeMenu() {
	fmt.Println("1 - Iniciar Monitoramento")
	fmt.Println("2 - Exibir Logs")
	fmt.Println("0 - Sair do programa")
	fmt.Printf("Rotina:")
}

func leComando() int {
	var comandoLido int
	fmt.Scan(&comandoLido)
	fmt.Println("O comando escolhido foi", comandoLido)
	fmt.Println("")

	return comandoLido
}

func testaSite(site string) {
	resp, err := http.Get(site)

	if err != nil {
		log.Fatal("Ocorreu um erro:", err)
	}

	if resp.StatusCode == 200 {
		fmt.Println("Site:", site, "Foi carregado com sucesso", resp.StatusCode)
		registraLog(site, true)
	} else {
		fmt.Println("Site:", site, "esta com problemas. Status Code:", resp.StatusCode)
		registraLog(site, false)
	}
}

func iniciarMonitoramento() {
	fmt.Println("Iniciando o Monitoramento...")

	sites := leSitesDoArquivo()

	fmt.Println(sites)

	for i := 0; i < MONITORAMENTOS; i++ {
		for i, site := range sites {
			fmt.Println("Testando Site", i, ":", site)
			testaSite(site)
		}
		time.Sleep(DELAY * time.Second)
		fmt.Println("")
	}
	fmt.Println("")
}

func leSitesDoArquivo() []string {
	var sites []string

	arquivo, err := os.Open("sites.txt")
	if err != nil {
		log.Fatal("Erro ao abrir o arquivo .txt", err)
	}
	defer arquivo.Close() // fechar o arquivo

	leitor := bufio.NewReader(arquivo)

	for {
		linha, err := leitor.ReadString('\n')
		linha = strings.TrimSpace(linha)

		// Adicionando site a site no slice sites
		sites = append(sites, linha)

		if err == io.EOF {
			break
		}
	}

	fmt.Println(sites)

	return sites
}

func registraLog(site string, status bool) {

	arquivo, err := os.OpenFile("log.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		fmt.Println(err)
	}
	defer arquivo.Close()

	arquivo.WriteString(time.Now().Format("02/01/2006 15:04:05") + " - " + site + "- online: " + strconv.FormatBool(status) + "\n")
}

func imprimeLogs() {

	// ioutil já fecha o arquivo, sem necessidade de escrever arquivo.Close()
	arquivo, err := ioutil.ReadFile("log.txt")
	if err != nil {
		fmt.Println("Erro ao abrir o arquivo de LOG", err)
	}
	fmt.Println(string(arquivo))
}
