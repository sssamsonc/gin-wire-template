package text_menu_controller

import (
	"gin-wire-template/repositories/text_menu_repository"
)

type Controller struct {
	textMenuRepo *text_menu_repository.Repository
}

func NewController(textMenuRepo *text_menu_repository.Repository) *Controller {
	return &Controller{
		textMenuRepo: textMenuRepo,
	}
}
