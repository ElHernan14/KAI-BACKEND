package users

import (
	"context"
	"fmt"
	usersdto "kai-back/internal/modules/users/dto"
	contextutil "kai-back/internal/shared/context"
	errorHandler "kai-back/internal/shared/errors"
	helper "kai-back/internal/shared/helpers"
	"kai-back/internal/shared/response"
	validatorx "kai-back/internal/shared/validator"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ServicePort interface {
	GetMe(ctx context.Context, userID uuid.UUID) (*usersdto.MeResponse, error)
	UpdateMe(ctx context.Context, userID uuid.UUID, req usersdto.UpdateMeRequest) error
	ChangePassword(ctx context.Context, userID uuid.UUID, req usersdto.ChangePasswordRequest) error
	UpdateProfilePhoto(ctx context.Context, userID uuid.UUID) error
	GetUserProfile(
		ctx context.Context,
		userID uuid.UUID,
	) (*usersdto.UserProfileResponse, error)
	UpdateUserProfile(
		ctx context.Context,
		userID uuid.UUID,
		req usersdto.UpdateUserProfileRequest,
	) error
	UpdateUserConfiguration(
		ctx context.Context,
		userID uuid.UUID,
		req usersdto.UpdateUserConfigurationRequest,
	) error
}

type Controller struct {
	service ServicePort
}

func NewController(service ServicePort) *Controller {
	return &Controller{service: service}
}

func (ctrl *Controller) Me(c *gin.Context) {

	userIDString, exists := contextutil.GetUserID(c.Request.Context())
	if !exists {
		_ = c.Error(
			errorHandler.NewAppError(
				http.StatusUnauthorized,
				"usuario no autenticado",
			),
		)
		return
	}

	userID, err := uuid.Parse(userIDString)
	if err != nil {
		_ = c.Error(err)
		return
	}

	user, err := ctrl.service.GetMe(
		c.Request.Context(),
		userID,
	)

	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(
		http.StatusOK,
		response.Success(user),
	)
}

func (ctrl *Controller) UpdateMe(c *gin.Context) {
	var req usersdto.UpdateMeRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		_ = c.Error(errorHandler.NewAppError(
			http.StatusBadRequest,
			"request json invalido",
		))
		return
	}

	if message, hasError := validatorx.ValidateStruct(req); hasError {
		_ = c.Error(errorHandler.NewAppError(
			http.StatusBadRequest,
			message,
		))
		return
	}

	userIDString, exists := contextutil.GetUserID(c.Request.Context())
	if !exists {
		_ = c.Error(
			errorHandler.NewAppError(
				http.StatusUnauthorized,
				"usuario no autenticados",
			),
		)
		return
	}

	userID, err := uuid.Parse(userIDString)
	if err != nil {
		_ = c.Error(err)
		return
	}

	errService := ctrl.service.UpdateMe(
		c.Request.Context(),
		userID,
		req,
	)
	if errService != nil {
		_ = c.Error(errService)
		return
	}

	c.JSON(
		http.StatusOK,
		response.SuccessWithCode(
			http.StatusOK,
			"perfil actualizado correctamente",
		),
	)
}

