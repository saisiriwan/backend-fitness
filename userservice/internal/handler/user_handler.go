package handler

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"strconv" // (เพิ่ม import นี้ สำหรับ GetUserByID)

	"users/internal/models"
	"users/internal/service"

	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

// (ตัวแปร global สำหรับเก็บ Config)
var (
	googleOauthConfig *oauth2.Config
	// (ย้าย JWT_SECRET มาไว้ที่นี่ และอ่านจาก .env)
	JWT_SECRET []byte
)

// init() จะทำงาน *ก่อน* main() เสมอ
func init() {
	// 1. ตั้งค่า Google OAuth
	clientID := os.Getenv("GOOGLE_CLIENT_ID")
	clientSecret := os.Getenv("GOOGLE_CLIENT_SECRET")
	redirectURL := os.Getenv("GOOGLE_REDIRECT_URL")

	if clientID == "" || clientSecret == "" {
		log.Println("WARNING: GOOGLE_CLIENT_ID or GOOGLE_CLIENT_SECRET is not set. Google OAuth will not work.")
	}

	googleOauthConfig = &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURL:  redirectURL,
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email", "https://www.googleapis.com/auth/userinfo.profile"},
		Endpoint:     google.Endpoint,
	}

	// 2. ตั้งค่า JWT Secret (จาก .env)
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		log.Println("WARNING: JWT_SECRET is not set. Using default (unsafe) secret.")
		JWT_SECRET = []byte("your_very_secret_key_should_be_long")
	} else {
		JWT_SECRET = []byte(secret)
	}
}

type UserHandler struct {
	userService service.UserService
}

func NewUserHandler(us service.UserService) *UserHandler {
	return &UserHandler{userService: us}
}

// ----------------------------------------------------
// (แก้ไข) ฟังก์ชัน CRUD เดิม (เติม Logic ให้สมบูรณ์)
// ----------------------------------------------------

func (h *UserHandler) GetAllUsers(c *gin.Context) {
	users, err := h.userService.GetAllUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get users"})
		return
	}
	c.JSON(http.StatusOK, users)
}

func (h *UserHandler) GetUserByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	user, err := h.userService.GetUserByID(id)
	if err != nil {
		if err.Error() == "user not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user"})
		}
		return
	}
	c.JSON(http.StatusOK, user)
}

func (h *UserHandler) CreateUser(c *gin.Context) {
	var req struct {
		Name  string `json:"name" binding:"required"`
		Email string `json:"email" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	// (หมายเหตุ: ฟังก์ชัน CreateUser เก่านี้ ไม่มีการ Hash Password)
	// (เราควรใช้ RegisterUser แทน)
	user, err := h.userService.CreateUser(req.Name, req.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, user)
}

func (h *UserHandler) UpdateUser(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	var req struct {
		Name  string `json:"name" binding:"required"`
		Email string `json:"email" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	user, err := h.userService.UpdateUser(id, req.Name, req.Email)
	if err != nil {
		if err.Error() == "user not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}
	c.JSON(http.StatusOK, user)
}

func (h *UserHandler) DeleteUser(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	err = h.userService.DeleteUser(id)
	if err != nil {
		if err.Error() == "user not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user"})
		}
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}

// ----------------------------------------------------
// (ฟังก์ชัน Auth ที่มี Logic จริง)
// ----------------------------------------------------

func (h *UserHandler) Register(c *gin.Context) {
	var req service.RegisterRequest // (ใช้ service.RegisterRequest)
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	// เรียก Service (Logic การ Hash อยู่ใน Service แล้ว)
	user, err := h.userService.RegisterUser(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, user) // ส่ง User ที่สร้างเสร็จกลับไป
}

func (h *UserHandler) Login(c *gin.Context) {
	var req service.LoginRequest // (ใช้ service.LoginRequest)
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	// +++ เพิ่มส่วนนี้เพื่อความปลอดภัยสูงสุด +++
	// ป้องกันกรณีคนพยายาม Login บัญชี Google โดยไม่ใส่รหัสผ่าน
	if req.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Password is required"})
		return
	}
	// ++++++++++++++++++++++++++++++++++++++

	// เรียก Service
	accessToken, err := h.userService.LoginUser(req)
	if err != nil {
		// (ปรับ Error Message ให้ผู้ใช้เข้าใจง่ายขึ้น ไม่ควรส่ง raw error)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	// ตั้งค่า httpOnly Cookie (จากเอกสาร Auth Part 1)
	c.SetCookie("access_token", accessToken, 15*60, "/", "localhost", false, true) // httpOnly=true

	c.JSON(http.StatusOK, gin.H{"message": "Login successful"})
}

