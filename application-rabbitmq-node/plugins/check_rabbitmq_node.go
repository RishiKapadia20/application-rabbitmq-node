package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/opsview/go-plugin"
	"time"
	"strings"
)

type apiResponse struct {
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

var opts struct {
	Hostname string `short:"H" long:"hostname" description:"Host" default:"localhost"`
	Port     string `short:"P" long:"port" description:"Port" required:"true"`
	Mode     string `short:"m" long:"mode" description:"Mode" required:"true"`
	Warning  string `short:"w" long:"warning" description:"Warning"`
	Critical string `short:"c" long:"critical" description:"Critical"`
	Node     string `short:"n" long:"Node" description:"Node name"`
	Username string `short:"u" long:"Username" description:"opsview username" required:"true"`
	Password string `short:"p" long:"Password" description:"opsview password" required:"true"`
}

func main() {
	check := checkPlugin()
	if err := check.ParseArgs(&opts); err != nil {
		check.ExitUnknown("Error parsing arguments: %s", err)
	}
	defer check.Final()
	check.AllMetricsInOutput = true

	if opts.Node == "" {
		opts.Node = searchNodes(check, opts.Hostname, opts.Port, opts.Username, opts.Password)
	}

	result := fetch(check, opts.Hostname, opts.Port, opts.Node, opts.Username, opts.Password)
	values := search(check, result)

	switch opts.Mode {
	case "mem_used":
		memUsed, memUsedUOM := convertBytes(values.Mem_Used, "b",2)
		check.AddMetric("Mem_Used", memUsed, memUsedUOM, opts.Warning, opts.Critical)
	case "fd_used":
		check.AddMetric("Fd_Used", values.Fd_Used, "", opts.Warning, opts.Critical)
	case "sockets_used":
		check.AddMetric("Sockets_Used", values.Sockets_Used, "", opts.Warning, opts.Critical)
	case "fd_left":
		check.AddMetric("Fd_left", (values.Fd_Total - values.Fd_Used), "", opts.Warning, opts.Critical)
	case "fd_used_percent":
		check.AddMetric("Fd_Used_Percent", removeDecimal(percentUsed(values.Fd_Total, values.Fd_Used)), "%", opts.Warning, opts.Critical)
	case "disk_free":
		check.AddMetric("Disk_Free", removeDecimal(values.Disk_Free/1024), "kb", opts.Warning, opts.Critical)
	case "disk_free_alarm":
		findResult(check, values.Disk_Free_Alarm, "Disk_Free_Alarm", true)
	case "running":
		findResult(check, values.Running, "Running", false)
	case "sockets_left":
		check.AddMetric("Sockets_Left", (values.Sockets_Total - values.Sockets_Used), "", opts.Warning, opts.Critical)
	case "sockets_used_percent":
		check.AddMetric("Sockets_Used_Percent", removeDecimal(percentUsed(values.Sockets_Total, values.Sockets_Used)), "%", opts.Warning, opts.Critical)
	case "mem_alarm":
		findResult(check, values.Mem_Alarm, "Memory_Alarm", true)
	case "io_read_count":
		check.AddMetric("Io_Read_Count", values.Io_Read_Count, "", opts.Warning, opts.Critical)
	case "io_write_count":
		check.AddMetric("Io_Write_Count", values.Io_Write_Count, "", opts.Warning, opts.Critical)
	case "io_write_avg_time":
		check.AddMetric("Io_Write_Avg_Time", roundto2DP(values.Io_Write_Avg_Time), "ms", opts.Warning, opts.Critical)
	case "io_read_avg_time":
		check.AddMetric("Io_Read_Avg_Time", roundto2DP(values.Io_Read_Avg_Time), "ms", opts.Warning, opts.Critical)
	case "io_seek_avg_time":
		check.AddMetric("Io_Seek_Avg_Time", roundto2DP(values.Io_Seek_Avg_Time), "ms", opts.Warning, opts.Critical)
	case "io_sync_avg_time":
		check.AddMetric("Io_Sync_Avg_Time", roundto2DP(values.Io_Sync_Avg_Time), "ms", opts.Warning, opts.Critical)
	case "context_switch":
		check.AddMetric("Context_Switch", values.Context_Switches_Details.Rate, "", opts.Warning, opts.Critical)
	case "node_summary":
		//Running
		findResult(check, values.Running, "Running", false)
		//mem_used
		memUsed, memUsedUOM := convertBytes(values.Mem_Used, "b",2)
		check.AddMetric("Mem_Used", memUsed, memUsedUOM, opts.Warning, opts.Critical)
		//Mem Alarm
		alarmCheck(check, values.Mem_Alarm, "Memory_Alarm")
		//Disk Alarm
		alarmCheck(check, values.Disk_Free_Alarm, "Disk_Free_Alarm")
	default:
		check.ExitUnknown("No check specified")
	}
}

func checkPlugin() *plugin.Plugin {
	check := plugin.New("check_rabbitmq_node", "v2.0.0")

	check.Preamble = `Copyright (C) 2003-2018 Opsview Limited. All rights reserved.
This plugin tests the stats of an rabbitmq_node.`

	check.Description = `Check for RabbitMQ Node status

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
	node_summary: Provides a summary of metrics which include running, mem_used, mem_alarm and disk_free_alarm
  -H URL, --hostname=URL. URL prefix for API access. Required.
  -P PORT, --port=PORT. Port for API access. If not specified default port 15672 will be used.
  -n NODE, --node=NODE. Name of the node being monitored.
  -w WARNING, --warning=WARNING. Warning level.
  -c CRITICAL, --critical=CRITICAL. Critical level.`

	return check
}

func fetch(check *plugin.Plugin, URL string, Port string, Node string, Username string, Password string) []byte {
	client := &http.Client{Timeout: time.Duration(10 * time.Second)}

	req, err := http.NewRequest("GET", "http://"+URL+":"+Port+"/api/nodes/"+Node, nil)
	if err != nil {
		check.ExitUnknown("Cannot create connection request")
	}
	req.SetBasicAuth(Username, Password)
	resp, err := client.Do(req)
	if err != nil {
		check.ExitUnknown("Could Not Connect to " + URL + ":" + Port)
	}

	if resp.StatusCode != 200 {
		check.ExitCritical("Could Not Connect to " + URL + " Response Code: " + string(resp.Status))
	}

	bodyText, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		check.ExitUnknown("Error reading output of" + Node + ": " + string(bodyText))
	}
	return bodyText
}

