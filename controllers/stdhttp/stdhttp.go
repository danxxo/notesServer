package stdhttp

import (
	"encoding/json"
	"net/http"
	storage "notesServer/gate/storage"
	"notesServer/models/dto"
	errorLogger "notesServer/pkg/errorLogger"
	"strconv"

	"github.com/pkg/errors"
)

type Controller struct {
	DB     storage.Storage
	Srv    *http.Server
	Logger *errorLogger.ErrorLogger
}

func NewController(addr string, db storage.Storage, logger *errorLogger.ErrorLogger) *Controller {

	controller := &Controller{
		Srv: &http.Server{
			Addr: addr,
		},
		DB:     db,
		Logger: logger,
	}

	http.HandleFunc("/add", controller.NoteAdd)
	http.HandleFunc("/get", controller.NoteGetByID)
	http.HandleFunc("/update", controller.NoteUpdate)
	http.HandleFunc("/delete", controller.NoteDeleteByPhone)
	http.HandleFunc("/configure", controller.NoteConfigure)

	return controller
}

func (c *Controller) NoteConfigure(w http.ResponseWriter, r *http.Request) {
	response := dto.Response{}

	defer responseWriteAndReturn(w, &response)

	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	config := dto.Config{}
	err := json.NewDecoder(r.Body).Decode(&config)
	if err != nil {
		err = errors.Wrap(err, "stdhttp.NoteConfigure(): json.NewDecoder(r.Body).Decode(&note)")
		c.Logger.LogError(err)
		response.Wrap("error", nil, err)
		return
	}

	storageNum, err := strconv.Atoi(config.StorageNum)
	if err != nil {
		err = errors.Wrap(err, "stdhttp.NoteConfigure(): strconv.Atoi(config.StorageNum)")
		c.Logger.LogError(err)
		response.Wrap("error", nil, err)
		return
	}

	if storageNum == 0 {
		c.DB = storage.NewMap()
	} else if storageNum == 1 {
		c.DB = storage.NewList()
	} else {
		err = errors.New("stdhttp.NoteConfigure(): Invalid configure number option, must be 0 or 1")
		c.Logger.LogError(err)
		response.Wrap("error", nil, err)
	}

	response.Wrap("OK", nil, nil)
}

func (c *Controller) NoteAdd(w http.ResponseWriter, r *http.Request) {

	response := dto.Response{}

	defer responseWriteAndReturn(w, &response)

	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	note := dto.Note{}
	err := json.NewDecoder(r.Body).Decode(&note)
	if err != nil {
		err = errors.Wrap(err, "stdhttp.NoteAdd(): json.NewDecoder(r.Body).Decode(&note)")
		c.Logger.LogError(err)
		response.Wrap("error", nil, err)
		return
	}

	if note.AuthorFirstName == "" || note.AuthorLastName == "" || note.Note == "" {
		err = errors.New("stdhttp.NoteAdd(): Not all Note fields filled")
		c.Logger.LogError(err)
		response.Wrap("error", nil, err)
		return
	}

	_, err = c.DB.Add(note)
	if err != nil {
		err = errors.Wrap(err, "stdhttp.NoteAdd(): c.DB.Add(note)")
		c.Logger.LogError(err)
		response.Wrap("error", nil, err)
		return
	}

	response.Wrap("OK", nil, nil)
}

func (c *Controller) NoteGetByID(w http.ResponseWriter, r *http.Request) {

	response := dto.Response{}

	defer responseWriteAndReturn(w, &response)

	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	note := dto.Note{}

	err := json.NewDecoder(r.Body).Decode(&note)
	if err != nil {
		err = errors.Wrap(err, "stdhttp.NoteGetByID(): json.NewDecoder(r.Body).Decode(&note)")
		c.Logger.LogError(err)
		response.Wrap("error", nil, err)
		return
	}

	id, err := strconv.Atoi(note.Id)
	if err != nil {
		err = errors.Wrap(err, "stdhttp.NoteGetByID(): strconv.Atoi(note.Id)")
		c.Logger.LogError(err)
		response.Wrap("error", nil, err)
		return
	}

	noteGet, err := c.DB.GetByIndex(int64(id))
	if err != nil {
		err = errors.Wrap(err, "stdhttp.NoteGetByID(): c.DB.GetByIndex(int64(id))")
		c.Logger.LogError(err)
		response.Wrap("error", nil, err)
		return
	}
	noteStruct := noteGet.(dto.Note)

	notesJSON, err := json.Marshal(&noteStruct)
	if err != nil {
		err = errors.Wrap(err, "stdhttp.NoteGetByID(): json.Marshal(&notes)")
		c.Logger.LogError(err)
		response.Wrap("error", nil, err)
		return
	}

	response.Wrap("OK", notesJSON, nil)
}

