package cmd

import (
	"fmt"

	"github.com/ripx80/brave/exit"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile string
	cfgMove bool
)

var rootCmd = &cobra.Command{
	Use:   "wgc",
	Short: "wireguard composer",
	Long:  `wireguard composer is a tool to manage your wg setups`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	defer exit.Safe()
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		exit.Exit(1)
	}
}

func init() {
	defer exit.Safe()
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is wgc.json)")
	//rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	//rootCmd.PersistentFlags().BoolVarP(&cfgMove, "move", "m", false, "move files to dst")
	rootCmd.AddCommand(connectCmd)
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		viper.AddConfigPath(".")
		viper.SetConfigName("wgc.json")
	}
	viper.AutomaticEnv()
	// problems with genkey
	if err := viper.ReadInConfig(); err != nil {
		return
		//fmt.Println("canot read config file:", viper.ConfigFileUsed())
	}
}
