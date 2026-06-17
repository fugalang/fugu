package parser

import (
	"fugu/pkg/parser/action/ast"
	"fugu/pkg/token"
	"testing"
)

func TestParserExpr(t *testing.T) {
	tests := []struct {
		name              string
		input             []byte
		expectedRoot      ast.NodeKind
		expectedRootOp    token.TokenKind
		expectedRight     ast.NodeKind
		expectedRightOp   token.TokenKind
		checkRightOp      bool
		checkRightAsPower bool
		expectedRoots     int
	}{
		{
			name:            "Тест приоритета умножения над сложением",
			input:           []byte("2 + 3 * 4"),
			expectedRoot:    ast.KindAdditiveExpr,
			expectedRootOp:  token.INCREASE,
			expectedRight:   ast.KindMultiplicativeExpr,
			expectedRightOp: token.MULTIPLY,
			checkRightOp:    true,
			expectedRoots:   1,
		},
		{
			name:              "Тест правой связности степени",
			input:             []byte("2 ^ 3 ^ 4"),
			expectedRoot:      ast.KindPowerExpr,
			expectedRootOp:    token.DEGREE,
			expectedRight:     ast.KindPowerExpr,
			checkRightAsPower: true,
			expectedRoots:     1,
		},
	}

	for _, tt := range tests {
		ps := New(tt.input, "main.fg")
		ps.Run()

		if ps.report.IsUse {
			t.Errorf("[%s] Парсер вернул ошибку", tt.name)
			continue
		}

		if len(ps.ast.Nodes) == 0 {
			t.Errorf("[%s] Парсер не создал ast узлы", tt.name)
			continue
		}

		if len(ps.Roots) != tt.expectedRoots {
			t.Errorf("[%s] Неверное количество корней. Ожидалось: %d, получено: %d",
				tt.name, tt.expectedRoots, len(ps.Roots))
			continue
		}

		root := ps.ast.Nodes[ps.Roots[len(ps.Roots)-1]]
		if root.Kind != tt.expectedRoot {
			t.Errorf("[%s] Неверный тип корня. Ожидался: %s, получен: %s",
				tt.name, tt.expectedRoot, root.Kind)
			continue
		}

		if token.TokenKind(root.Data3) != tt.expectedRootOp {
			t.Errorf("[%s] Неверный оператор корня. Ожидался: %s, получен: %s",
				tt.name, tt.expectedRootOp, token.TokenKind(root.Data3))
			continue
		}

		right := ps.ast.Nodes[root.Data2]
		if right.Kind != tt.expectedRight {
			t.Errorf("[%s] Неверный тип правого узла. Ожидался: %s, получен: %s",
				tt.name, tt.expectedRight, right.Kind)
			continue
		}

		if tt.checkRightOp && token.TokenKind(right.Data3) != tt.expectedRightOp {
			t.Errorf("[%s] Неверный оператор правого узла. Ожидался: %s, получен: %s",
				tt.name, tt.expectedRightOp, token.TokenKind(right.Data3))
			continue
		}

		if tt.checkRightAsPower && right.Kind != ast.KindPowerExpr {
			t.Errorf("[%s] Правая часть степени должна быть тоже степенью. Получен: %s",
				tt.name, right.Kind)
			continue
		}
	}
}
