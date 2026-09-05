package httpapi_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/MrKyomoto/ieat/internal/catalog"
	"github.com/MrKyomoto/ieat/internal/config"
	"github.com/MrKyomoto/ieat/internal/httpapi"
	"github.com/MrKyomoto/ieat/internal/mockdata"
)

func TestMockDataSupportsLoginAndCanteenCatalog(t *testing.T) {
	store, err := mockdata.New("ieat-dev-only")
	if err != nil {
		t.Fatal(err)
	}
	router := httpapi.NewRouter(config.Config{
		WebOrigin:  "http://localhost:5173",
		SessionTTL: 24 * time.Hour,
	}, httpapi.Dependencies{Auth: store, Catalog: store})

	login := httptest.NewRequest(http.MethodPost, "/api/v1/auth/session", strings.NewReader(`{
		"email":"student@mail.ustc.edu.cn",
		"password":"ieat-dev-only"
	}`))
	login.Header.Set("Origin", "http://localhost:5173")
	loginRecorder := httptest.NewRecorder()
	router.ServeHTTP(loginRecorder, login)
	if loginRecorder.Code != http.StatusOK {
		t.Fatalf("login status = %d, want %d; body = %s", loginRecorder.Code, http.StatusOK, loginRecorder.Body.String())
	}

	catalogRequest := httptest.NewRequest(http.MethodGet, "/api/v1/catalog/canteens", nil)
	for _, cookie := range loginRecorder.Result().Cookies() {
		catalogRequest.AddCookie(cookie)
	}
	catalogRecorder := httptest.NewRecorder()
	router.ServeHTTP(catalogRecorder, catalogRequest)
	if catalogRecorder.Code != http.StatusOK {
		t.Fatalf("catalog status = %d, want %d; body = %s", catalogRecorder.Code, http.StatusOK, catalogRecorder.Body.String())
	}

	var canteens []catalog.Canteen
	if err := json.NewDecoder(catalogRecorder.Body).Decode(&canteens); err != nil {
		t.Fatal(err)
	}
	if len(canteens) != 1 || canteens[0].Name != "示例食堂" {
		t.Fatalf("canteens = %#v, want one 示例食堂", canteens)
	}
}
