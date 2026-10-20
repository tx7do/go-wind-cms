package password

import (
	"testing"
)

func TestValidateComplexity(t *testing.T) {
	const minLen = 8

	// 合规：长度达标 + 四类中至少三类
	pass := []string{
		"Abcdefg1", // 小写+大写+数字（三类）
		"Abcdefg!", // 大写+小写+符号（三类）
		"ABCDEFG1", // 大写+数字（两类）——修正：应为两类，拒绝；此行仅示意
	}
	_ = pass

	passCases := []string{
		"Abcdefg1",   // 小写+大写+数字
		"Abcdefg!",   // 大写+小写+符号
		"12345678aB", // 数字+小写+大写
		"Abcdef1!@#", // 四类齐全
	}
	for _, c := range passCases {
		if err := ValidateComplexity(c, minLen); err != nil {
			t.Errorf("ValidateComplexity(%q) = %v, want nil", c, err)
		}
	}

	// 不合规：任一维度不足
	weakCases := []string{
		"",          // 空
		"Abc1",      // 长度不足
		"abcdefgh",  // 只有小写（一类）
		"12345678",  // 只有数字（一类）
		"ABCDEFG1",  // 大写+数字（两类）
		"abcdefg!",  // 小写+符号（两类）
		"12345678a", // 数字+小写（两类）
		"abcABC",    // 两类
	}
	for _, c := range weakCases {
		if err := ValidateComplexity(c, minLen); err == nil {
			t.Errorf("ValidateComplexity(%q) = nil, want error", c)
		}
	}

	// minLen 透传：调用方配置 12 位时 8 位口令拒绝
	if err := ValidateComplexity("Abcdefg1", 12); err == nil {
		t.Error("ValidateComplexity with minLen=12 should reject 8-char password")
	}

	// minLen<=0 回落 8
	if err := ValidateComplexity("Abcdefg1", 0); err != nil {
		t.Errorf("ValidateComplexity with minLen=0 should fallback to 8: %v", err)
	}
}
