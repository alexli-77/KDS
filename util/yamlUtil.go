package util

import (
//	"strings"
//	"errors"
	"io/ioutil"
	"gopkg.in/yaml.v2"
//	"reflect"
	"fmt"
//	"strconv"
	"github.com/descheduler-controller/model"
//	"os"
//	"encoding/json"
//	"reflect"
)

func GetDataFromYaml() (string,string){
	var policyyaml model.Config
	config, err := ioutil.ReadFile("./config/config.yaml")
	if err != nil {
		fmt.Print(err)
	}
	yaml.Unmarshal(config,&policyyaml)

	policyfilepath := policyyaml.PolicyFilepath
	templatefilepath := policyyaml.TemplateFilepath
	//fmt.Println(policyyaml.Strategies.LowNodeUtilization.Params.NodeResourceUtilizationThresholds.Thresholds.Cpu)
	return policyfilepath,templatefilepath
}
