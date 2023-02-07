/*
Copyright 2022 Nokia.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"context"
	"fmt"
	"os"
	"strconv"

	fnrunv1alpha1 "github.com/fnrunner/fnruntime/apis/fnrun/v1alpha1"
	"github.com/fnrunner/fnsvcsdk/grpcserver"
	"github.com/fnrunner/fnsvcsdk/healthhandler"
	"github.com/fnrunner/ipam-injector-service/internal/servicehandler"
	"github.com/henderiw-k8s-lcnc/discovery/discovery"
	"github.com/henderiw-k8s-lcnc/discovery/registrator"
	"go.uber.org/zap/zapcore"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
)

func main() {

	opts := zap.Options{
		Development: true,
		TimeEncoder: zapcore.ISO8601TimeEncoder,
	}
	ctrl.SetLogger(zap.New(zap.UseFlagOptions(&opts)))
	l := ctrl.Log.WithName("ipamInjectorService")

	ctx := context.Background()

	reg, err := registrator.New(ctx, ctrl.GetConfigOrDie(), &registrator.Options{
		ServiceDiscovery:          discovery.ServiceDiscoveryTypeK8s,
		ServiceDiscoveryNamespace: "ipam",
	})
	if err != nil {
		l.Error(err, "Cannot create registrator")
		os.Exit(1)
	}

	hh := healthhandler.New()
	sh := servicehandler.New(ctx, reg)

	address := fmt.Sprintf(":%s", strconv.Itoa(fnrunv1alpha1.FnGRPCServerPort))
	if os.Getenv("FN_SERVICE_PORT") != "" {
		address = fmt.Sprintf(":%s", os.Getenv("FN_SERVICE_PORT"))
	}
	l.Info("grpc server", "address", address)

	s := grpcserver.New(grpcserver.Config{
		Address:  address,
		Insecure: true,
	},
		grpcserver.WithServiceApplyResourceHandler(sh.ApplyResource),
		grpcserver.WithServiceDeleteResourceHandler(sh.DeleteResource),
		grpcserver.WithWatchHandler(hh.Watch),
		grpcserver.WithCheckHandler(hh.Check),
	)

	if err := s.Start(); err != nil {
		l.Error(err, "cannot start grpcserver")
		os.Exit(1)
	}
}
