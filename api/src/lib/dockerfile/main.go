package dockerfile
import (
	"github.com/docker/docker/builder/dockerfile/parser"
	"io"
	"regexp"
	"strconv"
)

func ExtractExposedPort(rwc io.Reader) (string, error) {
	x, _ := parser.Parse(rwc)

	for i :=0; i< len(x.AST.Children); i++  {
		row := x.AST.Children[i]

		switch row.Value {
			case "expose":
				if _, err := strconv.Atoi(row.Next.Value); err == nil {
					return row.Next.Value, nil
				}
			case "env":
				if row.Next.Value != "LISTEN_PORT" {
					continue
				}
				r := regexp.MustCompile(`([\d]{4})`)
				results := r.FindStringSubmatch(row.Next.Next.Value)
				if len(results) != 2 {
					continue
				}
				if _, err := strconv.Atoi( results[1]); err == nil {
					return  results[1], nil
				}
				continue
		}
	}

	return "", nil
}