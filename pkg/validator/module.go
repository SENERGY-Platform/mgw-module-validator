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
	"github.com/SENERGY-Platform/mgw-modfile-lib/modfile"
	"github.com/SENERGY-Platform/mgw-modfile-lib/v1/v1dec"
	"github.com/SENERGY-Platform/mgw-modfile-lib/v1/v1gen"
	module_lib "github.com/SENERGY-Platform/mgw-module-lib/model"
	"gopkg.in/yaml.v3"
)

var mfDecoders = make(modfile.Decoders)
var mfGenerators = make(modfile.Generators)

func init() {
	mfDecoders.Add(v1dec.GetDecoder)
	mfGenerators.Add(v1gen.GetGenerator)
}

func getModule(dirPath string) (*module_lib.Module, error) {
	file, err := openModFile(dirPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	yd := yaml.NewDecoder(file)
	mf := modfile.New(mfDecoders, mfGenerators)
	err = yd.Decode(&mf)
	if err != nil {
		return nil, err
	}
	mod, err := mf.GetModule()
	if err != nil {
		return nil, err
	}
	return mod, nil
}
