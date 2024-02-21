package service

import "github.com/gin-gonic/gin"

// DeleteUser
// @Tags 管理员方法
// @Summary 删除用户
// @Accept multipart/form-data
// @Param file formData file true "File to upload"
// @Param document formData string true "document"
// @Param time formData string true "time"
// @Success 200 {string} json "{"code":"0"}"
// @Router /delete-user [post]
func DeleteUser(c *gin.Context) {

}
