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

package servicehandler

import (
	"context"
	"encoding/json"

	"github.com/fnrunner/fnproto/pkg/service/servicepb"
	ipamv1alpha1 "github.com/nokia/k8s-ipam/apis/ipam/v1alpha1"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

func (r *subServer) ApplyResource(ctx context.Context, req *servicepb.FunctionServiceRequest) (*servicepb.FunctionServiceResponse, error) {
	r.l.Info("service apply", "req", req)

	cr := &ipamv1alpha1.IPAllocation{}
	if err := json.Unmarshal([]byte(req.Resource), cr); err != nil {
		r.l.Error(err, "cannot unmarshal service apply req")
		return nil, err
	}

	resp, err := r.ipamclient.AllocateIPPrefix(ctx, cr, nil)
	if err != nil {
		r.l.Error(err, "cannot allocate prefix")
		return nil, err
	}
	b, err := json.Marshal(resp)
	if err != nil {
		r.l.Error(err, "cannot marshal response")
		return nil, err
	}
	r.l.Info("service apply success", "resp", string(b))
	return &servicepb.FunctionServiceResponse{Resource: string(b)}, nil

}

func (r *subServer) DeleteResource(ctx context.Context, req *servicepb.FunctionServiceRequest) (*emptypb.Empty, error) {
	r.l.Info("service delete", "req", req)

	cr := &ipamv1alpha1.IPAllocation{}
	if err := json.Unmarshal([]byte(req.Resource), cr); err != nil {
		r.l.Error(err, "cannot unmarshal service delete req")
		return nil, err
	}

	if err := r.ipamclient.DeAllocateIPPrefix(ctx, cr, nil); err != nil {
		r.l.Error(err, "cannot deallocate prefix")
		return nil, err
	}
	r.l.Info("service delete success")
	return &emptypb.Empty{}, nil
}
