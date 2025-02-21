package utils

import (
	"testing"

	ctypes "github.com/docker/engine-api/types/container"
	"github.com/cmosh/interlock/ext"
)

func TestContextRoot(t *testing.T) {
	testContext := "/context"

	cfg := &ctypes.Config{
		Labels: map[string]string{
			ext.InterlockContextRootLabel: testContext,
		},
	}

	context := ContextRoot(cfg)

	if context != testContext {
		t.Fatalf("expected %s; received %s", testContext, context)
	}
}

func TestContextRootRewrite(t *testing.T) {
	cfg := &ctypes.Config{
		Labels: map[string]string{
			ext.InterlockContextRootRewriteLabel: "true",
		},
	}

	rewrite := ContextRootRewrite(cfg)

	if !rewrite {
		t.Fatal("expected context root rewrite")
	}
}
