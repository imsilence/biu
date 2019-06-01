package cmd

import (
	"biu/config"
	"biu/plugin"
	"bufio"
	"html/template"
	"net/http"
	"os"
	"path/filepath"

	"biu/assets"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var address string

var apiCmd = &cobra.Command{
	Use:   "api",
	Short: "Biu Framwork web api server",
	Long:  `Biu Framwork web api server`,
	RunE: func(cmd *cobra.Command, args []string) error {
		logDir := filepath.Join(config.HOME, "logs")
		if verbose {
			logDir = filepath.Join(config.CWD, "logs")
			logrus.SetLevel(logrus.DebugLevel)
			gin.SetMode(gin.DebugMode)
		} else {
			gin.SetMode(gin.ReleaseMode)
		}

		os.MkdirAll(logDir, os.ModePerm)

		file, err := os.Create(filepath.Join(logDir, "api.log"))
		if err != nil {
			return err
		}
		defer file.Close()
		writer := bufio.NewWriter(file)
		defer writer.Flush()
		logrus.SetOutput(writer)

		manager := plugin.Manager{Timeout: timeout}

		manager.Load()

		r := gin.Default()

		tpl := template.New("")
		for _, name := range assets.AssetNames() {
			tpl, err = tpl.New(name).Parse(string(assets.MustAsset(name)))
			if err != nil {
				return err
			}
		}

		r.SetHTMLTemplate(tpl)

		r.GET("/task/", func(c *gin.Context) {
			c.HTML(http.StatusOK, "templates/task.html", nil)
		})

		r.POST("/task/", func(c *gin.Context) {

			pocName := c.DefaultPostForm("p", "*")
			host := c.DefaultPostForm("H", "")
			cidr := c.DefaultPostForm("c", "")

			logrus.WithFields(logrus.Fields{
				"poc":  pocName,
				"host": host,
				"cidr": cidr,
			}).Debug("create task")

			pocs := manager.Search(pocName)
			targets := manager.ParseTargets(host, cidr)

			logrus.WithFields(logrus.Fields{
				"pocs":    pocs,
				"targets": targets,
			}).Debug("execute task")

			results := manager.Execute(targets, pocs, worker)
			rs := make([]map[string]string, 0)
			for result := range results {
				if r, ok := result.(plugin.Result); ok && !r.Safety {
					rs = append(rs, map[string]string{
						"url":   r.Target.Raw,
						"pid":   r.Plugin.POC.ID,
						"pname": r.Plugin.POC.Name,
						"pdesc": r.Plugin.POC.Desc,
					})
				}
			}
			c.JSON(http.StatusOK, gin.H{
				"results": rs,
			})
		})
		logrus.Info("Listen on: ", address)
		r.Run(address)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(apiCmd)
	apiCmd.Flags().StringVarP(&address, "listen", "l", ":8080", "web server listen addr")
}
