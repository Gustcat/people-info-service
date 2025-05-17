package response

type Status string

const (
	StatusOK    Status = "ok"
	StatusError Status = "error"
)

type Response[T any] struct {
	Status Status `json:"status"`
	Data   *T     `json:"data,omitempty"`
	Error  string `json:"error,omitempty"`
}

func OK[T any](data *T) Response[T] {
	return Response[T]{
		Status: StatusOK,
		Data:   data,
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
