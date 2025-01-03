package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"sync"

	"github.com/jcelliott/lumber"
)

type (
	Logger interface {
		Fatal(string,...interface{})
		Error(string,...interface{})
		Warn(string,...interface{})
		Info(string,...interface{})
		Debug(string,...interface{})
		Trace(string,...interface{})
	}

	Driver struct {
		mutex sync.Mutex
		mutexes map[string]*sync.Mutex
		dir string 
		log Logger
	}
)

type Options struct {
	Logger
}

func New(dir string,options *Options)(*Driver,error) {

	dir = filepath.Clean(dir)
	opts := Options{}
	
	if options != nil {
		opts = *options
	}
	
	if opts.Logger == nil {
		opts.Logger = lumber.NewConsoleLogger(lumber.INFO)
	}
	
	driver := Driver{
		mutexes: make(map[string]*sync.Mutex),
		dir: dir,
		log: opts.Logger,
	}

	if _,err := os.Stat(dir); err == nil {
		opts.Logger.Debug("Using %s (database already exists)\n",dir)
		return &driver,nil
	}
	opts.Logger.Debug("Creating a database at the path %s\n",dir)
	return &driver,os.MkdirAll(dir,0755)
}

func (d *Driver) Write(collection string, resources string,v interface{}) error {
	if collection == "" {
		fmt.Errorf("Missing the name of the collection")
	}

	if resources == "" {
		fmt.Errorf("Missing the name of the file")
	}

	mutex := d.getOrCreateMutex(collection)
	mutex.Lock()
	defer mutex.Unlock()

	dir := filepath.Join(d.dir,collection)
	fnlPath := filepath.Join(dir,resources + ".json")
	tmpPath := fnlPath + ".tmp"

	if err := os.MkdirAll(dir,0755); err != nil {
		return err
	}

	b,err := json.MarshalIndent(v,"","\t")

	if err != nil {
		return err
	}

	b = append(b, byte('\n'))
	if err := ioutil.WriteFile(tmpPath,b,0664); err != nil {
		return err
	}

	return os.Rename(tmpPath,fnlPath)
}

func (d *Driver) Read(collection string,resource string,v interface{}) error {

	if collection == "" {
		fmt.Errorf("Incorrect name of the collection")
	}

	if resource == "" {
		fmt.Errorf("Incorrect name of the file")
	}

	record := filepath.Join(d.dir,collection,resource)

	if _,err := stat(record); err != nil {
		return err
	}

	b, err := ioutil.ReadFile(record)

	if err != nil {
		return err
	}
	
	return json.Unmarshal(b,&v)
}

func (d *Driver) ReadAll(collection string)([]string , error) {

	if collection == "" {
		return nil,fmt.Errorf("Incorrect name of the collection")
	}

	dir := filepath.Join(d.dir,collection)

	if _,err := stat(dir); err != nil {
		return nil,err
	}

	files,_ := ioutil.ReadDir(dir)
	var records []string

	for _,file := range files{
		b,err := ioutil.ReadFile(filepath.Join(dir,file.Name()))

		if err != nil {
			return nil,err
		}

		records = append(records, string(b))
	}

	return records,nil
}

func (d *Driver) Delete(collection string, resource string) error {

	if collection == "" {
		fmt.Errorf("Incorrect name of the collection")
	}

	if resource == "" {
		fmt.Errorf("Incorrect name of the file")
	}

	path := filepath.Join(collection,resource)
	mutex := d.mutexes[collection]
	mutex.Lock()
	defer mutex.Unlock()

	dir := filepath.Join(d.dir,path)

	switch fi,err := stat(dir); {
		case fi == nil, err != nil:
		return fmt.Errorf("Unable to find the path %s\n",dir)

		case fi.Mode().IsDir():
		return os.RemoveAll(dir)
		
		case fi.Mode().IsRegular():
		return os.RemoveAll(dir+".json")
	}

	return nil
}

