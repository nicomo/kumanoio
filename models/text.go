package models

import (
	"encoding/json"
	"time"

	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/pop/nulls"

	"github.com/gobuffalo/uuid"
	"github.com/gobuffalo/validate"
)

// Text is the base struct for content on our site
type Text struct {
	ID          uuid.UUID  `json:"id" db:"id"`
	CreatedAt   time.Time  `json:"created_at" db:"created_at"`
	PublishedAt nulls.Time `json:"published_at" db:"published_at"`
	UpdatedAt   time.Time  `json:"updated_at" db:"updated_at"`
	Title       string     `json:"title" db:"title"`
	Content     string     `json:"content" db:"content"`
	Author      User       `belongs_to:"user"`
	AuthorID    uuid.UUID  `json:"author_id" db:"author_id"`
	Draft       bool       `json:"draft" db:"draft"`
	StarredBy   Users      `many_to_many:"stars" db:"-"`
}

// String is not required by pop and may be deleted
func (t Text) String() string {
	jt, _ := json.Marshal(t)
	return string(jt)
}

// Texts is not required by pop and may be deleted
type Texts []Text

// String is not required by pop and may be deleted
func (t Texts) String() string {
	jt, _ := json.Marshal(t)
	return string(jt)
}

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
// This method is not required and may be deleted.
func (t *Text) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateCreate gets run every time you call "pop.ValidateAndCreate" method.
// This method is not required and may be deleted.
func (t *Text) ValidateCreate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateUpdate gets run every time you call "pop.ValidateAndUpdate" method.
// This method is not required and may be deleted.
func (t *Text) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}
