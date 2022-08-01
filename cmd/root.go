// Copyright 2022 Akiomi Kamakura
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"fmt"
	"os"

	"github.com/akiomik/get-old-tweets/config"
	"github.com/spf13/cobra"
)

var (
	text string
	userAgent string
)

var rootCmd = &cobra.Command{
	Use:     "get-old-tweets",
	Short:   "get-old-tweets " + config.Version,
	Version: config.Version,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Hello, " + text)
	},
}

func init() {
	rootCmd.Flags().StringVarP(&text, "text", "t", "", "query text to be matched")
	rootCmd.Flags().StringVarP(&userAgent, "user-agent", "", "", "user-agent for request")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println("Error:", err.Error())
		os.Exit(1)
	}
}
