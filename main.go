package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

//Cria as estruturas que serão utilizadas nas funções
type Person struct {
	ID        string   `json:"id,omitempty"`
	Firstname string   `json:"firstname,omitempty"`
	Lastname  string   `json:"lastname,omitempty"`
	Address   *Address `json:"address,omitempty"`
}

type Address struct {
	City  string `json:"city,omitempty"`
	State string `json:"state,omitempty"`
}

//cria um slice de Person
var people []Person

//Cria as funções que serão chamadas pelas requisições

func GetPeople(w http.ResponseWriter, r *http.Request) {
	//devolve na resposta do get o slice de people
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(people)
	log.Println("Devolvendo consulta de todos os usuarios")
}

func GetPerson(w http.ResponseWriter, r *http.Request) {
	//verifica a informação passada através da requisição e verifica se a informação corresponde ao ID de algum Person então passa isso na resposta
	params := mux.Vars(r)
	for _, item := range people {
		if item.ID == params["id"] {
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(item)
			log.Println("Devolvendo do usuario :", item.ID, "-", item.Firstname)
			return
		} else if item.ID != params["id"] {
			w.WriteHeader(http.StatusNoContent)
			json.NewEncoder(w).Encode("error:person not found")
			log.Println("usuario não encontrado")
			return
		}
	}
}
func CreatePerson(w http.ResponseWriter, r *http.Request) {
	//verifica a informação passada na requisição,converte para a variavel person, adiciona ao slice e sem seguida retorna o slice na resposta
	param := mux.Vars(r)
	var person Person
	_ = json.NewDecoder(r.Body).Decode(&person)
	person.ID = param["id"]
	people = append(people, person)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(people)
	log.Println("Cadastrando usuario: ", person.ID, "-", person.Firstname)
}
func DeletePerson(w http.ResponseWriter, r *http.Request) {
	//deleta uma people do slice através do id e em seguida retorna o slice no response
	params := mux.Vars(r)
	for index, item := range people {
		if item.ID == params["id"] {
			//valido se o id passado correspode a algume dentro meu slice, então realizo um append no slice adicionando todos os elementos exceto o atual, deletando o mesmo do slice
			people = append(people[:index], people[index+1:]...)
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(people)
			log.Println("Removendo Usuario: ", item.ID, "-", item.Firstname)
			break
		} else {
			w.WriteHeader(http.StatusNoContent)
			json.NewEncoder(w).Encode("error:person not found")
			log.Println("usuario não encontrado")
		}

	}
}
func main() {

	//cria um router para direcionar as requisições
	router := mux.NewRouter()
	log.Println("Servidor Rodando na porta 8080")

	//popula o slice de Person
	people = append(people, Person{ID: "1", Firstname: "Mauricio", Lastname: "Fernandes", Address: &Address{City: "São Paulo", State: "SP"}})
	people = append(people, Person{ID: "2", Firstname: "Marcos", Lastname: "Santos", Address: &Address{City: "Parana", State: "PR"}})
	people = append(people, Person{ID: "3", Firstname: "Maria", Lastname: "Souza", Address: &Address{City: "São Bernando", State: "SP"}})

	//cria as rotas e as funçoes que elas irão chamar
	router.HandleFunc("/people", GetPeople).Methods("GET")
	router.HandleFunc("/people/{id}", GetPerson).Methods("GET")
	router.HandleFunc("/people/{id}", CreatePerson).Methods("POST")
	router.HandleFunc("/people/{id}", DeletePerson).Methods("DELETE	")
	log.Fatal(http.ListenAndServe(":8080", router))

}
