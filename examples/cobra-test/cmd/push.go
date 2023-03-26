/*
Copyright © 2021 Thomas Meitz

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"context"
	"time"

	"github.com/Masterminds/log-go"
	"github.com/oscarzhou/ksqldb-go"
	"github.com/oscarzhou/ksqldb-go/net"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// pushCmd represents the push command
var pushCmd = &cobra.Command{
	Use:   "push",
	Short: "push dogs example",
}

func init() {
	pushCmd.Run = push
	rootCmd.AddCommand(pushCmd)
}

func push(cmd *cobra.Command, args []string) {
	setLogger()
	host := viper.GetString("host")
	user := viper.GetString("username")
	password := viper.GetString("password")

	options := net.Options{
		Credentials: net.Credentials{Username: user, Password: password},
		BaseUrl:     host,
		AllowHTTP:   true,
	}

	kcl, err := ksqldb.NewClientWithOptions(options)
	if err != nil {
		log.Fatal(err)
	}
	defer kcl.Close()

	// You don't need to parse your ksql statement; Client.Pull parses it for you
	// if parsing is enabled (default)
	// you can disable parsing with `kcl.EnableParseSQL(false)`
	query := "select rowtime, id, name, dogsize, age from dogs emit changes;"

	rowChannel := make(chan ksqldb.Row)
	headerChannel := make(chan ksqldb.Header, 1)

	// This Go routine will handle rows as and when they
	// are sent to the channel
	go func() {
		var dataTs float64
		var id string
		var name string
		var dogSize string
		var age string
		for row := range rowChannel {
			if row != nil {
				// Should do some type assertions here
				dataTs = row[0].(float64)
				id = row[1].(string)
				name = row[2].(string)
				dogSize = row[3].(string)
				age = row[4].(string)

				// Handle the timestamp
				t := int64(dataTs)
				ts := time.Unix(t/1000, 0).Format(time.RFC822)

				log.Infof("🐾 New dog at %v: '%v' is %v and %v (id %v)\n", ts, name, dogSize, age, id)
			}
		}

	}()

	ctx, cancel := context.WithTimeout(context.TODO(), 120*time.Second)
	defer cancel()

	e := kcl.Push(ctx, ksqldb.QueryOptions{Sql: query}, rowChannel, headerChannel)

	if e != nil {
		log.Fatal(e)
	}
}
