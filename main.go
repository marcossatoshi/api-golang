package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// "Pessoa" (para quem veio de OO e ficar mais fácil)
// omitempty é uma tag pra seu arquivo json não criar propriedades que não foram declaradas e/ou não possuem valor
type Pessoa struct {
	ID        string    `json:"id,omitempty"`
	Nome      string    `json:"nome,omitempty"`
	Sobrenome string    `json:"sobrenome,omitempty"`
	Endereco  *Endereco `json:"endereco,omitempty"`
}

//Outra 'classe' que pode ser acessada com * no exemplo assim
type Endereco struct {
	Cidade string `json:"cidade,omitempty"`
	Estado string `json:"estado,omitempty"`
	Bairro string `json:"bairro,omitempty"`
}

//Ao invés de criar um banco de dados usamos um array para armazenar as informações
var pessoas []Pessoa

// GetPessoas mostra todos os contatos da variável pessoas
func GetPessoas(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(pessoas)
}

// GetPessoa mostra apenas um contato
func GetPessoa(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	for _, item := range pessoas {
		if item.ID == params["id"] {
			json.NewEncoder(w).Encode(item)
			return
		}
	}
	json.NewEncoder(w).Encode(&Pessoa{})
}

// CreatePessoa cria um novo contato
func CreatePessoa(w http.ResponseWriter, r *http.Request) {
	var pessoa Pessoa
	var id int64
	_ = json.NewDecoder(r.Body).Decode(&pessoa)
	id, _ = strconv.ParseInt(pessoas[len(pessoas)-1].ID, 10, 8)
	id++
	pessoa.ID = strconv.FormatInt(id, 10)
	pessoas = append(pessoas, pessoa)
	json.NewEncoder(w).Encode(pessoas)
}

// UpdatePessoa atualiza um contato
func UpdatePessoa(w http.ResponseWriter, r *http.Request) {
	var pessoa Pessoa
	params := mux.Vars(r)
	for _, item := range pessoas {
		if item.ID == params["id"] {
			_ = json.NewDecoder(r.Body).Decode(&pessoa)
			item.Nome = pessoa.Nome
			item.Sobrenome = pessoa.Sobrenome
			item.Endereco = pessoa.Endereco
			pessoas = pessoas[:len(pessoas)-1]
			pessoas = append(pessoas, item)
		}
	}
}

// DeletePessoa deleta um contato
func DeletePessoa(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	for index, item := range pessoas {
		if item.ID == params["id"] {
			pessoas = append(pessoas[:index], pessoas[index+1:]...)
			break
		}
		json.NewEncoder(w).Encode(pessoas)
	}
}

// função principal para executar a api
func main() {
	router := mux.NewRouter()
	pessoas = append(pessoas, Pessoa{ID: "1", Nome: "Joao", Sobrenome: "Silva", Endereco: &Endereco{Cidade: "Cidade X", Estado: "Estado X", Bairro: "Bairro X"}})
	pessoas = append(pessoas, Pessoa{ID: "2", Nome: "Maria", Sobrenome: "Silva", Endereco: &Endereco{Cidade: "Cidade Z", Estado: "Estado Y", Bairro: "Bairro Y"}})
	router.HandleFunc("/contato", GetPessoas).Methods("GET")
	router.HandleFunc("/contato/{id}", GetPessoa).Methods("GET")
	router.HandleFunc("/contato", CreatePessoa).Methods("POST")
	router.HandleFunc("/contato/{id}", UpdatePessoa).Methods("PUT")
	router.HandleFunc("/contato/{id}", DeletePessoa).Methods("DELETE")
	log.Fatal(http.ListenAndServe(":8000", router))
}
