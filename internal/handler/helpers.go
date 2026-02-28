package handler

import (
	"errors"
	"net/http"

	"github.com/CackSocial/cack-backend/internal/dto"
	"github.com/CackSocial/cack-backend/pkg/response"
	ucerrors "github.com/CackSocial/cack-backend/internal/usecase/errors"
	"github.com/gin-gonic/gin"
)

func getUserID(c *gin.Context) string {
	userID, _ := c.Get("userID")
	if id, ok := userID.(string); ok {
		return id
	}
	return ""
}

func getPagination(c *gin.Context) (int, int) {
	var q dto.PaginationQuery
	if err := c.ShouldBindQuery(&q); err != nil {
		q.Page = 1
		q.Limit = 20
	}
	if q.Page < 1 {
		q.Page = 1
	}
	if q.Limit < 1 {
		q.Limit = 20
	}
	return q.Page, q.Limit
}

func handleError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, ucerrors.ErrUserNotFound), errors.Is(err, ucerrors.ErrPostNotFound), errors.Is(err, ucerrors.ErrCommentNotFound):
		response.Error(c, http.StatusNotFound, err.Error())
	case errors.Is(err, ucerrors.ErrInvalidCredentials), errors.Is(err, ucerrors.ErrUnauthorized):
		response.Error(c, http.StatusUnauthorized, err.Error())
	case errors.Is(err, ucerrors.ErrUsernameTaken):
		response.Error(c, http.StatusConflict, err.Error())
	case errors.Is(err, ucerrors.ErrSelfFollow), errors.Is(err, ucerrors.ErrAlreadyFollowing), errors.Is(err, ucerrors.ErrAlreadyLiked):
		response.Error(c, http.StatusBadRequest, err.Error())
	default:
		response.Error(c, http.StatusInternalServerError, "internal server error")
	}
}
