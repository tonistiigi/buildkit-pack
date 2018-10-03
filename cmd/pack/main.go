package main

import (
	"github.com/moby/buildkit/frontend/gateway/grpcclient"
	"github.com/moby/buildkit/util/appcontext"
	"github.com/sirupsen/logrus"
	pack "github.com/tonistiigi/buildkit-pack"
)

func main() {
	if err := grpcclient.RunFromEnvironment(appcontext.Context(), pack.Build); err != nil {
		logrus.Errorf("fatal error: %+v", err)
		panic(err)
	}
}
