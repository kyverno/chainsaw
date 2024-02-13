package mutate

// import (
// 	"context"
// 	"reflect"

// 	"github.com/jmespath-community/go-jmespath/pkg/binding"
// 	"github.com/kyverno/kyverno-json/pkg/engine/template"
// 	reflectutils "github.com/kyverno/kyverno-json/pkg/utils/reflect"
// )

// type projection struct {
// 	foreach     bool
// 	foreachName string
// 	binding     string
// 	result      any
// }

// func project(ctx context.Context, key any, value any, bindings binding.Bindings, opts ...template.Option) (*projection, error) {
// 	expression := parseExpression(ctx, key)
// 	if expression != nil {
// 		if expression.engine != "" {
// 			projected, err := template.Execute(ctx, expression.statement, value, bindings, opts...)
// 			if err != nil {
// 				return nil, err
// 			}
// 			return &projection{
// 				foreach:     expression.foreach,
// 				foreachName: expression.foreachName,
// 				binding:     expression.binding,
// 				result:      projected,
// 			}, nil
// 		} else {
// 			if reflectutils.GetKind(value) == reflect.Map {
// 				mapValue := reflect.ValueOf(value).MapIndex(reflect.ValueOf(expression.statement))
// 				var value any
// 				if mapValue.IsValid() {
// 					value = mapValue.Interface()
// 				}
// 				return &projection{
// 					foreach:     expression.foreach,
// 					foreachName: expression.foreachName,
// 					binding:     expression.binding,
// 					result:      value,
// 				}, nil
// 			}
// 		}
// 	}
// 	if reflectutils.GetKind(value) == reflect.Map {
// 		mapValue := reflect.ValueOf(value).MapIndex(reflect.ValueOf(key))
// 		var value any
// 		if mapValue.IsValid() {
// 			value = mapValue.Interface()
// 		}
// 		return &projection{
// 			result: value,
// 		}, nil
// 	}
// 	// TODO is this an error ?
// 	return &projection{
// 		result: value,
// 	}, nil
// }
