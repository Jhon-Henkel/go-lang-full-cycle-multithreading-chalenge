package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

const baseUriBrasilApi = "https://brasilapi.com.br/api/cep/v1/"
const baseUriViaCep = "https://viacep.com.br/ws/"

type AddressData struct {
	Cep     string `json:"cep"`
	City    string `json:"city"`
	State   string `json:"state"`
	Service string
}

type AddressDataViaCep struct {
	Cep   string `json:"cep"`
	City  string `json:"localidade"`
	State string `json:"uf"`
}

func main() {
	cep := os.Args[1]
	if len(cep) != 8 {
		fmt.Println("CEP inválido, um cep válido deve ser no formato 12345678")
		os.Exit(1)
	}
	chanelBrasilApi := make(chan AddressData)
	chanelViaCep := make(chan AddressData)

	go getBrasilApiAddressInfo(cep, chanelBrasilApi)
	go getViaCepAddressInfo(cep, chanelViaCep)

	select {
	case brasilApi := <-chanelBrasilApi:
		pintData(brasilApi)
	case viaCep := <-chanelViaCep:
		pintData(viaCep)
	case <-time.After(time.Second):
		fmt.Println("Timeout")
	}
}

func getBrasilApiAddressInfo(cep string, chanel chan AddressData) {
	req, _ := http.Get(baseUriBrasilApi + cep)
	defer req.Body.Close()
	res, _ := io.ReadAll(req.Body)
	data := AddressData{}
	json.Unmarshal(res, &data)
	data.Service = "BrasilApi"
	chanel <- data
}

func getViaCepAddressInfo(cep string, chanel chan AddressData) {
	req, _ := http.Get(baseUriViaCep + cep + "/json")
	defer req.Body.Close()
	res, _ := io.ReadAll(req.Body)
	data := AddressDataViaCep{}
	json.Unmarshal(res, &data)
	chanel <- populateViaCepAddressData(data)
}

func populateViaCepAddressData(viaCepData AddressDataViaCep) AddressData {
	data := AddressData{}
	data.Cep = viaCepData.Cep
	data.City = viaCepData.City
	data.State = viaCepData.State
	data.Service = "ViaCep"
	return data
}

func pintData(data AddressData) {
	fmt.Println("CEP:     ", data.Cep)
	fmt.Println("Cidade:  ", data.City)
	fmt.Println("Estado:  ", data.State)
	fmt.Println("Serviço: ", data.Service)
}
