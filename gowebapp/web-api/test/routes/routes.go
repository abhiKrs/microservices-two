package routes

// import (
// 	// "net/http"
// 	"web-api/test/controller"

// 	"github.com/go-chi/chi/v5"
// )

// // type DummyRoute struct {

// // }

// // func DummyRoute() chi.Router {
// // 	routes := chi.NewRouter()

// // 	d := controller.DummyController{}
// // 	routes.Get("/", d.sayHello)
// // 	return routes
// // }

// type DummyResource struct{}

// func (rs DummyResource) Routes() chi.Router {
// 	r := chi.NewRouter()

// 	r.Get("/", controller.DummyController{}.List)    // GET /posts - Read a list of posts.
// 	r.Post("/", controller.DummyController{}.Create) // POST /posts - Create a new post.

// 	// r.Route("/{id}", func(r chi.Router) {
// 	// 	r.Use(PostCtx)
// 	// 	r.Get("/", rs.Get)       // GET /posts/{id} - Read a single post by :id.
// 	// 	r.Put("/", rs.Update)    // PUT /posts/{id} - Update a single post by :id.
// 	// 	r.Delete("/", rs.Delete) // DELETE /posts/{id} - Delete a single post by :id.
// 	// })

// 	r.Get("/:id", controller.DummyController{}.Post)   // GET /posts/{id} - Read a single post by :id.
// 	r.Put("/:id", controller.DummyController{}.Update) // PUT /posts/{id} - Update a single post by :id.
// 	r.Delete("/:id", controller.DummyController{}.Delete)

// 	return r
// }
