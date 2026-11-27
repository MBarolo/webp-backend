package controller

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/mbarolo/webp_backend/forms"
	"github.com/mbarolo/webp_backend/services"
	"github.com/mbarolo/webp_backend/utils"
)

// GetProfile godoc
// @Summary      Obtener perfil
// @Description  Obtener información del usuario autenticado
// @Tags         users
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  map[string]interface{}
// @Failure      401  {object}  map[string]interface{}
// @Router       /users/profile [get]
func GetProfile(c *gin.Context) {
	user, err := services.GetCurrentUser(c.Request)
	if err != nil {
		utils.JsonResponse(c, http.StatusUnauthorized, "Error al obtener el usuario autenticado: "+err.Error(), nil)
		return
	}

	utils.JsonResponse(c, http.StatusOK, "Perfil obtenido", user)
}

// GetAllUsers godoc
// @Summary      Obtener todos los usuarios
// @Description  Listar todos los usuarios del sistema (admin)
// @Tags         admin
// @Accept       json
// @Produce      json
// @Security     BasicAuth
// @Success      200  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /admin/users [get]
func GetAllUsers(c *gin.Context) {
	users, err := services.GetAllUsers()
	if err != nil {
		utils.JsonResponse(c, http.StatusInternalServerError, "Error al obtener los usuarios: "+err.Error(), nil)
		return
	}

	utils.JsonResponse(c, http.StatusOK, "Usuarios obtenidos", users)
}

// GetUserById godoc
// @Summary      Obtener usuario por ID
// @Description  Obtener información de un usuario específico (admin)
// @Tags         admin
// @Accept       json
// @Produce      json
// @Security     BasicAuth
// @Param        id   path      int  true  "ID del usuario"
// @Success      200  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /admin/users/{id} [get]
func GetUserById(c *gin.Context) {
	idStr := c.Param("id")
	if idStr == "" {
		utils.JsonResponse(c, http.StatusBadRequest, "Parametro id no encontrado", nil)
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		utils.JsonResponse(c, http.StatusInternalServerError, "Error al convertir el parametro id: "+err.Error(), nil)
		return
	}

	user, err := services.GetUserById(int64(id))
	if err != nil {
		utils.JsonResponse(c, http.StatusInternalServerError, "Error al obtener el usuario: "+err.Error(), nil)
		return
	}

	utils.JsonResponse(c, http.StatusOK, "Usuario obtenido", user)
}

// UpdateProfile godoc
// @Summary      Actualizar perfil
// @Description  Modificar información del usuario autenticado
// @Tags         users
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        user  body      forms.UserForm  true  "Datos actualizados del usuario"
// @Success      200   {object}  map[string]interface{}
// @Failure      400   {object}  map[string]interface{}
// @Failure      401   {object}  map[string]interface{}
// @Failure      500   {object}  map[string]interface{}
// @Router       /users/profile [patch]
func UpdateProfile(c *gin.Context) {
	user, err := services.GetCurrentUser(c.Request)
	if err != nil {
		utils.JsonResponse(c, http.StatusUnauthorized, "Error al obtener el usuario autenticado: "+err.Error(), nil)
		return
	}

	var userForm *forms.UserForm
	if err := json.NewDecoder(c.Request.Body).Decode(&userForm); err != nil {
		utils.JsonResponse(c, http.StatusBadRequest, "Error al decodificar el cuerpo de la solicitud: "+err.Error(), nil)
		return
	}

	updatedUser, err := services.UpdateUser(user.Id, userForm)
	if err != nil {
		utils.JsonResponse(c, http.StatusInternalServerError, "Error al actualizar el usuario: "+err.Error(), nil)
		return
	}

	utils.JsonResponse(c, http.StatusOK, "Usuario actualizado correctamente", updatedUser)
}

// UpdateUser godoc
// @Summary      Actualizar usuario
// @Description  Modificar información de un usuario específico (admin)
// @Tags         admin
// @Accept       json
// @Produce      json
// @Security     BasicAuth
// @Param        id    path      int             true  "ID del usuario"
// @Param        user  body      forms.UserForm  true  "Datos actualizados del usuario"
// @Success      200   {object}  map[string]interface{}
// @Failure      400   {object}  map[string]interface{}
// @Failure      500   {object}  map[string]interface{}
// @Router       /admin/users/{id} [patch]
func UpdateUser(c *gin.Context) {
	idStr := c.Param("id")
	if idStr == "" {
		utils.JsonResponse(c, http.StatusBadRequest, "Parametro id no encontrado", nil)
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		utils.JsonResponse(c, http.StatusInternalServerError, "Error al convertir el parametro id: "+err.Error(), nil)
		return
	}

	var userForm *forms.UserForm
	if err := json.NewDecoder(c.Request.Body).Decode(&userForm); err != nil {
		utils.JsonResponse(c, http.StatusBadRequest, "Error al decodificar el cuerpo de la solicitud: "+err.Error(), nil)
		return
	}

	updatedUser, err := services.UpdateUser(int64(id), userForm)
	if err != nil {
		utils.JsonResponse(c, http.StatusInternalServerError, "Error al actualizar el usuario: "+err.Error(), nil)
		return
	}

	utils.JsonResponse(c, http.StatusOK, "Usuario actualizado correctamente", updatedUser)
}