func (d *Driver) getOrCreateMutex(collection string) *sync.Mutex {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	m,ok := d.mutexes[collection]

	if !ok {
		m = &sync.Mutex{}
		d.mutexes[collection] = m
	}

	return m
}

func stat(path string)(f os.FileInfo,err error) {

	if f,err = os.Stat(path); os.IsNotExist(err) {
		f,err = os.Stat(path + ".json")
	}
	return
}

type (
	// Adress struct {
	// 	City string
	// 	Province string
	// 	Country string
	// 	PostalCode string
	// }
	
	//  Employee struct {
	// 	Name string
	// 	Age json.Number
	// 	Contact string
	// 	Adress Adress
	// }

	Employee struct {
		Name       string `json:"name"`
		Age        string `json:"age"`
		Contact    string `json:"contact"`
		City       string `json:"city"`
		Province   string `json:"province"`
		Country    string `json:"country"`
		PostalCode string `json:"postalCode"`
	}
) 

func databaseHandler(w http.ResponseWriter, r *http.Request) {

	fmt.Println("Reached here")
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read the body", http.StatusInternalServerError)
		return
	}

	// fmt.Println(string(body))

	// var employees []Employee
	var requestData map[string]interface{}

	dir := "./"
	db,err := New(dir,nil)

	if err != nil {
		http.Error(w, "Failed to create the database", http.StatusInternalServerError)
		return
	}

	// err = json.Unmarshal(body, &employees)
	err = json.Unmarshal(body, &requestData)

	if err != nil {
		http.Error(w, "Failed to unmarshal the data", http.StatusInternalServerError)
		return
	}

	// fmt.Println("The collection name is ",requestData["collectionName"])
	// fmt.Println("The employees is ",requestData["employees"])

	collectionName := requestData["collectionName"].(string)
	employees := requestData["employees"].([]interface{})

	fmt.Println(collectionName)

	for _,value := range employees {
		emp,ok := value.(map[string]interface{})

		if !ok {
			http.Error(w, "Invalid employee data", http.StatusBadRequest)
			return
		}

		db.Write(collectionName,emp["name"].(string),
				Employee{
					Name: emp["name"].(string),	
					Age: emp["age"].(string),
					Contact: emp["contact"].(string),
					City: emp["city"].(string),
					Province: emp["province"].(string),
					Country: emp["country"].(string),
					PostalCode: emp["postalCode"].(string),
				})
		
	}

	// for _,value := range employees {
	// 		db.Write("employees",value.Name,
	// 			Employee{
	// 				Name: value.Name,
	// 				Age: value.Age,
	// 				Contact: value.Contact,
	// 				City: value.City,
	// 				Province: value.Province,
	// 				Country: value.Country,
	// 				PostalCode: value.PostalCode,
	// 			})
	// }

}

func main() {

	// for _,value := range employees {
	// 	db.Write("employees",value.Name,
	// 		Employee{
	// 			Name: value.Name,
	// 			Age: value.Age,
	// 			Contact: value.Contact,
	// 			Adress: value.Adress,
	// 		})
	// }

	// records, err := db.ReadAll("employees")
	// if err != nil {
	// 	fmt.Println("The error is ",err)
	// }

	// allEmployees := []Employee{}

	// for _,f := range records {
	// 	tmp := Employee{}

	// 	err := json.Unmarshal([]byte(f),&tmp)
	// 	if err != nil {
	// 		fmt.Println("The error is ",err)
	// 	}
	// 	allEmployees = append(allEmployees,tmp)
	// }
	// fmt.Println(allEmployees)

	fileServer := http.FileServer(http.Dir("./src"))

	http.Handle("/",fileServer)
	http.HandleFunc("/api/employees",databaseHandler)

	if err := http.ListenAndServe(":8080",nil); err != nil {
		fmt.Println("There was an error in serving the files")
	}

	// if err := db.Delete("employees","Nahan"); err != nil {
	// 	fmt.Println("The error is ", err)
	// }

	// if err := db.Delete("employees",""); err != nil {
	// 	fmt.Println("The error is ", err)
	// }

}