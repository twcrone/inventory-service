package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
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

func productHandler(w http.ResponseWriter, r *http.Request) {
	urlPathSegments := strings.Split(r.URL.Path, "product/")
	productId, err := strconv.Atoi(urlPathSegments[len(urlPathSegments)-1])
	fmt.Println("Looking for ", productId)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	product, listItemIndex := findProductById(productId)
	if product == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	switch r.Method {
	case http.MethodGet:
		productJson, err := json.Marshal(product)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(productJson)
	case http.MethodPut:
		var updatedProduct Product
		bodyBytes, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		err = json.Unmarshal(bodyBytes, &updatedProduct)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		product = &updatedProduct
		productList[listItemIndex] = *product
		w.WriteHeader(http.StatusOK)
	}
}

func findProductById(id int) (*Product, int) {
	var found *Product = nil
	for i := range productList {
		product := productList[i]
		if product.Id == id {
			found = &product
			break
		}
	}
	return found, 0
}

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
	http.HandleFunc("/product/", productHandler)
	http.HandleFunc("/products", productsHandler)
	http.ListenAndServe(":5000", nil)
}
