package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type AssetLinksHandler struct {
}

func NewAssetLinks() *AssetLinksHandler {
	return &AssetLinksHandler{}
}

func (handler *AssetLinksHandler) SendManifest(c *gin.Context) {
	c.JSON(http.StatusOK,
		[]gin.H{{
			"relation": []string{"delegate_permission/common.handle_all_urls"},
			"target": gin.H{
				"namespace":                "android_app",
				"package_name":             "com.example",
				"sha256_cert_fingerprints": []string{"14:6D:E9:83:C5:73:06:50:D8:EE:B9:95:2F:34:FC:64:16:A0:83:42:E6:1D:BE:A8:8A:04:96:B2:3F:CF:44:E5"},
			},
		}})
}
