package web_test

import (
	"testing"

	"github.com/l3lackShark/binance-rest-listener/web"
	"github.com/stretchr/testify/assert"
)

func TestRequestFiveMinAVG(t *testing.T) {
	_, err := web.RequestFiveMinAVG()
	assert.NoError(t, err)
}
