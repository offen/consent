package consent

import (
	"fmt"
	"os"

	esbuild "github.com/evanw/esbuild/pkg/api"
)

func minifyJS(js string) (string, error) {
	_, isDevelopment := os.LookupEnv("DEVELOPMENT")
	result := esbuild.Transform(js, esbuild.TransformOptions{
		MinifyWhitespace:  !isDevelopment,
		MinifyIdentifiers: !isDevelopment,
		MinifySyntax:      !isDevelopment,
		Target:            esbuild.ES2015,
		Format:            esbuild.FormatIIFE,
	})

	if len(result.Errors) != 0 {
		return "", fmt.Errorf("minifyJS: error minifying script: %s", result.Errors[0].Text)
	}
	return string(result.Code), nil
}
