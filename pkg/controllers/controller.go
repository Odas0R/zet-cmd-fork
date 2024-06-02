package controllers

import (
	"net/http"
	"path"

	"github.com/a-h/templ"
	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3" // Import the SQLite driver
	"github.com/odas0r/zet/pkg/database"
	"github.com/odas0r/zet/pkg/domain/workspace"
	wq "github.com/odas0r/zet/pkg/domain/workspace/sqlite"
	"github.com/odas0r/zet/pkg/domain/zettel"
	zq "github.com/odas0r/zet/pkg/domain/zettel/sqlite"
	"github.com/odas0r/zet/pkg/view"
)

type Controller struct {
	workspaceRepo workspace.Repository
	zettelRepo    zettel.Repository
}

func NewController(db *database.Database) (*Controller, error) {
	workspaceRepo, err := wq.New(db)
	if err != nil {
		return nil, err
	}
	zettelRepo, err := zq.New(db)
	if err != nil {
		return nil, err
	}

	return &Controller{
		workspaceRepo: workspaceRepo,
		zettelRepo:    zettelRepo,
	}, nil
}

func (c *Controller) renderError(w http.ResponseWriter, r *http.Request, err error) {
	component := view.ErrorMessage(err.Error())
	templ.Handler(component).ServeHTTP(w, r)
}

func (c *Controller) HandleHome(w http.ResponseWriter, r *http.Request) {
	workspaces, err := c.workspaceRepo.FindAllWorkspaces()
	if err != nil {
		c.renderError(w, r, err)
		return
	}

	var component templ.Component
	if len(workspaces) == 0 {
		component = view.CreateWorkspaceForm()
	} else {
		component = view.ListWorkspaces(workspaces)
	}

	templ.Handler(component).ServeHTTP(w, r)
}

func (c *Controller) HandleListWorkspaces(w http.ResponseWriter, r *http.Request) {
	workspaces, err := c.workspaceRepo.FindAllWorkspaces()
	if err != nil {
		c.renderError(w, r, err)
		return
	}
	component := view.ListWorkspaces(workspaces)
	templ.Handler(component).ServeHTTP(w, r)
}

func (c *Controller) HandleCreateWorkspaceForm(w http.ResponseWriter, r *http.Request) {
	component := view.CreateWorkspaceForm()
	templ.Handler(component).ServeHTTP(w, r)
}

func (c *Controller) HandleCreateWorkspace(w http.ResponseWriter, r *http.Request) {
	path := r.FormValue("path")
	wrk, err := workspace.New(path)
	if err != nil {
		c.renderError(w, r, err)
		return
	}
	if err := c.workspaceRepo.Save(wrk); err != nil {
		c.renderError(w, r, err)
		return
	}
	c.HandleListWorkspaces(w, r)
}

func (c *Controller) HandleEditWorkspaceForm(w http.ResponseWriter, r *http.Request) {
	id := path.Base(r.URL.Path)
	wrkID, err := uuid.Parse(id)
	if err != nil {
		c.renderError(w, r, err)
		return
	}
	wrk, err := c.workspaceRepo.FindWorkspaceByID(wrkID)
	if err != nil {
		c.renderError(w, r, err)
		return
	}
	component := view.EditWorkspaceForm(wrk)
	templ.Handler(component).ServeHTTP(w, r)
}

func (c *Controller) HandleEditWorkspace(w http.ResponseWriter, r *http.Request) {
	id := path.Base(r.URL.Path)
	wrkID, err := uuid.Parse(id)
	if err != nil {
		c.renderError(w, r, err)
		return
	}
	wrk, err := c.workspaceRepo.FindWorkspaceByID(wrkID)
	if err != nil {
		c.renderError(w, r, err)
		return
	}
	wrk.SetPath(r.FormValue("path"))
	if err := c.workspaceRepo.Save(wrk); err != nil {
		c.renderError(w, r, err)
		return
	}
	c.HandleListWorkspaces(w, r)
}

func (c *Controller) HandleDeleteWorkspace(w http.ResponseWriter, r *http.Request) {
	id := path.Base(r.URL.Path)
	wrkID, err := uuid.Parse(id)
	if err != nil {
		c.renderError(w, r, err)
		return
	}
	if err := c.workspaceRepo.Delete(wrkID); err != nil {
		c.renderError(w, r, err)
		return
	}
	c.HandleListWorkspaces(w, r)
}

func (c *Controller) HandleListZettels(w http.ResponseWriter, r *http.Request) {
	zettels, err := c.zettelRepo.FindAllZettels() // You'll need to implement this method in the repository
	if err != nil {
		c.renderError(w, r, err)
		return
	}
	component := view.ListZettels(zettels)
	templ.Handler(component).ServeHTTP(w, r)
}

func (c *Controller) HandleCreateZettelForm(w http.ResponseWriter, r *http.Request) {
	component := view.CreateZettelForm()
	templ.Handler(component).ServeHTTP(w, r)
}

func (c *Controller) HandleCreateZettel(w http.ResponseWriter, r *http.Request) {
	title := r.FormValue("title")
	content := r.FormValue("content")
	kind := r.FormValue("kind")

	zett, err := zettel.New(title, content, zettel.Kind(kind))
	if err != nil {
		c.renderError(w, r, err)
		return
	}
	if err := c.zettelRepo.Save(zett); err != nil {
		c.renderError(w, r, err)
		return
	}
	c.HandleListZettels(w, r)
}

func (c *Controller) HandleEditZettelForm(w http.ResponseWriter, r *http.Request) {
	id := path.Base(r.URL.Path)
	zettID, err := uuid.Parse(id)
	if err != nil {
		c.renderError(w, r, err)
		return
	}
	zett, err := c.zettelRepo.FindByID(zettID)
	if err != nil {
		c.renderError(w, r, err)
		return
	}
	component := view.EditZettelForm(zett)
	templ.Handler(component).ServeHTTP(w, r)
}

func (c *Controller) HandleEditZettel(w http.ResponseWriter, r *http.Request) {
	id := path.Base(r.URL.Path)
	zettID, err := uuid.Parse(id)
	if err != nil {
		c.renderError(w, r, err)
		return
	}
	zett, err := c.zettelRepo.FindByID(zettID)
	if err != nil {
		c.renderError(w, r, err)
		return
	}
	zett.SetTitle(r.FormValue("title"))
	zett.SetBody(r.FormValue("content"))
	zett.SetKind(zettel.Kind(r.FormValue("kind")))
	if err := c.zettelRepo.Save(zett); err != nil {
		c.renderError(w, r, err)
		return
	}
	c.HandleListZettels(w, r)
}

func (c *Controller) HandleDeleteZettel(w http.ResponseWriter, r *http.Request) {
	id := path.Base(r.URL.Path)
	zettID, err := uuid.Parse(id)
	if err != nil {
		c.renderError(w, r, err)
		return
	}
	if err := c.zettelRepo.Delete(zettID); err != nil {
		c.renderError(w, r, err)
		return
	}
	c.HandleListZettels(w, r)
}
