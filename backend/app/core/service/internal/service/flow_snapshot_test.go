package service

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	commentV1 "go-wind-cms/api/gen/go/comment/service/v1"
	contentV1 "go-wind-cms/api/gen/go/content/service/v1"
	identityV1 "go-wind-cms/api/gen/go/identity/service/v1"
	mediaV1 "go-wind-cms/api/gen/go/media/service/v1"
)

// ── 事件快照构造（Flow 条件判定的数据源，字段名即条件 field 取值域）──

func TestPostEventSnapshot(t *testing.T) {
	assert.Nil(t, PostEventSnapshot(nil, "DRAFT"), "nil 实体 → nil 快照")

	post := &contentV1.Post{
		Status:      contentV1.Post_POST_STATUS_PUBLISHED.Enum(),
		AuthorId:    intPtr(11),
		CategoryIds: []int64{3, 1, 2},
		Translations: []*contentV1.PostTranslation{
			{Title: strPtr(""), Slug: strPtr("")},
			{Title: strPtr("标题甲"), Slug: strPtr("post-a")},
			{Title: strPtr("标题乙"), Slug: strPtr("post-b")},
		},
		CustomFields: map[string]string{"cover": "media:1", "rating": "5"},
	}

	snap := PostEventSnapshot(post, "DRAFT")
	assert.Equal(t, "POST_STATUS_PUBLISHED", snap["status"])
	assert.Equal(t, "11", snap["author_id"])
	assert.Equal(t, "DRAFT", snap["prev_status"])
	assert.Equal(t, "标题甲", snap["title"], "取首个非空翻译")
	assert.Equal(t, "post-a", snap["slug"])
	assert.Equal(t, "3|1|2", snap["category_id"], "多分类按 | 拼接")
	assert.Equal(t, "media:1", snap["cf.cover"], "自定义字段带 cf. 前缀")
	assert.Equal(t, "5", snap["cf.rating"])
}

func TestPostEventSnapshotNoPrevStatus(t *testing.T) {
	post := &contentV1.Post{Status: contentV1.Post_POST_STATUS_DRAFT.Enum()}
	snap := PostEventSnapshot(post, "")
	_, ok := snap["prev_status"]
	assert.False(t, ok, "prevStatus 为空不写键")
}

func TestPageEventSnapshot(t *testing.T) {
	assert.Nil(t, PageEventSnapshot(nil, ""))

	page := &contentV1.Page{
		Status: contentV1.Page_PAGE_STATUS_PUBLISHED.Enum(),
		Translations: []*contentV1.PageTranslation{
			{Title: strPtr("关于我们"), Slug: strPtr("about")},
		},
		CustomFields: map[string]string{"hero": "media:9"},
	}
	snap := PageEventSnapshot(page, "PAGE_POST_STATUS_DRAFT")
	assert.Equal(t, "PAGE_POST_STATUS_DRAFT", snap["prev_status"])
	assert.Equal(t, "关于我们", snap["title"])
	assert.Equal(t, "about", snap["slug"])
	assert.Equal(t, "media:9", snap["cf.hero"])
}

func TestCommentEventSnapshot(t *testing.T) {
	assert.Nil(t, CommentEventSnapshot(nil))

	comment := &commentV1.Comment{
		Status:     commentV1.Comment_STATUS_PENDING.Enum(),
		Content:    strPtr("评论内容"),
		ObjectId:   int64Ptr(77),
		AuthorId:   int64Ptr(5),
		AuthorType: commentV1.Comment_AUTHOR_TYPE_USER.Enum(),
	}
	snap := CommentEventSnapshot(comment)
	assert.Equal(t, "STATUS_PENDING", snap["status"])
	assert.Equal(t, "评论内容", snap["content"])
	assert.Equal(t, "77", snap["object_id"])
	assert.Equal(t, "5", snap["author_id"])
	assert.Equal(t, "AUTHOR_TYPE_USER", snap["authorType"])
}

func TestCommentEventSnapshotContentTruncation(t *testing.T) {
	// 超长评论截断到 200 字符（按字节），防 payload 膨胀。
	long := make([]rune, 500)
	for i := range long {
		long[i] = '评'
	}
	comment := &commentV1.Comment{Content: strPtr(string(long))}
	snap := CommentEventSnapshot(comment)
	assert.Len(t, snap["content"], flowSnapshotContentLimit)
}

