package main

import (
	"context"
	"testing"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	testclient "k8s.io/client-go/kubernetes/fake"
)

func TestOptions_namespaceExist(t *testing.T) {
	objects := getTestNamespaces()
	namespaces = []string{"ns-01"}

	tests := []struct {
		name    string
		want    bool
		wantErr bool
	}{
		{
			name:    "check ns-01 exists",
			want:    true,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := testclient.NewSimpleClientset(objects...)
			o := Options{
				client: client,
			}
			got, err := o.namespaceExist(context.Background())
			if (err != nil) != tt.wantErr {
				t.Errorf("Options.namespaceExist() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Options.namespaceExist() = %v, want %v", got, tt.want)
			}
		})
	}
}

func getTestNamespaces() []runtime.Object {

	return []runtime.Object{
		&v1.Namespace{
			ObjectMeta: metav1.ObjectMeta{
				Name: "ns-01",
			},
		},
		&v1.Namespace{
			ObjectMeta: metav1.ObjectMeta{
				Name: "ns-02",
			},
		},
	}
}
