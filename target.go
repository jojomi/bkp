package bkp

import "fmt"

type Target struct {
	Name     string `json:"name"`
	Filename string
	Type     string `json:"type"`
	Path     string `json:"path"`
	Password string `json:"password"`
}

func (t *Target) IsReady() bool {
	switch t.Type {
	case "local":
		ex, err := exists(t.Path)
		return err == nil && ex
	}
	return true
}

func (t *Target) String() string {
	return fmt.Sprintf("%s [%s]", t.Name, t.Type)
}
