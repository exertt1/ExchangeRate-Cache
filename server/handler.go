package server

import (
	core "Excnahge-Cacher/core/cache"
	"net/http"
)

type HTTPHandler struct {
	Cache *core.Cache
}

type AllCoursesResponse struct {
	Course map[string]float64
}

func (h *HTTPHandler) GetAllCourses(w http.Response, r *http.Request) {

}

func (h *HTTPHandler) GetCourse(fromValue, ToValue string) string {

}
