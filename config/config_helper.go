package config

import (
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/chaos-io/chaos/config/reader"
	"github.com/chaos-io/chaos/config/source"
	"github.com/chaos-io/chaos/config/source/file"
)

var supportedFileSuffixes = map[string]bool{
	"json": true,
	// "toml": true,
	"xml":  true,
	"yaml": true,
}

func init() {
	_ = Load(defaultSources()...)
}

func defaultSources() []source.Source {
	workDir, _ := os.Getwd()
	dirs := []string{
		filepath.Join(workDir, "conf"),
		filepath.Join(workDir, "config"),
		filepath.Join(workDir, "configs"),
	}

	if configPath := os.Getenv("CONFIG_PATH"); len(configPath) > 0 {
		dirs = append([]string{configPath}, dirs...)
	}

	if strings.Contains(workDir, "/cmd/") || strings.HasSuffix(workDir, "/cmd") {
		dirs = append(dirs, []string{"../configs", "../../configs"}...)
	}

	var sources []source.Source
	for _, dir := range dirs {
		if sources = newFileSources(dir, os.Getenv("DEPLOY_ENV")); len(sources) > 0 {
			break
		}
	}

	// sources = append(sources, env.NewSource())
	return sources
}

func newFileSources(dir string, env string) []source.Source {
	var sources []source.Source

	files, _ := os.ReadDir(dir)
	for _, f := range files {
		if f.IsDir() {
			ss := newFileSources(f.Name(), env)
			sources = append(sources, ss...)
		} else {
			segments := strings.Split(f.Name(), ".")
			suffix := ""
			if len(segments) >= 2 {
				suffix = segments[len(segments)-1]
			}
			if !supportedFileSuffixes[suffix] {
				continue
			}
			p := path.Join(dir, f.Name())
			if len(env) > 0 {
				name := strings.Join(segments[:len(segments)-1], ".")
				if strings.HasSuffix(name, env) {
					sources = append(sources, file.NewSource(file.WithPath(p)))
				}
			} else {
				sources = append(sources, file.NewSource(file.WithPath(p)))
			}
		}
	}

	return sources
}

// GetValue a value from the config
func GetValue(path ...string) reader.Value {
	if len(path) == 1 {
		segments := strings.Split(path[0], ".")
		return Get(segments...)
	}

	return Get(path...)
}

type watchCloser struct {
	exit chan struct{}
}

func (w watchCloser) Close() error {
	fmt.Println("close")
	w.exit <- struct{}{}
	return nil
}

func WatchFunc(handle func(reader.Value), paths ...string) (io.Closer, error) {
	_path := make([]string, 0, len(paths))
	for _, v := range paths {
		_path = append(_path, strings.Split(v, ".")...)
	}

	exit := make(chan struct{})
	w, err := Watch(_path...)
	if err != nil {
		return nil, err
	}

	go func() {
		for {
			v, err := w.Next()
			// if err == err_code.WatchStoppedError {
			//	return
			// }
			if err != nil {
				continue
			}

			log.Printf("file changed, %s", string(v.Bytes()))

			// if v.Empty() {
			//	continue
			// }

			if handle != nil {
				handle(v)
			}
		}
	}()

	go func() {
		select {
		case <-exit:
			_ = w.Stop()
		}
	}()

	return watchCloser{exit: exit}, nil
}
