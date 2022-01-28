package describe

import (
	"net/http"
	"testing"

	"github.com/aziontech/azion-cli/pkg/httpmock"
	"github.com/aziontech/azion-cli/pkg/testutils"
	"github.com/stretchr/testify/require"
)

var successResponse string = `
{
    "results":{
        "id":1337,
        "name":"SUUPA_FUNCTION",
        "language":"javascript",
        "code":"async function handleRequest(request) {return new Response(\"Hello World!\",{status:200})}",
        "json_args":{"a":1,"b":2},"function_to_run":"",
        "initiator_type":"edge_application",
        "active":true,
        "last_editor":"testando@azion.com",
        "modified":"2022-01-26T12:31:09.865515Z",
        "reference_count":0
    },
    "schema_version":3
}
`

func TestDescribe(t *testing.T) {
	t.Run("describe a function", func(t *testing.T) {
		mock := &httpmock.Registry{}

		mock.Register(
			httpmock.REST("GET", "edge_functions/123"),
			httpmock.JSONFromString(successResponse),
		)

		f, stdout, _ := testutils.NewFactory(mock)

		cmd := NewCmd(f)

		cmd.SetArgs([]string{"123"})

		err := cmd.Execute()
		require.NoError(t, err)

		require.Equal(t, `ID: 1337
Name: SUUPA_FUNCTION
Language: javascript
Reference Count: 0
Modified at: 2022-01-26T12:31:09.865515Z
Initiator Type: edge_application
Last Editor: testando@azion.com
Function to run: 
JSON Args: {"a":1,"b":2}
`, stdout.String())

	})

	t.Run("with code", func(t *testing.T) {
		mock := &httpmock.Registry{}

		mock.Register(
			httpmock.REST("GET", "edge_functions/123"),
			httpmock.JSONFromString(successResponse),
		)

		f, stdout, _ := testutils.NewFactory(mock)

		cmd := NewCmd(f)

		cmd.SetArgs([]string{"123", "--with-code"})

		err := cmd.Execute()
		require.NoError(t, err)

		require.Equal(t, `ID: 1337
Name: SUUPA_FUNCTION
Language: javascript
Reference Count: 0
Modified at: 2022-01-26T12:31:09.865515Z
Initiator Type: edge_application
Last Editor: testando@azion.com
Function to run: 
JSON Args: {"a":1,"b":2}
Code:
async function handleRequest(request) {return new Response("Hello World!",{status:200})}
`, stdout.String())

		t.Run("not found", func(t *testing.T) {
			mock := &httpmock.Registry{}

			mock.Register(
				httpmock.REST("GET", "edge_functions/1234"),
				httpmock.StatusStringResponse(http.StatusNotFound, "Not Found"),
			)

			f, _, _ := testutils.NewFactory(mock)

			cmd := NewCmd(f)

			cmd.SetArgs([]string{"1234", "--with-code"})

			err := cmd.Execute()

			require.Error(t, err)
		})

	})
}