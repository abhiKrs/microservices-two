// // main_test.go

package main

// import (
// 	"bytes"
// 	"encoding/json"
// 	"net/http"
// 	"net/http/httptest"
// 	"os"

// 	// // "strconv"
// 	"testing"
// 	"web-api/app/server"
// 	// sourceSchema "web-api/app/source/schema"
// 	// // "github.com/logfire-sh/web-api"
// 	// "github.com/stretchr/testify/require"
// )

// var a = server.New()

// func TestMain(m *testing.M) {
// 	a.Init("testing")
// 	code := m.Run()
// 	clearTable()
// 	os.Exit(code)
// }

// func clearTable() {
// 	a.DB().Exec("DELETE FROM users")
// 	a.DB().Exec("DELETE FROM credentials")
// 	a.DB().Exec("DELETE FROM magic_links")
// 	a.DB().Exec("DELETE FROM profiles")
// 	a.DB().Exec("DELETE FROM emails")
// 	// a.DB().Exec("ALTER SEQUENCE products_id_seq RESTART WITH 1")
// }

// //
// // )

// func executeRequest(req *http.Request) *httptest.ResponseRecorder {
// 	rr := httptest.NewRecorder()
// 	a.Router.ServeHTTP(rr, req)

// 	return rr
// }

// func checkResponseCode(t *testing.T, expected, actual int) {
// 	if expected != actual {
// 		t.Errorf("Expected response code %d. Got %d\n", expected, actual)
// 	}
// }

// func TestEmptyTable(t *testing.T) {
// 	// clearTable()

// 	req, _ := http.NewRequest("GET", "/auth", nil)
// 	response := executeRequest(req)

// 	checkResponseCode(t, http.StatusOK, response.Code)

// 	if body := response.Body.String(); body != "Hello From Auth!!!" {
// 		t.Errorf("Expected an empty array. Got %s", body)
// 	}
// }

// func TestCreateAuth(t *testing.T) {

// 	clearTable()

// 	var jsonStr = []byte(`{"email":"rahulprakash@contient.com"}`)
// 	req, _ := http.NewRequest("POST", "/auth/signup", bytes.NewBuffer(jsonStr))
// 	req.Header.Set("Content-Type", "application/json")

// 	response := executeRequest(req)
// 	checkResponseCode(t, http.StatusCreated, response.Code)

// 	var m map[string]interface{}
// 	json.Unmarshal(response.Body.Bytes(), &m)

// 	if m["email"] != "rahulprakash@constient.com" {
// 		t.Errorf("Expected product name to be 'rahulprakash@constient.com'. Got '%v'", m["email"])
// 	}

// 	if m["isSuccessful"] != true {
// 		t.Errorf("Expected isSuccessful to be 'true'. Got '%v'", m["isSuccessful"])
// 	}

// 	// the id is compared to 1.0 because JSON unmarshaling converts numbers to
// 	// floats, when the target is a map[string]interface{}
// 	// if m["expiryTime"] < time. {
// 	//     t.Errorf("Expected product ID to be '1'. Got '%v'", m["id"])
// 	// }
// }
