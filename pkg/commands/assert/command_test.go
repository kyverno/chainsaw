package assert

import (
	"bytes"
	"context"
	"path"
	"testing"
	"time"

	fakeClient "github.com/kyverno/chainsaw/pkg/client/testing"
	fakeNamespacer "github.com/kyverno/chainsaw/pkg/runner/namespacer/testing"
	"github.com/spf13/cobra"
	testify "github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	ctrlclient "sigs.k8s.io/controller-runtime/pkg/client"
)

func Test_preRunE(t *testing.T) {
	tests := []struct {
		name        string
		args        []string
		initialOpt  *options
		expectedOpt *options
		wantErr     bool
		errMsg      string
	}{
		{
			name: "No args and no filePath set",
			args: []string{},
			initialOpt: &options{
				filePath: "",
			},
			wantErr: true,
			errMsg:  "either a file path as an argument or the --file flag must be provided",
		},
		{
			name: "Args provided and filePath not set",
			args: []string{"./path/to/config.yaml"},
			initialOpt: &options{
				filePath: "",
			},
			expectedOpt: &options{
				filePath: "./path/to/config.yaml",
			},
			wantErr: false,
		},
		{
			name: "filePath already set, no args",
			args: []string{},
			initialOpt: &options{
				filePath: "./already/set/path.yaml",
			},
			expectedOpt: &options{
				filePath: "./already/set/path.yaml",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := &cobra.Command{}
			err := preRunE(tt.initialOpt, cmd, tt.args)
			if tt.wantErr {
				testify.Error(t, err)
				testify.Equal(t, tt.errMsg, err.Error())
			} else {
				testify.NoError(t, err)
				testify.Equal(t, tt.expectedOpt.filePath, tt.initialOpt.filePath)
			}
		})
	}
}

func Test_runE(t *testing.T) {
	basePath := path.Join("..", "..", "..", "testdata", "commands", "assert")
	tests := []struct {
		name       string
		setupFunc  func() *cobra.Command
		opts       options
		client     *fakeClient.FakeClient
		nspacer    *fakeNamespacer.FakeNamespacer
		wantErrMsg string
		wantErr    bool
	}{
		{
			name: "Success case - file input",
			setupFunc: func() *cobra.Command {
				cmd := &cobra.Command{}
				cmd.Args = cobra.RangeArgs(0, 1)
				cmd.SilenceUsage = true
				cmd.SetOut(bytes.NewBufferString(""))
				return cmd
			},
			opts: options{
				filePath:  path.Join(basePath, "assert.yaml"),
				noColor:   true,
				namespace: "default",
				timeout:   metav1.Duration{Duration: 5 * time.Second},
			},
			client: &fakeClient.FakeClient{
				GetFn: func(ctx context.Context, call int, key ctrlclient.ObjectKey, obj ctrlclient.Object, opts ...ctrlclient.GetOption) error {
					obj.(*unstructured.Unstructured).Object = map[string]any{
						"apiVersion": "v1",
						"kind":       "ConfigMap",
						"metadata": map[string]any{
							"name": "quick-start",
						},
						"data": map[string]any{
							"foo": "bar",
						},
					}
					return nil
				},
			},
			nspacer: &fakeNamespacer.FakeNamespacer{
				ApplyFn: func(obj ctrlclient.Object, call int) error {
					return nil
				},
				GetNamespaceFn: func(call int) string {
					return "default"
				},
			},
			wantErr: false,
		},
		{
			name: "Failure case - file input",
			setupFunc: func() *cobra.Command {
				cmd := &cobra.Command{}
				cmd.Args = cobra.RangeArgs(0, 1)
				cmd.SilenceUsage = true
				cmd.SetOut(bytes.NewBufferString(""))
				return cmd
			},
			opts: options{
				filePath:  path.Join(basePath, "assert.yaml"),
				noColor:   true,
				namespace: "default",
				timeout:   metav1.Duration{Duration: 5 * time.Second},
			},
			client: &fakeClient.FakeClient{
				GetFn: func(ctx context.Context, call int, key ctrlclient.ObjectKey, obj ctrlclient.Object, opts ...ctrlclient.GetOption) error {
					obj.(*unstructured.Unstructured).Object = map[string]any{
						"apiVersion": "v1",
						"kind":       "ConfigMap",
						"metadata": map[string]any{
							"name": "quick-start",
						},
						"data": map[string]any{
							"bar": "foo",
						},
					}
					return nil
				},
			},
			nspacer: &fakeNamespacer.FakeNamespacer{
				ApplyFn: func(obj ctrlclient.Object, call int) error {
					return nil
				},
				GetNamespaceFn: func(call int) string {
					return "default"
				},
			},
			wantErrMsg: "assertion failed: ------------------------\n" +
				"v1/ConfigMap/quick-start\n" +
				"------------------------\n" +
				"* data.foo: Invalid value: \"null\": Expected value: \"bar\"\n\n" +
				"--- expected\n" +
				"+++ actual\n" +
				"@@ -1,6 +1,5 @@\n" +
				" apiVersion: v1\n" +
				"-data:\n" +
				"-  foo: bar\n" +
				"+data: {}\n" +
				" kind: ConfigMap\n" +
				" metadata:\n" +
				"   name: quick-start\n",
			wantErr: true,
		},
		{
			name: "Success case - stdin input",
			setupFunc: func() *cobra.Command {
				cmd := &cobra.Command{}
				cmd.Args = cobra.RangeArgs(0, 1)
				cmd.SilenceUsage = true
				cmd.SetOut(bytes.NewBufferString(""))
				cmd.SetIn(bytes.NewBufferString(`apiVersion: v1
kind: ConfigMap
metadata:
  name: quick-start
data:
  foo: bar`))
				return cmd
			},
			opts: options{
				filePath:  "-",
				noColor:   true,
				namespace: "default",
				timeout:   metav1.Duration{Duration: 5 * time.Second},
			},
			client: &fakeClient.FakeClient{
				GetFn: func(ctx context.Context, call int, key ctrlclient.ObjectKey, obj ctrlclient.Object, opts ...ctrlclient.GetOption) error {
					obj.(*unstructured.Unstructured).Object = map[string]any{
						"apiVersion": "v1",
						"kind":       "ConfigMap",
						"metadata": map[string]any{
							"name": "quick-start",
						},
						"data": map[string]any{
							"foo": "bar",
						},
					}
					return nil
				},
			},
			nspacer: &fakeNamespacer.FakeNamespacer{
				ApplyFn: func(obj ctrlclient.Object, call int) error {
					return nil
				},
				GetNamespaceFn: func(call int) string {
					return "default"
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := tt.setupFunc()
			err := runE(tt.opts, cmd, tt.client, tt.nspacer)
			if tt.wantErr {
				testify.Error(t, err)
				if err != nil {
					testify.Equal(t, tt.wantErrMsg, err.Error())
				}
			} else {
				testify.NoError(t, err)
			}
		})
	}
}
