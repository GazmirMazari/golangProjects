package main

import (
	"context"
	"log"
	"net/http"
	"os/signal"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/thedevsaddam/renderer"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var rnd *renderer.Render
var db *mgo.Database

const (
	hostName       string = "localhost:4000"
	dbName         string = "demo_todo"
	collectionName string = "todo"
	port           string = ":9000"
)

type (
	todoModel struct {
		ID        bson.ObjectId `bson:" _id, omitempty"`
		Title     string        `bson:"title"`
		Completed bool          `bson:"completed"`
		CreatedAt time.Time     `bson:"created_at"`
	}

	todo struct {
		ID        string    `json:id`
		Title     string    "json: title"
		Completed string    "json: completed"
		Createdat time.Time "json: created_at"
	}
)

func init() {
	rnd = renderer.New()
	sess, err := mgo.Dial(hostName)
	checkError(err)
	sess.SetMode(mgo.Monotonic, true)
	db = sess.DB(dbName)
}

func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}

}



func toHandler() http.Handler {
	rg := chi.NewRouter()
	rg.Group(func(r chi.Router){
		r.Get("/", fetchTodos)
		r.Post("/", createTodo)
		r.Put("/{id}", updateTodo)
		r.Delete("/{id}", deleteTodo)
	})

	return rg
}

func homdeHandler(w http.ResponseWriter,  r *http.Request){
	err := rnd.Template(w, http.StatusOK, []string{"/static/home.tpl"}, nil)
	checkError(err)
}

func fetchTodos(w http.ResponseWriter, r *http.Request){
	todos := []models.Todo{}

	if err:= db.C(collectionName).Find(bson.M{}.All(&todos)); err!= nil{
		rnd.JSON(w, http.StatusProcessing, renderer.M{
			"message": "Failed to fetch todos",
			"erro:" :err,
		})

		return
	}
	todoList : []todo{}

	for _, t :=range todos{
		todoList = append(todoList, todo{
			ID: t.ID.Hex(),
			Title: t.Title,
			Completed : t.Completed,
			CreatedAt: t.CreatedAt,
		})
	}

	rnd.JSON(w, http.StatusProcessing, renderer.M{
		"Message": "Failed to fetch todo",
		"error": err,
	})
}


func createTodo(w http.)

func main() {
	stopChan := make(chan os.Signal)
	signal.Notify(stopChan, os.Interrupt)
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Get("/", homdeHandler)
	r.Mount("/todo", todoHandler())

	srv := &http.Server{
		Addr: port,
		Handler: r,
		ReadTimeout: 60 * time.Second,
		WriteTimeout: 60 * time.Second,
		IdleTimeout: 60 * time.Second,
	}

	go func() {
		log.Println("Listening on port", port)
		if err:=srv.ListenAndServe(); err!=nil {
			log.Fatal(err)
		}
	}()

	<-stopChan

	log.Println("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(),)
	srv.Shutdown(ctx)
	defer cancel(
		log.Println("Server stopped", err)
	)

}
