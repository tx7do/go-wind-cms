package data

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func lp(id uint32, typ, method, value, reason string) EffectivePolicy {
	return EffectivePolicy{TargetID: id, Type: typ, Method: method, Value: value, Reason: reason}
}

func TestMatchLoginPolicyIP(t *testing.T) {
	now := time.Date(2026, 10, 20, 12, 0, 0, 0, time.Local)

	// 黑名单：精确 IP 命中
	blocked, reason := MatchLoginPolicy([]EffectivePolicy{lp(0, "BLACKLIST", "IP", "10.0.0.5", "")}, 7, "10.0.0.5", "", "", now)
	require.True(t, blocked)
	require.Contains(t, reason, "ip blacklist")

	// CIDR 命中
	blocked, _ = MatchLoginPolicy([]EffectivePolicy{lp(0, "BLACKLIST", "IP", "10.0.0.0/24", "")}, 7, "10.0.0.99", "", "", now)
	require.True(t, blocked)

	// 不在网段放行
	blocked, _ = MatchLoginPolicy([]EffectivePolicy{lp(0, "BLACKLIST", "IP", "10.0.0.0/24", "")}, 7, "10.1.0.99", "", "", now)
	require.False(t, blocked)

	// 非法 CIDR 不命中（配置错误不阻断登录）
	blocked, _ = MatchLoginPolicy([]EffectivePolicy{lp(0, "BLACKLIST", "IP", "not-a-cidr", "")}, 7, "10.0.0.5", "", "", now)
	require.False(t, blocked)

	// 定向条目只约束目标用户
	blocked, _ = MatchLoginPolicy([]EffectivePolicy{lp(7, "BLACKLIST", "IP", "10.0.0.5", "")}, 8, "10.0.0.5", "", "", now)
	require.False(t, blocked, "定向条目对其他用户不生效")
	blocked, _ = MatchLoginPolicy([]EffectivePolicy{lp(7, "BLACKLIST", "IP", "10.0.0.5", "")}, 7, "10.0.0.5", "", "", now)
	require.True(t, blocked)
}

func TestMatchLoginPolicyWhitelist(t *testing.T) {
	now := time.Date(2026, 10, 20, 12, 0, 0, 0, time.Local)

	whitelist := []EffectivePolicy{lp(0, "WHITELIST", "IP", "192.168.1.0/24", "")}

	// 白名单命中放行
	blocked, _ := MatchLoginPolicy(whitelist, 7, "192.168.1.10", "", "", now)
	require.False(t, blocked)

	// 未命中任何白名单 → 拒绝
	blocked, reason := MatchLoginPolicy(whitelist, 7, "10.0.0.5", "", "", now)
	require.True(t, blocked)
	require.Contains(t, reason, "not in ip whitelist")

	// 白名单只存在但定向他人：对当前用户无约束
	blocked, _ = MatchLoginPolicy([]EffectivePolicy{lp(9, "WHITELIST", "IP", "192.168.1.0/24", "")}, 7, "10.0.0.5", "", "", now)
	require.False(t, blocked)
}

func TestMatchLoginPolicyTime(t *testing.T) {
	at := func(h, m int) time.Time {
		return time.Date(2026, 10, 20, h, m, 0, 0, time.Local)
	}

	window := []EffectivePolicy{lp(0, "BLACKLIST", "TIME", "22:00-06:00", "")}
	require.True(t, matchTimeWindow(at(23, 30), "22:00-06:00"), "跨午夜窗口内")
	require.True(t, matchTimeWindow(at(2, 0), "22:00-06:00"))
	require.False(t, matchTimeWindow(at(12, 0), "22:00-06:00"))
	require.False(t, matchTimeWindow(at(6, 0), "22:00-06:00"), "终点开区间")
	require.True(t, matchTimeWindow(at(6, 0), "06:00-18:00"), "起点闭区间")

	blocked, _ := MatchLoginPolicy(window, 7, "", "", "", at(23, 30))
	require.True(t, blocked)
}

func TestMatchLoginPolicyDevice(t *testing.T) {
	now := time.Now()
	policies := []EffectivePolicy{lp(0, "BLACKLIST", "DEVICE", "dev-123", "")}

	require.True(t, blockedHelper(policies, 7, "", "dev-123", now))
	require.False(t, blockedHelper(policies, 7, "", "dev-456", now))
	require.False(t, blockedHelper(policies, 7, "", "", now), "空 deviceId 不判定")
}

func TestMatchLoginPolicyRegion(t *testing.T) {
	now := time.Now()
	policies := []EffectivePolicy{lp(0, "BLACKLIST", "REGION", "四川省", "")}

	// 省精确相等（大小写不敏感）命中
	require.True(t, blockedHelper(policies, 7, "", "", now.Add(0), "四川省"))
	require.True(t, blockedHelper([]EffectivePolicy{lp(0, "BLACKLIST", "REGION", "sichuan", "")},
		7, "", "", now, "SICHUAN"), "大小写不敏感")

	// 他省不命中
	require.False(t, blockedHelper(policies, 7, "", "", now, "广东省"))

	// 解析失败（空归属地）不判定：fail-open
	require.False(t, blockedHelper(policies, 7, "", "", now, ""))

	// REGION 白名单：命中放行、未命中拒绝
	wl := []EffectivePolicy{lp(0, "WHITELIST", "REGION", "四川省", "")}
	require.False(t, blockedHelper(wl, 7, "", "", now, "四川省"))
	require.True(t, blockedHelper(wl, 7, "", "", now, "广东省"))
}

// blockedHelper 跑一遍 MatchLoginPolicy 只回 blocked。
func blockedHelper(policies []EffectivePolicy, userId uint32, clientIP, deviceId string, now time.Time, region ...string) bool {
	r := ""
	if len(region) > 0 {
		r = region[0]
	}
	blocked, _ := MatchLoginPolicy(policies, userId, clientIP, deviceId, r, now)
	return blocked
}

func TestValidateLoginPolicyValue(t *testing.T) {
	require.NoError(t, ValidateLoginPolicyValue("IP", "10.0.0.1"))
	require.NoError(t, ValidateLoginPolicyValue("IP", "10.0.0.0/24"))
	require.Error(t, ValidateLoginPolicyValue("IP", "not-an-ip"))
	require.Error(t, ValidateLoginPolicyValue("IP", "10.0.0.0/99"))
	require.NoError(t, ValidateLoginPolicyValue("TIME", "22:00-06:00"))
	require.Error(t, ValidateLoginPolicyValue("TIME", "22:00-22:00"), "起止相等拒绝")
	require.Error(t, ValidateLoginPolicyValue("TIME", "25:00-06:00"))
	require.NoError(t, ValidateLoginPolicyValue("DEVICE", "dev-123"))
	require.NoError(t, ValidateLoginPolicyValue("REGION", "四川省"))
	require.Error(t, ValidateLoginPolicyValue("IP", ""), "空值拒绝")
}
