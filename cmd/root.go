package cmd

import (
	"bufio"
	"encoding/csv"
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
var target string
var cidr string
var pocName string
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

		logDir := filepath.Join(config.HOME, "logs")
		reportDir := filepath.Join(config.HOME, "reports")

		if verbose {
			logDir = filepath.Join(config.CWD, "logs")
			reportDir = filepath.Join(config.CWD, "reports")
			logrus.SetLevel(logrus.DebugLevel)
		}

		os.MkdirAll(logDir, os.ModePerm)
		os.MkdirAll(reportDir, os.ModePerm)

		file, err := os.Create(filepath.Join(logDir, "biu.log"))
		if err != nil {
			return err
		}
		defer file.Close()
		writer := bufio.NewWriter(file)
		defer writer.Flush()
		logrus.SetOutput(writer)

		manager := plugin.Manager{Timeout: timeout}

		manager.Load()

		pocs := manager.Search(pocName)

		fmt.Printf("Total Poc: %5d,    Query POCs: %5d\n", len(manager.POCs), len(pocs))

		if query {
			tpl := "| %-8s | %-50s | %-100s |"
			fmt.Println(strings.Repeat("=", len(fmt.Sprintf(tpl, "", "", ""))))
			fmt.Printf(tpl, "ID", "Name", "Desc")
			fmt.Println()
			fmt.Println(strings.Repeat("-", len(fmt.Sprintf(tpl, "", "", ""))))
			if len(pocs) > 0 {
				for _, p := range pocs {
					fmt.Printf(tpl, p.ID, p.Name, p.Desc)
					fmt.Println()
				}
			} else {
				fmt.Printf("| %-[1]*[2]s |\n", len(fmt.Sprintf(tpl, "", "", ""))-4, "Empty")
			}

			fmt.Println(strings.Repeat("=", len(fmt.Sprintf(tpl, "", "", ""))))
		} else {
			reportPath := filepath.Join(reportDir, fmt.Sprintf("%s.csv", time.Now().Format("2006-01-02_15-04-05")))
			reportFile, err := os.Create(reportPath)
			if err != nil {
				return err
			}
			defer reportFile.Close()
			reportWriter := csv.NewWriter(reportFile)
			defer reportWriter.Flush()
			reportWriter.Write([]string{"URL", "POC ID", "POC Name", "POC Desc", "Request"})

			targets := manager.ParseTargets(target, cidr)
			fmt.Printf("Total Target: %5d\n", len(targets))

			logrus.WithFields(logrus.Fields{
				"targets": targets,
			}).Debug("check targets")

			results := manager.Execute(targets, pocs, worker)
			tpl := "| %-45s | %-60s | %-100s |"
			fmt.Println(strings.Repeat("=", len(fmt.Sprintf(tpl, "", "", ""))))
			fmt.Printf("| %-[1]*[2]s |\n", len(fmt.Sprintf(tpl, "", "", ""))-4, "Results:")
			fmt.Println(strings.Repeat("-", len(fmt.Sprintf(tpl, "", "", ""))))
			empty := true

			for result := range results {
				if r, ok := result.(plugin.Result); ok && !r.Safety {
					reportWriter.Write([]string{r.Target.Raw, r.Plugin.POC.ID, r.Plugin.POC.Name, r.Plugin.POC.Desc, r.Plugin.Request.String()})
					empty = false
					fmt.Printf(tpl, r.Plugin.POC.Name, r.Target.Raw, r.Plugin.Request)
					fmt.Println()
				}
			}
			if empty {
				fmt.Printf("| %-133s |\n", "Empty")
			}
			fmt.Println(strings.Repeat("=", len(fmt.Sprintf(tpl, "", "", ""))))
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
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "verbose run log")
	rootCmd.PersistentFlags().DurationVarP(&timeout, "timeout", "T", 3*time.Second, "http request timeout")
	rootCmd.Flags().StringVarP(&target, "target", "t", "", "target addr(127.0.0.1)")
	rootCmd.Flags().StringVarP(&cidr, "cidr", "c", "", "target cidr addr(192.168.0.0/24)")
	rootCmd.Flags().StringVarP(&pocName, "poc", "p", "*", "poc")
	rootCmd.Flags().BoolVarP(&query, "query", "q", false, "query poc")

}
