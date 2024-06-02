package controllers

import (
	"net/http"

	"github.com/odas0r/zet/pkg/domain/workspace"
)

type Controller struct {
	Workspace workspace.Workspace
}

func NewController(ws workspace.Workspace) *Controller {
	return &Controller{Workspace: ws}
}

func (c *Controller) HandleHome(w http.ResponseWriter, r *http.Request) {
	// TODO
}
