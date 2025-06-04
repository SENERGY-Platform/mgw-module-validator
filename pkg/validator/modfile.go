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

package validator

import (
	"errors"
	"os"
	"path"
	"regexp"
)

var re = regexp.MustCompile(`^Modfile\.(?:yml|yaml)$`)
var NoModfileErr = errors.New("missing modfile")

type OpenModfileErr struct {
	err error
}

func (e *OpenModfileErr) Error() string {
	return e.err.Error()
}

func (e *OpenModfileErr) Unwrap() error {
	return e.err
}

func openModFile(dirPath string) (*os.File, error) {
	dirEntries, err := os.ReadDir(dirPath)
	if err != nil {
		return nil, &OpenModfileErr{err: err}
	}
	for _, entry := range dirEntries {
		if !entry.IsDir() {
			eName := entry.Name()
			if re.MatchString(eName) {
				f, err := os.Open(path.Join(dirPath, eName))
				if err != nil {
					return nil, &OpenModfileErr{err: err}
				}
				return f, nil
			}
		}
	}
	return nil, &OpenModfileErr{err: NoModfileErr}
}
