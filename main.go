package main

import (
	"encoding/json"
	"fmt"
	"github.com/twcrone/inventoryservice/database"
	"github.com/twcrone/inventoryservice/service"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

type Note struct {
	Id   int    `json: "id"`
	Note string `json: "message"`
}

var productList []Note
var nextId = 1

func init() {
	database.SetupDatabase()
	var data = `[
    {
        "Id": 1,
        "Message": "Hello",
        "Age": 50,
        "Name": "Todd"
    },
    {
        "Id": 2,
        "Message": "Second one",
        "Age": 22,
        "Name": "Number Two"
    }
]`

	err := json.Unmarshal([]byte(data), &productList)
	if err != nil {
		log.Fatal(err)
	}
}

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
		var updatedProduct Note
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
		fmt.Println("Updating product at index ", listItemIndex)
		productList[listItemIndex] = *product
		w.WriteHeader(http.StatusOK)
	}
}

func findProductById(id int) (*Note, int) {
	var found *Note = nil
	var index = -1
	for i := range productList {
		product := productList[i]
		if product.Id == id {
			found = &product
			index = i
			break
		}
	}
	return found, index
}

func notesHandler(w http.ResponseWriter, r *http.Request) {
	//	switch r.Method {
	//	case http.MethodGet:
	notes := service.GetAllNotes()
	notesJson, err := json.Marshal(notes)
	if err != nil {
		log.Fatal(err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(notesJson)
	//	}
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
		var newProduct Note
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
	http.HandleFunc("/notes", notesHandler)
	http.HandleFunc("/product/", productHandler)
	http.HandleFunc("/products", productsHandler)
	http.ListenAndServe(":5000", nil)
}
