package static

import (
	_ "embed"
)

//go:embed connect.html
var ConnectHTML string

//go:embed uploads.html
var UploadHTML string

//go:embed client.html
var UiClientHTML string

//go:embed server.html
var UiServerHTML string

//go:embed index.html
var UiIndexHTML string
