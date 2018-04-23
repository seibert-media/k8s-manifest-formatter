package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"runtime"
	"github.com/golang/glog"
	"github.com/ghodss/yaml"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/client-go/kubernetes/scheme"
	"bytes"
	k8s_runtime "k8s.io/apimachinery/pkg/runtime"
	"io"
	io_util "github.com/bborbe/io/util"
)

const (
	parameterPath     = "path"
	parameterWrite    = "write"
	parameterValidate = "validate"
)

var (
	pathPtr     = flag.String(parameterPath, "", "path")
	writePtr    = flag.Bool(parameterWrite, false, "write formates content back to file")
	validatePtr = flag.Bool(parameterValidate, false, "validate content is already formated")
)

func main() {
	defer glog.Flush()
	glog.CopyStandardLogTo("info")
	flag.Parse()
	runtime.GOMAXPROCS(runtime.NumCPU())

	if len(*pathPtr) == 0 {
		fmt.Fprintf(os.Stderr, "parameter %s missing\n", parameterPath)
		os.Exit(1)
	}
	path, err := io_util.NormalizePath(*pathPtr)
	if err != nil {
		glog.Exitf("normalize path: %s failed: %v", path, err)
	}
	file, err := os.Open(path)
	if err != nil {
		glog.Exitf("open file %s failed: %v", path, err)
	}
	defer file.Close()
	fileInfo, err := file.Stat()
	if err != nil {
		glog.Exitf("get file stat failed: %v", path, err)
	}
	content, err := ioutil.ReadAll(file)
	if err != nil {
		glog.Exitf("read %s failed: %v", path, err)
	}
	formatedContent, err := formatYaml(content)
	if err != nil {
		glog.Exitf("format yaml %s failed: %v", path, err)
	}
	if *writePtr {
		if err := ioutil.WriteFile(path, formatedContent, fileInfo.Mode()); err != nil {
			glog.Exitf("write file failed: %v", err)
		}
		glog.V(0).Infof("write file completed")
	} else if *validatePtr {
		if bytes.Compare(content, formatedContent) != 0 {
			fmt.Fprintf(os.Stderr, "content of %s is not formatted\n", path)
			os.Exit(1)
		}
		glog.V(0).Infof("content is formatted")
	} else {
		if _, err := io.Copy(os.Stdout, bytes.NewBuffer(formatedContent)); err != nil {
			glog.Exitf("print content failed")
		}
	}
}

func formatYaml(content []byte) ([]byte, error) {
	content, err := yaml.YAMLToJSON(content)
	if err != nil {
		return nil, fmt.Errorf("yaml to json failed: %v", err)
	}
	obj, err := kind(content)
	if err != nil {
		return nil, fmt.Errorf("create object by content failed: %v", err)
	}
	if obj, _, err = unstructured.UnstructuredJSONScheme.Decode(content, nil, obj); err != nil {
		return nil, fmt.Errorf("unmarshal to object failed: %v", err)
	}
	return yaml.Marshal(obj)
}

func kind(content []byte) (k8s_runtime.Object, error) {
	_, kind, err := unstructured.UnstructuredJSONScheme.Decode(content, nil, nil)
	if err != nil {
		return nil, fmt.Errorf("unmarshal to unknown failed: %v", err)
	}
	if kind.Kind == "Secret" {
		return nil, nil
	}
	obj, err := scheme.Scheme.New(*kind)
	if err != nil {
		return nil, fmt.Errorf("create object failed: %v", err)
	}
	return obj, nil
}
