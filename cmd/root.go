package cli

import (
	"fmt"
	"os"
	"strings"

	"github.com/Tinkoff/node-shell/pkg/plugin"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"k8s.io/cli-runtime/pkg/genericclioptions"
)

var (
	KubernetesConfigFlags *genericclioptions.ConfigFlags
)

func RootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:           "node-shell",
		Short:         "nsh",
		Long:          `.`,
		SilenceErrors: true,
		SilenceUsage:  true,
		PreRun: func(cmd *cobra.Command, args []string) {
			viper.BindPFlags(cmd.Flags())
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			p, err := plugin.NewPlugin(cmd.Flags())
			if err != nil {
				return fmt.Errorf("failed to initialize plugin: %w", err)
			}

			if err = p.RunPlugin(KubernetesConfigFlags); err != nil {
				return fmt.Errorf("failed to run plugin: %w", err)
			}

			return nil
		},
	}

	cobra.OnInitialize(initConfig)

	cmd.Flags().String("debug-image", "", "set an image for debug container")
	cmd.Flags().String("nodename", "", "set a name of a node to run debug pod on")
	cmd.Flags().String("ips", "image-pull-secret", "set a name of image pull secret to pull debug image")
	cmd.Flags().String("cpu", "500m", "set cpu request and limit for debug container")
	cmd.Flags().String("mem", "256Mi", "set memory request and limit for debug container")
	cmd.Flags().StringSlice("caps", []string{"NET_ADMIN", "SYS_ADMIN", "SYS_PTRACE"}, "set capabilities to add to debug container")
	cmd.MarkFlagRequired("debug-image")

	KubernetesConfigFlags = genericclioptions.NewConfigFlags(false)
	KubernetesConfigFlags.AddFlags(cmd.Flags())

	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	return cmd
}

func InitAndExecute() {
	if err := RootCmd().Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func initConfig() {
	viper.AutomaticEnv()
}
