package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"trip_app/api"
	"trip_app/internal/domain"
	"trip_app/internal/handler"
	"trip_app/internal/middleware"
	"trip_app/internal/repository"
	"trip_app/internal/security"
	"trip_app/internal/usecase"
	"trip_app/test/mock"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	testDB          *gorm.DB
	testServer      *echo.Echo
	jwtSecret       = "test-secret-key-for-testing-only"
	mockEmailSender *mock.MockEmailSender
)

// setupTestDB はテスト用DBへの接続とマイグレーションを実行
func setupTestDB(t *testing.T) {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "postgres://user:password@localhost:5432/trip_app?sslmode=disable"
	}

	var err error
	testDB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	require.NoError(t, err, "Failed to connect to test database")

	err = testDB.AutoMigrate(
		&domain.User{},
		&domain.Trip{},
		&domain.Member{},
		&domain.Member{},
		&domain.Member{},
		&domain.Schedule{},
		&domain.ShareToken{},
	)
	require.NoError(t, err, "Failed to migrate database")
}

// cleanupTestDB はテスト後に全テーブルをクリーンアップ
func cleanupTestDB(t *testing.T) {
	testDB.Exec("TRUNCATE TABLE schedules, share_tokens, trips, users RESTART IDENTITY CASCADE")
}

// setupTestServer はテスト用HTTPサーバーを構築（全層を初期化）
func setupTestServer(t *testing.T) {
	userRepo := repository.NewUserRepository(testDB)
	tripRepo := repository.NewTripRepository(testDB)
	scheduleRepo := repository.NewScheduleRepository(testDB)
	shareTokenRepo := repository.NewShareTokenRepository(testDB)
	publicTripRepo := repository.NewPublicTripRepository(testDB)

	passwordGenerator := security.NewPasswordGenerator()
	tokenGenerator := security.NewTokenGenerator()
	authTokenGenerator := security.NewAuthTokenGenerator(jwtSecret)
	mockEmailSender = mock.NewMockEmailSender()

	userValidator := usecase.NewUserUsecaseValidator()
	scheduleUsecaseValidator := usecase.NewScheduleUsecaseValidator()
	userHandlerValidator := handler.NewUserHandlerValidator()
	scheduleHandlerValidator := handler.NewScheduleHandlerValidator()

	userUsecase := usecase.NewUserUsecase(userRepo, userValidator, passwordGenerator, tokenGenerator, authTokenGenerator, mockEmailSender)
	tripUsecase := usecase.NewTripUsecase(tripRepo, tokenGenerator)
	scheduleUsecase := usecase.NewScheduleUsecase(scheduleRepo, scheduleUsecaseValidator)
	shareTokenUsecase := usecase.NewShareTokenUsecase(shareTokenRepo, tokenGenerator)
	publicTripUsecase := usecase.NewPublicTripUsecase(publicTripRepo, tokenGenerator)

	h := handler.NewHandler(
		userUsecase,
		tripUsecase,
		scheduleUsecase,
		shareTokenUsecase,
		publicTripUsecase,
		userHandlerValidator,
		scheduleHandlerValidator,
	)

	tripOwnershipMiddleware := middleware.TripOwnershipMiddleware(tripUsecase)
	authMiddleware := middleware.AuthMiddleware(jwtSecret)
	shareTokenOwnershipMiddleware := middleware.ShareTokenOwnershipMiddleware(publicTripUsecase)

	e := echo.New()
	wrapper := &api.ServerInterfaceWrapper{Handler: h}

	e.POST("/login", wrapper.LoginUser)
	e.POST("/signup", wrapper.CreateUser)
	e.POST("/users/verify/:verificationToken", wrapper.VerifyUser)

	publicTripGroup := e.Group("/public/trips/:shareToken")
	publicTripGroup.Use(shareTokenOwnershipMiddleware)
	publicTripGroup.GET("", wrapper.GetPublicTripByShareToken)
	publicTripGroup.PUT("", wrapper.UpdatePublicTripByShareToken)
	publicTripGroup.GET("/details", wrapper.GetTripDetailsForPublicTrip)
	publicTripGroup.GET("/schedules", wrapper.GetSchedulesForPublicTrip)
	publicTripGroup.POST("/schedules", wrapper.AddScheduleToPublicTrip)
	publicTripGroup.GET("/schedules/:scheduleId", wrapper.GetScheduleForPublicTrip)
	publicTripGroup.PATCH("/schedules/:scheduleId", wrapper.UpdateScheduleForPublicTrip)
	publicTripGroup.DELETE("/schedules/:scheduleId", wrapper.DeleteScheduleForPublicTrip)

	authRequired := e.Group("")
	authRequired.Use(authMiddleware)
	authRequired.POST("/logout", wrapper.LogoutUser)
	authRequired.GET("/me", wrapper.GetMe)
	authRequired.PUT("/me/password", wrapper.ChangePassword)
	authRequired.GET("/trips", wrapper.GetUserTrips)
	authRequired.POST("/trips", wrapper.CreateUserTrip)

	tripOwnerGroup := authRequired.Group("/trips/:tripId")
	tripOwnerGroup.Use(tripOwnershipMiddleware)
	tripOwnerGroup.GET("", wrapper.GetUserTrip)
	tripOwnerGroup.PUT("", wrapper.UpdateUserTrip)
	tripOwnerGroup.DELETE("", wrapper.DeleteUserTrip)
	tripOwnerGroup.GET("/details", wrapper.GetTripDetails)
	tripOwnerGroup.GET("/schedules", wrapper.GetSchedulesForTrip)
	tripOwnerGroup.POST("/schedules", wrapper.AddScheduleToTrip)
	tripOwnerGroup.GET("/schedules/:scheduleId", wrapper.GetScheduleForTrip)
	tripOwnerGroup.PATCH("/schedules/:scheduleId", wrapper.UpdateScheduleForTrip)
	tripOwnerGroup.DELETE("/schedules/:scheduleId", wrapper.DeleteScheduleForTrip)
	tripOwnerGroup.POST("/share", wrapper.CreateShareLinkForTrip)

	testServer = e
}

