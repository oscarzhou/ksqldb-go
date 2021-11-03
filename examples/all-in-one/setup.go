/*
Copyright © 2021 Robin Moffat & Contributors

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

package main

import (
	"fmt"
	"time"

	"github.com/Masterminds/log-go"
	"github.com/thmeitz/ksqldb-go"
)

func setup() (*ksqldb.Client, error) {

	options := ksqldb.Options{
		Credentials: ksqldb.Credentials{Username: ksqlDBUser, Password: ksqlDBPW},
		BaseUrl:     ksqlDBServer,
	}

	client, err := ksqldb.NewClient(options, log.Current)
	if err != nil {
		log.Fatal(err)
	}
	// create the dummy data connector
	// IF NOT EXISTS didnt works - it throws an error - this is a bug in the ksql-rest-api
	if err := client.Execute(`
		CREATE SOURCE CONNECTOR DOGS WITH (
		'connector.class'               = 'io.mdrogalis.voluble.VolubleSourceConnector',
		'key.converter'                 = 'org.apache.kafka.connect.storage.StringConverter',
		'value.converter'               = 'org.apache.kafka.connect.json.JsonConverter',
		'value.converter.schemas.enable'='false',
		'genkp.dogs.with'               = '#{Internet.uuid}',
		'genv.dogs.name.with'           = '#{Dog.name}',
		'genv.dogs.dogsize.with'        = '#{Dog.size}',
		'genv.dogs.age.with'            = '#{Dog.age}',
		'topic.dogs.throttle.ms'        = 1000 
		);
		`); err != nil {
		return nil, fmt.Errorf("error creating the source connector.\n%v", err)
	}

	// this is a bit lame but without doing the cool stuff with CommandId etc
	// it's the easiest way to make sure the topic exists before continuing
	time.Sleep(5 * time.Second)

	// create the DOGS stream
	if err := client.Execute(`
	CREATE STREAM IF NOT EXISTS DOGS (ID STRING KEY, 
						NAME STRING, 
						DOGSIZE STRING, 
						AGE STRING) 
				  WITH (KAFKA_TOPIC='dogs', 
				  VALUE_FORMAT='JSON', PARTITIONS=1);
	`); err != nil {
		return nil, fmt.Errorf("error creating the dogs stream.\n%v", err)
	}

	// this is a bit lame but without doing the cool stuff with CommandId etc
	// it's the easiest way to make sure the stream exists before continuing
	time.Sleep(5 * time.Second)

	// create the DOGS_BY_SIZE table
	if err := client.Execute(`
	CREATE TABLE IF NOT EXISTS DOGS_BY_SIZE AS 
		SELECT DOGSIZE AS DOG_SIZE, COUNT(*) AS DOGS_CT 
		FROM DOGS WINDOW TUMBLING (SIZE 15 MINUTE) 
		GROUP BY DOGSIZE;
	`); err != nil {
		return nil, fmt.Errorf("error creating the dogs stream.\n%v", err)
	}
	// This is a bit lame but without doing the cool stuff with CommandId etc
	// it's the easiest way to make sure the table exists before continuing
	time.Sleep(10 * time.Second)

	return client, nil
}
