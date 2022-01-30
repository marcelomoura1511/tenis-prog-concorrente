package main

import (
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"time"
	"sync"
)

var qtdPontos = 0
var players = [2]string{"Jogador 1", "Jogador 2"}
var pontuacao = [2]int{0,0} 

var wg sync.WaitGroup

var encerrar = false
var saque = true
func partida(jogadorIndex int, games chan int) {
	//esperando que as goroutines terminem
	defer wg.Done();
	for true {
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
		game := <-games
		//verificando se o jogo acabou e fechando o canal
		if pontuacao[oponenteIndex] >= qtdPontos && pontuacao[oponenteIndex]-pontuacao[jogadorIndex]>1  {
			close(games)
			encerrar=true
			break
		}
		game++
		games <- game
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

	if erro != nil {
		fmt.Println("Erro, informe um número inteiro para limite de pontos do game.")
		os.Exit(2)
	}

	//gerar números aleatórios a cada execução
	rand.Seed(time.Now().UnixNano())

	//qtd de pontos para finalizar a partida
	qtdPontos = num
	//canal do game
	games := make(chan int)

	wg.Add(2)

	go partida(0, games)
	games <- 1
	go partida(1, games)


	wg.Wait()
}

