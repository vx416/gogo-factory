package codegen

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAst(t *testing.T) {
	node, err := ParseFile("./user.go")
	require.NoError(t, err)

	fm := ParseFileMeta(node, "User")
	require.NoError(t, err)

	res, err := GetTempalte(false, fm)
	require.NoError(t, err)
	t.Log(res)

}
