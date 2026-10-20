package service

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCheckQuotaNotExceeded(t *testing.T) {
	// 计数型（USER_LIMIT / API_CALL）：used >= limit 即超
	require.Error(t, checkQuotaNotExceeded("USER_LIMIT", 10, 10, 0), "达到上限即拒绝")
	require.Error(t, checkQuotaNotExceeded("API_CALL", 5, 9, 0), "超过上限拒绝")
	require.NoError(t, checkQuotaNotExceeded("USER_LIMIT", 10, 9, 0), "未达上限放行")

	// 容量型（STORAGE）：used+本次上传 > limit 即超
	require.Error(t, checkQuotaNotExceeded("STORAGE", 100, 60, 41), "60+41>100 拒绝")
	require.NoError(t, checkQuotaNotExceeded("STORAGE", 100, 60, 40), "60+40=100 不超（正好用满放行）")
	require.NoError(t, checkQuotaNotExceeded("STORAGE", 100, 100, 0), "无本次增量的纯查询不拒")
}

func TestQuotaWatermarkExceeded(t *testing.T) {
	// 80% 水位
	require.False(t, quotaWatermarkExceeded(79, 100))
	require.True(t, quotaWatermarkExceeded(80, 100))
	require.True(t, quotaWatermarkExceeded(101, 100), "超限必命中")
	// limit==0：有用量即事实超限
	require.True(t, quotaWatermarkExceeded(1, 0))
	require.False(t, quotaWatermarkExceeded(0, 0))
}

func TestRatioPct(t *testing.T) {
	require.EqualValues(t, 80, ratioPct(80, 100))
	require.EqualValues(t, 150, ratioPct(150, 100), "不封顶：超限如实报百分比（视觉封顶在 FE Progress）")
	require.EqualValues(t, 100, ratioPct(1, 0), "配 0 上限且有用量视为 100%")
	require.EqualValues(t, 0, ratioPct(0, 0))
}

func TestWatermarkCopy(t *testing.T) {
	require.Contains(t, watermarkTitle("zh-cn", "租户甲"), "租户甲")
	require.Contains(t, watermarkTitle("en-US", "Tenant A"), "Tenant A")
	require.Contains(t, watermarkTitle("unknown", "T"), "T", "未识别语言回落中文")

	body := watermarkBody("zh-cn", "租户甲", []QuotaWatermarkHit{
		{QuotaType: "USER_LIMIT", Used: 9, Limit: 10, RatioPct: 90},
	})
	require.Contains(t, body, "USER_LIMIT")
	require.Contains(t, body, "90%")
}
