package cleaner

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/kyverno/chainsaw/pkg/client"
	tclient "github.com/kyverno/chainsaw/pkg/client/testing"
	"github.com/kyverno/chainsaw/pkg/model"
	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	kerror "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/utils/ptr"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name    string
		timeout time.Duration
		delay   *time.Duration
		want    Cleaner
	}{{
		name:    "with timeout",
		timeout: time.Minute,
		delay:   nil,
		want: &cleaner{
			timeout:     time.Minute,
			delay:       nil,
			propagation: metav1.DeletePropagationBackground,
		},
	}, {
		name:    "with delay",
		timeout: time.Minute,
		delay:   ptr.To(10 * time.Second),
		want: &cleaner{
			timeout:     time.Minute,
			delay:       ptr.To(10 * time.Second),
			propagation: metav1.DeletePropagationBackground,
		},
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := New(tt.timeout, tt.delay, metav1.DeletePropagationBackground)
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_cleaner_Empty(t *testing.T) {
	tests := []struct {
		name    string
		entries []cleanupEntry
		want    bool
	}{{
		name:    "without entries",
		entries: nil,
		want:    true,
	}, {
		name:    "without entries",
		entries: []cleanupEntry{},
		want:    true,
	}, {
		name:    "with entries",
		entries: []cleanupEntry{{}},
		want:    false,
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &cleaner{
				entries: tt.entries,
			}
			got := c.Empty()
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_cleaner_Add(t *testing.T) {
	entry := cleanupEntry{
		object: &unstructured.Unstructured{},
		client: &tclient.FakeClient{},
	}
	tests := []struct {
		name    string
		entries []cleanupEntry
		client  client.Client
		object  client.Object
		want    []cleanupEntry
	}{{
		name:   "empty",
		client: entry.client,
		object: entry.object,
		want:   []cleanupEntry{entry},
	}, {
		name:    "with entries",
		entries: []cleanupEntry{entry},
		client:  entry.client,
		object:  entry.object,
		want:    []cleanupEntry{entry, entry},
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &cleaner{
				entries: tt.entries,
			}
			c.Add(tt.client, tt.object)
			assert.Equal(t, tt.want, c.entries)
		})
	}
}

func Test_cleaner_Run(t *testing.T) {
	obj := &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: "foo",
		},
	}
	tests := []struct {
		name    string
		entries []cleanupEntry
		want    []error
	}{{
		name: "empty",
	}, {
		name: "ok",
		entries: []cleanupEntry{{
			object: obj,
			client: &tclient.FakeClient{
				DeleteFn: func(ctx context.Context, call int, obj client.Object, opts ...client.DeleteOption) error {
					return nil
				},
				GetFn: func(ctx context.Context, call int, key client.ObjectKey, obj client.Object, opts ...client.GetOption) error {
					return kerror.NewNotFound(corev1.Resource("namespace"), "foo")
				},
			},
		}},
		want: nil,
	}, {
		name: "not found",
		entries: []cleanupEntry{{
			object: obj,
			client: &tclient.FakeClient{
				DeleteFn: func(ctx context.Context, call int, obj client.Object, opts ...client.DeleteOption) error {
					return kerror.NewNotFound(corev1.Resource("namespace"), "foo")
				},
			},
		}},
		want: nil,
	}, {
		name: "delete error",
		entries: []cleanupEntry{{
			object: obj,
			client: &tclient.FakeClient{
				DeleteFn: func(ctx context.Context, call int, obj client.Object, opts ...client.DeleteOption) error {
					return errors.New("dummy")
				},
			},
		}},
		want: []error{errors.New("dummy")},
	}, {
		name: "get error",
		entries: []cleanupEntry{{
			object: obj,
			client: &tclient.FakeClient{
				DeleteFn: func(ctx context.Context, call int, obj client.Object, opts ...client.DeleteOption) error {
					return nil
				},
				GetFn: func(ctx context.Context, call int, key client.ObjectKey, obj client.Object, opts ...client.GetOption) error {
					return errors.New("dummy")
				},
			},
		}},
		want: []error{errors.New("dummy")},
	}, {
		name: "failed",
		entries: []cleanupEntry{{
			object: obj,
			client: &tclient.FakeClient{
				DeleteFn: func(ctx context.Context, call int, obj client.Object, opts ...client.DeleteOption) error {
					return nil
				},
				GetFn: func(ctx context.Context, call int, key client.ObjectKey, obj client.Object, opts ...client.GetOption) error {
					return nil
				},
			},
		}},
		want: []error{context.DeadlineExceeded},
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &cleaner{
				delay:   ptr.To(time.Second),
				timeout: 1 * time.Second,
				entries: tt.entries,
			}
			got := c.Run(context.TODO(), &model.StepReport{})
			assert.Equal(t, tt.want, got)
		})
	}
}
