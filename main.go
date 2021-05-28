package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"   
	"strings" 
)

var (
	host string
	port string
)

func init() {
	flag.StringVar(&host, "host", "localhost", "Host on which to run")
	flag.StringVar(&port, "port", "8080", "Port on which to run")
}
 
type Suppliers struct {
    Suppliers []Supplier `json:"suppliers"` 
}

type Supplier struct {
	Suppid   string  `json:"suppid"`
	Name     string  `json:"name"`
	Status   string  `json:"status"`
	Addr1    string  `json:"addr1"`
	Addr2    string  `json:"addr2"`
	City     string  `json:"city"`
	State    string  `json:"state"`
	Zip      string  `json:"zip"`
	Phone    string  `json:"phone"`
}

func doNothing(w http.ResponseWriter, r *http.Request) {}
 
func forbidden(w http.ResponseWriter, r *http.Request) {
	// see http://golang.org/pkg/net/http/#pkg-constants
	w.WriteHeader(http.StatusForbidden)
	w.Write([]byte("403 HTTP status code returned!"))
}

func find() Suppliers {
	
	suppliers := Suppliers{}
	data, err := ioutil.ReadFile("data/suppliers.json")
 
	if err != nil {
        log.Fatal(err)
    }

	fmt.Println("Successfully Opened suppliers.json")

	err = json.Unmarshal(data, &suppliers)
    if err != nil {
        log.Fatal(err)
    }

	fmt.Println(suppliers)
	
	return suppliers
}

func one(id string) Supplier {
	
	suppliers := Suppliers{}
	supplier := Supplier{}
	data, err := ioutil.ReadFile("data/suppliers.json")
 
	if err != nil {
        log.Fatal(err)
    }

	fmt.Println("Successfully Opened suppliers.json")

	err = json.Unmarshal(data, &suppliers)
    if err != nil {
        log.Fatal(err)
    } 

	for i := range suppliers.Suppliers {  
		//fmt.Println("|",suppliers.Suppliers[i].Catid , "|==|" , id,"|")
		if(string(suppliers.Suppliers[i].Suppid) == id) {
			//fmt.Println(suppliers.Suppliers[i])
			supplier = suppliers.Suppliers[i] 
		}
    }

	return supplier
}

func findAll(w http.ResponseWriter, r *http.Request) {
		
	if r.Method == "POST" {
	
	} else if r.Method == "GET" {
		fmt.Println("Endpoint Hit: findAll")

		suppliers := find() 
	
		output, err := json.Marshal(suppliers.Suppliers)
		//output, err := json.MarshalIndent(suppliers.Suppliers, "", "    ")
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		w.Header().Set("content-type", "application/json")
		w.Write(output)
		fmt.Println(string(output))

		fmt.Println("Endpoint Hit: findAll") 
	} else {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte("403 HTTP status code returned!")) 
	}
}


func findOne(w http.ResponseWriter, r *http.Request) {  
	fmt.Println("Endpoint Hit: findOne")
	id := strings.TrimPrefix(r.URL.Path, "/suppliers/") 
	suppliers := one(id)  
	
	output, err := json.Marshal(suppliers)
	//output, err := json.MarshalIndent(suppliers.Suppliers, "", "    ")
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Header().Set("content-type", "application/json")
	w.Write(output)
	fmt.Println(string(output)) 
	fmt.Println("Endpoint Hit: findOne")
}

func handleRequests() {
	http.HandleFunc("/favicon.ico", doNothing) 
	http.HandleFunc("/suppliers", findAll)
	http.HandleFunc("/suppliers/", findOne)
	address := ":" + port

	log.Println("Starting server on address", address)
	err := http.ListenAndServe(address, nil)
	if err != nil {
		panic(err)
	}
}
 
func main() {
	flag.Parse()
	handleRequests()
}