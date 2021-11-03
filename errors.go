/*
Copyright © 2021 Robin Moffat & Contributors
Copyright © 2021 Thomas Meitz <thme219@gmail.com>

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

package ksqldb

import (
	"errors"
	"fmt"
)

var (
	ErrNotFound = errors.New("no result found")
)

type ClientError struct {
	Message string
}

func (e *ClientError) Error() string {
	return fmt.Sprintf("ISPConfig Client error: %v", e.Message)
}

func NewClientError(message string) *ClientError {
	return &ClientError{Message: message}
}
