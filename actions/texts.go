package actions

import (
	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/pop"
	"github.com/nicomo/kumano/models"
	"github.com/pkg/errors"
)

// This file is generated by Buffalo. It offers a basic structure for
// adding, editing and deleting a page. If your model is more
// complex or you need more than the basic implementation you need to
// edit this file.

// Following naming logic is implemented in Buffalo:
// Model: Singular (Text)
// DB Table: Plural (texts)
// Resource: Plural (Texts)
// Path: Plural (/texts)
// View Template Folder: Plural (/templates/texts/)

// TextsResource is the resource for the Text model
type TextsResource struct {
	buffalo.Resource
}

// List gets all Texts. This function is mapped to the path
// GET /texts
func (v TextsResource) List(c buffalo.Context) error {
	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return errors.WithStack(errors.New("no transaction found"))
	}

	texts := &models.Texts{}

	// Paginate results. Params "page" and "per_page" control pagination.
	// Default values are "page=1" and "per_page=20".
	q := tx.PaginateFromParams(c.Params())

	// Retrieve all Texts from the DB
	if err := q.All(texts); err != nil {
		return errors.WithStack(err)
	}

	// Add the paginator to the context so it can be used in the template.
	c.Set("pagination", q.Paginator)

	return c.Render(200, r.Auto(c, texts))
}

// Show gets the data for one Text. This function is mapped to
// the path GET /texts/{text_id}
func (v TextsResource) Show(c buffalo.Context) error {
	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return errors.WithStack(errors.New("no transaction found"))
	}

	// Allocate an empty Text
	text := &models.Text{}

	// To find the Text the parameter text_id is used.
	if err := tx.Find(text, c.Param("text_id")); err != nil {
		return c.Error(404, err)
	}

	return c.Render(200, r.Auto(c, text))
}

// New renders the form for creating a new Text.
// This function is mapped to the path GET /texts/new
func (v TextsResource) New(c buffalo.Context) error {
	return c.Render(200, r.Auto(c, &models.Text{}))
}

// Create adds a Text to the DB. This function is mapped to the
// path POST /texts
func (v TextsResource) Create(c buffalo.Context) error {
	// Allocate an empty Text
	text := &models.Text{}

	// Bind text to the html form elements
	if err := c.Bind(text); err != nil {
		return errors.WithStack(err)
	}

	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return errors.WithStack(errors.New("no transaction found"))
	}

	// Validate the data from the html form
	verrs, err := tx.ValidateAndCreate(text)
	if err != nil {
		return errors.WithStack(err)
	}

	if verrs.HasAny() {
		// Make the errors available inside the html template
		c.Set("errors", verrs)

		// Render again the new.html template that the user can
		// correct the input.
		return c.Render(422, r.Auto(c, text))
	}

	// If there are no errors set a success message
	c.Flash().Add("success", "Text was created successfully")

	// and redirect to the texts index page
	return c.Render(201, r.Auto(c, text))
}

// Edit renders a edit form for a Text. This function is
// mapped to the path GET /texts/{text_id}/edit
func (v TextsResource) Edit(c buffalo.Context) error {
	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return errors.WithStack(errors.New("no transaction found"))
	}

	// Allocate an empty Text
	text := &models.Text{}

	if err := tx.Find(text, c.Param("text_id")); err != nil {
		return c.Error(404, err)
	}

	return c.Render(200, r.Auto(c, text))
}

// Update changes a Text in the DB. This function is mapped to
// the path PUT /texts/{text_id}
func (v TextsResource) Update(c buffalo.Context) error {
	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return errors.WithStack(errors.New("no transaction found"))
	}

	// Allocate an empty Text
	text := &models.Text{}

	if err := tx.Find(text, c.Param("text_id")); err != nil {
		return c.Error(404, err)
	}

	// Bind Text to the html form elements
	if err := c.Bind(text); err != nil {
		return errors.WithStack(err)
	}

	verrs, err := tx.ValidateAndUpdate(text)
	if err != nil {
		return errors.WithStack(err)
	}

	if verrs.HasAny() {
		// Make the errors available inside the html template
		c.Set("errors", verrs)

		// Render again the edit.html template that the user can
		// correct the input.
		return c.Render(422, r.Auto(c, text))
	}

	// If there are no errors set a success message
	c.Flash().Add("success", "Text was updated successfully")

	// and redirect to the texts index page
	return c.Render(200, r.Auto(c, text))
}

// Destroy deletes a Text from the DB. This function is mapped
// to the path DELETE /texts/{text_id}
func (v TextsResource) Destroy(c buffalo.Context) error {
	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return errors.WithStack(errors.New("no transaction found"))
	}

	// Allocate an empty Text
	text := &models.Text{}

	// To find the Text the parameter text_id is used.
	if err := tx.Find(text, c.Param("text_id")); err != nil {
		return c.Error(404, err)
	}

	if err := tx.Destroy(text); err != nil {
		return errors.WithStack(err)
	}

	// If there are no errors set a flash message
	c.Flash().Add("success", "Text was destroyed successfully")

	// Redirect to the texts index page
	return c.Render(200, r.Auto(c, text))
}
