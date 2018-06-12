// Copyright (C) 2003-2018 Opsview Limited. All rights reserved
// 
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//  http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	opspacks "github.com/webb249/ops/opspacks"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"time"
)

type PublicKey struct {
	Mem_Used                 float64 `json:"mem_used"`
	Fd_Used                  float64 `json:"fd_used"`
	Sockets_Used             float64 `json:"sockets_used"`
	Disk_Free                float64 `json:"disk_free"`
	Io_Read_Count            float64 `json:"Io_read_count"`
	Io_Read_Avg_Time         float64 `json:"io_read_avg_time"`
	Io_Write_Count           float64 `json:"io_write_count"`
	Io_Write_Avg_Time        float64 `json:"io_write_avg_time"`
	Io_Sync_Count            float64 `json:"io_sync_count"`
	Io_Sync_Avg_Time         float64 `json:"io_sync_avg_time"`
	Io_Seek_Count            float64 `json:"io_seek_count"`
	Io_Seek_Avg_Time         float64 `json:"io_seek_avg_time"`
	Context_Switches         float64 `json:"context_switches"`
	Context_Switches_Details struct {
		Rate float64 `json:"rate"`
	} `json:"context_switches_details"`
	Fd_Total        float64 `json:"fd_total"`
	Sockets_Total   float64 `json:"sockets_total"`
	Mem_Alarm       bool    `json:"mem_alarm"`
	Disk_Free_Alarm bool    `json:"disk_free_alarm"`
	Running         bool    `json:"running"`
}

func main() {
	check := opspacks.NewCheck("RabbitMQ")
	defer check.Finish()
	HelpMenu := flag.Bool("h", false, "Print help screen")
	URL := flag.String("H", "127.0.0.1", "Host address")
	Port := flag.String("P", "15672", "Connection port")
	Username := flag.String("u", "guest", "Username login")
	Password := flag.String("p", "guest", "Password login")
	Node := flag.String("n", "", "Name of the node being monitered")
	Warning := flag.Int("w", 0, "Warning level")
	Critical := flag.Int("c", 0, "Critical level")
	Mode := flag.String("m", "", "Run this check")

	flag.Parse()

	if *HelpMenu == true {
		helpText(*HelpMenu)
		os.Exit(int(opspacks.UNKNOWN))
	}
	flagCheck(*URL, *Port, *Mode, *Username, *Password, *Node)

	result := fetch(*URL, *Port, *Node, *Username, *Password)
	values := search(result)

	switch *Mode {
	case "mem_used":
		returnValue := removeDecimal(values.Mem_Used / 1024)
		findPerf(check, *Warning, *Critical, "Mem_Used", returnValue, "kb", true, false)
	case "fd_used":
		returnValue := values.Fd_Used
		findPerf(check, *Warning, *Critical, "Fd_Used", returnValue, "", false, true)
	case "sockets_used":
		returnValue := values.Sockets_Used
		findPerf(check, *Warning, *Critical, "Sockets_Used", returnValue, "", false, true)
	case "fd_left":
		returnValue := values.Fd_Total - values.Fd_Used
		findPerf(check, *Warning, *Critical, "Fd_Left", returnValue, "", true, true)
	case "fd_used_percent":
		returnValue := removeDecimal(percentUsed(values.Fd_Total, values.Fd_Used))
		findPerf(check, *Warning, *Critical, "Fd_Used_Percent", returnValue, "%", false, true)
	case "disk_free":
		returnValue := values.Disk_Free
		returnValue = removeDecimal(returnValue / 1024)
		findPerf(check, *Warning, *Critical, "Disk_Free", returnValue, "kb", true, false)
	case "disk_free_alarm":
		returnValue := values.Disk_Free_Alarm
		findResult(check, returnValue, "Disk_Free_Alarm", false)
	case "running":
		returnValue := values.Running
		findResult(check, returnValue, "Running", true)
	case "sockets_left":
		returnValue := values.Sockets_Total - values.Sockets_Used
		findPerf(check, *Warning, *Critical, "Sockets_Left", returnValue, "", true, true)
	case "sockets_used_percent":
		returnValue := removeDecimal(percentUsed(values.Sockets_Total, values.Sockets_Used))
		findPerf(check, *Warning, *Critical, "Sockets_Used_Percent", returnValue, "%", false, true)
	case "mem_alarm":
		returnValue := values.Mem_Alarm
		findResult(check, returnValue, "Memory_Alarm", false)
	case "io_read_count":
		returnValue := values.Io_Read_Count
		findPerf(check, *Warning, *Critical, "Io_Read_Count", returnValue, "", false, false)
	case "io_write_count":
		returnValue := values.Io_Write_Count
		findPerf(check, *Warning, *Critical, "Io_Write_Count", returnValue, "", false, false)
	case "io_write_avg_time":
		returnValue := values.Io_Write_Avg_Time
		findPerf(check, *Warning, *Critical, "Io_Write_Avg_Time", returnValue, "ms", false, true)
	case "io_read_avg_time":
		returnValue := values.Io_Read_Avg_Time
		findPerf(check, *Warning, *Critical, "Io_Read_Avg_Time", returnValue, "ms", false, true)
	case "io_seek_avg_time":
		returnValue := values.Io_Seek_Avg_Time
		findPerf(check, *Warning, *Critical, "Io_Seek_Avg_Time", returnValue, "ms", false, false)
	case "io_sync_avg_time":
		returnValue := values.Io_Sync_Avg_Time
		findPerf(check, *Warning, *Critical, "Io_Sync_Avg_Time", returnValue, "ms", false, true)
	case "context_switch":
		returnValue := values.Context_Switches_Details.Rate
		findPerf(check, *Warning, *Critical, "Context_Switch", returnValue, "", false, true)
	default:
		opspacks.Exit(opspacks.UNKNOWN, "No check specified")
	}
}

