package service

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/tx7do/go-utils/trans"
	"github.com/tx7do/go-utils/translator"
	"github.com/tx7do/kratos-bootstrap/bootstrap"

	adminV1 "go-wind-cms/api/gen/go/admin/service/v1"
	translatorV1 "go-wind-cms/api/gen/go/translator/service/v1"

	"go-wind-cms/pkg/translator/mdprotect"
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
	// 编辑器提交的是 Markdown 正文，先遮蔽代码块/行内代码/URL，
	// 防止引擎翻译过程中破坏代码与超链接结构；纯文本不受影响
	content, restore := mdprotect.Protect(req.GetContent())

	targetContent, err := s.translator.Translate(content, req.GetSourceLanguage(), req.GetTargetLanguage())
	if err != nil {
		s.log.Errorf("translator.Translate err: %+v", err)
		return nil, adminV1.ErrorInternalServerError("翻译失败")
	}

	return &translatorV1.TranslateResponse{
		TranslatedContent: trans.Ptr(restore(targetContent)),
		RawContent:        req.Content,
	}, nil
}
