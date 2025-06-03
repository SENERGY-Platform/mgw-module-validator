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

package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/SENERGY-Platform/mgw-module-validator/pkg/validator"
	"os"
)

const (
	textOutputFormat = "text"
	jsonOutputFormat = "json"
)

func main() {
	var targetPath string
	var outputFormat string
	var outputPath string
	var multiple bool
	var dependencies bool

	flag.StringVar(&targetPath, "target", "", "target path")
	flag.StringVar(&targetPath, "t", "", "target path")
	flag.StringVar(&outputFormat, "format", textOutputFormat, "output format")
	flag.StringVar(&outputFormat, "f", textOutputFormat, "output format")
	flag.StringVar(&outputPath, "output", "", "output file path")
	flag.StringVar(&outputPath, "o", "", "output file path")
	flag.BoolVar(&multiple, "multiple", false, "validate multiple modules")
	flag.BoolVar(&multiple, "m", false, "validate multiple modules")
	flag.BoolVar(&dependencies, "dependencies", false, "check dependencies")
	flag.BoolVar(&dependencies, "d", false, "check dependencies")
	flag.Parse()

	if targetPath == "" {
		targetPath = flag.Arg(0)
	}
	if targetPath == "" {
		fmt.Println("no target path specified")
		os.Exit(1)
	}

	if multiple {
		reports, err := validator.ValidateMany(targetPath, dependencies)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		switch outputFormat {
		case textOutputFormat:
			var str string
			lenReports := len(reports)
			for i, report := range reports {
				str += report.String() + "\n"
				if i < lenReports-1 {
					str += "\n"
				}
			}
			if err := writeOutputString(outputPath, str); err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		case jsonOutputFormat:
			b, err := json.MarshalIndent(reports, "", "  ")
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			if err = writeOutputBytes(outputPath, b); err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		default:
			fmt.Println("invalid output format")
		}
		for _, report := range reports {
			if len(report.Errs) > 0 {
				os.Exit(1)
			}
		}
	} else {
		report := validator.Validate(targetPath)
		switch outputFormat {
		case textOutputFormat:
			if err := writeOutputString(outputPath, report.String()+"\n"); err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		case jsonOutputFormat:
			b, err := json.MarshalIndent(report, "", "  ")
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			if err = writeOutputBytes(outputPath, b); err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		default:
			fmt.Println("invalid output format")
		}
		if len(report.Errs) > 0 {
			os.Exit(1)
		}
	}
}

func writeOutputBytes(filePath string, bytes []byte) error {
	if filePath != "" {
		file, err := os.Create(filePath)
		if err != nil {
			return err
		}
		defer file.Close()
		_, err = file.Write(bytes)
	} else {
		_, err := os.Stdout.Write(bytes)
		if err != nil {
			return err
		}
	}
	return nil
}

func writeOutputString(filePath, str string) error {
	if filePath != "" {
		file, err := os.Create(filePath)
		if err != nil {
			return err
		}
		defer file.Close()
		_, err = file.WriteString(str)
	} else {
		_, err := os.Stdout.WriteString(str)
		if err != nil {
			return err
		}
	}
	return nil
}
