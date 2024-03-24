package service

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// @Tags 公共方法
// @Summary 下载指定路径的文件
// @Accept multipart/form-data
//
//	@Param Authorization header string false "Authentication header"
//
// @Param filePath formData string true "filePath"
// @Success 200 {string} json "{"code":"0"}"
// @Router /download-file [post]
func DownloadFile(c *gin.Context) {
	// todo 所有文件都能请求，有点危险啊，
	// 应该检查参数，保证只能请求files目录下的文件,并且该目录下的文件应该全部不重要
	filePath := c.PostForm("filePath")
	if filePath == "" {
		c.JSON(http.StatusOK, gin.H{
			"code": 1,
			"msg":  "The path to the file cannot be empty",
		})
		return
	}
	c.File(filePath)
}