// makeRequest はHTTPリクエストを送信してレスポンスを返す
func makeRequest(t *testing.T, method, path string, body interface{}, token string) *httptest.ResponseRecorder {
	var reqBody io.Reader
	// リクエストボディがあればJSON形式に変換
	if body != nil {
		jsonBody, err := json.Marshal(body)
		require.NoError(t, err)
		reqBody = bytes.NewBuffer(jsonBody)
	}

	req := httptest.NewRequest(method, path, reqBody)
	// HTTPリクエストオブジェクトを作成
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	if token != "" {
		req.Header.Set(echo.HeaderAuthorization, "Bearer "+token)
	}

	rec := httptest.NewRecorder()
	testServer.ServeHTTP(rec, req)
	return rec
}

// TestMain はテスト全体のエントリーポイント
func TestMain(m *testing.M) {
	code := m.Run()
	os.Exit(code)
}

// TestScenario_BasicUserFlow はユーザー登録・認証・ログインの基本フローをテスト
func TestScenario_BasicUserFlow(t *testing.T) {
	setupTestDB(t)
	defer cleanupTestDB(t)
	setupTestServer(t)

	// ユーザー登録
	signupReq := map[string]interface{}{
		"name":  "testuser",
		"email": "test@example.com",
	}
	rec := makeRequest(t, http.MethodPost, "/signup", signupReq, "")
	assert.Equal(t, http.StatusCreated, rec.Code)

	var signupResp map[string]interface{}
	err := json.Unmarshal(rec.Body.Bytes(), &signupResp)
	require.NoError(t, err)
	assert.Contains(t, signupResp, "id")

	// メール認証（モックから初回パスワードとトークンを取得）
	verificationToken := mockEmailSender.GetLastToken()
	initialPassword := mockEmailSender.GetLastPassword()
	require.NotEmpty(t, verificationToken, "Verification token should be captured by mock")
	require.NotEmpty(t, initialPassword, "Initial password should be captured by mock")

	rec = makeRequest(t, http.MethodPost, "/users/verify/"+verificationToken, nil, "")
	assert.Equal(t, http.StatusOK, rec.Code)

	// ログイン（初回パスワードを使用）
	loginReq := map[string]interface{}{
		"email":    "test@example.com",
		"password": initialPassword,
	}
	rec = makeRequest(t, http.MethodPost, "/login", loginReq, "")
	assert.Equal(t, http.StatusOK, rec.Code)

	var loginResp map[string]interface{}
	err = json.Unmarshal(rec.Body.Bytes(), &loginResp)
	require.NoError(t, err)
	token := loginResp["token"].(string)
	assert.NotEmpty(t, token)

	// /me エンドポイントでユーザー情報取得
	rec = makeRequest(t, http.MethodGet, "/me", nil, token)
	assert.Equal(t, http.StatusOK, rec.Code)

	var meResp map[string]interface{}
	err = json.Unmarshal(rec.Body.Bytes(), &meResp)
	require.NoError(t, err)
	assert.Equal(t, "testuser", meResp["name"])
	assert.Equal(t, "test@example.com", meResp["email"])
}

