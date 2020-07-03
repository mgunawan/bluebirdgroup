package error

import (
	"fmt"
	"net/http"
	"testing"
)

func TestError(t *testing.T) {
	h := NewHTTPError(http.StatusInternalServerError, "ErrServerDown", "Server is likely down")
	err := h.Error()
	if err == "" {
		t.Fail()
	}
	fmt.Println(err)
}
