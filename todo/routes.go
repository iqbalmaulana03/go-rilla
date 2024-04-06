package todo

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/oklog/ulid/v2"
)

func writeMessage(w http.ResponseWriter, status int, msg string) {
	var j struct {
		Msg string `json:"message"`
	}

	j.Msg = msg

	w.Header().Add("content-type", "application/json")
	w.WriteHeader(status)

	json.NewEncoder(w).Encode(j)
}

func wrietError(w http.ResponseWriter, status int, err error) {
	writeMessage(w, status, err.Error())
}

func ListItemsHandler(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	resp, err := listItem(ctx)

	if err != nil {
		wrietError(w, http.StatusInternalServerError, err)
	}

	w.Header().Add("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

func GetItemHandler(w http.ResponseWriter, req *http.Request) {

	ctx := req.Context()

	vars := mux.Vars(req)

	itemId := vars["itemId"]

	id, err := ulid.Parse(itemId)

	if err != nil {
		wrietError(w, http.StatusBadRequest, err)
		return
	}

	var resp TodoItem

	item, err := findItem(ctx, id)

	if err != nil {
		if err == ErrTodoNotFound {
			writeMessage(w, http.StatusNotFound, err.Error())
			return
		}

		wrietError(w, http.StatusInternalServerError, err)
		return
	}

	resp = item

	w.Header().Add("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

func CreateItemHandler(w http.ResponseWriter, req *http.Request) {
	var requestBody struct {
		Title string `json:"title"`
	}

	if err := json.NewDecoder(req.Body).Decode(&requestBody); err != nil {
		wrietError(w, http.StatusBadRequest, err)
		return
	}

	ctx := req.Context()

	id, err := createItem(ctx, requestBody.Title)

	if err != nil {
		wrietError(w, http.StatusBadRequest, err)
		return
	}

	var resp struct {
		Id string `json:"id"`
	}

	resp.Id = id.String()

	w.Header().Add("content-type", "application/json")
	w.WriteHeader(http.StatusCreated)

	json.NewEncoder(w).Encode(resp)
}

func MakeItemDoneHandler(w http.ResponseWriter, req *http.Request) {
	if err := req.ParseForm(); err != nil {
		wrietError(w, http.StatusBadRequest, err)
		return
	}

	ctx := req.Context()

	idStr := req.FormValue("id")

	id, err := ulid.Parse(idStr)

	if err != nil {
		wrietError(w, http.StatusBadRequest, err)
		return
	}

	err = makeItemDone(ctx, id)

	if err != nil {
		wrietError(w, http.StatusInternalServerError, err)
		return
	}

	w.WriteHeader(http.StatusOK)
}
