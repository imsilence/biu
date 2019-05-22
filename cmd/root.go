package cmd

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"biu/config"
	"biu/plugin"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var worker int
var host string
var cidr string
var pluginName string
var query bool
var verbose bool
var timeout time.Duration

var rootCmd = &cobra.Command{
	Use:   "biu",
	Short: "A brief description of your application",
	Long:  `Security Scan Framework For Enterprise Intranet Based Services`,
	RunE: func(cmd *cobra.Command, args []string) error {
		start := time.Now()
		fmt.Printf("Start Time: %s\n", start.Format("2006-01-02 15:04:05"))
		if verbose {
			logrus.SetLevel(logrus.DebugLevel)
		}

		file, err := os.Create(filepath.Join(config.HOME, "logs", "biu.log"))
		if err != nil {
			return err
		}
		defer file.Close()
		writer := bufio.NewWriter(file)
		defer writer.Flush()
		logrus.SetOutput(writer)

		manager := plugin.Manager{Timeout: timeout}

		manager.Load()

		plugins := manager.Search(pluginName)

		fmt.Printf("Total Plugin: %5d,    Query Plugins: %5d\n", len(manager.Plugins), len(plugins))

		if query {
			tpl := "| %010s.%-50s |"
			fmt.Println(strings.Repeat("=", len(fmt.Sprintf(tpl, "", ""))))
			for _, p := range plugins {
				fmt.Printf(tpl, p.Id, p.Name)
				fmt.Println()
			}
			fmt.Println(strings.Repeat("=", len(fmt.Sprintf(tpl, "", ""))))
		} else {
			targets := manager.ParseTargets(host, cidr)

			fmt.Printf("Total Target: %5d\n", len(targets))

			logrus.WithFields(logrus.Fields{
				"targets": targets,
			}).Debug("check targets")

			results := manager.Execute(targets, plugins, worker)
			tpl := "| %-50s | %80s |"
			fmt.Println(strings.Repeat("=", len(fmt.Sprintf(tpl, "", ""))))
			fmt.Printf("| %-133s |\n", "Results:")
			fmt.Println(strings.Repeat("-", len(fmt.Sprintf(tpl, "", ""))))
			for result := range results {
				if r, ok := result.(plugin.Result); ok && !r.Safety {
					fmt.Printf(tpl, r.Plugin.Name, r.Target.Raw)
					fmt.Println()
				}
			}
			fmt.Println(strings.Repeat("=", len(fmt.Sprintf(tpl, "", ""))))

		}

		end := time.Now()
		fmt.Printf("Over Time: %s / Total Time: %s \n", end.Format("2006-01-02 15:04:05"), end.Sub(start))
		return nil
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().IntVarP(&worker, "worker", "w", 10000, "worker")
	rootCmd.Flags().StringVarP(&host, "host", "H", "", "target host addr")
	rootCmd.Flags().StringVarP(&cidr, "cidr", "c", "", "target cidr addr")
	rootCmd.Flags().StringVarP(&pluginName, "plugin", "p", "*", "plugin")
	rootCmd.Flags().BoolVarP(&query, "query", "q", false, "query plugin")
	rootCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "verbose run log")
	rootCmd.Flags().DurationVarP(&timeout, "timeout", "T", 3*time.Second, "http request timeout")
}
