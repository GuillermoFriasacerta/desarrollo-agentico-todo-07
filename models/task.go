package models

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/gobuffalo/pop/v6"
	"github.com/gobuffalo/validate/v3"
	"github.com/gofrs/uuid"
)

// TitleMaxLength es el número máximo de caracteres permitido en el título.
const TitleMaxLength = 200

// DueDateLayout es el formato de la fecha límite en los formularios, el de <input type="date">.
const DueDateLayout = "2006-01-02"

// dueDateDisplayLayout es el formato de la fecha límite en la lista.
const dueDateDisplayLayout = "02/01/2006"

// ErrInvalidDueDate es el error de una fecha límite que no es un día válido.
var ErrInvalidDueDate = errors.New("La fecha límite no es válida.")

// Task representa una tarea de la lista.
type Task struct {
	ID        uuid.UUID  `json:"id" db:"id" form:"-"`
	Title     string     `json:"title" db:"title" form:"title"`
	Completed bool       `json:"completed" db:"completed" form:"-"`
	DueDate   *time.Time `json:"due_date" db:"due_date" form:"-"`
	CreatedAt time.Time  `json:"created_at" db:"created_at" form:"-"`
	UpdatedAt time.Time  `json:"updated_at" db:"updated_at" form:"-"`
}

// Tasks es una colección de Task.
type Tasks []Task

// Validate se ejecuta en cada llamada a los métodos "pop.Validate*".
func (t *Task) Validate(tx *pop.Connection) (*validate.Errors, error) {
	verrs := validate.NewErrors()

	title := strings.TrimSpace(t.Title)
	if title == "" {
		verrs.Add("title", "El título es obligatorio.")
	} else if len([]rune(title)) > TitleMaxLength {
		verrs.Add("title", fmt.Sprintf("El título no puede superar los %d caracteres.", TitleMaxLength))
	}

	return verrs, nil
}

// ParseDueDate convierte la fecha límite de un formulario (AAAA-MM-DD) en un día.
// Una fecha vacía o solo con espacios deja la tarea sin fecha límite.
func ParseDueDate(value string) (*time.Time, error) {
	value = strings.TrimSpace(value)
	if value == "" {
		return nil, nil
	}

	day, err := time.Parse(DueDateLayout, value)
	if err != nil {
		return nil, ErrInvalidDueDate
	}
	return &day, nil
}

// DueDateValue devuelve la fecha límite en el formato de los formularios, o "" si no tiene.
func (t Task) DueDateValue() string {
	if t.DueDate == nil {
		return ""
	}
	return t.DueDate.Format(DueDateLayout)
}

// DueDateText devuelve la fecha límite como se lee en la lista (DD/MM/AAAA), o "" si no tiene.
func (t Task) DueDateText() string {
	if t.DueDate == nil {
		return ""
	}
	return t.DueDate.Format(dueDateDisplayLayout)
}
