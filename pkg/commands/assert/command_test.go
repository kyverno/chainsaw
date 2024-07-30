package assert

import (
	"bytes"
	"context"
	"io"
	"os"
	"path"
	"path/filepath"
	"time"

	"github.com/kyverno/chainsaw/pkg/client"
	tclient "github.com/kyverno/chainsaw/pkg/client/testing"
	"github.com/kyverno/chainsaw/pkg/commands/root"
	fakeNamespacer "github.com/kyverno/chainsaw/pkg/engine/namespacer/testing"
	"github.com/kyverno/chainsaw/pkg/testing"
	"github.com/spf13/cobra"
	testify "github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

func Test_Execute(t *testing.T) {
	basePath := path.Join("..", "..", "..", "testdata", "commands", "assert")
	tests := []struct {
		name    string
		args    []string
		wantErr bool
		out     string
	}{{
		name: "help",
		args: []string{
			"assert",
			"--help",
		},
		out:     filepath.Join(basePath, "help.txt"),
		wantErr: false,
	}, {
		name:    "no args and no flags",
		args:    []string{"assert"},
		wantErr: true,
	}, {
		name: "unknow flag",
		args: []string{
			"assert",
			"--foo",
		},
		wantErr: true,
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := root.Command()
			cmd.AddCommand(Command())
			testify.NotNil(t, cmd)
			cmd.SetArgs(tt.args)
			out := bytes.NewBufferString("")
			cmd.SetOut(out)
			err := cmd.Execute()
			if tt.wantErr {
				testify.Error(t, err)
			} else {
				testify.NoError(t, err)
			}
			actual, err := io.ReadAll(out)
			testify.NoError(t, err)
			if tt.out != "" {
				expected, err := os.ReadFile(tt.out)
				testify.NoError(t, err)
				testify.Equal(t, string(expected), string(actual))
			}
		})
	}
}

