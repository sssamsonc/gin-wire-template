package repositories

import (
	"gin-wire-template/repositories/text_menu_repository"
	"github.com/google/wire"
)

var RepositorySet = wire.NewSet(
	text_menu_repository.NewRepository,
)
