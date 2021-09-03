package main

import (
	"fmt"
	"github.com/justinas/nosurf"
	_ "github.com/justinas/nosurf"
	"net/http"
)

func WriteToConcole(next http.Handler) http.Handler {
	 return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		  fmt.Println("Hit The page")
	   next.ServeHTTP(w,r)
	 })
}

//NoSurf adds CSRF protection to all Post Request
func NoSurf(next http.Handler) http.Handler {
   	   crsfHanlder := nosurf.New(next)
   	   crsfHanlder.SetBaseCookie(http.Cookie{
   	   	HttpOnly: true,
   	   	 Path: "/", // for perPage bases, the entire site
   	   	 Secure: app.InProduction ,
   	   	 SameSite: http.SameSiteLaxMode,
	   })
  return crsfHanlder
}


//loads and saves the session on every session
//tell the webapp to remember state using session
func SessionLoad(next http.Handler) http.Handler {
	return session.LoadAndSave(next)
}

