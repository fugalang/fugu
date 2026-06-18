package action

import (
	"fugu/pkg/token"
	"testing"
)

func TestActionMap(t *testing.T) {
	tests := []struct {
		name     string
		state    int
		tk       token.TokenKind
		expected ActionType
		next     int
	}{
		{
			name:     "Тест группы литералов",
			state:    0,
			tk:       token.INTEGER,
			expected: Shift,
			next:     1,
		},
		{
			name:     "Тест группы арифметики",
			state:    1,
			tk:       token.INCREASE,
			expected: Reduce,
			next:     1,
		},
	}

	for _, tt := range tests {
		act := Action(tt.state, tt.tk)

		if act.Typ != tt.expected {
			t.Errorf("[%s] Неверный тип действия. Ожидался: %d, получен: %d",
				tt.name, tt.expected, act.Typ)
			continue
		}

		if act.State != tt.next {
			t.Errorf("[%s] Неверное состояние. Ожидалось: %d, получено: %d",
				tt.name, tt.next, act.State)
			continue
		}
	}
}
