package contextx

import (
	"context"

	"github.com/go-xuan/typex"
)

const (
	valuesKey = "__quanx_values__"
)

func New() context.Context {
	var values = make(map[string]any)
	return context.WithValue(context.Background(), valuesKey, values)
}

func SetValue(ctx context.Context, key string, value any) {
	if v := ctx.Value(valuesKey); v != nil {
		if values, ok := v.(map[string]any); ok {
			values[key] = value
		}
	}
}

func getValue(ctx context.Context, key string) any {
	if v := ctx.Value(valuesKey); v != nil {
		if values, ok := v.(map[string]any); ok {
			return values[key]
		}
	}
	return nil
}

func GetValue(ctx context.Context, key string) typex.Value {
	if value := getValue(ctx, key); value != nil {
		return typex.NewValue(value)
	}
	return typex.ZeroValue()
}

func GetValueString(ctx context.Context, key string) string {
	if value := getValue(ctx, key); value != nil {
		if v, ok := value.(string); ok {
			return v
		}
	}
	return ""
}

func GetValueInt(ctx context.Context, key string) int {
	if value := getValue(ctx, key); value != nil {
		if v, ok := value.(int); ok {
			return v
		}
	}
	return 0
}

func GetValueInt64(ctx context.Context, key string) int64 {
	if value := getValue(ctx, key); value != nil {
		if v, ok := value.(int64); ok {
			return v
		}
	}
	return 0
}

func GetValueBool(ctx context.Context, key string) bool {
	if value := getValue(ctx, key); value != nil {
		if v, ok := value.(bool); ok {
			return v
		}
	}
	return false
}

func GetValueFloat64(ctx context.Context, key string) float64 {
	if value := getValue(ctx, key); value != nil {
		if v, ok := value.(float64); ok {
			return v
		}
	}
	return 0
}