func TestMediaEventSnapshot(t *testing.T) {
	assert.Nil(t, MediaEventSnapshot(nil))

	asset := &mediaV1.MediaAsset{
		Type:     mediaV1.MediaAsset_ASSET_TYPE_IMAGE.Enum(),
		MimeType: strPtr("image/jpeg"),
		Filename: strPtr("cat.jpg"),
		Size:     uint64Ptr(1024),
	}
	snap := MediaEventSnapshot(asset)
	assert.Equal(t, "ASSET_TYPE_IMAGE", snap["asset_type"])
	assert.Equal(t, "image/jpeg", snap["mime_type"])
	assert.Equal(t, "cat.jpg", snap["file_name"])
	assert.Equal(t, "1024", snap["size_bytes"])
}

func TestUserEventSnapshot(t *testing.T) {
	assert.Nil(t, UserEventSnapshot(nil))

	user := &identityV1.User{
		Status:   identityV1.User_NORMAL.Enum(),
		Username: strPtr("alice"),
		Email:    strPtr("alice@example.com"),
	}
	snap := UserEventSnapshot(user)
	assert.Equal(t, "NORMAL", snap["status"])
	assert.Equal(t, "alice", snap["username"])
	assert.Equal(t, "alice@example.com", snap["email"])
}

// ── parseFlowConditions 与条件算子 ──

func TestParseFlowConditions(t *testing.T) {
	conds, err := parseFlowConditions("")
	require.NoError(t, err)
	assert.Nil(t, conds, "空串 = 无条件")

	conds, err = parseFlowConditions("   ")
	require.NoError(t, err)
	assert.Nil(t, conds)

	conds, err = parseFlowConditions(`[{"field":"status","op":"eq","value":"PUBLISHED"}]`)
	require.NoError(t, err)
	require.Len(t, conds, 1)
	assert.Equal(t, "status", conds[0].Field)
	assert.Equal(t, "eq", conds[0].Op)
	assert.Equal(t, "PUBLISHED", conds[0].Value)

	_, err = parseFlowConditions("not-json")
	assert.Error(t, err)
}

func TestFlowConditionOps(t *testing.T) {
	cases := []struct {
		op   string
		got  string
		want string
	}{
		{"eq", "PUBLISHED", "PUBLISHED"},
		{"ne", "DRAFT", "PUBLISHED"},
		{"contains", "go-wind-cms 是一款", "wind"},
		{"in", "PUBLISHED", "PUBLISHED|SCHEDULED"},
		{"in", "SCHEDULED", "PUBLISHED|SCHEDULED"},
		{"in", "PUBLISHED", "PUBLISHED | SCHEDULED"}, // 候选带空格：TrimSpace 后命中
	}
	for _, c := range cases {
		match := flowConditionOps[c.op]
		assert.Equal(t, true, match(c.got, c.want), "%s(%q, %q)", c.op, c.got, c.want)
	}

	negatives := []struct {
		op   string
		got  string
		want string
	}{
		{"eq", "DRAFT", "PUBLISHED"},
		{"ne", "PUBLISHED", "PUBLISHED"},
		{"contains", "other", "wind"},
		{"in", "DRAFT", "PUBLISHED|SCHEDULED"},
	}
	for _, c := range negatives {
		match := flowConditionOps[c.op]
		assert.False(t, match(c.got, c.want), "%s(%q, %q) 应不命中", c.op, c.got, c.want)
	}

	// 未知算子不在表内（validateConditions 据此拒绝非法 op）。
	_, ok := flowConditionOps["regex"]
	assert.False(t, ok)
}

// ── access key 辅助 ──

func TestRandomToken(t *testing.T) {
	ak, err := randomToken("AK", 16)
	require.NoError(t, err)
	assert.True(t, len(ak) > 2)
	assert.True(t, len(ak) >= 2+32, "16 字节 hex = 32 字符")
	assert.Contains(t, ak[:2], "AK")

	sk1, err := randomToken("SK", 24)
	require.NoError(t, err)
	sk2, err := randomToken("SK", 24)
	require.NoError(t, err)
	assert.NotEqual(t, sk1, sk2, "两次生成不得相同")

	// 前缀忠实保留。
	sk, err := randomToken("gwc_sk_", 32)
	require.NoError(t, err)
	assert.True(t, len(sk) >= len("gwc_sk_")+64)
}

func TestHashSecret(t *testing.T) {
	// SHA-256 hex 摘要：SK 只存摘要，secret 不落明文。
	h1 := hashSecret("secret-value")
	assert.Len(t, h1, 64)
	assert.Equal(t, h1, hashSecret("secret-value"))
	assert.NotEqual(t, h1, hashSecret("other-value"))
}

// ── 辅助构造 ──

func intPtr(v int64) *int64      { return &v }
func int64Ptr(v int64) *int64    { return &v }
func uint64Ptr(v uint64) *uint64 { return &v }
