package main

import (
	"context"
	"testing"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes"
	testclient "k8s.io/client-go/kubernetes/fake"
)

func getTestNamespaces() []runtime.Object {
	return []runtime.Object{
		&v1.Namespace{
			ObjectMeta: metav1.ObjectMeta{
				Name: "cert-manager",
			},
		},
		&v1.Namespace{
			ObjectMeta: metav1.ObjectMeta{
				Name: "default",
			},
		},
	}
}

func Test_doExpectedNamespacesExist(t *testing.T) {
	type args struct {
		ctx                context.Context
		client             kubernetes.Interface
		expectedNamespaces []string
		currentEnv         string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "The correct namespaces exist",
			args: args{
				ctx:                context.TODO(),
				client:             testclient.NewSimpleClientset(getTestNamespaces()...),
				expectedNamespaces: []string{"cert-manager", "default"},
				currentEnv:         "whatever",
			},
			wantErr: false,
		},
		{
			name: "The correct namespaces don't exist",
			args: args{
				ctx:                context.TODO(),
				client:             testclient.NewSimpleClientset(getTestNamespaces()...),
				expectedNamespaces: []string{"cert-manager", "default", "test"},
				currentEnv:         "whatever",
			},
			wantErr: true,
		},
		{
			name: "is not prod",
			args: args{
				ctx:                context.TODO(),
				client:             testclient.NewSimpleClientset(getTestNamespaces()...),
				expectedNamespaces: []string{"velero", "default"},
				currentEnv:         "whatever",
			},
			wantErr: false,
		},
		{
			name: "is prod",
			args: args{
				ctx:                context.TODO(),
				client:             testclient.NewSimpleClientset(getTestNamespaces()...),
				expectedNamespaces: []string{"velero", "default"},
				currentEnv:         "live",
			},
			wantErr: true,
		},
		{
			name: "is non-live",
			args: args{
				ctx:                context.TODO(),
				client:             testclient.NewSimpleClientset(getTestNamespaces()...),
				expectedNamespaces: []string{"overprovision", "default"},
				currentEnv:         "manager",
			},
			wantErr: false,
		},
		{
			name: "Bad client",
			args: args{
				ctx:                context.TODO(),
				client:             testclient.NewSimpleClientset(),
				expectedNamespaces: []string{"cert-manager", "default"},
				currentEnv:         "whatever",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := doExpectedNamespacesExist(tt.args.ctx, tt.args.client, tt.args.expectedNamespaces, tt.args.currentEnv); (err != nil) != tt.wantErr {
				t.Errorf("doExpectedNamespacesExist() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
