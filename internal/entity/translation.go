// Package entity defines main entities for business logic (services), data base mapping and
// HTTP response objects if suitable. Each logic group entities in own file.
package entity

// Translation -.
type Translation struct {
	Source      string `example:"auto"                               json:"source"`
	Destination string `example:"en"                                 json:"destination"`
	Original    string `example:"текст для перевода"                 json:"original"`
	Translation string `example:"text for translation"               json:"translation"`
}