// TestScenario_TripManagementFlow は旅行のCRUD操作をテスト
func TestScenario_TripManagementFlow(t *testing.T) {
	setupTestDB(t)
	defer cleanupTestDB(t)
	setupTestServer(t)

	// ユーザー作成とログイン
	token := createAndLoginUser(t, "tripuser", "trip@example.com", "password123")

	// 旅行を作成
	tripReq := map[string]interface{}{
		"title":     "沖縄旅行",
		"startDate": "2025-10-01",
		"endDate":   "2025-10-05",
		"members":   []interface{}{},
	}
	rec := makeRequest(t, http.MethodPost, "/trips", tripReq, token)
	assert.Equal(t, http.StatusCreated, rec.Code)

	var tripResp map[string]interface{}
	err := json.Unmarshal(rec.Body.Bytes(), &tripResp)
	require.NoError(t, err)
	tripID := tripResp["id"].(string)

	// 旅行一覧を取得
	rec = makeRequest(t, http.MethodGet, "/trips", nil, token)
	assert.Equal(t, http.StatusOK, rec.Code)

	var tripsResp []map[string]interface{}
	err = json.Unmarshal(rec.Body.Bytes(), &tripsResp)
	require.NoError(t, err)
	assert.Len(t, tripsResp, 1)
	assert.Equal(t, "沖縄旅行", tripsResp[0]["title"])

	// 単一旅行の詳細を取得
	rec = makeRequest(t, http.MethodGet, fmt.Sprintf("/trips/%s", tripID), nil, token)
	assert.Equal(t, http.StatusOK, rec.Code)

	// 旅行情報を更新
	updateReq := map[string]interface{}{
		"title":     "沖縄旅行（更新）",
		"startDate": "2025-10-01",
		"endDate":   "2025-10-06",
		"members":   []interface{}{},
	}
	rec = makeRequest(t, http.MethodPut, fmt.Sprintf("/trips/%s", tripID), updateReq, token)
	assert.Equal(t, http.StatusOK, rec.Code)

	// 旅行を削除
	rec = makeRequest(t, http.MethodDelete, fmt.Sprintf("/trips/%s", tripID), nil, token)
	assert.Equal(t, http.StatusNoContent, rec.Code)

	// 削除確認（一覧が空になることを確認）
	rec = makeRequest(t, http.MethodGet, "/trips", nil, token)
	assert.Equal(t, http.StatusOK, rec.Code)
	err = json.Unmarshal(rec.Body.Bytes(), &tripsResp)
	require.NoError(t, err)
	assert.Len(t, tripsResp, 0)
}

