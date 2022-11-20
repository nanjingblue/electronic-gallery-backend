package service

import (
	"electronic-album/global"
	"electronic-album/internal/serializer"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/google/uuid"
	"mime"
	"path/filepath"
)

type UploadTokenService struct {
	Filename string `form:"filename" json:"filename" binding:"required"`
}

// GetToken 创建token
func (svc *Service) GetToken(param *UploadTokenService) serializer.Response {
	client, err := oss.New(global.OSSSetting.END_POINT, global.OSSSetting.ACCESS_KEY_ID, global.OSSSetting.ACCESS_KEY_SECRET)
	if err != nil {
		return serializer.Response{
			Code:  500,
			Msg:   "OSS配置错误",
			Error: err.Error(),
		}
	}

	// 获取存储空间。
	bucket, err := client.Bucket(global.OSSSetting.BUCKET)
	if err != nil {
		return serializer.Response{
			Code:  500,
			Msg:   "OSS配置错误",
			Error: err.Error(),
		}
	}

	// 获取扩展名
	ext := filepath.Ext(param.Filename)

	// 带可选参数的签名直传。
	options := []oss.Option{
		oss.ContentType(mime.TypeByExtension(ext)),
	}

	key := "upload/cover/" + uuid.Must(uuid.NewRandom()).String() + ext
	// 签名直传。
	signedPutURL, err := bucket.SignURL(key, oss.HTTPPut, 600, options...)
	if err != nil {
		return serializer.Response{
			Code:  500,
			Msg:   "OSS配置错误",
			Error: err.Error(),
		}
	}
	// 查看图片
	signedGetURL, err := bucket.SignURL(key, oss.HTTPGet, 600)
	if err != nil {
		return serializer.Response{
			Code:  500,
			Msg:   "OSS配置错误",
			Error: err.Error(),
		}
	}

	return serializer.Response{
		Code: 200,
		Data: map[string]string{
			"key": key,
			"put": signedPutURL,
			"get": signedGetURL,
		},
		Msg: "获取上传cover token 成功",
	}
}