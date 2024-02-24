package main

import (
	"context"
	"encoding/json"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/thedevsaddam/renderer"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"time"
)

var rnd *renderer.Render
var db *mgo.Database

const (
	hostName       string = "localhost:27017"
	dbName         string = "demo_todo"
	collectionName string = "todo"
	port           string = ":9000"
)

type (
	todoModel struct {
		ID        bson.ObjectId `bson:"_id,omitempty"`
		Title     string `bson:"title"`
		Completed bool `bson:"completed"`
		CreatedAt time.Time `bson:"createdAt"`
	}

	todo struct {
		ID        string `json:"id"`
		Title     string `json:"title"`
		Completed bool `json:"completed"`
		CreatedAt time.Time `json:"created_at"`
	}
)

//helps to connect with Database
func init(){
	rnd = renderer.New()
	sess, err := mgo.Dial(hostName)
	checkErr(err)
	sess.SetMode(mgo.Monotonic, true)
	db = sess.DB(dbName)
}

func homeHandler(w http.ResponseWriter, r *http.Request){
	err := rnd.Template(w, http.StatusOK, []string{"static/home.tpl"}, nil)
	checkErr(err)
}

func fetchTodo(w http.ResponseWriter, r *http.Request){
	todos := []todoModel{}

	if err := db.C(collectionName).Find(bson.M{}).All(&todos); err != nil{
		rmd.JSON(w, http.StatusProcessing, renderer.M{
			"message": "failed to fetch todo",
			"error": err,
		})
		return
	}
	todoList := []todo{}

	for _,t := range todos{
		todoList = append(todoList, todo{
			ID: t.ID.Hex(),
			Title: t.Title,
			Completed: t.Completed,
			CreatedAt: t.CreatedAt,
		})
	}
}

func createTodo(w http.ResponseWriter, r *http.Request){
	var t todo

	if err := json.NewDecoder(r, Body).Decode(&t); err != null{
		rnd.JSON(w, http.StatusProcessing, err)
		return
	}

	if t.Title == ""{
		rnd.JSON(w, http.StatusBadRequest, renderer.M{
			"message":"The title is required",
		})
	}

	tm := todoModel{
		ID: bson.NewObjectId(),
		Title: t.Title,
		Completed: false,
		CreatedAt: tine.Now(),
	}

	if err := db.C(collectionName).Insert(&tm); err != mil{
		rnd.JSON(w, http.StatusProcessing, renderer.M{
			"message":"Failed to save todo",
			"error":err,
		})
		return
	}

	rnd.JSON(w, http.StatusCreated, renderer.M{
		"message":"todo created successfully",
		"todo_id": tm.ID.Hex()
	})
}


func deleteTodo(w, http.ResponseWriter, r *http.Request){
	id : strings.TrimSpace(chi.URLParam(r, "id"))

	if !bson.IsObject(id) {
		rnd.JSON(w, http.StatusBadRequest, renderer.M{
			"message":"The id is invalid",	
		})
		return
	}

	if err := db.C(collectionName).RemoveId(bson.ObjectIdHex(id)); err!= nil{
		rnd.JSON(w, http.StatusProcessing, renderer.M{
            "message":"Failed to delete todo",
            "error":err,
        })
        return
	}

	rnd.JSON(w, http.StatusOK, renderer.M{
        "message":"todo deleted successfully",
    })

}

func main(){
	stopChan := make(chan os.Signal)
	signal.Notify(stopChan, os.Interrupt)
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Get("/", homeHandler)
	r.Mount("/todo", todoHandlers())

	srv := &http.Server{
		Addr: port,
		Handler: r,
		ReadTimeout: 60*time.Second,
		WriteTimeout: 60*time.Second,
		IdleTimeout: 60*time.Second,
	}

	go func(){
		log.Println("Listening to Port", port)
		if err := srv.ListenAndServe(); err != nil{
			log.Printf("listen:%\n", err)
		}
	}()

	<-stopChan
	log.Println("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(),5*time.Second)
	srv.Shutdown(ctx)
	defer cancel(
		log.Println("server gracefully stopped")
	)

}

func todoHandlers() http.Handler{
	rg := chi.NewRouter()
	rg.Group(func(r chi.Router){
		r.Get("/", fetchTodos)
		r.Post("/", createTodo)
		r.Put("/{id}", updateTodo)
		r.Delete("{/id}", deleteTodo)
	})
}


func checkErr(err error){
	if err != nil{
		log.Fatal(err)
	}
}