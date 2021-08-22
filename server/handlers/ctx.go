package handlers

import (
	"context"
	"net/http"
	"net/url"

	"github.com/anfimovoleh/ms-users/db"
	"github.com/anfimovoleh/ms-users/email"
	"github.com/go-chi/jwtauth"
)

type CtxKey int

const (
	webAppCtxKey = iota
	httpCtxKey
	emailClientCtxKey
	dbCtxKey
	jwtCtxKey
)

func CtxWebApp(webApp *url.URL) func(context.Context) context.Context {
	return func(ctx context.Context) context.Context {
		return context.WithValue(ctx, webAppCtxKey, webApp)
	}
}

func WebApp(r *http.Request) *url.URL {
	return r.Context().Value(webAppCtxKey).(*url.URL)
}

func CtxHTTP(http *url.URL) func(context.Context) context.Context {
	return func(ctx context.Context) context.Context {
		return context.WithValue(ctx, httpCtxKey, http)
	}
}

func HTTP(r *http.Request) *url.URL {
	return r.Context().Value(httpCtxKey).(*url.URL)
}

func CtxEmailClient(emailClient email.Client) func(context.Context) context.Context {
	return func(ctx context.Context) context.Context {
		return context.WithValue(ctx, emailClientCtxKey, emailClient)
	}
}

func EmailClient(r *http.Request) email.Client {
	return r.Context().Value(emailClientCtxKey).(email.Client)
}
func CtxDB(db *db.DB) func(context.Context) context.Context {
	return func(ctx context.Context) context.Context {
		return context.WithValue(ctx, dbCtxKey, db)
	}
}

func DB(r *http.Request) *db.DB {
	return r.Context().Value(dbCtxKey).(*db.DB)
}

func CtxJWT(entry *jwtauth.JWTAuth) func(context.Context) context.Context {
	return func(ctx context.Context) context.Context {
		return context.WithValue(ctx, jwtCtxKey, entry)
	}
}

func JWT(r *http.Request) *jwtauth.JWTAuth {
	return r.Context().Value(jwtCtxKey).(*jwtauth.JWTAuth)
}
