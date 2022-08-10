package controller

import (
	"sync"

	"tdgames/storage"
)

// InfoController is an implementation of controller.Info
type InfoController struct {
	stg *storage.Storage
	m   sync.Mutex
}

// NewInfoController constructor of InfoController struct
func NewInfoController(stg *storage.Storage) *InfoController {
	return &InfoController{stg: stg}
}

// GetLastTimestamp returns time of the latest post
func (icon *InfoController) GetLastTimestamp() uint64 {
	return icon.stg.GetLastTimestamp()
}

// SetLastTimestamp sets time of the latest post
func (icon *InfoController) SetLastTimestamp(tsp uint64) {

}
