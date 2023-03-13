package controller

// import (
// 	// "io"
// 	"encoding/json"
// 	"net/http"
// 	pgsql "web-api/app/database"
// 	"web-api/test/testModels"

// 	"github.com/go-chi/chi/v5"
// )

// type DummyController struct{}

// // func (d *DummyController) sayHello(w http.ResponseWriter, r *http.Request) {
// // 	w.Write([]byte("Hello World!"))
// // }

// func (rs DummyController) List(w http.ResponseWriter, r *http.Request) {
// 	w.Header().Set("Content-Type", "application/json")
// 	// w.Header().st

// 	dummy := []testModels.Dummy{}
// 	pgsql.DB.Find(&dummy)
// 	// w.Write([]byte("Hello World!"))
// 	json.NewEncoder(w).Encode(dummy)

// }

// // Request Handler - POST /Dummy - Create a new post.
// func (rs DummyController) Create(w http.ResponseWriter, r *http.Request) {

// 	w.Write([]byte("Hello World!"))

// 	// resp, err := http.Post("https://jsonplaceholder.typicode.com/Dummy", "application/json", r.Body)

// 	// if err != nil {
// 	//   http.Error(w, err.Error(), http.StatusInternalServerError)
// 	//   return
// 	// }

// 	// defer resp.Body.Close()

// 	// w.Header().Set("Content-Type", "application/json")

// 	// if _, err := io.Copy(w, resp.Body); err != nil {
// 	//   http.Error(w, err.Error(), http.StatusInternalServerError)
// 	//   return
// 	// }
// }
// func (rs DummyController) Post(w http.ResponseWriter, r *http.Request) {

// 	w.Header().Set("Content-Type", "application/json")

// 	var dummy testModels.Dummy

// 	pgsql.DB.Where("id = ?", chi.URLParam(r, "id")).First(&dummy)
// 	// w.Write([]byte("Hello World!"))
// 	json.NewEncoder(w).Encode(&dummy)

// }

// func (rs DummyController) Update(w http.ResponseWriter, r *http.Request) {
// 	w.Header().Set("Content-Type", "application/json")
// 	// w.Header().st
// 	var dummy testModels.Dummy

// 	pgsql.DB.Where("id = ?", chi.URLParam(r, "id")).Delete(&dummy)
// 	// w.Write([]byte("Hello World!"))
// 	json.NewEncoder(w).Encode(&dummy)

// }

// func (rs DummyController) Delete(w http.ResponseWriter, r *http.Request) {

// 	w.Header().Set("Content-Type", "application/json")
// 	// w.Header().st
// 	var dummy testModels.Dummy

// 	pgsql.DB.Where("id = ?", chi.URLParam(r, "id")).First(&dummy)
// 	// w.Write([]byte("Hello World!"))
// 	json.NewEncoder(w).Encode(&dummy)
// }
