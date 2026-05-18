package out 

import "context"

type Encoder interface {
	EncodeVariants(ctx context.Context, inputPath string, outputDir string) error 
}