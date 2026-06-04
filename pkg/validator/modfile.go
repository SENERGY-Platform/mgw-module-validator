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
	"io/fs"
	"os"
	"regexp"

	modfile_lib "github.com/SENERGY-Platform/mgw-modfile-lib"
	module_lib "github.com/SENERGY-Platform/mgw-module-lib/model"
)

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

var regExp = regexp.MustCompile(`^Modfile\.(?:yml|yaml)$`)

func getModule(dirPath string) (module_lib.Module, error) {
	fSys := os.DirFS(dirPath)
	mfPath, err := findFile(fSys, regExp.MatchString)
	if err != nil {
		return module_lib.Module{}, &OpenModfileErr{err: err}
	}
	if mfPath == "" {
		return module_lib.Module{}, &OpenModfileErr{err: NoModfileErr}
	}
	file, err := fSys.Open(mfPath)
	if err != nil {
		return module_lib.Module{}, &OpenModfileErr{err: err}
	}
	defer file.Close()
	return modfile_lib.Decode(file)
}

func findFile(fSys fs.FS, match func(v string) bool) (string, error) {
	var filePath string
	err := fs.WalkDir(fSys, ".", func(currentPath string, dirEntry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !dirEntry.IsDir() && match(dirEntry.Name()) {
			filePath = currentPath
			return fs.SkipAll
		}
		return nil
	})
	if err != nil {
		return "", err
	}
	return filePath, nil
}
