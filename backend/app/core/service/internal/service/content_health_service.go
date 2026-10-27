package service

import (
	"context"
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	paginationV1 "github.com/tx7do/go-crud/api/gen/go/pagination/v1"
	"github.com/tx7do/go-utils/trans"
	"github.com/tx7do/kratos-bootstrap/bootstrap"

	appViewer "go-wind-cms/pkg/entgo/viewer"

	"go-wind-cms/app/core/service/internal/data"
	"go-wind-cms/pkg/task"
)

// 内容健康度巡检（B8）：每日扫描已发布文章——
//   1. 翻译缺摘要（影响 SEO 与前台卡片）
//   2. 正文外链可达性（HEAD 探测；超时/非 2xx 记疑似死链）
// 问题按作者聚合发站内信（internalMessages 缝）。旁路观测，不改内容数据。

const (
	// contentHealthMaxPosts 单轮最多检查的文章数（防首日冷启动扫全库拖垮 worker）。
	contentHealthMaxPosts = 50
	// contentHealthMaxLinksPerPost 每篇最多探测的外链数。
	contentHealthMaxLinksPerPost = 10
	// contentHealthHTTPTimeout 单链接探测超时。
	contentHealthHTTPTimeout = 5 * time.Second
	// contentHealthMaxIssuesPerNotice 单条站内信最多列的问题行（超出给省略提示）。
	contentHealthMaxIssuesPerNotice = 20
)

// externalLinkPattern 正文里的外链（HTML href 与 Markdown [](...)），scheme 限 http(s)。
var externalLinkPattern = regexp.MustCompile(`(?:href="|href='|\]\()(https?://[^"'\)\s>]+)`)

// ContentHealthService 内容健康度巡检服务。
type ContentHealthService struct {
	postRepo *data.PostRepo

	// internalMessages 站内信缝（构造后注册，解双向引用——与租户水位告警同款）
	internalMessages internalMessageSender

	httpClient *http.Client
	log        *log.Helper
}

func NewContentHealthService(ctx *bootstrap.Context, postRepo *data.PostRepo) *ContentHealthService {
	return &ContentHealthService{
		log:        ctx.NewLoggerHelper("content-health/service/core-service"),
		postRepo:   postRepo,
		httpClient: &http.Client{Timeout: contentHealthHTTPTimeout},
	}
}

// RegisterInternalMessages 注册站内信投递缝。
func (s *ContentHealthService) RegisterInternalMessages(sender internalMessageSender) {
	s.internalMessages = sender
}

// healthIssue 单条问题（作者 ID → 行列表的聚合粒度）。
type healthIssue struct {
	postID uint32
	title  string
	detail string
}

