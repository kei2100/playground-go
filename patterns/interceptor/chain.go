package interceptor

import "context"

// Handler type
type Handler func(ctx context.Context, req interface{}) (res interface{}, err error)

// Interceptor type
type Interceptor func(ctx context.Context, handler Handler) (res interface{}, err error)

// ChainInterceptors chains multiple interceptors
func ChainInterceptors(interceptors []Interceptor, index int, handler Handler) Handler {
	if index == len(interceptors) {
		return handler
	}
	return func(ctx context.Context, req interface{}) (res interface{}, err error) {
		return interceptors[index](ctx, ChainInterceptors(interceptors, index+1, handler))
	}
}

type server struct {
	interceptors []Interceptor
}

func (s *server) HandleFunc(pattern string, handler Handler) {
	handler = ChainInterceptors(s.interceptors, 0, handler)
	var req interface{}
	handler(context.Background(), req)
}
