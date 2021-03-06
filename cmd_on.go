// Copyright 2019 Aporeto Inc.
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//     http://www.apache.org/licenses/LICENSE-2.0
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.aporeto.io/remod/internal/remod"
)

var cmdOn = &cobra.Command{
	Use:   "on",
	Short: "Activate remod on the current branch",
	Long: `This command will turn on the development mode for the current branch.

It will create a backup of your current 'go.mod' file (stripped from any
eventual development replacement remains), and will apply the replacements
from your 'remod.dev' file.

It will also align the git filters from '.git/config' if they are missing.

Note that everytime you switch to a new branch, or pull a change, you must
run 'remod on' to realign remod with the latest changes.
`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return viper.BindPFlags(cmd.Flags())
	},
	RunE: func(cmd *cobra.Command, args []string) error {

		if err := remod.Off(); err != nil {
			return err
		}

		if err := remod.GitConfig(); err != nil {
			return err
		}

		if err := remod.On(); err != nil {
			return err
		}

		return remod.GitAdd()
	},
}
