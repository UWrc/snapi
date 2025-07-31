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

var version string = "0.0.3"
var url string = "https://uwconnect.uw.edu/yavin.do"
var KeyFile string
var RecordNumber string
var CI string
var state string
var AssignedTo string
var CommentWatch string
var NoteWatch string

type Payload struct {
	RecordNumber string `json:"RecordNumber"`
	CI           string `json:"CI"`
	State        string `json:"State"`
	WorkNotes    string `json:"WorkNotes"`
	AssignedTo   string `json:"AssignedTo"`
	CommentWatch   string `json:"WatchList"`
	NoteWatch   string `json:"WorkNotesList"`
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "snapi [flags] [command]",
	Short: fmt.Sprintf("snapi v%s: A command line tool to interact with the ServiceNow API.", version),
	Long:  ``,
	Args:  cobra.ArbitraryArgs,
	Run: func(cmd *cobra.Command, args []string) {
		// Read in API credentials.
		username := viper.Get("SNAPI_USERNAME").(string)
		password := viper.Get("SNAPI_PASSWORD").(string)

		WorkNotes := ""
		if len(args) > 0 {
			WorkNotes = GetWorkNotes(args)
		}

		CIs := map[string]string{
			"hyak":  "Shared HPC Cluster (Hyak)",
			"kopah": "Kopah",
			"lolo":  "Shared Central File System (lolo)",
		}

		states := map[string]string{
			//"n":    "new",
			//"new":  "new",
			"o":    "open",
			"open": "open",
			//"h":        "on hold",
			//"hold":     "on hold",
			"r":        "resolved",
			"resolve": "resolved",
			"resolved": "resolved",
		}

		// Create the JSON payload.
		// https://uwconnect.uw.edu/kb_view.do?sysparm_article=KB0025022
		data := Payload{
			RecordNumber: RecordNumber,
			CI:           CIs[CI],
			State:        states[state],
			WorkNotes:    WorkNotes,
			AssignedTo:   AssignedTo,
			CommentWatch: CommentWatch,
			NoteWatch: 	  NoteWatch,
		}
		//fmt.Printf("%v\n", data)
		//os.Exit(1)
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
			log.Fatalf("Error: %s\n%s\n", resp.Status, body)
		} else {
			fmt.Printf("%s updated [%s].\n", RecordNumber, strings.ToUpper(states[state]))
		}
		//body, _ := io.ReadAll(resp.Body)
		//fmt.Printf("Error: %s\n%s\n", resp.Status, body)
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

	rootCmd.PersistentFlags().StringVarP(&KeyFile, "key", "k", "", "API key file.")
	rootCmd.PersistentFlags().StringVarP(&CI, "configuration-item", "c", "hyak", "Configuration item (required).")
	rootCmd.PersistentFlags().StringVarP(&AssignedTo, "assigned-to", "a", "", "A single netID or email address for the primary contact working on the record.")
	rootCmd.PersistentFlags().StringVarP(&CommentWatch, "watch-list", "w", "", "A comma-separated list of email addresses to add to the watch list for this record. This is for all customer facing communications.")
	rootCmd.PersistentFlags().StringVarP(&NoteWatch, "note-list", "n", "", "A comma-separated list of email addresses to add to the work note watch list for this record. This is for all internal (i.e., non-customer) facing communications and notes.")
	rootCmd.PersistentFlags().StringVarP(&state, "state", "s", "open", "The state of the record. Valid values are (o)pen or (r)esolved.")

	rootCmd.Flags().StringVarP(&RecordNumber, "record", "r", "", "Service Now record number (required). Only REQs, CHGs, and INCs supported.")
	rootCmd.MarkFlagRequired("record")
}

func GetCredentials() {
	if KeyFile != "" {
		_, err := os.Stat(KeyFile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading config file: %s\n", err)
			os.Exit(1)
		} else {
			viper.SetConfigFile(KeyFile)
		}
	} else {
		viper.AddConfigPath(".")
		home, _ := os.UserHomeDir()
		viper.AddConfigPath(home)
		viper.SetConfigName(".snapi")
	}
	viper.SetConfigType("dotenv")

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
	WorkNotes := fmt.Sprintf("Ran on <code>%s</code> by <code>%s</code> in <code>%s</code>.<br /><br /><b>Command: <code>%s</code></b><br /><br /><code>--- STDOUT ---</code><br /><br /><pre>%s</pre><br /><br /><code>--- STDERR ---</code><br /><br /><pre>%s</pre><br />", hostname, user.Username, cwd, strings.Join(args, " "), capturedOutput, capturedError)

	return WorkNotes
}
