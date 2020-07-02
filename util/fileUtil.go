package util

import (
	"os"
	"github.com/golang/glog"
	"fmt"
	"bufio"
	"io"
	"strings"
	"github.com/descheduler-controller/pkg/api"
	"k8s.io/api/core/v1"
	"strconv"
)


func Modifyfile(fileName string,thresholds api.ResourceThresholds,targetThresholds api.ResourceThresholds) {
	in, err := os.Open(fileName)
	if err != nil {
		fmt.Println("open file fail:", err)
		os.Exit(-1)
	}
	defer in.Close()
	policyPath,templatePath := GetDataFromYaml()
	glog.Warningf("path : %s ", templatePath)
	out, err := os.OpenFile(policyPath, os.O_RDWR|os.O_CREATE, 0766)
	if err != nil {
		fmt.Println("Open write file fail:", err)
		os.Exit(-1)
	}
	defer out.Close()

	br := bufio.NewReader(in)
	index := 1
	for {
		line, _, err := br.ReadLine()
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Println("read err:", err)
			os.Exit(-1)
		}
		if -1 != strings.Index(string(line),"thresholdscpu") {
			thresholdscpu := strings.Replace(string(line), "thresholdscpu", strconv.FormatFloat(float64(thresholds[v1.ResourceCPU]),'f',6,64), -1)
			_, err = out.WriteString(thresholdscpu + "\n")
			if err != nil {
				fmt.Println("write to file fail:", err)
				os.Exit(-1)
			}
		} else if -1 != strings.Index(string(line),"thresholdsmemory") {
			thresholdscpu := strings.Replace(string(line), "thresholdsmemory", strconv.FormatFloat(float64(thresholds[v1.ResourceMemory]),'f',6,64), -1)
			_, err = out.WriteString(thresholdscpu + "\n")
			if err != nil {
				fmt.Println("write to file fail:", err)
				os.Exit(-1)
			}
		} else if -1 != strings.Index(string(line),"thresholdspods") {
			thresholdscpu := strings.Replace(string(line), "thresholdspods", strconv.FormatFloat(float64(thresholds[v1.ResourcePods]),'f',6,64), -1)
			_, err = out.WriteString(thresholdscpu + "\n")
			if err != nil {
				fmt.Println("write to file fail:", err)
				os.Exit(-1)
			}
		} else if -1 != strings.Index(string(line),"targetThresholdscpu") {
			thresholdscpu := strings.Replace(string(line), "targetThresholdscpu", strconv.FormatFloat(float64(targetThresholds[v1.ResourceCPU]),'f',6,64), -1)
			_, err = out.WriteString(thresholdscpu + "\n")
			if err != nil {
				fmt.Println("write to file fail:", err)
				os.Exit(-1)
			}
		} else if -1 != strings.Index(string(line),"targetThresholdsmemory") {
			thresholdscpu := strings.Replace(string(line), "targetThresholdsmemory",  strconv.FormatFloat(float64(targetThresholds[v1.ResourceMemory]),'f',6,64), -1)
			_, err = out.WriteString(thresholdscpu + "\n")
			if err != nil {
				fmt.Println("write to file fail:", err)
				os.Exit(-1)
			}
		} else if -1 != strings.Index(string(line),"targetThresholdspods") {
			thresholdscpu := strings.Replace(string(line), "targetThresholdspods",  strconv.FormatFloat(float64(targetThresholds[v1.ResourcePods]),'f',6,64), -1)
			_, err = out.WriteString(thresholdscpu + "\n")
			if err != nil {
				fmt.Println("write to file fail:", err)
				os.Exit(-1)
			}
		}else {
			_, err = out.WriteString(string(line) + "\n")
			if err != nil {
				fmt.Println("write to file fail:", err)
				os.Exit(-1)
			}
		}

//		fmt.Println("done ", index)
		index++
	}
	fmt.Println("FINISH!")

}
