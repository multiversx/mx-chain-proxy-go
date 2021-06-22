package txcost

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRemoveLatestArgumentFromDataField(t *testing.T) {
	t.Parallel()

	dataField := removeLatestArgumentFromDataField("function@arg1@arg2@arg3@shouldBeRemoved")
	require.Equal(t, "function@arg1@arg2@arg3", dataField)

	dataField = removeLatestArgumentFromDataField("@new@arg1@arg2@arg3@shouldBeRemoved")
	require.Equal(t, "@new@arg1@arg2@arg3", dataField)

	dataField = removeLatestArgumentFromDataField("1@2@3")
	require.Equal(t, "1@2", dataField)

	dataField = removeLatestArgumentFromDataField("")
	require.Equal(t, "", dataField)

	dataField = removeLatestArgumentFromDataField("the-field")
	require.Equal(t, "the-field", dataField)
}
