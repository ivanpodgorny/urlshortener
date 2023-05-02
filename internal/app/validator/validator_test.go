package validator

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsURL(t *testing.T) {
	tests := []struct {
		name  string
		url   string
		valid bool
	}{
		{
			name:  "корректный URL",
			url:   "https://ya.ru/",
			valid: true,
		},
		{
			name:  "проверка присутствия схемы в URL",
			url:   "ya.ru",
			valid: false,
		},
		{
			name:  "проверка присутствия hostname в URL",
			url:   "/path/test",
			valid: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := IsURL(tt.url)
			if !tt.valid {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestLength(t *testing.T) {
	var (
		l         = 10
		validator = Length(l)
	)
	tests := []struct {
		name  string
		str   string
		valid bool
	}{
		{
			name:  "максимальная длина",
			str:   strings.Repeat("a", l),
			valid: true,
		},
		{
			name:  "длина меньше заданной",
			str:   strings.Repeat("a", l-1),
			valid: true,
		},
		{
			name:  "длина больше заданной",
			str:   strings.Repeat("a", l+1),
			valid: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator(tt.str)
			if !tt.valid {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestSize(t *testing.T) {
	validator := Size[int](2)
	tests := []struct {
		name  string
		arr   []int
		valid bool
	}{
		{
			name:  "максимальная длина",
			arr:   []int{1, 2},
			valid: true,
		},
		{
			name:  "длина меньше заданной",
			arr:   []int{1},
			valid: true,
		},
		{
			name:  "длина больше заданной",
			arr:   []int{1, 2, 3},
			valid: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator(tt.arr)
			if !tt.valid {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func BenchmarkIsURL(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = IsURL("https://ya.ru/")
	}
}
