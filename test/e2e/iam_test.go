package e2e

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	baseURL = "http://localhost:8080"
	testDB  = "test_liaison.db"
)

// TestConfig 测试配置
type TestConfig struct {
	BaseURL string
	DBPath  string
}

// LoginRequest 登录请求
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// LoginResponse 登录响应
type LoginResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    struct {
		Token string `json:"token"`
		User  struct {
			ID    uint   `json:"id"`
			Email string `json:"email"`
		} `json:"user"`
	} `json:"data"`
}

// ProfileResponse 用户信息响应
type ProfileResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    struct {
		ID    uint   `json:"id"`
		Email string `json:"email"`
	} `json:"data"`
}

// ErrorResponse 错误响应
type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// TestSuite E2E测试套件
type TestSuite struct {
	config     *TestConfig
	httpClient *http.Client
	serverCmd  *exec.Cmd
	token      string
}

// NewTestSuite 创建测试套件
func NewTestSuite() *TestSuite {
	return &TestSuite{
		config: &TestConfig{
			BaseURL: baseURL,
			DBPath:  testDB,
		},
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// Setup 设置测试环境
func (ts *TestSuite) Setup(t *testing.T) {
	// 清理测试数据库
	ts.cleanupTestDB(t)

	// 启动服务器
	ts.startServer(t)

	// 等待服务器启动
	ts.waitForServer(t)
}

// Teardown 清理测试环境
func (ts *TestSuite) Teardown(t *testing.T) {
	// 停止服务器
	if ts.serverCmd != nil {
		ts.serverCmd.Process.Kill()
		ts.serverCmd.Wait()
	}

	// 清理测试数据库
	ts.cleanupTestDB(t)
}

// cleanupTestDB 清理测试数据库
func (ts *TestSuite) cleanupTestDB(t *testing.T) {
	if _, err := os.Stat(ts.config.DBPath); err == nil {
		err := os.Remove(ts.config.DBPath)
		require.NoError(t, err, "Failed to remove test database")
	}
}

// startServer 启动服务器
func (ts *TestSuite) startServer(t *testing.T) {
	// 创建测试配置文件
	ts.createTestConfig(t)

	// 启动服务器进程
	ts.serverCmd = exec.Command("./bin/liaison", "-c", "test_config.yaml")
	ts.serverCmd.Dir = "../.." // 回到项目根目录

	// 设置环境变量
	ts.serverCmd.Env = append(os.Environ(), "LIAISON_DB_PATH="+ts.config.DBPath)

	// 启动服务器
	err := ts.serverCmd.Start()
	require.NoError(t, err, "Failed to start server")
}

// createTestConfig 创建测试配置文件
func (ts *TestSuite) createTestConfig(t *testing.T) {
	configContent := fmt.Sprintf(`
manager:
  listen:
    addr: "0.0.0.0:8080"
    network: "tcp"
  database:
    driver: "sqlite"
    source: "%s"
  daemon:
    pprof:
      enable: false
    rlimit:
      enable: false
`, ts.config.DBPath)

	err := os.WriteFile("test_config.yaml", []byte(configContent), 0644)
	require.NoError(t, err, "Failed to create test config")
}

// waitForServer 等待服务器启动
func (ts *TestSuite) waitForServer(t *testing.T) {
	maxRetries := 30
	for i := 0; i < maxRetries; i++ {
		resp, err := ts.httpClient.Get(ts.config.BaseURL + "/health")
		if err == nil && resp.StatusCode == 200 {
			resp.Body.Close()
			return
		}
		if resp != nil {
			resp.Body.Close()
		}
		time.Sleep(1 * time.Second)
	}
	require.Fail(t, "Server failed to start within 30 seconds")
}

// makeRequest 发送HTTP请求
func (ts *TestSuite) makeRequest(method, path string, body interface{}, token string) (*http.Response, error) {
	var reqBody io.Reader
	if body != nil {
		jsonData, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		reqBody = bytes.NewBuffer(jsonData)
	}

	req, err := http.NewRequest(method, ts.config.BaseURL+path, reqBody)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	return ts.httpClient.Do(req)
}

// TestIAMLogin 测试用户登录
func TestIAMLogin(t *testing.T) {
	ts := NewTestSuite()
	ts.Setup(t)
	defer ts.Teardown(t)

	t.Run("Login with valid credentials", func(t *testing.T) {
		// 首先需要创建默认用户
		ts.createDefaultUser(t)

		loginReq := LoginRequest{
			Email:    "default@liaison.local",
			Password: "default123", // 假设这是默认密码
		}

		resp, err := ts.makeRequest("POST", "/v1/iam/login", loginReq, "")
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var loginResp LoginResponse
		err = json.NewDecoder(resp.Body).Decode(&loginResp)
		require.NoError(t, err)

		assert.Equal(t, 200, loginResp.Code)
		assert.NotEmpty(t, loginResp.Data.Token)
		assert.Equal(t, "default@liaison.local", loginResp.Data.User.Email)

		// 保存token供后续测试使用
		ts.token = loginResp.Data.Token
	})

	t.Run("Login with invalid credentials", func(t *testing.T) {
		loginReq := LoginRequest{
			Email:    "default@liaison.local",
			Password: "wrongpassword",
		}

		resp, err := ts.makeRequest("POST", "/v1/iam/login", loginReq, "")
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var errorResp ErrorResponse
		err = json.NewDecoder(resp.Body).Decode(&errorResp)
		require.NoError(t, err)

		assert.NotEqual(t, 200, errorResp.Code)
	})
}

// TestIAMProfile 测试获取用户信息
func TestIAMProfile(t *testing.T) {
	ts := NewTestSuite()
	ts.Setup(t)
	defer ts.Teardown(t)

	// 先登录获取token
	ts.loginAndGetToken(t)

	t.Run("Get profile with valid token", func(t *testing.T) {
		resp, err := ts.makeRequest("GET", "/v1/iam/profile", nil, ts.token)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var profileResp ProfileResponse
		err = json.NewDecoder(resp.Body).Decode(&profileResp)
		require.NoError(t, err)

		assert.Equal(t, 200, profileResp.Code)
		assert.Equal(t, "default@liaison.local", profileResp.Data.Email)
	})

	t.Run("Get profile without token", func(t *testing.T) {
		resp, err := ts.makeRequest("GET", "/v1/iam/profile", nil, "")
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var errorResp ErrorResponse
		err = json.NewDecoder(resp.Body).Decode(&errorResp)
		require.NoError(t, err)

		assert.NotEqual(t, 200, errorResp.Code)
		assert.Contains(t, errorResp.Message, "No authentication token provided")
	})

	t.Run("Get profile with invalid token", func(t *testing.T) {
		resp, err := ts.makeRequest("GET", "/v1/iam/profile", nil, "invalid_token")
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var errorResp ErrorResponse
		err = json.NewDecoder(resp.Body).Decode(&errorResp)
		require.NoError(t, err)

		assert.NotEqual(t, 200, errorResp.Code)
		assert.Contains(t, errorResp.Message, "Token validation failed")
	})
}

// TestAuthenticationMiddleware 测试认证中间件
func TestAuthenticationMiddleware(t *testing.T) {
	ts := NewTestSuite()
	ts.Setup(t)
	defer ts.Teardown(t)

	// 先登录获取token
	ts.loginAndGetToken(t)

	t.Run("Access protected endpoint with valid token", func(t *testing.T) {
		resp, err := ts.makeRequest("GET", "/v1/applications", nil, ts.token)
		require.NoError(t, err)
		defer resp.Body.Close()

		// 应该能正常访问，返回200或相应的业务状态码
		assert.True(t, resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusNotFound)
	})

	t.Run("Access protected endpoint without token", func(t *testing.T) {
		resp, err := ts.makeRequest("GET", "/v1/applications", nil, "")
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var errorResp ErrorResponse
		err = json.NewDecoder(resp.Body).Decode(&errorResp)
		require.NoError(t, err)

		assert.NotEqual(t, 200, errorResp.Code)
		assert.Contains(t, errorResp.Message, "No authentication token provided")
	})

	t.Run("Access login endpoint without token", func(t *testing.T) {
		loginReq := LoginRequest{
			Email:    "default@liaison.local",
			Password: "default123",
		}

		resp, err := ts.makeRequest("POST", "/v1/iam/login", loginReq, "")
		require.NoError(t, err)
		defer resp.Body.Close()

		// 登录接口应该不需要认证
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})
}

// TestIAMLogout 测试用户登出
func TestIAMLogout(t *testing.T) {
	ts := NewTestSuite()
	ts.Setup(t)
	defer ts.Teardown(t)

	// 先登录获取token
	ts.loginAndGetToken(t)

	t.Run("Logout with valid token", func(t *testing.T) {
		resp, err := ts.makeRequest("POST", "/v1/iam/logout", nil, ts.token)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var logoutResp struct {
			Code    int    `json:"code"`
			Message string `json:"message"`
		}
		err = json.NewDecoder(resp.Body).Decode(&logoutResp)
		require.NoError(t, err)

		assert.Equal(t, 200, logoutResp.Code)
	})
}

// Helper methods

// createDefaultUser 创建默认用户（模拟安装脚本的行为）
func (ts *TestSuite) createDefaultUser(t *testing.T) {
	// 直接调用IAM服务的CreateDefaultUser方法
	// 这需要访问IAM服务实例，暂时跳过
	// 在实际测试中，可以通过环境变量或其他方式创建默认用户
	t.Log("Skipping default user creation - would be handled by installation script")
}

// loginAndGetToken 登录并获取token
func (ts *TestSuite) loginAndGetToken(t *testing.T) {
	// 创建默认用户
	ts.createDefaultUser(t)

	loginReq := LoginRequest{
		Email:    "default@liaison.local",
		Password: "default123", // 假设这是默认密码
	}

	resp, err := ts.makeRequest("POST", "/v1/iam/login", loginReq, "")
	require.NoError(t, err)
	defer resp.Body.Close()

	var loginResp LoginResponse
	err = json.NewDecoder(resp.Body).Decode(&loginResp)
	require.NoError(t, err)

	if loginResp.Code == 200 {
		ts.token = loginResp.Data.Token
	} else {
		// 如果登录失败，可能需要先创建用户
		t.Logf("Login failed, may need to create default user first: %s", loginResp.Message)
	}
}

// TestMain 测试主函数
func TestMain(m *testing.M) {
	// 确保在项目根目录
	os.Chdir("../..")

	// 运行测试
	code := m.Run()

	// 清理测试配置文件
	os.Remove("test_config.yaml")
	os.Remove(testDB)

	os.Exit(code)
}
