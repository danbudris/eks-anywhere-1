package v1alpha1_test

import (
	"testing"

	"github.com/aws/eks-anywhere/release/api/v1alpha1"
)

func TestImageVersionedImage(t *testing.T) {
	tests := []struct {
		testName string
		URI      string
		want     string
	}{
		{
			testName: "full uri",
			URI:      "public.ecr.aws/l0g8r8j6/kubernetes-sigs/kind/node:v1.20.4-eks-d-1-20-1-eks-a-0.0.1.build.38",
			want:     "public.ecr.aws/l0g8r8j6/kubernetes-sigs/kind/node:v1.20.4-eks-d-1-20-1-eks-a-0.0.1.build.38",
		},
		{
			testName: "full uri with port",
			URI:      "public.ecr.aws:8484/l0g8r8j6/kubernetes-sigs/kind/node:v1.20.4-eks-d-1-20-1-eks-a-0.0.1.build.38",
			want:     "public.ecr.aws:8484/l0g8r8j6/kubernetes-sigs/kind/node:v1.20.4-eks-d-1-20-1-eks-a-0.0.1.build.38",
		},
	}
	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			i := v1alpha1.Image{
				URI: tt.URI,
			}
			if got := i.VersionedImage(); got != tt.want {
				t.Errorf("Image.VersionedImage() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestImageImage(t *testing.T) {
	tests := []struct {
		testName string
		URI      string
		want     string
	}{
		{
			testName: "full uri",
			URI:      "public.ecr.aws/l0g8r8j6/kubernetes-sigs/kind/node:v1.20.4-eks-d-1-20-1-eks-a-0.0.1.build.38",
			want:     "public.ecr.aws/l0g8r8j6/kubernetes-sigs/kind/node",
		},
		{
			testName: "full uri with port",
			URI:      "public.ecr.aws:8484/l0g8r8j6/kubernetes-sigs/kind/node:v1.20.4-eks-d-1-20-1-eks-a-0.0.1.build.38",
			want:     "public.ecr.aws:8484/l0g8r8j6/kubernetes-sigs/kind/node",
		},
		{
			testName: "no tag",
			URI:      "public.ecr.aws/l0g8r8j6/kubernetes-sigs/kind/node",
			want:     "public.ecr.aws/l0g8r8j6/kubernetes-sigs/kind/node",
		},
	}
	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			i := v1alpha1.Image{
				URI: tt.URI,
			}
			if got := i.Image(); got != tt.want {
				t.Errorf("Image.Image() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestImageTag(t *testing.T) {
	tests := []struct {
		testName string
		URI      string
		want     string
	}{
		{
			testName: "full uri",
			URI:      "public.ecr.aws/l0g8r8j6/kubernetes-sigs/kind/node:v1.20.4-eks-d-1-20-1-eks-a-0.0.1.build.38",
			want:     "v1.20.4-eks-d-1-20-1-eks-a-0.0.1.build.38",
		},
		{
			testName: "full uri with port",
			URI:      "public.ecr.aws:8484/l0g8r8j6/kubernetes-sigs/kind/node:v1.20.4-eks-d-1-20-1-eks-a-0.0.1.build.38",
			want:     "v1.20.4-eks-d-1-20-1-eks-a-0.0.1.build.38",
		},
		{
			testName: "no tag",
			URI:      "public.ecr.aws/l0g8r8j6/kubernetes-sigs/kind/node",
			want:     "",
		},
		{
			testName: "empty tag",
			URI:      "public.ecr.aws/l0g8r8j6/kubernetes-sigs/kind/node:",
			want:     "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			i := v1alpha1.Image{
				URI: tt.URI,
			}
			if got := i.Tag(); got != tt.want {
				t.Errorf("Image.Tag() = %v, want %v", got, tt.want)
			}
		})
	}
}
