package collab

import (
	"context"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/reearth/ygo/crdt"
	"github.com/reearth/ygo/provider/websocket"
	"github.com/tx7do/kratos-bootstrap/bootstrap"

	authenticationV1 "go-wind-cms/api/gen/go/authentication/service/v1"
	contentV1 "go-wind-cms/api/gen/go/content/service/v1"
)

// Service 协同编辑服务。
type Service struct {
	authClient  authenticationV1.AuthenticationServiceClient
	docStateCli contentV1.CollabDocStateServiceClient
	postClient  contentV1.PostServiceClient
	log         *log.Helper
	ws          *websocket.Server
	roomsMu     sync.Mutex
	liveRooms   map[string]struct{} // 本进程已加载的房间名（供按前缀踢房）
}

// New 构造协同服务：upgrade 前经 ValidateToken 准入；传入 docStateClient 时
// 房间状态持久化到 core（重启可恢复）；传入 postClient 时 post 房间由服务端
// 从内容权威种子化（HTML→yjs 受控转换，仅当恢复后文档仍为空）。
func New(ctx *bootstrap.Context, authClient authenticationV1.AuthenticationServiceClient, docStateClient contentV1.CollabDocStateServiceClient, postClient contentV1.PostServiceClient, allowedOrigins []string) *Service {
	s := &Service{
		authClient:  authClient,
		docStateCli: docStateClient,
		postClient:  postClient,
		log:         ctx.NewLoggerHelper("collab/admin-service"),
		liveRooms:   make(map[string]struct{}),
	}

	var srv *websocket.Server
	if docStateClient != nil {
		srv = websocket.NewServerWithPersistence(newGRPCPersistence(docStateClient))
	} else {
		s.log.Warnf("collab doc state client is nil: room persistence DISABLED (rooms are memory-only and lost on restart)")
		srv = websocket.NewServer()
	}
	// 原生 y-websocket 客户端帧（非 Hocuspocus docName 前缀帧）；
	// 房间名取自 URL 路径末段（如 /post:12:zh-CN）。
	srv.AuthFunc = s.authenticate
	srv.AllowedOrigins = allowedOrigins
	srv.OnLoadDocument = func(_ context.Context, room string, doc *crdt.Doc) error {
		s.trackRoom(room, true)
		s.seedRoom(room, doc)
		return nil
	}
	srv.OnUnloadDocument = func(_ context.Context, room string) {
		s.trackRoom(room, false)
	}
	s.ws = srv
	return s
}

// Handler 返回 ws 挂载点（http.Handler）。
func (s *Service) Handler() http.Handler {
	return s.ws
}

// Shutdown 优雅停机（等待在途持久化写）。
func (s *Service) Shutdown(ctx context.Context) error {
	return s.ws.Shutdown(ctx)
}

