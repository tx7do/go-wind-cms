package service

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/tx7do/go-utils/trans"
	"github.com/tx7do/go-utils/translator"
	"github.com/tx7do/kratos-bootstrap/bootstrap"

	adminV1 "go-wind-cms/api/gen/go/admin/service/v1"
	translatorV1 "go-wind-cms/api/gen/go/translator/service/v1"
)

type TranslatorService struct {
	adminV1.TranslatorServiceHTTPServer

	log *log.Helper

	translator translator.Translator
}

func NewTranslatorService(
	ctx *bootstrap.Context,
	translator translator.Translator,
) *TranslatorService {
	return &TranslatorService{
		log:        ctx.NewLoggerHelper("translator/service/admin-service"),
		translator: translator,
	}
}

func (s *TranslatorService) Translate(_ context.Context, req *translatorV1.TranslateRequest) (*translatorV1.TranslateResponse, error) {
	targetContent, err := s.translator.Translate(req.GetContent(), req.GetSourceLanguage(), req.GetTargetLanguage())
	if err != nil {
		return nil, adminV1.ErrorInternalServerError("翻译失败: %v", err)
	}

	return &translatorV1.TranslateResponse{
		TranslatedContent: trans.Ptr(targetContent),
		RawContent:        req.Content,
	}, nil
}
