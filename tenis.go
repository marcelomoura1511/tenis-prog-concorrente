package main

import (
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"time"
	"sync"
)

var qtdPontosGame = 0
var qtdGames = 0
var qtdSets = 0
var players = [2]string{"Jogador 1", "Jogador 2"}
var pontuacao = [2]int{0,0} 
var games = [2]int{0,0} 
var sets = [2]int{0,0} 

var wg sync.WaitGroup

var encerrar = false
var saque = true
func partida(jogadorIndex int, passos chan int) {
	//esperando que as goroutines terminem
	defer wg.Done();
	for true {
		fmt.Print()
		//verificando se a partida deve ser encerrada
		if encerrar {
			fmt.Printf("%s venceu a partida.\n", players[jogadorIndex])
			break
		}
		//gerando o número aleatório
		num := rand.Intn(60)
		//pegando o index do openente
		oponenteIndex := oponente(jogadorIndex);
		//verificando se o número é primo, caso sim, o jogador errou
		if isPrimo(num) {
			//verificando se errou o saque ou a rebatida
			if saque {
				fmt.Printf("%s errou o saque.\n", players[jogadorIndex])
			}else{
				fmt.Printf("%s falhou.\n", players[jogadorIndex])
				saque=true
			}
			//printando o placar
			fmt.Printf("%s pontuou.\n", players[oponenteIndex])
			pontuacao[oponenteIndex]++
			fmt.Print("\n******** Placar do Game ********\n")
			fmt.Printf("%s | %d x %d | %s \n\n", players[0], pontuacao[0], pontuacao[1], players[1])
		} else {
			if saque{
				fmt.Printf("%s sacou.\n", players[jogadorIndex])
				saque=false
			}else{
				fmt.Printf("%s rebateu.\n", players[jogadorIndex])
			}
		}
		passo := <-passos
		if pontuacao[oponenteIndex] >= qtdPontosGame && pontuacao[oponenteIndex]-pontuacao[jogadorIndex]>1  {
			games[oponenteIndex]++
			gamesjogados := games[jogadorIndex]+games[oponenteIndex]
			fmt.Printf("%s ganhou o %dº game.\n", players[oponenteIndex], gamesjogados)
			fmt.Print("\n******** Placar de Games ********\n")
			fmt.Printf("%s | %d x %d | %s \n\n", players[0], games[0], games[1], players[1])
			pontuacao = [2]int{0, 0}
		}
		if games[oponenteIndex] >= qtdGames && games[oponenteIndex]-games[jogadorIndex]>1  {
			sets[oponenteIndex]++
			setjogados:=sets[jogadorIndex]+sets[oponenteIndex]
			fmt.Printf("%s ganhou o %dº set.\n", players[oponenteIndex], setjogados)
			fmt.Print("\n******** Placar de Sets ********\n")
			fmt.Printf("%s | %d x %d | %s \n\n", players[0], sets[0], sets[1], players[1])
			games = [2]int{0, 0}
		}
		if sets[oponenteIndex]==qtdSets{
			close(passos)
			encerrar=true
			break
		}
		passo++
		passos <- passo
	}
	// wg.Done()
}
//func de pegar o index do oponente
func oponente(jogadorIndex int) int{
	if jogadorIndex == 0{
		return 1
	}
	return 0
}

func isPrimo(num int) bool{
	if num > 1 {
		for i := 2; i < (num/2); i++ {
		   if num % i == 0 {
			   return false
		   }
		}
		return true
	}
	return false
}

func main() {
	args := os.Args[1:]

	if len(args) < 1 {
		fmt.Println("Erro, digite um número para limite de pontos do game")
		os.Exit(2)
	}

	num, erro:= strconv.Atoi(os.Args[1])
	num1, erro1:= strconv.Atoi(os.Args[2])
	num2, erro2:= strconv.Atoi(os.Args[3])

	if erro != nil {
		fmt.Println("Erro, informe um número inteiro para limite de pontos do game.")
		os.Exit(2)
	}
	if num < 2{
		fmt.Println("Erro, informe um número inteiro maior ou igual a 2 para limite de pontos.")
		os.Exit(2)
	}
	if erro1 != nil {
		fmt.Println("Erro, informe um número inteiro para limite de games.")
		os.Exit(2)
	}

	if num1 < 1{
		fmt.Println("Erro, informe um número inteiro maior ou igual a 1 para limite de games.")
		os.Exit(2)
	}
	if erro2 != nil {
		fmt.Println("Erro, informe um número inteiro para limite de sets.")
		os.Exit(2)
	}

	if num2 < 1{
		fmt.Println("Erro, informe um número inteiro maior ou igual a 1 para limite de sets.")
		os.Exit(2)
	}

	//gerar números aleatórios a cada execução
	rand.Seed(time.Now().UnixNano())

	//qtd de pontos para finalizar a partida
	qtdPontosGame = num
	qtdGames = num1
	qtdSets = num2
	fmt.Printf("%d, %d, %d", qtdPontosGame, qtdGames, qtdSets)
	//canal do game
	passos := make(chan int)

	wg.Add(2)

	go partida(0, passos)
	passos <- 1
	go partida(1, passos)


	wg.Wait()
}