func Test_preRunE(t *testing.T) {
	tests := []struct {
		name        string
		args        []string
		initialOpt  *options
		expectedOpt *options
		wantErr     bool
		errMsg      string
	}{{
		name: "No args and no filePath set",
		args: []string{},
		initialOpt: &options{
			assertPath: "",
		},
		wantErr: true,
		errMsg:  "either a file path as an argument or the --file flag must be provided",
	}, {
		name: "Args provided and filePath not set",
		args: []string{"./path/to/config.yaml"},
		initialOpt: &options{
			assertPath: "",
		},
		expectedOpt: &options{
			assertPath: "./path/to/config.yaml",
		},
		wantErr: false,
	}, {
		name: "filePath already set, no args",
		args: []string{},
		initialOpt: &options{
			assertPath: "./already/set/path.yaml",
		},
		expectedOpt: &options{
			assertPath: "./already/set/path.yaml",
		},
		wantErr: false,
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := &cobra.Command{}
			err := preRunE(tt.initialOpt, cmd, tt.args)
			if tt.wantErr {
				testify.Error(t, err)
				testify.Equal(t, tt.errMsg, err.Error())
			} else {
				testify.NoError(t, err)
				testify.Equal(t, tt.expectedOpt.assertPath, tt.initialOpt.assertPath)
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
		client     *tclient.FakeClient
		nspacer    *fakeNamespacer.FakeNamespacer
		wantErrMsg string
		wantErr    bool
	}{{
		name: "Success case - file input",
		setupFunc: func() *cobra.Command {
			cmd := &cobra.Command{}
			cmd.Args = cobra.RangeArgs(0, 1)
			cmd.SilenceUsage = true
			cmd.SetOut(bytes.NewBufferString(""))
			return cmd
		},
		opts: options{
			assertPath: path.Join(basePath, "assert.yaml"),
			noColor:    true,
			namespace:  "default",
			timeout:    metav1.Duration{Duration: 5 * time.Second},
		},
		client: &tclient.FakeClient{
			GetFn: func(ctx context.Context, call int, key client.ObjectKey, obj client.Object, opts ...client.GetOption) error {
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
			ApplyFn: func(call int, client client.Client, obj client.Object) error {
				return nil
			},
			GetNamespaceFn: func(call int) string {
				return "default"
			},
		},
		wantErr: false,
	}, {
		name: "Failure case - file input",
		setupFunc: func() *cobra.Command {
			cmd := &cobra.Command{}
			cmd.Args = cobra.RangeArgs(0, 1)
			cmd.SilenceUsage = true
			cmd.SetOut(bytes.NewBufferString(""))
			return cmd
		},
		opts: options{
			assertPath: path.Join(basePath, "assert.yaml"),
			noColor:    true,
			namespace:  "default",
			timeout:    metav1.Duration{Duration: 5 * time.Second},
		},
		client: &tclient.FakeClient{
			GetFn: func(ctx context.Context, call int, key client.ObjectKey, obj client.Object, opts ...client.GetOption) error {
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
			ApplyFn: func(call int, client client.Client, obj client.Object) error {
				return nil
			},
			GetNamespaceFn: func(call int) string {
				return "default"
			},
		},
		wantErrMsg: "assertion failed\n" +
			"------------------------\n" +
			"v1/ConfigMap/quick-start\n" +
			"------------------------\n" +
			"* data.foo: Required value: field not found in the input object\n\n" +
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
	}, {
		name: "Failure case - Non-exist file input",
		setupFunc: func() *cobra.Command {
			cmd := &cobra.Command{}
			cmd.Args = cobra.RangeArgs(0, 1)
			cmd.SilenceUsage = true
			cmd.SetOut(bytes.NewBufferString(""))
			return cmd
		},
		opts: options{
			assertPath: path.Join(basePath, "non-exist-file.yaml"),
			noColor:    true,
			namespace:  "default",
			timeout:    metav1.Duration{Duration: 5 * time.Second},
		},
		client: &tclient.FakeClient{
			GetFn: func(ctx context.Context, call int, key client.ObjectKey, obj client.Object, opts ...client.GetOption) error {
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
			ApplyFn: func(call int, client client.Client, obj client.Object) error {
				return nil
			},
			GetNamespaceFn: func(call int) string {
				return "default"
			},
		},
		wantErrMsg: "failed to load file '../../../testdata/commands/assert/non-exist-file.yaml': " +
			"no files found matching path: ../../../testdata/commands/assert/non-exist-file.yaml",
		wantErr: true,
	}, {
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
			assertPath: "-",
			noColor:    true,
			namespace:  "default",
			timeout:    metav1.Duration{Duration: 5 * time.Second},
		},
		client: &tclient.FakeClient{
			GetFn: func(ctx context.Context, call int, key client.ObjectKey, obj client.Object, opts ...client.GetOption) error {
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
			ApplyFn: func(call int, client client.Client, obj client.Object) error {
				return nil
			},
			GetNamespaceFn: func(call int) string {
				return "default"
			},
		},
		wantErr: false,
	}, {
		name: "Failure case - can't read from stdin",
		setupFunc: func() *cobra.Command {
			cmd := &cobra.Command{}
			cmd.Args = cobra.RangeArgs(0, 1)
			cmd.SilenceUsage = true
			cmd.SetOut(bytes.NewBufferString(""))
			cmd.SetIn(&testing.ErrReader{})
			return cmd
		},
		opts: options{
			assertPath: "-",
			noColor:    true,
			namespace:  "default",
			timeout:    metav1.Duration{Duration: 5 * time.Second},
		},
		client: &tclient.FakeClient{
			GetFn: func(ctx context.Context, call int, key client.ObjectKey, obj client.Object, opts ...client.GetOption) error {
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
			ApplyFn: func(call int, client client.Client, obj client.Object) error {
				return nil
			},
			GetNamespaceFn: func(call int) string {
				return "default"
			},
		},
		wantErrMsg: "failed to read from stdin: error reading from stdin",
		wantErr:    true,
	}}
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
