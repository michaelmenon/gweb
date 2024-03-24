package gweb

import (
	"errors"
	"net/http"
	"strings"
)

// Use ... add a middleware for the Group
func (wg *WebGroup) Use(f WebHandler) {
	wg.middlewares = append(wg.middlewares, f)
}
func (wg *WebGroup) Get(pattern string, f WebHandler) error {
	if !strings.HasPrefix(pattern, "/") {
		return errors.New(InvalidPath)
	}
	wg.w.addRoutes(http.MethodGet+" "+wg.pattern+pattern, f, wg)
	return nil
}

func (wg *WebGroup) Post(pattern string, f WebHandler) error {
	if !strings.HasPrefix(pattern, "/") {
		return errors.New(InvalidPath)
	}
	wg.w.addRoutes(http.MethodPost+" "+wg.pattern+pattern, f, wg)
	return nil
}

func (wg *WebGroup) Put(pattern string, f WebHandler) error {
	if !strings.HasPrefix(pattern, "/") {
		return errors.New(InvalidPath)
	}
	wg.w.addRoutes(http.MethodPut+" "+wg.pattern+pattern, f, wg)
	return nil
}

func (wg *WebGroup) Patch(pattern string, f WebHandler) error {
	if !strings.HasPrefix(pattern, "/") {
		return errors.New(InvalidPath)
	}
	wg.w.addRoutes(http.MethodPatch+" "+wg.pattern+pattern, f, wg)
	return nil
}

func (wg *WebGroup) Delete(pattern string, f WebHandler) error {
	if !strings.HasPrefix(pattern, "/") {
		return errors.New(InvalidPath)
	}
	wg.w.addRoutes(http.MethodDelete+" "+wg.pattern+pattern, f, wg)
	return nil
}

func (wg *WebGroup) Options(pattern string, f WebHandler) error {
	if !strings.HasPrefix(pattern, "/") {
		return errors.New(InvalidPath)
	}
	wg.w.addRoutes(http.MethodOptions+" "+wg.pattern+pattern, f, wg)
	return nil
}
