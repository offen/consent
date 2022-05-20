package main

import (
	"fmt"
	"os"
	"path"

	esbuild "github.com/evanw/esbuild/pkg/api"
)

func main() {
	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	for _, b := range []struct {
		in  string
		out string
	}{
		{"proxy/proxy.js", "proxy/bundle.min.js"},
		{"client/client.js", "client/bundle.min.js"},
	} {
		if err := bundleJS(
			path.Join(cwd, b.in), path.Join(cwd, b.out),
		); err != nil {
			panic(err)
		}
	}
}

func bundleJS(in, out string) error {
	result := esbuild.Build(esbuild.BuildOptions{
		EntryPoints:       []string{in},
		Target:            esbuild.ES2015,
		Format:            esbuild.FormatIIFE,
		Outfile:           out,
		MinifyWhitespace:  true,
		MinifyIdentifiers: true,
		MinifySyntax:      true,
		Write:             true,
		Bundle:            true,
	})

	if numErrs := len(result.Errors); numErrs != 0 {
		err := fmt.Errorf(
			"bundleJS: %d error(s) bundling script `%s`: %s",
			numErrs, in, result.Errors[0].Text,
		)
		return err
	}
	return nil
}
