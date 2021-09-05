package registry

import (
	"github.com/stretchr/testify/suite"
	"testing"
)

type ExtractImageRefContextUtilsTestSuite struct {
	suite.Suite
}

func (suite *ExtractImageRefContextUtilsTestSuite) TestExtractRegistryAndRepositoryFromImageReferencePublicImageRef() {
	ctx, err := ExtractImageRefContext("redis")
	suite.Nil(err)
	suite.Equal("index.docker.io", ctx.Registry)
	suite.Equal("library/redis", ctx.Repository)
}

func (suite *ExtractImageRefContextUtilsTestSuite) TestExtractRegistryAndRepositoryFromImageReferenceTag() {
	ctx, err := ExtractImageRefContext("tomer.azurecr.io/redis:v1")
	suite.Nil(err)
	suite.Equal("tomer.azurecr.io", ctx.Registry)
	suite.Equal("redis", ctx.Repository)
}

func (suite *ExtractImageRefContextUtilsTestSuite) TestExtractImageRefContext_NoIdentifier() {
	ctx, err := ExtractImageRefContext("tomer.azurecr.io/redis")
	suite.Nil(err)
	suite.Equal("tomer.azurecr.io", ctx.Registry)
	suite.Equal("redis", ctx.Repository)
}

func (suite *ExtractImageRefContextUtilsTestSuite) TestExtractImageRefContext_Digest_Parsed() {
	imageRef := "tomer.azurecr.io/redis@sha256:4a1c4b21597c1b4415bdbecb28a3296c6b5e23ca4f9feeb599860a1dac6a0108"
	ctx, err := ExtractImageRefContext(imageRef)
	suite.Nil(err)
	suite.Equal("tomer.azurecr.io", ctx.Registry)
	suite.Equal("redis", ctx.Repository)
}

func (suite *ExtractImageRefContextUtilsTestSuite) TestExtractImageRefContext_DigestBadFormat_Err() {
	// The last 4 chars of the digest are deleted:
	imageRef := "tomer.azurecr.io/redis@sha256:4a1c4b21597c1b4415bdbecb28a3296c6b5e23ca4f9feeb599860a1dac6a"
	ctx, _ := ExtractImageRefContext(imageRef)
	//suite.Equal(reflect.TypeOf(&name.ErrBadName{}), reflect.TypeOf(err))
	suite.Nil(ctx)
}

func (suite *ExtractImageRefContextUtilsTestSuite) TestExtractImageRefContext_TagAndDigest_ParsedDigestIgnoreTag() {
	imageRef := "tomer.azurecr.io/redis:v1@sha256:4a1c4b21597c1b4415bdbecb28a3296c6b5e23ca4f9feeb599860a1dac6a0108"
	ctx, err := ExtractImageRefContext(imageRef)
	suite.Nil(err)
	suite.Equal("tomer.azurecr.io", ctx.Registry)
	suite.Equal("redis", ctx.Repository)
}

func TestExtractImageRefContext(t *testing.T) {
	suite.Run(t, new(ExtractImageRefContextUtilsTestSuite))
}