func (ctrl *Controller) UpdateProfilePhoto(c *gin.Context) {

	userIDString, exists := contextutil.GetUserID(c.Request.Context())
	if !exists {
		_ = c.Error(
			errorHandler.NewAppError(
				http.StatusUnauthorized,
				"usuario no autenticado",
			),
		)
		return
	}

	userID, err := uuid.Parse(userIDString)
	if err != nil {
		_ = c.Error(err)
		return
	}

	file, err := c.FormFile("foto")
	if err != nil {
		_ = c.Error(
			errorHandler.NewAppError(
				http.StatusBadRequest,
				"foto de perfil requerida",
			),
		)
		return
	}

	if file.Size > maxProfilePhotoMB*1024*1024 {
		_ = c.Error(
			errorHandler.NewAppError(
				http.StatusBadRequest,
				fmt.Sprintf("la foto no puede superar %dMB", maxProfilePhotoMB),
			),
		)
		return
	}

	ext := normalizedProfilePhotoExt(file.Filename)
	if !isAllowedProfilePhotoExt(ext) {
		_ = c.Error(
			errorHandler.NewAppError(
				http.StatusBadRequest,
				"formato de foto no permitido",
			),
		)
		return
	}

	if err := ctrl.service.UpdateProfilePhoto(c.Request.Context(), userID); err != nil {
		_ = c.Error(err)
		return
	}

	if err := os.MkdirAll(profileUploadDir, 0755); err != nil {
		_ = c.Error(
			errorHandler.NewAppError(
				http.StatusInternalServerError,
				"no se pudo preparar carpeta de uploads",
			),
		)
		return
	}

	if err := removeProfilePhotos(userID); err != nil {
		_ = c.Error(
			errorHandler.NewAppError(
				http.StatusInternalServerError,
				"no se pudo reemplazar foto de perfil",
			),
		)
		return
	}

	fileName := userID.String() + ext
	destination := filepath.Join(profileUploadDir, fileName)

	if err := c.SaveUploadedFile(file, destination); err != nil {
		_ = c.Error(
			errorHandler.NewAppError(
				http.StatusInternalServerError,
				"no se pudo guardar foto de perfil",
			),
		)
		return
	}

	profilePhotoURL := "/" + filepath.ToSlash(destination)

	c.JSON(
		http.StatusOK,
		response.Success(gin.H{
			"foto_perfil": profilePhotoURL,
		}),
	)
}

func (ctrl *Controller) ChangePassword(c *gin.Context) {

	var req usersdto.ChangePasswordRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		_ = c.Error(
			errorHandler.NewAppError(
				http.StatusBadRequest,
				"request json invalido",
			),
		)
		return
	}

	if message, hasError := validatorx.ValidateStruct(req); hasError {
		_ = c.Error(
			errorHandler.NewAppError(
				http.StatusBadRequest,
				message,
			),
		)
		return
	}

	userIDString, exists := contextutil.GetUserID(c.Request.Context())
	if !exists {
		_ = c.Error(
			errorHandler.NewAppError(
				http.StatusUnauthorized,
				"usuario no autenticado",
			),
		)
		return
	}

	userID, errParse := uuid.Parse(userIDString)
	if errParse != nil {
		_ = c.Error(errParse)
		return
	}

	err := ctrl.service.ChangePassword(
		c.Request.Context(),
		userID,
		req,
	)

	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(
		http.StatusOK,
		response.SuccessWithCode(
			http.StatusOK,
			"contraseña actualizada correctamente",
		),
	)
}

func (ctrl *Controller) GetProfile(c *gin.Context) {
	userID, err := helper.ValidateUserUUID(c)
	if err != nil {
		_ = c.Error(err)
		return
	}

	res, err := ctrl.service.GetUserProfile(c.Request.Context(), userID)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(
		http.StatusOK,
		response.Success(res),
	)
}

func (ctrl *Controller) PutUserProfile(c *gin.Context) {
	var req usersdto.UpdateUserProfileRequest

	if helper.BindAndValidate(c, &req) {
		return
	}

	userID, err := helper.ValidateUserUUID(c)
	if err != nil {
		_ = c.Error(err)
		return
	}

	err = ctrl.service.UpdateUserProfile(c.Request.Context(), userID, req)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(
		http.StatusOK,
		response.Success("Perfil actualizado exitosamente"),
	)
}

func (ctrl *Controller) PutUserConfigurationProfile(c *gin.Context) {
	var req usersdto.UpdateUserConfigurationRequest

	if helper.BindAndValidate(c, &req) {
		return
	}

	userID, err := helper.ValidateUserUUID(c)
	if err != nil {
		_ = c.Error(err)
		return
	}

	err = ctrl.service.UpdateUserConfiguration(c.Request.Context(), userID, req)

	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(
		http.StatusOK,
		response.Success("Configuración de perfil actualizada exitosamente"),
	)
}
