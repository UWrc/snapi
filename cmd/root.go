package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"os/user"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var url string = "https://uwconnect.uw.edu/yavin.do"
var KeyFile string
var RecordNumber string
var CI string

type Payload struct {
	RecordNumber string `json:"RecordNumber"`
	CI           string `json:"CI"`
	WorkNotes    string `json:"WorkNotes"`
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "snapi [flags] [command]",
	Short: "A command line tool to interact with the ServiceNow API.",
	Long:  ``,
	Args:  cobra.ArbitraryArgs,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			fmt.Println("Error: you must specify a command to execute.")
			os.Exit(1)
		}

		// Read in API credentials.
		username := viper.Get("SNAPI_USERNAME").(string)
		password := viper.Get("SNAPI_PASSWORD").(string)

		WorkNotes := GetWorkNotes(args)

		CIs := map[string]string{
			"hyak":  "Shared HPC Cluster (Hyak)",
			"kopah": "Kopah",
			"lolo":  "Shared Central File System (lolo)",
		}

		// Create the JSON payload.
		// https://uwconnect.uw.edu/kb_view.do?sysparm_article=KB0025022
		data := Payload{
			RecordNumber: RecordNumber,
			CI:           CIs[CI],
			WorkNotes:    WorkNotes,
		}
		fmt.Printf("%v\n", data)
		os.Exit(1)
		payloadBytes, _ := json.Marshal(data)
		req, _ := http.NewRequest("POST", url, bytes.NewBuffer(payloadBytes))

		req.Header.Set("Content-Type", "application/json")
		req.SetBasicAuth(username, password)

		// Send the API call
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			log.Fatalf("Error sending request: %v", err)
		}
		defer resp.Body.Close()

		if resp.Status != "200 OK" {
			fmt.Println("Error: Command was run, record was NOT updated.")
			body, _ := io.ReadAll(resp.Body)
			log.Fatalf("Error: %s\n%v\n", resp.Status, body)
		} else {
			fmt.Printf("%s updated.\n", RecordNumber)
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() {
	rootCmd.Flags().SetInterspersed(false)
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// This function runs before main() and sets up the application.
	cobra.OnInitialize(GetCredentials)

	rootCmd.PersistentFlags().StringVarP(&KeyFile, "key", "k", ".snapi", "config file")
	rootCmd.PersistentFlags().StringVarP(&CI, "configuration-item", "c", "hyak", "Configuration item (required).")

	rootCmd.Flags().StringVarP(&RecordNumber, "record", "r", "", "Service Now record number (required).")
	rootCmd.MarkFlagRequired("record")
}

func GetCredentials() {
	if KeyFile != "" {
		viper.SetConfigFile(KeyFile)
		viper.SetConfigType("dotenv")
	}

	err := viper.ReadInConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading config file: %s\n", err)
		os.Exit(1)
	}
}

func GetWorkNotes(args []string) string {
	// Get context of command being run.
	user, _ := user.Current()
	hostname, _ := os.Hostname()
	cwd, _ := os.Getwd()

	// STDOUT buffer
	var stdoutBuf bytes.Buffer
	stdoutWriter := io.MultiWriter(os.Stdout, &stdoutBuf)

	// STDERR buffer
	var stderrBuf bytes.Buffer
	stderrWriter := io.MultiWriter(os.Stderr, &stderrBuf)

	// Run command and save output.
	cmd := exec.Command(args[0], args[1:]...)
	cmd.Stdout = stdoutWriter
	cmd.Stderr = stderrWriter
	cmd.Run()

	// Capture outputs, make them HTML-friendly.
	capturedOutput := strings.ReplaceAll(stdoutBuf.String(), "\n", "<br />")
	capturedError := strings.ReplaceAll(stderrBuf.String(), "\n", "<br />")

	// Generate SN WorkNotes from the above information.
	WorkNotes := fmt.Sprintf("Command run on <code>%s</code> by <code>%s</code> in <code>%s</code>.<br /><br /><b><code>%s</code></b><br /><br /><code>--- STDOUT ---</code><br /><br /><pre>%s</pre><br /><br /><code>--- STDERR ---</code><br /><br /><pre>%s</pre><br />", hostname, user.Username, cwd, strings.Join(args, " "), capturedOutput, capturedError)

	return WorkNotes
}
