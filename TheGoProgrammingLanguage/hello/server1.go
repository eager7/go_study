package main
import (
	"fmt"
	"log"
	"net/http"
)

func main(){
	fmt.Println("server1 test")

	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe("localhost:8000", nil))
}

func handler(w http.ResponseWriter, r *http.Request){
	fmt.Fprintf(w, "Request is %q\n", r.URL.Path)
}