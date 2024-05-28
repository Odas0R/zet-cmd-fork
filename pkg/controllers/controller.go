package controllers

import (
    "fmt"
    "net/http"

    "github.com/a-h/templ"
    "github.com/google/uuid"
    "github.com/odas0r/zet/pkg/components"
    "github.com/odas0r/zet/pkg/domain/workspace"
    "github.com/odas0r/zet/pkg/domain/zettel"
)

type Controller struct {
    Workspace workspace.Workspace
}

func NewController(ws workspace.Workspace) *Controller {
    return &Controller{Workspace: ws}
}

func (c *Controller) HandleHome(w http.ResponseWriter, r *http.Request) {
    if c.Workspace.ID() == uuid.Nil {
        component := components.WorkspaceInvalid()
        templ.Handler(component).ServeHTTP(w, r)
        return
    }
    component := components.Home(c.Workspace.Zettels())
    templ.Handler(component).ServeHTTP(w, r)
}

func (c *Controller) HandleCreateForm(w http.ResponseWriter, r *http.Request) {
    if c.Workspace.ID() == uuid.Nil {
        component := components.WorkspaceInvalid()
        templ.Handler(component).ServeHTTP(w, r)
        return
    }
    component := components.CreateForm()
    templ.Handler(component).ServeHTTP(w, r)
}

func (c *Controller) HandleCreate(w http.ResponseWriter, r *http.Request) {
    if c.Workspace.ID() == uuid.Nil {
        component := components.WorkspaceInvalid()
        templ.Handler(component).ServeHTTP(w, r)
        return
    }
    if err := r.ParseForm(); err != nil {
        http.Error(w, "Invalid form", http.StatusBadRequest)
        return
    }
    title := r.Form.Get("title")
    content := r.Form.Get("content")
    if err := c.CreateZettel(title, content); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (c *Controller) HandleArchive(w http.ResponseWriter, r *http.Request) {
    if c.Workspace.ID() == uuid.Nil {
        component := components.WorkspaceInvalid()
        templ.Handler(component).ServeHTTP(w, r)
        return
    }
    idStr := r.PathValue("id")
    id, err := uuid.Parse(idStr)
    if err != nil {
        http.Error(w, "Invalid ID", http.StatusBadRequest)
        return
    }
    if err := c.ArchiveZettel(id); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (c *Controller) HandleInitializeForm(w http.ResponseWriter, r *http.Request) {
    component := components.InitializeForm()
    templ.Handler(component).ServeHTTP(w, r)
}

func (c *Controller) HandleInitialize(w http.ResponseWriter, r *http.Request) {
    if err := r.ParseForm(); err != nil {
        http.Error(w, "Invalid form", http.StatusBadRequest)
        return
    }
    path := r.Form.Get("path")
    if err := c.InitializeWorkspace(path); err != nil {
        component := components.WorkspaceInvalid()
        templ.Handler(component).ServeHTTP(w, r)
        return
    }
    http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (c *Controller) CreateZettel(title, content string) error {
    z, err := zettel.New(title, content, zettel.Permanent)
    if err != nil {
        return err
    }
    c.Workspace.SetZettels(append(c.Workspace.Zettels(), z))
    return nil
}

func (c *Controller) ArchiveZettel(id uuid.UUID) error {
    zettels := c.Workspace.Zettels()
    for i, z := range zettels {
        if z.ID() == id {
            c.Workspace.SetZettels(append(zettels[:i], zettels[i+1:]...))
            return nil
        }
    }
    return fmt.Errorf("zettel not found")
}

func (c *Controller) InitializeWorkspace(path string) error {
    ws, err := workspace.New(path)
    if err != nil {
        return err
    }
    c.Workspace = ws
    return nil
}
