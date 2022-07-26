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
				"package_name":             "com.romanillos.pack",
				"sha256_cert_fingerprints": []string{"AD:F4:9F:F3:B4:EB:23:50:24:3E:F0:AA:8A:E3:8F:EE:1A:79:39:33:36:28:16:4B:63:D4:9A:F3:DF:8F:72:5A"},
			},
		}})
}
