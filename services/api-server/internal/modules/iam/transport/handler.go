package transport

import (
	"net/http"

	"github.com/gin-gonic/gin"
	iamservice "github.com/liangluo/weouc2026/services/api-server/internal/modules/iam/service"
	iamtypes "github.com/liangluo/weouc2026/services/api-server/internal/modules/iam/types"
	"github.com/liangluo/weouc2026/services/api-server/internal/platform/auth"
	"github.com/liangluo/weouc2026/services/api-server/internal/platform/httpx"
)

type Handler struct {
	service *iamservice.Service
}

func NewHandler(service *iamservice.Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) LoginWithPassword(c *gin.Context) {
	var request iamtypes.AdminLoginRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		httpx.AbortWithError(c, httpx.BadRequest("登录参数格式错误", nil))
		return
	}

	response, err := h.service.LoginWithPassword(c.Request.Context(), request)
	if err != nil {
		httpx.AbortWithError(c, err)
		return
	}

	httpx.JSON(c, http.StatusOK, response)
}

func (h *Handler) LoginWithWeChat(c *gin.Context) {
	var request iamtypes.WeChatLoginRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		httpx.AbortWithError(c, httpx.BadRequest("登录参数格式错误", nil))
		return
	}

	response, err := h.service.LoginWithWeChat(c.Request.Context(), request)
	if err != nil {
		httpx.AbortWithError(c, err)
		return
	}

	httpx.JSON(c, http.StatusOK, response)
}

func (h *Handler) GetStudentProfile(c *gin.Context) {
	response, err := h.service.GetStudentProfile(c.Request.Context(), auth.PrincipalFromContext(c))
	if err != nil {
		httpx.AbortWithError(c, err)
		return
	}

	httpx.JSON(c, http.StatusOK, response)
}

func (h *Handler) SendCaptcha(c *gin.Context) {
	var request iamtypes.SendCaptchaRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		httpx.AbortWithError(c, httpx.BadRequest("验证码参数格式错误", nil))
		return
	}

	if err := h.service.SendCaptcha(c.Request.Context(), auth.PrincipalFromContext(c), request.StudentID); err != nil {
		httpx.AbortWithError(c, err)
		return
	}

	httpx.JSON(c, http.StatusOK, gin.H{"message": "验证码已发送"})
}

func (h *Handler) BindStudent(c *gin.Context) {
	var request iamtypes.BindStudentRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		httpx.AbortWithError(c, httpx.BadRequest("绑定参数格式错误", nil))
		return
	}

	response, err := h.service.BindStudent(c.Request.Context(), auth.PrincipalFromContext(c), request)
	if err != nil {
		httpx.AbortWithError(c, err)
		return
	}

	httpx.JSON(c, http.StatusOK, response)
}

func (h *Handler) UpdateStudent(c *gin.Context) {
	var request iamtypes.UpdateStudentRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		httpx.AbortWithError(c, httpx.BadRequest("更新参数格式错误", nil))
		return
	}
	if request.IsBound == nil || *request.IsBound {
		httpx.AbortWithError(c, httpx.BadRequest("当前仅支持提交 is_bound=false 进行解绑", nil))
		return
	}

	if err := h.service.UnbindStudent(c.Request.Context(), auth.PrincipalFromContext(c)); err != nil {
		httpx.AbortWithError(c, err)
		return
	}

	httpx.JSON(c, http.StatusOK, gin.H{"is_bound": false})
}
