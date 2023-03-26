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

package ksqldb_test

import (
	"testing"

	"github.com/oscarzhou/ksqldb-go"
	"github.com/stretchr/testify/require"
)

func TestResponseError_ErrorMessage(t *testing.T) {
	error := ksqldb.ResponseError{
		ErrType: "some error",
		ErrCode: 40001,
		Message: "Some error occured\nNext line",
	}
	require.Equal(t, "Some error occured Next line", error.Error())
}