func (h *UserHandler) Logout(c *gin.Context) {
	// ล้าง Cookie
	c.SetCookie("access_token", "", -1, "/", "localhost", false, true)
	c.JSON(http.StatusOK, gin.H{"message": "Logout successful"})
}

// (Logic ให้ GoogleLogin)
func (h *UserHandler) GoogleLogin(c *gin.Context) {
	state := "random-state-string-for-csrf-protection" // (ในระบบจริงควรสุ่มค่านี้)
	url := googleOauthConfig.AuthCodeURL(state)

	// สั่ง Redirect (เด้ง) เบราว์เซอร์ของผู้ใช้ไปหน้า Login ของ Google
	c.Redirect(http.StatusTemporaryRedirect, url)
}

// (Logic ให้ GoogleCallback)
func (h *UserHandler) GoogleCallback(c *gin.Context) {
	// 1. รับ "code" ที่ Google ส่งกลับมา
	code := c.Query("code")

	// 2. นำ "code" ไปแลกเป็น "Google Token"
	token, err := googleOauthConfig.Exchange(context.Background(), code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to exchange token"})
		return
	}

	// 3. นำ "Google Token" ไปขอข้อมูลโปรไฟล์ผู้ใช้
	response, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token.AccessToken)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user info"})
		return
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read user info"})
		return
	}

	// 4. (Flow 3) ค้นหา หรือ สร้าง User
	var googleUser struct {
		Email string `json:"email"`
		Name  string `json:"name"`
	}
	json.Unmarshal(body, &googleUser)

	// (a) ค้นหา User ด้วย Email
	user, err := h.userService.GetUserByEmail(googleUser.Email)
	if err != nil {
		// (b) ถ้า "ไม่เจอ" (Sign Up) -> สร้าง User ใหม่
		if err.Error() == "user not found" {
			// สร้าง User ใหม่ด้วย Role "client"
			newUser := models.User{
				Name:  googleUser.Name,
				Email: googleUser.Email,
				Role:  "trainer",
			}

			// (แก้ไข) ใช้ Password สั้นๆ ที่แน่นอน
			createdUser, err := h.userService.RegisterUser(service.RegisterRequest{
				FirstName: newUser.Name,
				LastName:  "", // (Google อาจจะไม่ได้แยกชื่อมาให้)
				Email:     newUser.Email,
				Password:  "google_user_placeholder_password",
			})
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create new user via Google"})
				return
			}
			user = createdUser
		} else {
			// ถ้า Error อื่น (เช่น DB ล่ม)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error on user lookup"})
			return
		}
	}

	// (c) ถ้า "เจอ" (Log In) -> เราได้ `user` มาแล้ว

	// 5. สร้าง JWT Token และตั้ง Cookie (เหมือน `Login` Handler)
	// (เราได้ปรับ Logic ใน service.LoginUser ให้รองรับ Password ว่างเปล่าแล้ว)
	accessToken, err := h.userService.LoginUser(service.LoginRequest{
		Email:    user.Email,
		Password: "", // (ส่ง Password ว่างเปล่าไป)
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to log in user after Google auth"})
		return
	}

	// 6. ตั้งค่า httpOnly Cookie
	c.SetCookie("access_token", accessToken, 15*60, "/", "localhost", false, true)

	// 7. (สำคัญ) Redirect กลับไปหน้า Frontend
	c.Redirect(http.StatusTemporaryRedirect, "http://localhost:3000/dashboard")
}

func (h *UserHandler) CheckAuth(c *gin.Context) {
	// (เราสามารถดึงข้อมูลที่ "ยาม" ส่งมาให้ได้)
	userID, _ := c.Get("user_id")
	userRole, _ := c.Get("role")

	c.JSON(http.StatusOK, gin.H{
		"message": "You are authenticated!",
		"user_id": userID,
		"role":    userRole,
	})
}
