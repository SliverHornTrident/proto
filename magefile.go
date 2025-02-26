//go:build mage

package main

import (
	"github.com/magefile/mage/sh"
	"github.com/pkg/errors"
	"os"
	"path/filepath"
	"strings"
)

var pwd string

func init() {
	pwd, _ = os.Getwd()
}

func Proto() error {
	files := make([]string, 0)
	err := filepath.Walk(pwd, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return errors.Wrapf(err, "[filepath:%s]读取失败!", path)
		}
		if info.IsDir() {
			return nil
		}
		if filepath.Ext(info.Name()) == ".proto" {
			relName, _ := filepath.Rel(pwd, path)
			if strings.Contains(relName, "third_party") {
				return nil
			}
			relName = strings.ReplaceAll(relName, `\`, "/")
			files = append(files, relName)
		}
		return nil
	})
	if err != nil {
		return err
	}
	args := []string{
		"--proto_path=.",
		"--proto_path=./third_party",
		"--go_out=paths=source_relative:.",
		"--go-http_out=paths=source_relative:.",
		"--go-grpc_out=paths=source_relative:.",
		"--go-errors_out=paths=source_relative:.",
		"--validate_out=paths=source_relative,lang=go:.",
		"--openapi_out=fq_schema_naming=true,default_response=false:.",
	}
	args = append(args, files...)
	err = sh.Run("protoc", args...)
	if err != nil {
		return errors.Wrap(err, "proto 生成失败!")
	}
	return nil
}

func Clear() error {
	err := filepath.Walk(pwd, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return errors.Wrapf(err, "[filepath:%s]读取失败!", path)
		}
		if info.IsDir() {
			return nil
		}
		switch {
		case strings.Contains(info.Name(), ".pb.go"), strings.Contains(info.Name(), ".pb.validate.go"):
			err = os.Remove(path)
			if err != nil {
				return errors.Wrapf(err, "[filepath:%s]删除失败!", path)
			}
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}
