package lexorank

import (
	"fmt"
	"testing"
)

func TestBetween(t *testing.T) {
	tests := []struct {
		a, b string
		desc string
	}{
		{
			a:    MinString,
			b:    MaxString,
			desc: "基本测试：最小值到最大值之间",
		},
		{
			a:    "yzzzzzzzzz",
			b:    MaxString,
			desc: "边界情况：yzzzzzzzzz到最大值之间",
		},
		{
			a:    "aaaaaaaaaa",
			b:    "aaaaaaaaab",
			desc: "相邻值测试",
		},
		{
			a:    "zzzzzzzzzy",
			b:    MaxString,
			desc: "另一个边界情况",
		},
		{
			a:    "zzzzzzzzzz",
			b:    "zzzzzzzzzz",
			desc: "相同值测试",
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			result := Between(tt.a, tt.b)
			fmt.Printf("Between(%q, %q) = %q (len=%d)\n", tt.a, tt.b, result, len(result))

			// 验证字符串长度为10
			if len(result) != 10 {
				t.Errorf("结果长度必须是10位，得到: %d", len(result))
			}

			// 验证结果不等于输入值
			if result == tt.a {
				t.Errorf("结果不能等于输入a: %s", result)
			}
			if result == tt.b {
				t.Errorf("结果不能等于输入b: %s", result)
			}

			// 验证结果在a和b之间（字典序）
			if result <= tt.a || result >= tt.b {
				t.Errorf("结果不在a和b之间: %s <= %s <= %s", tt.a, result, tt.b)
			}
		})
	}
}
