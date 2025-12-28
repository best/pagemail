package errors

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ProblemDetail struct {
	Type     string            `json:"type"`
	Title    string            `json:"title"`
	Status   int               `json:"status"`
	Detail   string            `json:"detail,omitempty"`
	Instance string            `json:"instance,omitempty"`
	Errors   []ValidationError `json:"errors,omitempty"`
}

type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func NewProblemDetail(status int, title, detail string) *ProblemDetail {
	typeURL := fmt.Sprintf("https://pagemail.app/errors/%s", slugify(title))
	return &ProblemDetail{
		Type:   typeURL,
		Title:  title,
		Status: status,
		Detail: detail,
	}
}

func BadRequest(detail string) *ProblemDetail {
	return NewProblemDetail(http.StatusBadRequest, "Bad Request", detail)
}

func Unauthorized(detail string) *ProblemDetail {
	return NewProblemDetail(http.StatusUnauthorized, "Unauthorized", detail)
}

func Forbidden(detail string) *ProblemDetail {
	return NewProblemDetail(http.StatusForbidden, "Forbidden", detail)
}

func NotFound(detail string) *ProblemDetail {
	return NewProblemDetail(http.StatusNotFound, "Not Found", detail)
}

func Conflict(detail string) *ProblemDetail {
	return NewProblemDetail(http.StatusConflict, "Conflict", detail)
}

func InternalError(detail string) *ProblemDetail {
	return NewProblemDetail(http.StatusInternalServerError, "Internal Server Error", detail)
}

func ValidationFailed(errors []ValidationError) *ProblemDetail {
	pd := NewProblemDetail(http.StatusBadRequest, "Validation Error", "One or more fields failed validation")
	pd.Errors = errors
	return pd
}

func (p *ProblemDetail) WithInstance(instance string) *ProblemDetail {
	p.Instance = instance
	return p
}

func (p *ProblemDetail) Respond(c *gin.Context) {
	p.Instance = c.Request.URL.Path
	c.JSON(p.Status, p)
}

func slugify(s string) string {
	result := ""
	for _, r := range s {
		if r >= 'a' && r <= 'z' {
			result += string(r)
		} else if r >= 'A' && r <= 'Z' {
			result += string(r + 32)
		} else if r == ' ' {
			result += "-"
		}
	}
	return result
}
