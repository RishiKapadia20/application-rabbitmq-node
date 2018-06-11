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

package rabbitmq

import "testing"
import "net"
import "fmt"
//import "bufio"
//import "strings" // only needed below for sample processing

/*
func main(){
  listener()
  testFetch()
}
8*/
func listenerServer() {

  fmt.Println("Launching server...")

  ln, err := net.Listen("tcp", ":15672")
  printError(err)
  // accept connection on port
  conn, err := ln.Accept()
  printError(err)
  // run loop forever (or until ctrl-c)
  for {
    // will listen for message to process ending in newline (\n)
    //message, _ := bufio.NewReader(conn).ReadString('\n')
    // sample process for string received
    newmessage := `[{"cluster_links":[],"mem_used":55214952,"mem_used_details":{"rate":420.8},"fd_used":23,"fd_used_details":{"rate":-0.4},"sockets_used":0,"sockets_used_details":{"rate":0.0},"proc_used":225,"proc_used_details":{"rate":0.0},"disk_free":13355036672,"disk_free_details":{"rate":0.0},"io_read_count":1,"io_read_count_details":{"rate":0.0},"io_read_bytes":1,"io_read_bytes_details":{"rate":0.0},"io_read_avg_time":0.032,"io_read_avg_time_details":{"rate":0.0},"io_write_count":0,"io_write_count_details":{"rate":0.0},"io_write_bytes":0,"io_write_bytes_details":{"rate":0.0},"io_write_avg_time":0.0,"io_write_avg_time_details":{"rate":0.0},"io_sync_count":0,"io_sync_count_details":{"rate":0.0},"io_sync_avg_time":0.0,"io_sync_avg_time_details":{"rate":0.0},"io_seek_count":0,"io_seek_count_details":{"rate":0.0},"io_seek_avg_time":0.0,"io_seek_avg_time_details":{"rate":0.0},"io_reopen_count":0,"io_reopen_count_details":{"rate":0.0},"mnesia_ram_tx_count":41,"mnesia_ram_tx_count_details":{"rate":0.0},"mnesia_disk_tx_count":13,"mnesia_disk_tx_count_details":{"rate":0.0},"msg_store_read_count":0,"msg_store_read_count_details":{"rate":0.0},"msg_store_write_count":0,"msg_store_write_count_details":{"rate":0.0},"queue_index_journal_write_count":0,"queue_index_journal_write_count_details":{"rate":0.0},"queue_index_write_count":0,"queue_index_write_count_details":{"rate":0.0},"queue_index_read_count":0,"queue_index_read_count_details":{"rate":0.0},"gc_num":41521158,"gc_num_details":{"rate":5.8},"gc_bytes_reclaimed":350905487240,"gc_bytes_reclaimed_details":{"rate":83347.2},"context_switches":173584554,"context_switches_details":{"rate":45.0},"io_file_handle_open_attempt_count":10,"io_file_handle_open_attempt_count_details":{"rate":0.0},"io_file_handle_open_attempt_avg_time":0.0587,"io_file_handle_open_attempt_avg_time_details":{"rate":0.0},"partitions":[],"os_pid":"99","fd_total":524288,"sockets_total":471767,"mem_limit":1658155827,"mem_alarm":false,"disk_free_limit":50000000,"disk_free_alarm":false,"proc_total":1048576,"rates_mode":"basic","uptime":2875667129,"run_queue":0,"processors":1,"exchange_types":[{"name":"topic","description":"AMQP topic exchange, as per the AMQP specification","enabled":true},{"name":"direct","description":"AMQP direct exchange, as per the AMQP specification","enabled":true},{"name":"headers","description":"AMQP headers exchange, as per the AMQP specification","enabled":true},{"name":"fanout","description":"AMQP fanout exchange, as per the AMQP specification","enabled":true}],"auth_mechanisms":[{"name":"PLAIN","description":"SASL PLAIN authentication mechanism","enabled":true},{"name":"AMQPLAIN","description":"QPid AMQPLAIN mechanism","enabled":true},{"name":"RABBIT-CR-DEMO","description":"RabbitMQ Demo challenge-response authentication mechanism","enabled":false}],"applications":[{"name":"amqp_client","description":"RabbitMQ AMQP Client","version":"3.6.5"},{"name":"asn1","description":"The Erlang ASN1 compiler version 4.0.3","version":"4.0.3"},{"name":"compiler","description":"ERTS  CXC 138 10","version":"7.0.1"},{"name":"crypto","description":"CRYPTO","version":"3.7"},{"name":"inets","description":"INETS  CXC 138 49","version":"6.3.2"},{"name":"kernel","description":"ERTS  CXC 138 10","version":"5.0.2"},{"name":"mnesia","description":"MNESIA  CXC 138 12","version":"4.14"},{"name":"mochiweb","description":"MochiMedia Web Server","version":"2.13.1"},{"name":"os_mon","description":"CPO  CXC 138 46","version":"2.4.1"},{"name":"public_key","description":"Public key infrastructure","version":"1.2"},{"name":"rabbit","description":"RabbitMQ","version":"3.6.5"},{"name":"rabbit_common","description":"","version":"3.6.5"},{"name":"rabbitmq_management","description":"RabbitMQ Management Console","version":"3.6.5"},{"name":"rabbitmq_management_agent","description":"RabbitMQ Management Agent","version":"3.6.5"},{"name":"rabbitmq_web_dispatch","description":"RabbitMQ Web Dispatcher","version":"3.6.5"},{"name":"ranch","description":"Socket acceptor pool for TCP protocols.","version":"1.2.1"},{"name":"sasl","description":"SASL  CXC 138 11","version":"3.0"},{"name":"ssl","description":"Erlang/OTP SSL application","version":"8.0.1"},{"name":"stdlib","description":"ERTS  CXC 138 10","version":"3.0.1"},{"name":"syntax_tools","description":"Syntax tools","version":"2.0"},{"name":"webmachine","description":"webmachine","version":"1.10.3"},{"name":"xmerl","description":"XML parser","version":"1.3.11"}],"contexts":[{"description":"RabbitMQ Management","path":"/","port":"15672"}],"log_file":"tty","sasl_log_file":"tty","db_dir":"/var/lib/rabbitmq/mnesia/rabbit@b44c00604528","config_files":["/etc/rabbitmq/rabbitmq.config"],"net_ticktime":60,"enabled_plugins":["rabbitmq_management"],"name":"rabbit@b44c00604528","type":"disc","running":true}]`
    // send new string back to client
    conn.Write([]byte(newmessage + "\n"))
  }
}

