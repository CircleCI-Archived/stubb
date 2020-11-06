package cmd

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	log "github.com/sirupsen/logrus"
)

var setDescriptionCmd = &cobra.Command{
	Use:   "description <image-name> <file>",
	Short: "Set the description for an image on Docker Hub",
	Run: func(cmd *cobra.Command, args []string) {

		if len(args) < 2 {
			log.Fatal("Please provide a text file.")
			os.Exit(1)
		}

		content, err := ioutil.ReadFile(args[1])
		if err != nil {
			log.Fatal(err)
		}

		// Escape file content for use in JSON
		content = []byte(strconv.Quote(string(content)))

		content = append([]byte("{\"full_description\": "), content[:len(content)-1]...)
		content = append(content, []byte("\"}")...)

		req, err := http.NewRequest("PATCH", "https://hub.docker.com/v2/repositories/"+args[0]+"/", bytes.NewBuffer(content))
		if err != nil {
			log.Fatal(err)
		}

		req.Header.Set("Content-Type", "application/json")

		if viper.Get("user") == nil || viper.Get("pass") == nil || len(viper.Get("user").(string)) == 0 || len(viper.Get("pass").(string)) == 0 {
			log.Fatal("This command requires Docker Hub credentials (DOCKER_USER and DOCKER_PASS), credentials to be set in your environment.")
		}

		resp, err := sendRequest(req, viper.Get("user").(string), viper.Get("pass").(string))
		if err != nil {
			log.Fatal(err)
		}
		defer resp.Body.Close()

		if resp.StatusCode >= 300 {
			log.Fatal("There was an error updating the description. Code " + resp.Status)
		} else {
			fmt.Printf("Successfully updated with code %d.\n", resp.StatusCode)
		}
	},
}

func init() {
	setCmd.AddCommand(setDescriptionCmd)
}