func (c *Controller) NoteUpdate(w http.ResponseWriter, r *http.Request) {

	response := dto.Response{}

	defer responseWriteAndReturn(w, &response)

	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	noteRequest := dto.Note{}

	err := json.NewDecoder(r.Body).Decode(&noteRequest)
	if err != nil {
		err = errors.Wrap(err, "stdhttp.NoteUpdate(): json.NewDecoder(r.Body).Decode(&noteRequest)")
		c.Logger.LogError(err)
		response.Wrap("error", nil, err)
		return
	}

	id := noteRequest.Id
	idInt, err := strconv.Atoi(id)
	if err != nil {
		err = errors.Wrap(err, "stdhttp.NoteUpdate(): strconv.Atoi(id)")
		c.Logger.LogError(err)
		response.Wrap("error", nil, err)
		return
	}

	noteToUpdateInterface, err := c.DB.GetByIndex(int64(idInt))
	if err != nil {
		err = errors.Wrap(err, "stdhttp.NoteUpdate(): c.DB.GetByIndex(int64(idInt))")
		c.Logger.LogError(err)
		response.Wrap("error", nil, err)
		return
	}
	noteToUpdate := noteToUpdateInterface.(dto.Note)

	if noteRequest.AuthorFirstName != "" {
		noteToUpdate.AuthorFirstName = noteRequest.AuthorFirstName
	}
	if noteRequest.AuthorLastName != "" {
		noteToUpdate.AuthorLastName = noteRequest.AuthorLastName
	}
	if noteRequest.Note != "" {
		noteToUpdate.Note = noteRequest.Note
	}

	err = c.DB.UpdateByIndex(int64(idInt), noteToUpdate)
	if err != nil {
		err = errors.Wrap(err, "stdhttp.NoteUpdate(): c.DB.UpdateByIndex(int64(idInt), noteToUpdate)")
		c.Logger.LogError(err)
		response.Wrap("error", nil, err)
		return
	}

	response.Wrap("OK", nil, nil)
}

func (c *Controller) NoteDeleteByPhone(w http.ResponseWriter, r *http.Request) {

	response := dto.Response{}

	defer responseWriteAndReturn(w, &response)

	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	note := dto.Note{}

	err := json.NewDecoder(r.Body).Decode(&note)
	if err != nil {
		err = errors.Wrap(err, "stdhttp.NoteDeleteByPhone(): json.NewDecoder(r.Body).Decode(&note)")
		c.Logger.LogError(err)
		response.Wrap("error", nil, err)
		return
	}

	id := note.Id
	idInt, err := strconv.Atoi(id)
	if err != nil {
		err = errors.Wrap(err, "stdhttp.NoteDeleteByPhone(): strconv.Atoi(id)")
		c.Logger.LogError(err)
		response.Wrap("error", nil, err)
		return
	}

	err = c.DB.RemoveByIndex(int64(idInt))
	if err != nil {
		err = errors.Wrap(err, "stdhttp.NoteDeleteByPhone(): c.DB.RemoveByIndex(int64(idInt))")
		c.Logger.LogError(err)
		response.Wrap("error", nil, err)
		return
	}
	response.Wrap("OK", nil, nil)
}

func responseWriteAndReturn(w http.ResponseWriter, response *dto.Response) {
	errEncode := json.NewEncoder(w).Encode(response)
	if errEncode != nil {
		w.WriteHeader(http.StatusPaymentRequired)
	}
}
