// Code generated by mockery v3.2.5. DO NOT EDIT.

package mocks

import (
	context "context"

	grpc "google.golang.org/grpc"
	mock "github.com/stretchr/testify/mock"
	pb "ride-service/pb/proto/ride"
)

// RideServiceClient is an autogenerated mock type for the RideServiceClient type
type RideServiceClient struct {
	mock.Mock
}

// CreateRide provides a mock function with given fields: ctx, in, opts
func (_m *RideServiceClient) CreateRide(ctx context.Context, in *pb.CreateRideRequest, opts ...grpc.CallOption) (*pb.CreateRideResponse, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *pb.CreateRideResponse
	if rf, ok := ret.Get(0).(func(context.Context, *pb.CreateRideRequest, ...grpc.CallOption) *pb.CreateRideResponse); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*pb.CreateRideResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *pb.CreateRideRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetRide provides a mock function with given fields: ctx, in, opts
func (_m *RideServiceClient) GetRide(ctx context.Context, in *pb.GetRideRequest, opts ...grpc.CallOption) (*pb.Ride, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *pb.Ride
	if rf, ok := ret.Get(0).(func(context.Context, *pb.GetRideRequest, ...grpc.CallOption) *pb.Ride); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*pb.Ride)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *pb.GetRideRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UpdateRide provides a mock function with given fields: ctx, in, opts
func (_m *RideServiceClient) UpdateRide(ctx context.Context, in *pb.UpdateRideRequest, opts ...grpc.CallOption) (*pb.UpdateRideResponse, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *pb.UpdateRideResponse
	if rf, ok := ret.Get(0).(func(context.Context, *pb.UpdateRideRequest, ...grpc.CallOption) *pb.UpdateRideResponse); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*pb.UpdateRideResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *pb.UpdateRideRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
