package docker

import (
	"bytes"
	"context"
	"os"
	"regexp"
	"strings"

	transformer "github.com/apenella/go-common-utils/transformer/string"
	"github.com/apenella/go-docker-builder/pkg/build"
	contextpath "github.com/apenella/go-docker-builder/pkg/build/context/path"
	"github.com/apenella/go-docker-builder/pkg/response"
)

func BuildImage(imagePath string, imageName string, imageTag string) error {
	w := bytes.NewBuffer(nil)

	res := response.NewDefaultResponse(
		response.WithTransformers(
			transformer.Prepend("buildPathContext"),
		),
		response.WithWriter(w),
	)

	dockerBuilder := build.NewDockerBuildCmd(cli).
		WithImageName(imageName).
		WithResponse(res)

	tag := imageName
	if len(imageTag) > 0 {
		tag = strings.Join([]string{imageName, imageTag}, ":")
	}
	dockerBuilder.AddTags(tag)
	dockerBuildContext := &contextpath.PathBuildContext{
		Path: imagePath,
	}

	if err := dockerBuilder.AddBuildContext(dockerBuildContext); err != nil {
		return err // errors.New("buildPathContext", "Error adding build docker context", err)
	}

	if err := dockerBuilder.Run(context.TODO()); err != nil {
		return err // errors.New("buildPathContext", fmt.Sprintf("Error building '%s'", imageName), err)
	}

	return nil
}

var loadImageRegexp = regexp.MustCompile(`(?m:^##load-image:.*$)`)

// ExtractLoadImage
//
// ##load-image:alpine
func ExtractLoadImage(buildFile string) string {
	content, err := os.ReadFile(buildFile)
	if err != nil {
		return ""
	}

	image := loadImageRegexp.Find(content)
	if len(image) > 0 {
		return strings.TrimSpace(strings.TrimPrefix(string(image), "##load-image:"))
	}
	return ""
}
