package app

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/evgenv123/go-shortener/model"
	"github.com/evgenv123/go-shortener/service"
	"github.com/go-chi/chi/v5"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

// MyHandlerGetID is for getting full URL from shortened
func MyHandlerGetID(w http.ResponseWriter, r *http.Request) {
	requestedID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "Wrong requested ID!", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), appConf.CtxTimeout*time.Second)
	defer cancel()
	obj, err2 := URLSvc.GetObjFromShortID(ctx, model.ShortID(requestedID))
	switch {
	case errors.Is(err2, service.ErrItemDeleted):
		http.Error(w, "Object deleted!", http.StatusGone)
		return
	case err2 != nil:
		http.Error(w, "Error finding object for short id!", http.StatusBadRequest)
		return
	default:
		http.Redirect(w, r, obj.LongURL, http.StatusTemporaryRedirect)
	}
}

// MyHandlerListUrls is for getting all URLS for specified user
func MyHandlerListUrls(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), appConf.CtxTimeout*time.Second)
	defer cancel()

	urls, err := URLSvc.GetUserURLs(ctx, r.Context().Value(ContextKeyUserID).(string))
	if err != nil {
		if errors.Is(err, service.ErrNoURLsForUser) {
			http.Error(w, "No links for user!", http.StatusNoContent)
			return
		} else {
			http.Error(w, "Error getting URLS for user!", http.StatusInternalServerError)
			return
		}
	}
	var result []OutputAllURLs
	// Appending urls to output array
	for _, v := range urls {
		result = append(result, OutputAllURLs{
			ShortURL:    URLSvc.GetFullLinkShortObj(ctx, &v),
			OriginalURL: v.LongURL},
		)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(result)
	if err != nil {
		http.Error(w, "Cannot write reply body!", http.StatusInternalServerError)
		return
	}
}

// MyHandlerPost is for shortening full URL and saving info to DB
func MyHandlerPost(w http.ResponseWriter, r *http.Request) {
	// reading original link body
	b, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Cannot read request body!", http.StatusInternalServerError)
		return
	}
	// Checking for valid URL
	_, err = url.ParseRequestURI(string(b))
	if err != nil {
		http.Error(w, "Wrong URL format!", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), appConf.CtxTimeout*time.Second)
	defer cancel()

	// Trying to shorten URL from request
	shortened, err := URLSvc.ShortenURL(ctx, string(b), r.Context().Value(ContextKeyUserID).(string))

	if err != nil {
		// If we receive duplicate error from SQL
		var myErr *service.DuplicateFullURLErr
		if errors.As(err, &myErr) {
			// Send StatusConflict and existing short url
			w.WriteHeader(http.StatusConflict)
		} else {
			http.Error(w, "Cannot shorten URL!"+err.Error(), http.StatusInternalServerError)
			return
		}
	} else {
		w.WriteHeader(http.StatusCreated)
	}
	// Writing reply
	_, err = w.Write([]byte(URLSvc.GetFullLinkShortObj(ctx, shortened)))
	if err != nil {
		http.Error(w, "Cannot write reply body!", http.StatusInternalServerError)
		return
	}
}

// MyHandlerShorten is a handler for /api/shorten endpoint
func MyHandlerShorten(w http.ResponseWriter, r *http.Request) {
	// reading original link body
	var input InputURL

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Cannot decode request body!", http.StatusInternalServerError)
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), appConf.CtxTimeout*time.Second)
	defer cancel()

	w.Header().Set("Content-Type", "application/json")

	// Trying to shorten URL from request
	shortened, err := URLSvc.ShortenURL(ctx, input.URL, r.Context().Value(ContextKeyUserID).(string))

	if err != nil {
		var myErr *service.DuplicateFullURLErr
		var myErr2 *service.InvalidURLError
		if errors.As(err, &myErr) {
			// Send StatusConflict and existing short url
			w.WriteHeader(http.StatusConflict)
		} else if errors.As(err, &myErr2) {
			http.Error(w, "Wrong URL format!", http.StatusBadRequest)
			return
		} else {
			http.Error(w, "Cannot shorten URL!"+err.Error(), http.StatusInternalServerError)
			return
		}
	} else {
		w.WriteHeader(http.StatusCreated)
	}

	err = json.NewEncoder(w).Encode(OutputShortURL{
		Result: URLSvc.GetFullLinkShortObj(ctx, shortened),
	})
	if err != nil {
		http.Error(w, "Cannot write reply body!", http.StatusInternalServerError)
		return
	}
}

