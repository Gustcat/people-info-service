package params

import (
	"github.com/gin-gonic/gin"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/Gustcat/people-info-service/internal/lib/response"
)

func ParseIDParam(c *gin.Context, log *slog.Logger) (int64, bool) {
	idStr := c.Param("id")

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		log.Error("Invalid id parameter", slog.String("parameter", idStr))
		c.AbortWithStatusJSON(http.StatusBadRequest, response.Error("invalid id parameter"))
		return 0, false
	}

	return id, true
}
