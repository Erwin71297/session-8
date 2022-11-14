package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/antonlindstrom/pgstore"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo"
)

const SESSION_ID = "id"

func main() {
	e := echo.New()

	// e.Use(echo.WrapMiddleware(context.ClearHandler))

	store := newPostGresStore()

	// http.HandleFunc("/index", func(w http.ResponseWriter, r *http.Request) {
	// 	w.Header().Set("Access-Control-Allow-Origin", "https://www.google.com")
	// 	w.Header().Set("Access-Control-Allow-Methods", "OPTIONS, GET, POST, PUT")
	// 	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, X-CSRF-Token")

	// 	if r.Method == "OPTIONS" {
	// 		w.Write([]byte("allowed"))
	// 		return
	// 	}

	// 	w.Write([]byte("hello"))
	// })

	e.GET("/set", func(c echo.Context) error {
		session, _ := store.Get(c.Request(), SESSION_ID)
		session.Values["message1"] = "hello"
		session.Values["message2"] = "world"
		session.Save(c.Request(), c.Response())

		return c.Redirect(http.StatusTemporaryRedirect, "/get")
	})

	e.GET("/get", func(c echo.Context) error {
		session, _ := store.Get(c.Request(), SESSION_ID)

		if len(session.Values) == 0 {
			return c.String(http.StatusOK, "empty result")
		}

		return c.String(http.StatusOK, fmt.Sprintf(
			"%s %s",
			session.Values["message1"],
			session.Values["message2"],
		))
	})

	e.GET("/delete", func(c echo.Context) error {
		session, _ := store.Get(c.Request(), SESSION_ID)
		session.Options.MaxAge = -1
		session.Save(c.Request(), c.Response())

		return c.Redirect(http.StatusTemporaryRedirect, "/get")
	})

	e.Start(":9000")
}

func newPostGresStore() *pgstore.PGStore {
	url := "postgres://postgres:password@127.0.0.1:5432/postgres?sslmode=disable"
	authKey := []byte("my-auth-key--very-secret")
	encryptionKey := []byte("my-encryption-key-very-secret123")

	store, err := pgstore.NewPGStore(url, authKey, encryptionKey)
	if err != nil {
		log.Println("Error", err)
		os.Exit(0)
	}

	return store
}

func newCookieStore() *sessions.CookieStore {
	authKey := []byte("my-auth-key--very-secret")
	encryptionKey := []byte("my-encryption-key-very-secret123")

	store := sessions.NewCookieStore(authKey, encryptionKey)
	store.Options.Path = "/"
	store.Options.MaxAge = 86400 * 7
	store.Options.HttpOnly = true

	return store
}
