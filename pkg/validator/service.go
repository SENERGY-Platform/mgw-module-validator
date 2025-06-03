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
	"fmt"
	module_lib "github.com/SENERGY-Platform/mgw-module-lib/model"
	"github.com/SENERGY-Platform/mgw-module-lib/util/sem_ver"
	module_lib_validation "github.com/SENERGY-Platform/mgw-module-lib/validation"
	"github.com/SENERGY-Platform/mgw-module-validator/pkg/models"
	"os"
	"path"
	"slices"
	"strings"
)

func ValidateMany(dirPath string, dependencies bool) ([]models.Report, error) {
	entries, err := os.ReadDir(dirPath)
	if err != nil {
		return nil, err
	}
	var reports []models.Report
	if dependencies {
		reportMap := make(map[string]models.Report)
		mods := make(map[string][2]string) // {modID:[modVer, dirName]}
		var deps [][2]string               // [depID, depVerRng]
		for _, entry := range entries {
			if entry.IsDir() {
				report, mod, err := validate(path.Join(dirPath, entry.Name()))
				if err != nil {
					fmt.Println(err)
					continue
				}
				reportMap[entry.Name()] = report
				if mod != nil {
					mods[mod.ID] = [2]string{mod.Version, entry.Name()}
					for depID, depVerRng := range mod.Dependencies {
						deps = append(deps, [2]string{depID, depVerRng})
					}
				}
			}
		}
		for _, dep := range deps {
			mod, ok := mods[dep[0]]
			if !ok {
				report := reportMap[mod[1]]
				report.Errs = append(report.Errs, fmt.Sprintf("missing dependency: %s", dep[0]))
				reportMap[mod[1]] = report
				continue
			}
			if ok, _ = sem_ver.InSemVerRange(dep[1], mod[0]); !ok {
				report := reportMap[mod[1]]
				report.Errs = append(report.Errs, fmt.Sprintf("dependency version not satisfied: %s available=%s required=%s", dep[0], mod[0], dep[1]))
				reportMap[mod[1]] = report
			}
		}
		for _, report := range reportMap {
			reports = append(reports, report)
		}
	} else {
		for _, entry := range entries {
			if entry.IsDir() {
				report, _, err := validate(path.Join(dirPath, entry.Name()))
				if err != nil {
					fmt.Println(err)
					continue
				}
				reports = append(reports, report)
			}
		}
	}
	slices.SortStableFunc(reports, func(a, b models.Report) int {
		return strings.Compare(a.ModID, b.ModID)
	})
	return reports, nil
}

func Validate(dirPath string) (models.Report, error) {
	report, _, err := validate(dirPath)
	if err != nil {
		return models.Report{}, err
	}
	return report, nil
}

func validate(dirPath string) (models.Report, *module_lib.Module, error) {
	mod, err := getModule(dirPath)
	if err != nil {
		return models.Report{}, nil, err
	}
	ri := models.Report{
		DirName: strings.TrimSuffix(path.Base(dirPath), path.Ext(dirPath)),
		ModID:   mod.ID,
		ModVer:  mod.Version,
		Status:  models.StatusPassed,
	}
	if err = module_lib_validation.Validate(mod); err != nil {
		ri.Errs = append(ri.Errs, err.Error())
		ri.Status = models.StatusFailed
	}
	return ri, mod, nil
}
