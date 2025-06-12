package response

import "io"

type StatusCode int

const (
	Ok             StatusCode = 200
	BadRequest     StatusCode = 400
	InternalServer StatusCode = 500
)

func WriteStatusLine(w io.Writer, statusCode StatusCode) error {
}