// TestScenario_ScheduleManagementFlow はスケジュールのCRUD操作をテスト
func TestScenario_ScheduleManagementFlow(t *testing.T) {
	setupTestDB(t)
	defer cleanupTestDB(t)
	setupTestServer(t)

	token := createAndLoginUser(t, "scheduleuser", "schedule@example.com", "password123")
	tripID := createTrip(t, token, "京都旅行", "2025-10-10", "2025-10-12")

	// スケジュールを作成
	scheduleReq := map[string]interface{}{
		"title":         "清水寺観光",
		"startDateTime": "2025-10-10T10:00:00Z",
		"endDateTime":   "2025-10-10T12:00:00Z",
		"memo":          "清水寺を観光する",
	}
	rec := makeRequest(t, http.MethodPost, fmt.Sprintf("/trips/%s/schedules", tripID), scheduleReq, token)
	assert.Equal(t, http.StatusCreated, rec.Code)

	var scheduleResp map[string]interface{}
	err := json.Unmarshal(rec.Body.Bytes(), &scheduleResp)
	require.NoError(t, err)
	scheduleID := scheduleResp["id"].(string)

	// スケジュール一覧を取得
	rec = makeRequest(t, http.MethodGet, fmt.Sprintf("/trips/%s/schedules", tripID), nil, token)
	assert.Equal(t, http.StatusOK, rec.Code)

	var schedulesResp []map[string]interface{}
	err = json.Unmarshal(rec.Body.Bytes(), &schedulesResp)
	require.NoError(t, err)
	assert.Len(t, schedulesResp, 1)
	assert.Equal(t, "清水寺観光", schedulesResp[0]["title"])

	// 単一スケジュールの詳細を取得
	rec = makeRequest(t, http.MethodGet, fmt.Sprintf("/trips/%s/schedules/%s", tripID, scheduleID), nil, token)
	assert.Equal(t, http.StatusOK, rec.Code)

	// スケジュールを更新
	updateReq := map[string]interface{}{
		"title":         "清水寺と金閣寺観光",
		"startDateTime": "2025-10-10T09:00:00Z",
		"endDateTime":   "2025-10-10T15:00:00Z",
		"memo":          "清水寺と金閣寺を観光する",
	}
	rec = makeRequest(t, http.MethodPatch, fmt.Sprintf("/trips/%s/schedules/%s", tripID, scheduleID), updateReq, token)
	assert.Equal(t, http.StatusOK, rec.Code)

	// スケジュールを削除
	rec = makeRequest(t, http.MethodDelete, fmt.Sprintf("/trips/%s/schedules/%s", tripID, scheduleID), nil, token)
	assert.Equal(t, http.StatusNoContent, rec.Code)
}

