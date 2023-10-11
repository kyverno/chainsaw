package test

import (
	"fmt"
	"os"
	"time"

	"github.com/kyverno/chainsaw/pkg/apis/v1alpha1"
	"github.com/kyverno/chainsaw/pkg/config"
	flagutils "github.com/kyverno/chainsaw/pkg/utils/flag"
	"github.com/spf13/cobra"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type options struct {
	config  string
	timeout metav1.Duration
}

func Command() *cobra.Command {
	var options options
	cmd := &cobra.Command{
		Use:   "test [flags]... [test directories]...",
		Short: "Stronger tool for e2e testing",
		RunE: func(cmd *cobra.Command, _ []string) error {
			out := cmd.OutOrStdout()
			var configuration v1alpha1.Configuration
			// if no config file was provided, give a chance to the default config name
			if options.config == "" {
				if _, err := os.Stat(config.DefaultFileName); err == nil {
					options.config = config.DefaultFileName
					fmt.Fprintf(out, "No configuration provided but found default file: %s\n", options.config)
				}
			}
			// try to load configuration file
			if options.config != "" {
				fmt.Fprintf(out, "Loading config (%s)...\n", options.config)
				config, err := config.Load(options.config)
				if err != nil {
					return err
				}
				configuration = *config
			} else {
				fmt.Fprintln(out, "Running without configuration")
			}
			// flags take precedence over configuration file
			flags := cmd.Flags()
			if flagutils.IsSet(flags, "duration") {
				configuration.Spec.Timeout = &options.timeout
			}
			// run tests
			fmt.Fprintln(out, "Running tests...")
			fmt.Fprintln(out, "- TODO")
			// done
			fmt.Fprintln(out, "Done.")
			return nil
		},
	}
	cmd.Flags().DurationVar(&options.timeout.Duration, "duration", 30*time.Second, "The duration to use as default for configuration.")
	cmd.Flags().StringVar(&options.config, "config", "", "Chainsaw configuration file.")
	// TODO: panic ?
	if err := cmd.MarkFlagFilename("config"); err != nil {
		panic(err)
	}
	return cmd
}
