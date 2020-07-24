package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

type Product struct {
	Id      int    `json: "id"`
	Message string `json: "message"`
	Age     int    `json: "age"`
	Name    string `json: "name"`
	surname string
}

var productList = []Product{}
var nextId = 1

func productsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		productsJson, err := json.Marshal(productList)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(productsJson)
	case http.MethodPost:
		var newProduct Product
		bodyBytes, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
		}
		err = json.Unmarshal(bodyBytes, &newProduct)
		newProduct.Id = nextId
		nextId++
		productList = append(productList, newProduct)
		w.WriteHeader(http.StatusCreated)
	}
}

func main() {
	http.HandleFunc("/products", productsHandler)
	http.ListenAndServe(":5000", nil)
}
