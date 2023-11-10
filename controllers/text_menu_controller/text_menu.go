package text_menu_controller

import (
	"gin-wire-template/models/text_menu"
	"gin-wire-template/utils/http_util"
	"github.com/gin-gonic/gin"
	"strconv"
	"strings"
)

// GetTextMenu
// @Summary		Get Text Menu
// @Description	get text menu items
// @Tags		Text Menu Demo
// @Accept		json
// @Produce		json
// @Param		item_types	query	[]int	false "//1 = WEB, 2 = APP, 3 = BOTH"
// @Success		200	{object}	 []text_menu.TextMenu
// @Failure		500	{object}	http_util.HTTPError
// @Router		/api/v1/text_menu [get]
func (controller *Controller) GetTextMenu(c *gin.Context) {
	var textMenu []text_menu.TextMenu

	var itemTypes []int
	itemType, _ := c.GetQuery("item_types")

	for _, d := range strings.Split(itemType, ",") {
		if i, err := strconv.Atoi(d); err == nil {
			itemTypes = append(itemTypes, i)
		}
	}

	cmsMenuList, err := controller.textMenuRepo.GetMenu(c, itemTypes)
	if err != nil {
		http_util.RenderErrorResponse(c, err.Error())
		return
	}
	if len(cmsMenuList) == 0 {
		http_util.RenderSuccessResponse(c, []text_menu.TextMenu{})
		return
	}

	http_util.RenderSuccessResponse(c, textMenu)
	return
}

// CreateTextMenu
// @Summary		Create Text Menu
// @Description	create text menu <br> TODO - please implement the logic for this endpoint
// @Tags		Text Menu Demo
// @Accept		json
// @Produce		json
// @Param		menu	body	text_menu.TextMenu	true "text menu json item"
// @Success		200	{object}	 http_util.HTTPSuccess
// @Failure		500	{object}	http_util.HTTPError
// @Router		/api/v1/text_menu [post]
func (controller *Controller) CreateTextMenu(c *gin.Context) {
	var menu text_menu.TextMenu

	if err := c.BindJSON(&menu); err != nil {
		http_util.RenderErrorResponse(c, err.Error())
		return
	}

	err := controller.textMenuRepo.CreateMenu(c, menu)
	if err != nil {
		http_util.RenderErrorResponse(c, err.Error())
		return
	}

	http_util.RenderSuccessResponse(c, nil)
	return
}

// UpdateTextMenu
// @Summary		Update Text Menu
// @Description	update text menu <br> TODO - please implement the logic for this endpoint
// @Tags		Text Menu Demo
// @Accept		json
// @Produce		json
// @Param		id		path	string				true	"item id"
// @Param		menu	body	text_menu.TextMenu	true	"text menu json item"
// @Success		200	{object}	 http_util.HTTPSuccess
// @Failure		500	{object}	http_util.HTTPError
// @Router		/api/v1/text_menu/{id} [put]
func (controller *Controller) UpdateTextMenu(c *gin.Context) {
	var menu text_menu.TextMenu
	if err := c.BindJSON(&menu); err != nil {
		http_util.RenderErrorResponse(c, "menu param error:"+err.Error())
		return
	}

	err := controller.textMenuRepo.UpdateMenu(c, c.Param("id"), menu)
	if err != nil {
		http_util.RenderErrorResponse(c, err.Error())
		return
	}

	http_util.RenderSuccessResponse(c, nil)
	return
}

// DeleteTextMenu
// @Summary		Delete Text Menu
// @Description	delete text menu <br> TODO - please implement the logic for this endpoint
// @Tags		Text Menu Demo
// @Accept		json
// @Produce		json
// @Param		id		path	string				true	"item id"
// @Success		200	{object}	 http_util.HTTPSuccess
// @Failure		500	{object}	http_util.HTTPError
// @Router		/api/v1/text_menu/{id} [delete]
func (controller *Controller) DeleteTextMenu(c *gin.Context) {
	err := controller.textMenuRepo.DeleteMenu(c, c.Param("id"))
	if err != nil {
		http_util.RenderErrorResponse(c, err.Error())
		return
	}

	http_util.RenderSuccessResponse(c, nil)
	return
}
