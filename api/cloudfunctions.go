package cloudfunctions

import (
	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	route "github.com/kikudesuyo/point-hub/app/routes/v1"

	// required by vendor dir deployment by Cloud Functions
	_ "github.com/GoogleCloudPlatform/functions-framework-go/funcframework"
)

func init() {
	functions.HTTP("RunHTTPServer", route.RunHTTPServer)
}
