package consent

import (
	"errors"
	"fmt"

	esbuild "github.com/evanw/esbuild/pkg/api"
)

func bundleJS(entry string) (string, error) {
	result := esbuild.Build(esbuild.BuildOptions{
		EntryPoints: []string{entry},
		Target:      esbuild.ES2015,
		Format:      esbuild.FormatIIFE,
		Write:       false,
		Bundle:      true,
	})

	if numErrs := len(result.Errors); numErrs != 0 {
		err := fmt.Errorf(
			"bundleJS: %d error(s) bundling script `%s`: %s",
			numErrs, entry, result.Errors[0].Text,
		)
		return "", err
	}
	for _, output := range result.OutputFiles {
		return string(output.Contents), nil
	}
	return "", errors.New("bundleJS: bundling operating did not yield output file")
}
