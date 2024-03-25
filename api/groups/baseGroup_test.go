package groups

import (
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/multiversx/mx-chain-proxy-go/data"
	"github.com/stretchr/testify/assert"
)

func TestCrudOperationsBaseGroup(t *testing.T) {
	t.Parallel()

	ginHandler := func(c *gin.Context) {}
	hd0 := &data.EndpointHandlerData{
		Path:    "path0",
		Handler: ginHandler,
		Method:  "GET",
	}
	hd1 := &data.EndpointHandlerData{
		Path:    "path1",
		Handler: ginHandler,
		Method:  "GET",
	}
	hd2 := &data.EndpointHandlerData{
		Path:    "path2",
		Handler: ginHandler,
		Method:  "GET",
	}
	hd3 := &data.EndpointHandlerData{
		Path:    "path3",
		Handler: ginHandler,
		Method:  "GET",
	}
	hd4 := &data.EndpointHandlerData{
		Path:    "path4",
		Handler: ginHandler,
		Method:  "GET",
	}

	bg := &baseGroup{
		endpoints: []*data.EndpointHandlerData{hd0, hd1, hd2},
	}

	// ensure the order is kept
	assert.Equal(t, hd0.Path, bg.endpoints[0].Path)
	assert.Equal(t, hd1.Path, bg.endpoints[1].Path)
	assert.Equal(t, hd2.Path, bg.endpoints[2].Path)

	err := bg.UpdateEndpoint(hd0.Path, *hd3)
	assert.NoError(t, err)
	assert.Equal(t, 3, len(bg.endpoints))

	// ensure the order
	assert.Equal(t, hd3.Path, bg.endpoints[0].Path)
	assert.Equal(t, hd1.Path, bg.endpoints[1].Path)
	assert.Equal(t, hd2.Path, bg.endpoints[2].Path)

	err = bg.AddEndpoint(hd4.Path, *hd4)
	assert.NoError(t, err)
	assert.Equal(t, 4, len(bg.endpoints))

	// ensure the order
	assert.Equal(t, hd3.Path, bg.endpoints[0].Path)
	assert.Equal(t, hd1.Path, bg.endpoints[1].Path)
	assert.Equal(t, hd2.Path, bg.endpoints[2].Path)
	assert.Equal(t, hd4.Path, bg.endpoints[3].Path)

	err = bg.RemoveEndpoint(hd2.Path)
	assert.NoError(t, err)
	assert.Equal(t, 3, len(bg.endpoints))

	// ensure the order
	assert.Equal(t, hd3.Path, bg.endpoints[0].Path)
	assert.Equal(t, hd1.Path, bg.endpoints[1].Path)
	assert.Equal(t, hd4.Path, bg.endpoints[2].Path)
}