// TestScenario_ShareLinkFlow は共有リンク機能をテスト
func TestScenario_ShareLinkFlow(t *testing.T) {
	setupTestDB(t)
	defer cleanupTestDB(t)
	setupTestServer(t)

	token := createAndLoginUser(t, "shareuser", "share@example.com", "password123")
	tripID := createTrip(t, token, "博多旅行", "2025-10-15", "2025-10-20")
	scheduleID := createSchedule(t, token, tripID, "太宰府天満宮観光", "2025-10-15")

	// 共有リンクを作成
	rec := makeRequest(t, http.MethodPost, fmt.Sprintf("/trips/%s/share", tripID), nil, token)
	assert.Equal(t, http.StatusCreated, rec.Code)

	var shareResp map[string]interface{}
	err := json.Unmarshal(rec.Body.Bytes(), &shareResp)
	require.NoError(t, err)
	shareToken := shareResp["shareToken"].(string)
	assert.NotEmpty(t, shareToken)

	// 共有トークンで旅行情報を取得（認証不要）
	rec = makeRequest(t, http.MethodGet, fmt.Sprintf("/public/trips/%s", shareToken), nil, "")
	assert.Equal(t, http.StatusOK, rec.Code)

	var publicTripResp map[string]interface{}
	err = json.Unmarshal(rec.Body.Bytes(), &publicTripResp)
	require.NoError(t, err)
	assert.Equal(t, "博多旅行", publicTripResp["title"])

	// 共有トークンで旅行詳細を取得（スケジュール込み）
	rec = makeRequest(t, http.MethodGet, fmt.Sprintf("/public/trips/%s/details", shareToken), nil, "")
	assert.Equal(t, http.StatusOK, rec.Code)

	var detailsResp map[string]interface{}
	err = json.Unmarshal(rec.Body.Bytes(), &detailsResp)
	require.NoError(t, err)
	schedules := detailsResp["schedules"].([]interface{})
	assert.Len(t, schedules, 1)

	// 共有トークンでスケジュール一覧を取得
	rec = makeRequest(t, http.MethodGet, fmt.Sprintf("/public/trips/%s/schedules", shareToken), nil, "")
	assert.Equal(t, http.StatusOK, rec.Code)

	// 共有トークンで新しいスケジュールを作成
	scheduleReq := map[string]interface{}{
		"title":         "中洲屋台巡り",
		"startDateTime": "2025-10-16T18:00:00Z",
		"endDateTime":   "2025-10-16T21:00:00Z",
		"memo":          "中洲の屋台で博多ラーメンを食べる",
	}
	rec = makeRequest(t, http.MethodPost, fmt.Sprintf("/public/trips/%s/schedules", shareToken), scheduleReq, "")
	assert.Equal(t, http.StatusCreated, rec.Code)

	// 共有トークンでスケジュールを更新
	updateReq := map[string]interface{}{
		"title": "太宰府天満宮と九州国立博物館観光",
	}
	rec = makeRequest(t, http.MethodPatch, fmt.Sprintf("/public/trips/%s/schedules/%s", shareToken, scheduleID), updateReq, "")
	assert.Equal(t, http.StatusOK, rec.Code)

	// 共有トークンで旅行情報を更新
	tripUpdateReq := map[string]interface{}{
		"title":     "博多旅行（更新）",
		"startDate": "2025-10-15",
		"endDate":   "2025-10-20",
	}
	rec = makeRequest(t, http.MethodPut, fmt.Sprintf("/public/trips/%s", shareToken), tripUpdateReq, "")
	assert.Equal(t, http.StatusOK, rec.Code)
}

// TestScenario_AuthorizationFlow は他人の旅行へのアクセス拒否をテスト
func TestScenario_AuthorizationFlow(t *testing.T) {
	setupTestDB(t)
	defer cleanupTestDB(t)
	setupTestServer(t)

	// ユーザー1とユーザー2を作成
	token1 := createAndLoginUser(t, "user1", "user1@example.com", "password123")
	token2 := createAndLoginUser(t, "user2", "user2@example.com", "password123")

	// ユーザー1が旅行作成
	tripID := createTrip(t, token1, "User1の旅行", "2025-10-25", "2025-10-28")

	// ユーザー2がユーザー1の旅行にアクセス試行（失敗するべき）
	rec := makeRequest(t, http.MethodGet, fmt.Sprintf("/trips/%s", tripID), nil, token2)
	assert.Equal(t, http.StatusForbidden, rec.Code)

	// ユーザー2が旅行削除試行（失敗するべき）
	rec = makeRequest(t, http.MethodDelete, fmt.Sprintf("/trips/%s", tripID), nil, token2)
	assert.Equal(t, http.StatusForbidden, rec.Code)

	// ユーザー1は正常にアクセス可能
	rec = makeRequest(t, http.MethodGet, fmt.Sprintf("/trips/%s", tripID), nil, token1)
	assert.Equal(t, http.StatusOK, rec.Code)
}

