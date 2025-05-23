package utils

import (
	"github.com/Digitalkeun-Creative/be-dzikra-pos-service/internal/module/product/dto"
	"github.com/microcosm-cc/bluemonday"
)

// SanitizeCreateOrUpdateArticleResponse sanitizes the CreateOrUpdateArticleResponse by removing any potentially harmful content
func SanitizeCreateOrUpdateProductResponse(resp dto.CreateOrUpdateProductResponse, policy *bluemonday.Policy) dto.CreateOrUpdateProductResponse {
	resp.Name = policy.Sanitize(resp.Name)
	resp.Description = policy.Sanitize(resp.Description)
	return resp
}
