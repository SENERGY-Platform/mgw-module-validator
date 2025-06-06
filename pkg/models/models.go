/*
 * Copyright 2025 InfAI (CC SES)
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *    http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package models

import (
	"fmt"
)

const (
	StatusFailed = "failed"
	StatusPassed = "passed"
)

type Report struct {
	DirName string   `json:"dir_name"`
	ModID   string   `json:"module_id"`
	ModVer  string   `json:"module_version"`
	Errs    []string `json:"errors"`
	Status  string   `json:"status"`
}

func (r *Report) String() string {
	var str string
	str += fmt.Sprintf("%s:\n", r.DirName)
	str += fmt.Sprintf("\tmodule_id: %s\n", r.ModID)
	str += fmt.Sprintf("\tmodule_version: %s\n", r.ModVer)
	str += fmt.Sprintf("\terrors:\n")
	for _, err := range r.Errs {
		str += fmt.Sprintf("\t\t%s\n", err)
	}
	str += fmt.Sprintf("\tstatus: %s", r.Status)
	return str
}
