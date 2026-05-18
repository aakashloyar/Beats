package http 
import (
	"github.com/aakashloyar/beats/encoding/internal/application/ports/in"
)
type Handler struct {
	service	in.EncodingService
}