func search(check *plugin.Plugin, result []byte) apiResponse {
	var pluginResponse apiResponse

	err := json.Unmarshal(result, &pluginResponse)
	if err != nil {
		check.ExitUnknown("Failed to parse JSON response")
	}
	return pluginResponse
}

func searchNodes(check *plugin.Plugin, HostName string, Port string, Username string, Password string) string {
	// Called if no node is provided for opts.Node
	// Returns node name as string

	client := &http.Client{Timeout: time.Duration(10 * time.Second)}

	req, err := http.NewRequest("GET", "http://"+HostName+":"+Port+"/api/nodes/", nil)
	if err != nil {
		check.ExitUnknown("Cannot create connection request")
	}
	req.SetBasicAuth(Username, Password)
	resp, err := client.Do(req)
	if err != nil {
		check.ExitUnknown("Could Not Connect to " + HostName + ":" + Port)
	}

	if resp.StatusCode != 200 {
		check.ExitCritical("Could Not Connect to " + HostName + " Response Code: " + string(resp.Status))
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		check.ExitUnknown("Error reading output of " + HostName + ": " + string(body))
	}

	type Node struct {
		NodeName string `json:"name"`
	}

	var newNode []Node

	err = json.Unmarshal(body, &newNode)
	if err != nil {
		check.ExitUnknown("Error %s", err)
	}

	if len(newNode) == 0 {
		check.ExitUnknown("No RabbitMQ Nodes found. Please check configuration.")
	} else if len(newNode) > 1 {
		nodes := ""
		for _, node := range newNode {
			nodes += node.NodeName + ", "
		}
		check.ExitUnknown("Multiple RabbitMQ Nodes found. Please specify node name (-n) from one of the following: " + nodes[:len(nodes)-2])
	}
	return newNode[0].NodeName
}

func findResult(check *plugin.Plugin, returnValue bool, returnMessage string, reverse bool) {
	returnStatus := plugin.OK

	if reverse == true {
		if returnValue == true {

			returnMessage = returnMessage + ": True"
			returnStatus = plugin.CRITICAL

		} else {
			check.AddMessage(returnMessage + ": False")
			returnMessage = returnMessage + ": False"
		}
	} else {
		if returnValue == false {
			returnMessage = returnMessage + ": False"
			returnStatus = plugin.CRITICAL

		} else {
			returnMessage = returnMessage + ": True"
		}
	}

	check.AddResult(returnStatus, returnMessage)
}

func alarmCheck(check *plugin.Plugin, returnValue bool, alarmName string) {
	if returnValue == true {
		check.ExitCritical(alarmName + ": True")
	}
}

func percentUsed(total float64, used float64) float64 {
	return (used / total) * 100
}

func removeDecimal(value float64) float64 {
	intValue := int(value)
	return float64(intValue)
}
func roundto2DP(input float64) string {
	Output := strconv.FormatFloat(input, 'f', 2, 64)
	return Output
}

func convertBytes(numberToConvert float64, startingUOM string, precision int) (string, string) {
	// Takes in a number that needs converting, the bytes UOM it is already in and requested precision of new value
	// Returns value and UOM, in form of most suitable UOM needed

	startingUOM = strings.ToUpper(startingUOM) // Ensure input is in uppercase

	units := []string{"B", "KB", "MB", "GB", "TB", "PB", "EB", "ZB", "YB"}

	var startingPoint int = 0 // Assume number is in bytes to begin with

	for i, unit := range units {
		// For all bytes units, find the index of the one that the value is already in
		if unit == startingUOM {
			startingPoint = i
		}
	}

	for _, unit := range units[startingPoint:] {
		// Starting at the index of the UOM the value is already in
		// Iterate over each UOM and divide by 1024 each time if needed

		if numberToConvert >= 1024 {
			// If >= 1024 then it can be shown better in the next UOM, so divide it
			numberToConvert /= 1024
		} else {
			// If < 1024, then highest UOM needed is found, so return value + UOM with specified precision
			newValue := strconv.FormatFloat(numberToConvert, 'f', precision, 64)
			return newValue, unit
		}
	}

	return strconv.FormatFloat(numberToConvert, 'f', precision, 64), startingUOM // Should never happen but returns input incase of errors
}