// TestScenario_PasswordChangeFlow はパスワード変更機能をテスト
func TestScenario_PasswordChangeFlow(t *testing.T) {
	setupTestDB(t)
	defer cleanupTestDB(t)
	setupTestServer(t)

	token := createAndLoginUser(t, "pwuser", "pw@example.com", "oldpassword")

	// 初回パスワードを取得（mockから）
	initialPassword := mockEmailSender.GetLastPassword()

	// パスワードを変更
	changeReq := map[string]interface{}{
		"currentPassword": initialPassword,
		"newPassword":     "newpassword123",
	}
	rec := makeRequest(t, http.MethodPut, "/me/password", changeReq, token)
	assert.Equal(t, http.StatusNoContent, rec.Code)

	// 初回パスワードでログイン試行（失敗するべき）
	loginReq := map[string]interface{}{
		"email":    "pw@example.com",
		"password": initialPassword,
	}
	rec = makeRequest(t, http.MethodPost, "/login", loginReq, "")
	assert.Equal(t, http.StatusUnauthorized, rec.Code)

	// 新しいパスワードでログイン試行（成功するべき）
	loginReq["password"] = "newpassword123"
	rec = makeRequest(t, http.MethodPost, "/login", loginReq, "")
	assert.Equal(t, http.StatusOK, rec.Code)
}

// ========================================
// ヘルパー関数
// ========================================

// createAndLoginUser はユーザー登録・認証・ログインを行いJWTトークンを返す
func createAndLoginUser(t *testing.T, username, email, password string) string {
	signupReq := map[string]interface{}{
		"name":  username,
		"email": email,
	}
	rec := makeRequest(t, http.MethodPost, "/signup", signupReq, "")
	require.Equal(t, http.StatusCreated, rec.Code)

	// メール認証（モックから初回パスワードとトークンを取得）
	verificationToken := mockEmailSender.GetLastToken()
	initialPassword := mockEmailSender.GetLastPassword()
	require.NotEmpty(t, verificationToken)
	require.NotEmpty(t, initialPassword)
	rec = makeRequest(t, http.MethodPost, "/users/verify/"+verificationToken, nil, "")
	require.Equal(t, http.StatusOK, rec.Code)

	// ログイン（初回パスワードを使用）
	loginReq := map[string]interface{}{
		"email":    email,
		"password": initialPassword,
	}
	rec = makeRequest(t, http.MethodPost, "/login", loginReq, "")
	require.Equal(t, http.StatusOK, rec.Code)

	var loginResp map[string]interface{}
	err := json.Unmarshal(rec.Body.Bytes(), &loginResp)
	require.NoError(t, err)
	return loginResp["token"].(string)
}

// createTrip は旅行を作成してIDを返す
func createTrip(t *testing.T, token, title, startDate, endDate string) string {
	tripReq := map[string]interface{}{
		"title":     title,
		"startDate": startDate,
		"endDate":   endDate,
		"members":   []interface{}{},
	}
	rec := makeRequest(t, http.MethodPost, "/trips", tripReq, token)
	require.Equal(t, http.StatusCreated, rec.Code)

	var tripResp map[string]interface{}
	err := json.Unmarshal(rec.Body.Bytes(), &tripResp)
	require.NoError(t, err)
	return tripResp["id"].(string)
}

// createSchedule はスケジュールを作成してIDを返す
func createSchedule(t *testing.T, token string, tripID string, title, date string) string {
	scheduleReq := map[string]interface{}{
		"title":         title,
		"startDateTime": date + "T10:00:00Z",
		"endDateTime":   date + "T12:00:00Z",
		"memo":          "テストスケジュール",
	}
	rec := makeRequest(t, http.MethodPost, fmt.Sprintf("/trips/%s/schedules", tripID), scheduleReq, token)
	require.Equal(t, http.StatusCreated, rec.Code)

	var scheduleResp map[string]interface{}
	err := json.Unmarshal(rec.Body.Bytes(), &scheduleResp)
	require.NoError(t, err)
	return scheduleResp["id"].(string)
}
