package controllers

import (
	"gin-wire-template/controllers/text_menu_controller"
	"github.com/google/wire"
)

type Controllers struct {
	TextMenuController *text_menu_controller.Controller
}

func NewControllers(
	textMenuController *text_menu_controller.Controller,
) *Controllers {
	return &Controllers{
		TextMenuController: textMenuController,
	}
}

var ControllerSet = wire.NewSet(
	NewControllers,
	text_menu_controller.NewController,
)
