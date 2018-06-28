package main

import (
	"bytes"
	"net/http"
	"os/exec"
	"strings"
	"time"
)

func main() {
	go startListener()

	time.Sleep(time.Second)

	runTests()
}

func runTests() {
	checks := []string{"mem_used", "fd_used", "sockets_used", "disk_free", "fd_left", "fd_used_percent", "running", "sockets_left", "sockets_used_percent", "mem_alarm", "io_read_count", "io_write_count", "io_write_avg_time", "io_read_avg_time", "io_seek_avg_time", "io_sync_avg_time", "context_switch"}
	warnings := []string{"0", "50"}
	criticals := []string{"0", "51"}
	expectedOk := []string{"OK: Mem_Used: 1 | Mem_Used=1kb;;;;", "OK: Fd_Used: 2 | Fd_Used=2;;;;", "OK: Sockets_Used: 3 | Sockets_Used=3;;;;", "OK: Disk_Free: 4 | Disk_Free=4kb;;;;", "OK: Fd_Left: 15 | Fd_Left=15;;;;", "OK: Fd_Used_Percent: 11 | Fd_Used_Percent=11%;;;;", "OK: Running: True", "OK: Sockets_Left: 16 | Sockets_Left=16;;;;", "OK: Sockets_Used_Percent: 15 | Sockets_Used_Percent=15%;;;;", "OK: Memory_Alarm: False", "OK: Io_Read_Count: 5 | Io_Read_Count=5;;;;", "OK: Io_Write_Count: 8 | Io_Write_Count=8;;;;", "OK: Io_Write_Avg_Time: 10 | Io_Write_Avg_Time=10ms;;;;", "OK: Io_Read_Avg_Time: 7 | Io_Read_Avg_Time=7ms;;;;", "OK: Io_Seek_Avg_Time: 14 | Io_Seek_Avg_Time=14ms;;;;", "OK: Io_Sync_Avg_Time: 12 | Io_Sync_Avg_Time=12ms;;;;", "OK: Context_Switch: 16 | Context_Switch=16;;;;"}
	expectedWarning := []string{"WARNING: Mem_Used: 1 | Mem_Used=1kb;;;;", "WARNING: Fd_Used: 2 | Fd_Used=2;;;;", "WARNING: Sockets_Used: 3 | Sockets_Used=3;;;;", "WARNING: Disk_Free: 4 | Disk_Free=4kb;;;;", "WARNING: Fd_Left: 15 | Fd_Left=15;;;;", "WARNING: Fd_Used_Percent: 11 | Fd_Used_Percent=11%;;;;", "WARNING: Running: True", "WARNING: Sockets_Left: 16 | Sockets_Left=16;;;;", "WARNING: Sockets_Used_Percent: 15 | Sockets_Used_Percent=15%;;;;", "WARNING: Memory_Alarm: False", "OK: Io_Read_Count: 5 | Io_Read_Count=5;;;;", "OK: Io_Write_Count: 8 | Io_Write_Count=8;;;;", "WARNING: Io_Write_Avg_Time: 10 | Io_Write_Avg_Time=10ms;;;;", "WARNING: Io_Read_Avg_Time: 7 | Io_Read_Avg_Time=7ms;;;;", "WARNING: Io_Seek_Avg_Time: 14 | Io_Seek_Avg_Time=14ms;;;;", "WARNING: Io_Sync_Avg_Time: 12 | Io_Sync_Avg_Time=12ms;;;;", "WARNING: Context_Switch: 16 | Context_Switch=16;;;;"}
	expectedCritical := []string{"CRITICAL: Mem_Used: 1 | Mem_Used=1kb;;;;", "CRITICAL: Fd_Used: 2 | Fd_Used=2;;;;", "CRITICAL: Sockets_Used: 3 | Sockets_Used=3;;;;", "CRITICAL: Disk_Free: 4 | Disk_Free=4kb;;;;", "CRITICAL: Fd_Left: 15 | Fd_Left=15;;;;", "CRITICAL: Fd_Used_Percent: 11 | Fd_Used_Percent=11%;;;;", "CRITICAL: Running: True", "CRITICAL: Sockets_Left: 16 | Sockets_Left=16;;;;", "CRITICAL: Sockets_Used_Percent: 15 | Sockets_Used_Percent=15%;;;;", "CRITICAL: Memory_Alarm: False", "OK: Io_Read_Count: 5 | Io_Read_Count=5;;;;", "OK: Io_Write_Count: 8 | Io_Write_Count=8;;;;", "CRITICAL: Io_Write_Avg_Time: 10 | Io_Write_Avg_Time=10ms;;;;", "CRITICAL: Io_Read_Avg_Time: 7 | Io_Read_Avg_Time=7ms;;;;", "CRITICAL: Io_Seek_Avg_Time: 14 | Io_Seek_Avg_Time=14ms;;;;", "CRITICAL: Io_Sync_Avg_Time: 12 | Io_Sync_Avg_Time=12ms;;;;", "CRITICAL: Context_Switch: 16 | Context_Switch=16;;;;"}

	for i := range checks {
		if (checks[i] == "disk_free") || (checks[i] == "fd_left") || (checks[i] == "sockets_left") {
			runRabbit(checks[i], warnings[0], criticals[0], expectedOk[i])
			runRabbit(checks[i], warnings[1], criticals[0], expectedWarning[i])
			runRabbit(checks[i], warnings[1], criticals[1], expectedCritical[i])
		} else if (checks[i] == "running") || (checks[i] == "mem_alarm") {
			runRabbit(checks[i], warnings[0], criticals[0], expectedOk[i])
		} else {
			runRabbit(checks[i], warnings[1], criticals[1], expectedOk[i])
			runRabbit(checks[i], warnings[0], criticals[1], expectedWarning[i])
			runRabbit(checks[i], warnings[0], criticals[0], expectedCritical[i])
		}
	}
}

func runRabbit(check string, warning string, critical string, expected string) {
	cmd := exec.Command("go", "run", "check_rabbitmq_node.go", "-m", check, "-H", "127.0.0.1", "-P", "3333", "-p", "guest", "-u", "guest", "-n", "fakeNode", "-w", warning, "-c", critical)
	var actual bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &actual
	cmd.Stderr = &stderr

	cmd.Run()

	actualString := strings.TrimSpace(actual.String())
	if expected == actualString {
		print("PASS " + actual.String() + "\n")
	} else {
		print("FAIL Expecting: " + expected + " Actual: " + actualString + "\n")
	}
}

func startListener() {
	print("Server start\n")
	http.HandleFunc("/", jsonSetup)
	http.ListenAndServe("127.0.0.1:3333", nil)
}

func jsonSetup(w http.ResponseWriter, r *http.Request) {
	js := `{"mem_used": 1024,"Fd_used": 2,"Sockets_used": 3,"Disk_free": 4099,"Io_read_count": 5,"Io_read_bytes": 6,"Io_read_avg_time": 7,"Io_write_count": 8,"Io_write_bytes": 9,"Io_write_avg_time": 10,"Io_sync_count": 11,"Io_sync_avg_time": 12,"Io_seek_count": 13,"Io_seek_avg_time": 14,"Context_switches": 15,"Context_switches_details": {"rate": 16},"Fd_total": 17,"Sockets_total": 19,"Mem_alarm": false,"Running": true}`

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(js))
}
