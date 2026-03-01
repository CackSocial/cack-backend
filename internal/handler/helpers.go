package handler

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/CackSocial/cack-backend/internal/dto"
	"github.com/CackSocial/cack-backend/internal/infrastructure/storage"
	ucerrors "github.com/CackSocial/cack-backend/internal/usecase/errors"
	"github.com/CackSocial/cack-backend/pkg/response"
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
	case errors.Is(err, ucerrors.ErrUserNotFound), errors.Is(err, ucerrors.ErrPostNotFound), errors.Is(err, ucerrors.ErrCommentNotFound), errors.Is(err, ucerrors.ErrRepostNotFound):
		response.Error(c, http.StatusNotFound, err.Error())
	case errors.Is(err, ucerrors.ErrInvalidCredentials), errors.Is(err, ucerrors.ErrUnauthorized):
		response.Error(c, http.StatusUnauthorized, err.Error())
	case errors.Is(err, ucerrors.ErrUsernameTaken):
		response.Error(c, http.StatusConflict, err.Error())
	case errors.Is(err, ucerrors.ErrSelfFollow), errors.Is(err, ucerrors.ErrAlreadyFollowing), errors.Is(err, ucerrors.ErrAlreadyLiked), errors.Is(err, ucerrors.ErrAlreadyBookmarked), errors.Is(err, ucerrors.ErrAlreadyReposted), errors.Is(err, ucerrors.ErrCannotRepost):
		response.Error(c, http.StatusBadRequest, err.Error())
	case errors.Is(err, storage.ErrFileTooLarge), errors.Is(err, storage.ErrInvalidFileType):
		response.Error(c, http.StatusBadRequest, err.Error())
	default:
		slog.Error("internal server error",
			"error", err,
			"method", c.Request.Method,
			"path", c.Request.URL.Path,
		)
		response.Error(c, http.StatusInternalServerError, "internal server error")
	}
}