func helpText(HelpMenu bool) {
	fmt.Printf(`Check for RabbitMQ Node status

  Usage: "check_rabbitmq_node -m <MODE> -H <URL> -P <PORT> -u <USERNAME> -p <PASSWORD> -n <NODE> -c <CRITICAL> -w <WARNING>"

  -h, --help
  Print this help screen
  -m Mode, -mode=MODE
  The collecter check to run, options are:
    mem_used: Total memory used
    mem_alarm: Returns critical if the memory alarm has gone off
    fd_left: Number of file descriptors remaining
    fd_used_percent: Percentage of file descriptors used
    disk_free: Disk free space in bytes
    disk_free_alarm: Returns critical if the disk free alarm has gone off
    running: Whether or not this node is up
    sockets_left: File descriptors available for use as sockets remaining
    sockets_used_percent: Percentage of file descriptors used as sockets
    io_read_count: Rate of read operations by the persister in the last statistics interval
    io_write_count: Rate of write operations by the persister in the last statistics interval
    io_write_avg_time: Average wait time (ms) for each disk write operation in the last statistics interval
    io_read_avg_time: Average wait time (ms) for each disk read operation in the last statistics interval
    io_seek_avg_time: Average wait time (ms) for each seek operation in the last statistics interval
    io_sync_avg_time: Average wait time (ms) for each fsync() operation in the last statistics interval
    context_switch: Rate at which context switching takes place on this node during last statistics interval
  -H URL, --hostname=URL. URL prefix for API access. Required.
  -P PORT, --port=PORT. Port for API access. If not specified default port 15672 will be used.
  -n NODE, --node=NODE. Name of the node being monitored.
  -w WARNING, --warning=WARNING. Warning level.
  -c CRITICAL, --critical=CRITICAL. Critical level.`)
}

func fetch(URL string, Port string, Node string, Username string, Password string) []byte {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr, Timeout: time.Duration(10 * time.Second)}

	req, err := http.NewRequest("GET", "http://"+URL+":"+Port+"/api/nodes/"+Node, nil)
	exitOnError(err, "Cannot create connection request")
	req.SetBasicAuth(Username, Password)
	resp, err := client.Do(req)
	exitOnError(err, "Could Not Connect to "+URL+":"+Port)

	if resp.StatusCode != 200 {
		opspacks.Exit(opspacks.CRITICAL, "Could Not Connect to "+URL+" Response Code: "+string(resp.Status))
	}

	bodyText, err := ioutil.ReadAll(resp.Body)
	exitOnError(err, "Error reading output of"+Node+": "+string(bodyText))

	return bodyText
}

func search(result []byte) PublicKey {
	var pluginResponse PublicKey

	err := json.Unmarshal(result, &pluginResponse)
	exitOnError(err, "Failed to parse JSON response")

	return pluginResponse
}

func findPerf(check *opspacks.Check, warning int, critical int, returnName string, returnValue float64, UOM string, reverse bool, compare bool) {
	returnMessage := ""
	returnStatus := opspacks.OK
	if compare == true {
		if reverse == false {
			if int(returnValue) >= critical {
				returnStatus = opspacks.CRITICAL
			} else if int(returnValue) >= warning {
				returnStatus = opspacks.WARNING
			}
		} else {
			if int(returnValue) <= critical {
				returnStatus = opspacks.CRITICAL
			} else if int(returnValue) <= warning {
				returnStatus = opspacks.WARNING
			}
		}
	}
	returnMessage = returnName + ": " + strconv.Itoa(int(returnValue))

	check.AddPerfData(returnName, UOM, returnValue)
	check.AddResult(returnStatus, returnMessage)
}

func findResult(check *opspacks.Check, returnValue bool, returnMessage string, reverse bool) {
	returnStatus := opspacks.OK

	if reverse == false {
		if returnValue == true {
			returnMessage = returnMessage + ": True"
			returnStatus = opspacks.CRITICAL

		} else {
			returnMessage = returnMessage + ": False"
		}
	} else {
		if returnValue == false {
			returnMessage = returnMessage + ": False"
			returnStatus = opspacks.CRITICAL

		} else {
			returnMessage = returnMessage + ": True"
		}
	}

	check.AddResult(returnStatus, returnMessage)
}

func percentUsed(total float64, used float64) float64 {
	return (used / total) * 100
}

func removeDecimal(value float64) float64 {
	intValue := int(value)
	return float64(intValue)
}

func flagCheck(URL string, Port string, Mode string, Username string, Password string, Node string) {
	if URL == "" {
		err := errors.New("Flag check error")
		exitOnError(err, "Hostname (-H) is a required argument")
	}
	if Port == "" {
		err := errors.New("Flag check error")
		exitOnError(err, "Port (-P) is a required argument")
	}
	if Mode == "" {
		err := errors.New("Flag check error")
		exitOnError(err, "Check (-m) is a required argument")
	}
	if Username == "" {
		err := errors.New("Flag check error")
		exitOnError(err, "Username (-u) is a required argument")
	}
	if Password == "" {
		err := errors.New("Flag check error")
		exitOnError(err, "Password (-p) is a required argument")
	}
	if Node == "" {
		err := errors.New("Flag check error")
		exitOnError(err, "Node (-n) is a required argument")
	}
}

func exitOnError(err error, message string) {
	if err != nil {
		opspacks.Exit(opspacks.CRITICAL, message)
	}
}
