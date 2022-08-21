package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/emersion/go-smtp"
	"github.com/manyou-io/smtp-post/server"
	"github.com/spf13/cobra"
)

var config *server.Config
var listenTls bool

var rootCmd = &cobra.Command{
	Use:   "smtp-post",
	Short: "Receive mail with SMTP and relay with HTTP",
}

var runCmd = &cobra.Command{
	Use:              "run [endpoint]",
	Short:            "Start the SMTP server",
	TraverseChildren: true,
	Args:             cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		config.Endpoint = args[0]

		s, err := config.CreateServer()
		if err != nil {
			log.Fatal(err)
			os.Exit(1)
		}

		log.Printf("server started on %s for %s\n", config.Addr, config.Domain)
		if err := listenAndServe(s); err != nil {
			log.Fatal(err)
			os.Exit(1)
		}
	},
}

func listenAndServe(s *smtp.Server) error {
	if listenTls {
		return s.ListenAndServeTLS()
	} else {
		return s.ListenAndServe()
	}
}

func init() {
	rootCmd.AddCommand(runCmd)
	config = &server.Config{}
	runCmd.Flags().StringVarP(&config.Addr, "bind", "b", ":1025", "host and port to listen on. For example: \":587\" or \"127.0.0.1:25\"")
	runCmd.Flags().StringVarP(&config.Domain, "domain", "d", "localhost", "hostname for greeting")
	runCmd.Flags().DurationVar(&config.ReadTimeout, "read-timeout", 10*time.Second, "read timeout in seconds")
	runCmd.Flags().DurationVar(&config.WriteTimeout, "write-timeout", 10*time.Second, "write timeout in seconds")
	runCmd.Flags().IntVar(&config.MaxMessageBytes, "max-size", 5*1024*1024, "maximum message size in bytes. Note that AWS Lambda has a 5MB payload size limit.")
	runCmd.Flags().IntVar(&config.MaxRecipients, "max-rcpt", 14, "maximum number of recipients per message")
	runCmd.Flags().StringVar(&config.CertFile, "tls-cert", "", "X509 certificate file for TLS")
	runCmd.Flags().StringVar(&config.KeyFile, "tls-key", "", "private key file for TLS")
	runCmd.Flags().StringVarP(&config.ApiKey, "api-key", "k", "", "value of X-Api-Key header")
	runCmd.Flags().StringVarP(&config.Username, "username", "u", "smtp-post", "username for SMTP authentication")
	runCmd.Flags().StringVarP(&config.Password, "password", "p", "smtp-post", "password for SMTP authentication")
	runCmd.Flags().BoolVar(&listenTls, "tls", false, "listen with TLS wrapper")
	runCmd.Flags().BoolVar(&config.AllowInsecureAuth, "allow-insecure-auth", false, "allow insecure authentication even if TLS is enabled")
	runCmd.MarkFlagsRequiredTogether("username", "password")
	runCmd.MarkFlagsRequiredTogether("tls-cert", "tls-key")
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
