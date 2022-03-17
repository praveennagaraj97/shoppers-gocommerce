package cmd

import (
	"context"

	"github.com/praveennagaraj97/shoppers-gocommerce/routes"
)

/*
Bootstraps the GlobalConfiguration with database connection and starts the app
*/
func Serve() {
	// global configuration which will be passed to all required modules as pointer ref.
	cfg := initializeApp()

	// close database connection when app stops serving.
	defer cfg.DatabaseClient.Disconnect(context.Background())
	// Router Initialization which will start app and listen to routes.
	r := &routes.Router{}
	r.InitializeRouter(cfg)
}