func TestFetch(t *testing.T){
  go listenerServer()
  expected := `[{"cluster_links":[],"mem_used":55214952,"mem_used_details":{"rate":420.8},"fd_used":23,"fd_used_details":{"rate":-0.4},"sockets_used":0,"sockets_used_details":{"rate":0.0},"proc_used":225,"proc_used_details":{"rate":0.0},"disk_free":13355036672,"disk_free_details":{"rate":0.0},"io_read_count":1,"io_read_count_details":{"rate":0.0},"io_read_bytes":1,"io_read_bytes_details":{"rate":0.0},"io_read_avg_time":0.032,"io_read_avg_time_details":{"rate":0.0},"io_write_count":0,"io_write_count_details":{"rate":0.0},"io_write_bytes":0,"io_write_bytes_details":{"rate":0.0},"io_write_avg_time":0.0,"io_write_avg_time_details":{"rate":0.0},"io_sync_count":0,"io_sync_count_details":{"rate":0.0},"io_sync_avg_time":0.0,"io_sync_avg_time_details":{"rate":0.0},"io_seek_count":0,"io_seek_count_details":{"rate":0.0},"io_seek_avg_time":0.0,"io_seek_avg_time_details":{"rate":0.0},"io_reopen_count":0,"io_reopen_count_details":{"rate":0.0},"mnesia_ram_tx_count":41,"mnesia_ram_tx_count_details":{"rate":0.0},"mnesia_disk_tx_count":13,"mnesia_disk_tx_count_details":{"rate":0.0},"msg_store_read_count":0,"msg_store_read_count_details":{"rate":0.0},"msg_store_write_count":0,"msg_store_write_count_details":{"rate":0.0},"queue_index_journal_write_count":0,"queue_index_journal_write_count_details":{"rate":0.0},"queue_index_write_count":0,"queue_index_write_count_details":{"rate":0.0},"queue_index_read_count":0,"queue_index_read_count_details":{"rate":0.0},"gc_num":41521158,"gc_num_details":{"rate":5.8},"gc_bytes_reclaimed":350905487240,"gc_bytes_reclaimed_details":{"rate":83347.2},"context_switches":173584554,"context_switches_details":{"rate":45.0},"io_file_handle_open_attempt_count":10,"io_file_handle_open_attempt_count_details":{"rate":0.0},"io_file_handle_open_attempt_avg_time":0.0587,"io_file_handle_open_attempt_avg_time_details":{"rate":0.0},"partitions":[],"os_pid":"99","fd_total":524288,"sockets_total":471767,"mem_limit":1658155827,"mem_alarm":false,"disk_free_limit":50000000,"disk_free_alarm":false,"proc_total":1048576,"rates_mode":"basic","uptime":2875667129,"run_queue":0,"processors":1,"exchange_types":[{"name":"topic","description":"AMQP topic exchange, as per the AMQP specification","enabled":true},{"name":"direct","description":"AMQP direct exchange, as per the AMQP specification","enabled":true},{"name":"headers","description":"AMQP headers exchange, as per the AMQP specification","enabled":true},{"name":"fanout","description":"AMQP fanout exchange, as per the AMQP specification","enabled":true}],"auth_mechanisms":[{"name":"PLAIN","description":"SASL PLAIN authentication mechanism","enabled":true},{"name":"AMQPLAIN","description":"QPid AMQPLAIN mechanism","enabled":true},{"name":"RABBIT-CR-DEMO","description":"RabbitMQ Demo challenge-response authentication mechanism","enabled":false}],"applications":[{"name":"amqp_client","description":"RabbitMQ AMQP Client","version":"3.6.5"},{"name":"asn1","description":"The Erlang ASN1 compiler version 4.0.3","version":"4.0.3"},{"name":"compiler","description":"ERTS  CXC 138 10","version":"7.0.1"},{"name":"crypto","description":"CRYPTO","version":"3.7"},{"name":"inets","description":"INETS  CXC 138 49","version":"6.3.2"},{"name":"kernel","description":"ERTS  CXC 138 10","version":"5.0.2"},{"name":"mnesia","description":"MNESIA  CXC 138 12","version":"4.14"},{"name":"mochiweb","description":"MochiMedia Web Server","version":"2.13.1"},{"name":"os_mon","description":"CPO  CXC 138 46","version":"2.4.1"},{"name":"public_key","description":"Public key infrastructure","version":"1.2"},{"name":"rabbit","description":"RabbitMQ","version":"3.6.5"},{"name":"rabbit_common","description":"","version":"3.6.5"},{"name":"rabbitmq_management","description":"RabbitMQ Management Console","version":"3.6.5"},{"name":"rabbitmq_management_agent","description":"RabbitMQ Management Agent","version":"3.6.5"},{"name":"rabbitmq_web_dispatch","description":"RabbitMQ Web Dispatcher","version":"3.6.5"},{"name":"ranch","description":"Socket acceptor pool for TCP protocols.","version":"1.2.1"},{"name":"sasl","description":"SASL  CXC 138 11","version":"3.0"},{"name":"ssl","description":"Erlang/OTP SSL application","version":"8.0.1"},{"name":"stdlib","description":"ERTS  CXC 138 10","version":"3.0.1"},{"name":"syntax_tools","description":"Syntax tools","version":"2.0"},{"name":"webmachine","description":"webmachine","version":"1.10.3"},{"name":"xmerl","description":"XML parser","version":"1.3.11"}],"contexts":[{"description":"RabbitMQ Management","path":"/","port":"15672"}],"log_file":"tty","sasl_log_file":"tty","db_dir":"/var/lib/rabbitmq/mnesia/rabbit@b44c00604528","config_files":["/etc/rabbitmq/rabbitmq.config"],"net_ticktime":60,"enabled_plugins":["rabbitmq_management"],"name":"rabbit@b44c00604528","type":"disc","running":true}]`
  response := fetch("127.0.0.1", "15672", "rabbitmq@fake", "guest", "guest")

  if expected != string(response){
    t.Error("Not expected text response: " + string(response))
  }else{
      print("Pass response: " + string(response))
  }
}

func printError(err error) {
	if err != nil {
		print(err)
  }
}