// AsyncContentHealthSweep 周期巡检 worker（asynq 订阅 content.health.sweep）。
func (s *ContentHealthService) AsyncContentHealthSweep(taskType string, taskData *task.ContentHealthSweepTaskData) error {
	ctx := appViewer.NewSystemViewerContext(context.Background())

	// 拉最近发布的已发布文章（publishTime 倒序，钳上限）
	resp, err := s.postRepo.List(ctx, &paginationV1.PagingRequest{
		Page:     trans.Ptr(uint32(1)),
		PageSize: trans.Ptr(uint32(contentHealthMaxPosts)),
		FilteringType: &paginationV1.PagingRequest_Query{
			Query: `{"status":"POST_STATUS_PUBLISHED"}`,
		},
		OrderBy: trans.Ptr("-publishTime"),
	})
	if err != nil {
		s.log.Errorf("content health sweep: list posts failed: %s", err.Error())
		return err // 瞬时查询失败交 asynq 退避重试
	}

	// 作者 → 问题行
	byAuthor := make(map[uint32][]string, 8)
	checkedLinks := 0
	for _, p := range resp.GetItems() {
		author := p.GetAuthorId()
		for _, tr := range p.GetTranslations() {
			if strings.TrimSpace(tr.GetSummary()) == "" {
				byAuthor[author] = append(byAuthor[author], fmt.Sprintf(
					"《%s》(#%d) %s 缺少摘要", firstNonEmptyMcp(tr.GetTitle(), "#"+strconv.FormatUint(uint64(p.GetId()), 10)),
					p.GetId(), tr.GetLanguageCode()))
			}
			// 外链探测（每篇限 10 条；全局预算防长尾）
			for _, link := range extractExternalLinks(tr.GetContent()) {
				if checkedLinks >= contentHealthMaxPosts*contentHealthMaxLinksPerPost/2 {
					break
				}
				checkedLinks++
				if status, lerr := s.probeLink(link); lerr != nil || status >= 400 {
					detail := "无法访问"
					if lerr != nil {
						detail = lerr.Error()
					} else {
						detail = fmt.Sprintf("HTTP %d", status)
					}
					byAuthor[author] = append(byAuthor[author], fmt.Sprintf(
						"《%s》(#%d) 疑似死链 %s（%s）",
						firstNonEmptyMcp(tr.GetTitle(), "#"+strconv.FormatUint(uint64(p.GetId()), 10)),
						p.GetId(), link, detail))
				}
			}
		}
	}

	if len(byAuthor) == 0 {
		s.log.Infof("content health sweep: all clean (%d posts, %d links)", len(resp.GetItems()), checkedLinks)
		return nil
	}

	// 按作者聚合发站内信（通知缝未注册则只留日志）
	for author, issues := range byAuthor {
		if author == 0 || s.internalMessages == nil {
			continue
		}
		lines := issues
		omitted := 0
		if len(lines) > contentHealthMaxIssuesPerNotice {
			omitted = len(lines) - contentHealthMaxIssuesPerNotice
			lines = lines[:contentHealthMaxIssuesPerNotice]
		}
		var b strings.Builder
		b.WriteString(fmt.Sprintf("你的 %d 篇已发布文章存在健康度问题：\n", countDistinctPosts(issues)))
		b.WriteString(strings.Join(lines, "\n"))
		if omitted > 0 {
			b.WriteString(fmt.Sprintf("\n…另有 %d 条问题未列出", omitted))
		}
		if _, nerr := s.internalMessages.SendSystemNotification(ctx,
			"内容健康度提醒", b.String(), []uint32{author}); nerr != nil {
			s.log.Errorf("content health sweep: notify author [%d] failed: %s", author, nerr.Error())
		}
	}
	s.log.Infof("content health sweep: %d post(s) with issues across %d author(s) (%d links checked)",
		countDistinctPosts(flattenIssues(byAuthor)), len(byAuthor), checkedLinks)
	return nil
}

// probeLink HEAD 探测（有些站点不支持 HEAD 时退化为 GET Range 首字节）。
func (s *ContentHealthService) probeLink(link string) (int, error) {
	req, err := http.NewRequest(http.MethodHead, link, nil)
	if err != nil {
		return 0, fmt.Errorf("invalid url")
	}
	resp, err := s.httpClient.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()
	// 405/501 = 不支持 HEAD：补一发 GET（只读状态码，不读体）
	if resp.StatusCode == http.StatusMethodNotAllowed || resp.StatusCode == http.StatusNotImplemented {
		getReq, gerr := http.NewRequest(http.MethodGet, link, nil)
		if gerr != nil {
			return 0, gerr
		}
		getResp, gerr2 := s.httpClient.Do(getReq)
		if gerr2 != nil {
			return 0, gerr2
		}
		defer getResp.Body.Close()
		return getResp.StatusCode, nil
	}
	return resp.StatusCode, nil
}

// extractExternalLinks 抽正文外链（去重保序）。
func extractExternalLinks(content string) []string {
	if content == "" {
		return nil
	}
	matches := externalLinkPattern.FindAllStringSubmatch(content, -1)
	seen := map[string]bool{}
	out := make([]string, 0, len(matches))
	for _, m := range matches {
		link := strings.TrimRight(m[1], ".,;:!?)")
		if link == "" || seen[link] {
			continue
		}
		seen[link] = true
		out = append(out, link)
	}
	return out
}

func firstNonEmptyMcp(a, b string) string {
	if a != "" {
		return a
	}
	return b
}

func countDistinctPosts(issueLines []string) int {
	set := map[string]bool{}
	for _, line := range issueLines {
		if idx := strings.Index(line, "(#"); idx >= 0 {
			if end := strings.Index(line[idx:], ")"); end > 0 {
				set[line[idx+1:idx+end]] = true
			}
		}
	}
	return len(set)
}

func flattenIssues(byAuthor map[uint32][]string) []string {
	var out []string
	for _, lines := range byAuthor {
		out = append(out, lines...)
	}
	return out
}
