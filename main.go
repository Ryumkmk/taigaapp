package main

import (
	"fmt"
	"html/template"
	"net/http"
	"os"

	"github.com/Knetic/govaluate"
)

func main() {
	//サーバ起動
	StartMainServer()
}

func generateHTML(writer http.ResponseWriter, data interface{}, filenames ...string) {
	var files []string
	for _, file := range filenames {
		files = append(files, fmt.Sprintf("views/%s.html", file))
	}

	templates := template.Must(template.ParseFiles(files...))
	templates.ExecuteTemplate(writer, "layout", data)
}

func StartMainServer() error {
	files := http.FileServer(http.Dir("views"))
	http.Handle("/static/", http.StripPrefix("/static/", files))

	http.HandleFunc("/", index)
	http.HandleFunc("/calculate", calculate)
	fmt.Println("Server started")
	port := os.Getenv("PORT")
	return http.ListenAndServe(":"+port, nil)
}

func index(w http.ResponseWriter, r *http.Request) {
	generateHTML(w, "Hello", "layout", "index")
}

func calculate(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		fmt.Println(err)
	}
	function := r.PostFormValue("txt")
	expression, err := govaluate.NewEvaluableExpression(function)
	if err != nil {
		fmt.Printf("エラー: %s\n", err)
		return
	}
	result, err := expression.Evaluate(nil)
	if err != nil {
		fmt.Printf("評価エラー: %s\n", err)
		return
	}
	rs := result.(float64)
	// fmt.Printf("結果: %.2f\n", rs)
	generateHTML(w, rs, "layout", "index")
}