// MyHandlerShortenBatch is a handler for /api/shorten/batch endpoint
func MyHandlerShortenBatch(w http.ResponseWriter, r *http.Request) {
	var input []InputBatch
	var output []OutputBatch
	// reading request body
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Cannot decode request body!", http.StatusInternalServerError)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), appConf.CtxTimeout*time.Second)
	defer cancel()

	w.Header().Set("Content-Type", "application/json")

	for i := 0; i < len(input); i++ {
		// Trying to shorten URL from request
		shortened, err := URLSvc.ShortenURL(ctx, input[i].OrigURL, r.Context().Value(ContextKeyUserID).(string))

		if err != nil {
			// If we receive duplicate error from SQL
			var myErr *service.DuplicateFullURLErr
			var myErr2 *service.InvalidURLError
			if errors.As(err, &myErr) {
				http.Error(w, "Duplicate Full URL "+myErr.FullURL, http.StatusBadRequest)
				return
			} else if errors.As(err, &myErr2) {
				http.Error(w, "Wrong URL format!", http.StatusBadRequest)
				return
			} else {
				http.Error(w, "Cannot shorten URL!"+err.Error(), http.StatusInternalServerError)
				return
			}
		}
		output = append(output, OutputBatch{
			CorrID:   input[i].CorrID,
			ShortURL: URLSvc.GetFullLinkShortObj(ctx, shortened),
		})
	}

	w.WriteHeader(http.StatusCreated)
	err := json.NewEncoder(w).Encode(output)
	if err != nil {
		http.Error(w, "Cannot write reply body!", http.StatusInternalServerError)
		return
	}
}

// MyHandlerPing is a handler for /ping endpoint
func MyHandlerPing(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), appConf.CtxTimeout*time.Second)
	defer cancel()
	if !URLSvc.Ping(ctx) {
		http.Error(w, "Cannot ping database!", http.StatusInternalServerError)
		return
	} else {
		w.WriteHeader(http.StatusOK)
	}
}

// MyHandlerDelete is a handler for /api/user/urls endpoint
func MyHandlerDelete(w http.ResponseWriter, r *http.Request) {
	var input []InputDelete

	// reading request body
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		log.Println("Cannot decode request body! " + err.Error())
		http.Error(w, "Cannot decode request body!", http.StatusInternalServerError)
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), appConf.CtxTimeout*time.Second)
	defer cancel()
	// Getting full objects from storage using provided model.ShortID
	var recvdObjects []model.ShortenedURL
	currentUserID := r.Context().Value(ContextKeyUserID).(string)
	for _, v := range input {
		if obj, err := URLSvc.GetObjFromShortID(ctx, model.ShortID(v)); err == nil {
			// Checking Authorization
			if obj.UserID != currentUserID {
				http.Error(w, "Only owner can delete its records!", http.StatusForbidden)
				return
			}
			recvdObjects = append(recvdObjects, *obj)
		} else {
			log.Println("Cannot get objects from DB! " + err.Error())
			http.Error(w, "Cannot get objects from DB! "+err.Error(), http.StatusInternalServerError)
			return
		}
	}
	if err := URLSvc.DeleteBatchURL(ctx, recvdObjects); err != nil {
		log.Println("Cannot delete objects! " + err.Error())
		http.Error(w, "Cannot delete objects! "+err.Error(), http.StatusInternalServerError)
		return
	} else {
		w.WriteHeader(http.StatusAccepted)
	}
}