// seedRoom 服务端权威种子化：post 房间在持久化恢复后仍为空时，
// 从 core 当前内容构造初始文档（HTML→yjs 受控转换）。
// page 房间不种子化（页面是 sections 结构，无单一 HTML 正文）。
// 客户端种子化因覆盖恢复态已废弃（localStorage 锁保留作历史客户端防御）。
func (s *Service) seedRoom(room string, doc *crdt.Doc) {
	targetType, targetID, locale, ok := parseRoom(room)
	if !ok || targetType != "post" {
		s.log.Infof("seed: skip room %s (parsed=%v type=%s)", room, ok, targetType)
		return
	}
	if s.postClient == nil {
		s.log.Warnf("seed: post client nil, skip %s", room)
		return
	}
	if doc.GetXmlFragment("default").Len() > 0 {
		return
	}
	if doc.GetXmlFragment("default").Len() > 0 {
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	post, err := s.postClient.Get(ctx, &contentV1.GetPostRequest{
		QueryBy: &contentV1.GetPostRequest_Id{Id: targetID},
	})
	if err != nil || post == nil {
		s.log.Warnf("seed: load post %d failed: %v", targetID, err)
		return
	}
	content := pickTranslationContent(post.GetTranslations(), locale)
	if strings.TrimSpace(content) == "" {
		s.log.Infof("seed: post %d has no content for locale %s", targetID, locale)
		return
	}

	// MARKDOWN 编辑器：内容即 markdown 串，种子化为 YText 'content'
	//（前端 MarkdownEditor 经 yCollab 绑定 YText 'content'，与富文本的
	// XmlFragment 'default' 是不同表示，按编辑器类型分流）。
	if post.GetEditorType() == contentV1.EditorType_EDITOR_TYPE_MARKDOWN {
		ytext := doc.GetText("content")
		// OnLoadDocument 每房间仅触发一次，无需判空重入
		if err := doc.TransactE(func(txn *crdt.Transaction) error {
			ytext.Insert(txn, 0, content, nil)
			return nil
		}); err != nil {
			s.log.Errorf("seed: ytext insert failed (room=%s): %v", room, err)
			return
		}
		s.log.Infof("collab markdown room seeded: %s (%d chars)", room, len(content))
		return
	}

	if err := seedDocFromHTML(doc, content); err != nil {
		s.log.Errorf("seed: convert html failed (room=%s): %v", room, err)
		return
	}
	s.log.Infof("collab room seeded from content: %s (%d chars)", room, len(content))
}

// parseRoom 解析房间名 post:{id}:{locale} / page:{id}:{locale}。
func parseRoom(room string) (targetType string, targetID uint32, locale string, ok bool) {
	parts := strings.SplitN(room, ":", 3)
	if len(parts) != 3 {
		return "", 0, "", false
	}
	id, err := strconv.ParseUint(parts[1], 10, 32)
	if err != nil || id == 0 {
		return "", 0, "", false
	}
	return parts[0], uint32(id), parts[2], true
}

// pickTranslationContent 挑选语言内容：精确匹配 locale，回退第一条。
func pickTranslationContent(translations []*contentV1.PostTranslation, locale string) string {
	if len(translations) == 0 {
		return ""
	}
	if locale != "" {
		for _, tr := range translations {
			if tr.GetLanguageCode() == locale {
				return tr.GetContent()
			}
		}
	}
	return translations[0].GetContent()
}

// trackRoom 记录/移除本进程活跃房间（供按前缀踢房）。
func (s *Service) trackRoom(room string, loaded bool) {
	s.roomsMu.Lock()
	defer s.roomsMu.Unlock()
	if loaded {
		s.liveRooms[room] = struct{}{}
	} else {
		delete(s.liveRooms, room)
	}
}

// CloseRoomsByPrefix 踢掉匹配前缀的活跃房间（force 断开所有连接；
// ygo 在 evict 前会 flush 待写状态——调用方应在踢房之后再清持久化行，
// 并预留短暂的 worker flush 窗口）。
func (s *Service) CloseRoomsByPrefix(prefix string) {
	s.roomsMu.Lock()
	s.log.Infof("CloseRoomsByPrefix: prefix=%s liveRooms=%d svc=%p", prefix, len(s.liveRooms), s)
	
	var matches []string
	for room := range s.liveRooms {
		if strings.HasPrefix(room, prefix) {
			matches = append(matches, room)
		}
	}
	s.roomsMu.Unlock()
	for _, room := range matches {
		if err := s.ws.CloseRoom(room, true); err != nil {
			s.log.Warnf("close collab room failed (room=%s): %v", room, err)
		} else {
			s.log.Infof("collab room closed: %s", room)
		}
	}
}

// authenticate ws upgrade 准入：query token → core ValidateToken（admin 端）。
// 浏览器 WebSocket 无法带 Authorization 头，token 走 query 参数；
// 内容级细粒度权限由保存链路（core Update）兜底，准入门槛与 SSE 同级。
func (s *Service) authenticate(r *http.Request) bool {
	token := r.URL.Query().Get("token")
	if token == "" {
		return false
	}
	resp, err := s.authClient.ValidateToken(context.Background(), &authenticationV1.ValidateTokenRequest{
		ClientType:    authenticationV1.ClientType_admin,
		Token:         token,
		TokenCategory: authenticationV1.TokenCategory_ACCESS,
	})
	if err != nil {
		s.log.Errorf("collab token authentication failed: %v", err)
		return false
	}
	if resp.GetIsBlocked() {
		s.log.Warnf("collab token is blocked")
		return false
	}
	return resp.GetIsValid()
}

// ParseAllowedOrigins 解析逗号分隔的允许来源列表（环境变量 COLLAB_ALLOWED_ORIGINS）。
func ParseAllowedOrigins(raw string) []string {
	var origins []string
	for _, o := range strings.Split(raw, ",") {
		if o = strings.TrimSpace(o); o != "" {
			origins = append(origins, o)
		}
	}
	return origins
}
