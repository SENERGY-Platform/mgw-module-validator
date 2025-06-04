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

type modWrapper struct {
	Module  *module_lib.Module
	DirName string
}

func ValidateMany(dirPath string, dependencies bool) ([]models.Report, error) {
	entries, err := os.ReadDir(dirPath)
	if err != nil {
		return nil, err
	}
	var reports []models.Report
	if dependencies {
		reportMap := make(map[string]models.Report)
		mods := make(map[string]modWrapper)
		for _, entry := range entries {
			if entry.IsDir() {
				report, mod, err := validate(path.Join(dirPath, entry.Name()))
				if err != nil {
					fmt.Println(err)
					continue
				}
				reportMap[entry.Name()] = report
				if mod != nil {
					mods[mod.ID] = modWrapper{
						Module:  mod,
						DirName: entry.Name(),
					}
				}
			}
		}
		for _, wrapper := range mods {
			for depID, depVerRng := range wrapper.Module.Dependencies {
				var errMsgs []string
				mod, ok := mods[depID]
				if !ok {
					errMsgs = append(errMsgs, fmt.Sprintf("missing dependency: %s", depID))
				} else {
					if ok, _ = sem_ver.InSemVerRange(depVerRng, mod.Module.Version); !ok {
						errMsgs = append(errMsgs, fmt.Sprintf("dependency version not satisfied: %s available=%s required=%s", depID, mod.Module.Version, depVerRng))
					}
				}
				if len(errMsgs) > 0 {
					report := reportMap[wrapper.DirName]
					report.Errs = append(report.Errs, errMsgs...)
					report.Status = models.StatusFailed
					reportMap[wrapper.DirName] = report
				}
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
