package response

type Status string

const (
	StatusOK    Status = "ok"
	StatusError Status = "error"
)

type Response[T any] struct {
	Status     Status      `json:"status"`
	Data       *T          `json:"data,omitempty"`
	Error      string      `json:"error,omitempty"`
	Pagination *Pagination `json:"pagination,omitempty"`
}

func OK[T any](data *T) Response[T] {
	return Response[T]{
		Status: StatusOK,
		Data:   data,
	}
}

func OKWithPagination[T any](data *T, p *Pagination) Response[T] {
	return Response[T]{
		Status:     StatusOK,
		Data:       data,
		Pagination: p,
	}
}

type Void struct{}

func Error(msg string) Response[Void] {
	return Response[Void]{
		Status: StatusError,
		Error:  msg,
	}
}

// TODO: написать валидатор
