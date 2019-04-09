package metricsbridge

import (
	"github.com/spf13/cobra"

	"github.com/openshift/openshift-azure/pkg/entrypoint/config"
)

// NewCommand returns the cobra command for "azure-controllers".
func NewCommand() *cobra.Command {
	cc := &cobra.Command{
		Use:  "metricsbridge",
		Long: "Start metrics-bridge application",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := configFromCmd(cmd)
			if err != nil {
				return err
			}
			return start(cfg)
		},
	}
	cc.Flags().String("config", "", "config file location")
	cobra.MarkFlagRequired(cc.Flags(), "config")

	return cc
}

func configFromCmd(cmd *cobra.Command) (*Config, error) {
	c := &Config{}
	var err error
	c.Common, err = config.CommonConfigFromCmd(cmd)
	if err != nil {
		return nil, err
	}
	return c, nil
}
