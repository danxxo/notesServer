package stdhttp

import (
	"encoding/json"
	"fmt"
	"net/http"
	storage "notesServer/gate/storage"
	"notesServer/models/dto"
	errorLogger "notesServer/pkg/errorLogger"
	"strconv"

	"github.com/pkg/errors"
)

type Controller struct {
	DB  storage.Storage
	Srv *http.Server
}

func NewController(addr string, db storage.Storage) *Controller {

	controller := &Controller{
		Srv: &http.Server{
			Addr: addr,
		},
		DB: db,
	}

	http.HandleFunc("/add", controller.NoteAdd)
	http.HandleFunc("/get", controller.NotesGet)
	http.HandleFunc("/update", controller.NoteUpdate)
	http.HandleFunc("/delete", controller.NoteDeleteByPhone)

	return controller
}

func (c *Controller) NoteAdd(w http.ResponseWriter, r *http.Request) {

	// creating logger
	logger, err := errorLogger.NewErrorLogger(".log")
	if err != nil {
		fmt.Println(err)
		return
	}

	response := dto.Response{}

	defer responseWriteAndReturn(w, &response)

	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	note := dto.Note{}
	err = json.NewDecoder(r.Body).Decode(&note)
	if err != nil {
		err = errors.Wrap(err, "stdhttp.NoteAdd()")
		logger.LogError(err)
		response.ErrorWrap(err)
		return
	}

	_, err = c.DB.Add(note)
	if err != nil {
		err = errors.Wrap(err, "stdhttp.NoteAdd()")
		logger.LogError(err)
		response.ErrorWrap(err)
		return
	}

	response.Wrap("OK", nil, nil)

	c.DB.Print()

}

func (c *Controller) NotesGet(w http.ResponseWriter, r *http.Request) {

	// creating logger
	logger, err := errorLogger.NewErrorLogger(".log")
	if err != nil {
		fmt.Println(err)
		return
	}

	response := dto.Response{}

	defer responseWriteAndReturn(w, &response)

	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	note := dto.Note{}

	err = json.NewDecoder(r.Body).Decode(&note)
	if err != nil {
		err = errors.Wrap(err, "stdhttp.NoteDeleteByPhone()")
		logger.LogError(err)
		response.ErrorWrap(err)
		return
	}

	ids := c.DB.GetAllByValueSelectedFields(note)
	var notes dto.Notes
	for _, v := range ids {
		note, _ := c.DB.GetByIndex(v)
		noteStruct := note.(dto.Note)
		notes = append(notes, noteStruct)
	}

	notesJSON, err := json.Marshal(&notes)
	if err != nil {
		err = errors.Wrap(err, "stdhttp.NotesGet()")
		logger.LogError(err)
		response.ErrorWrap(err)
		return
	}

	response.Wrap("OK", notesJSON, nil)
}

func (c *Controller) NoteUpdate(w http.ResponseWriter, r *http.Request) {

	// creating logger
	logger, err := errorLogger.NewErrorLogger(".log")
	if err != nil {
		fmt.Println(err)
		return
	}

	response := dto.Response{}

	defer responseWriteAndReturn(w, &response)

	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	note := dto.Note{}

	err = json.NewDecoder(r.Body).Decode(&note)
	if err != nil {
		err = errors.Wrap(err, "stdhttp.NoteDeleteByPhone()")
		logger.LogError(err)
		response.ErrorWrap(err)
		return
	}

	id := note.Id
	idInt, err := strconv.Atoi(id)
	if err != nil {
		err = errors.Wrap(err, "stdhttp.NoteDeleteByPhone")
		logger.LogError(err)
		response.ErrorWrap(err)
		return
	}

	c.DB.RemoveByIndex(int64(idInt))
	response.Wrap("OK", nil, nil)
}

func (c *Controller) NoteDeleteByPhone(w http.ResponseWriter, r *http.Request) {

	// creating logger
	logger, err := errorLogger.NewErrorLogger(".log")
	if err != nil {
		fmt.Println(err)
		return
	}

	response := dto.Response{}

	defer responseWriteAndReturn(w, &response)

	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	note := dto.Note{}

	err = json.NewDecoder(r.Body).Decode(&note)
	if err != nil {
		err = errors.Wrap(err, "stdhttp.NoteDeleteByPhone()")
		logger.LogError(err)
		response.ErrorWrap(err)
		return
	}

	id := note.Id
	idInt, err := strconv.Atoi(id)
	if err != nil {
		err = errors.Wrap(err, "stdhttp.NoteDeleteByPhone")
		logger.LogError(err)
		response.ErrorWrap(err)
		return
	}

	c.DB.RemoveByIndex(int64(idInt))
	response.Wrap("OK", nil, nil)
}

func responseWriteAndReturn(w http.ResponseWriter, response *dto.Response) {
	errEncode := json.NewEncoder(w).Encode(response)
	if errEncode != nil {
		w.WriteHeader(http.StatusPaymentRequired)
	}
}
