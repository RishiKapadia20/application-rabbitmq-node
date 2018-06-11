# RabbitMQ Opspack

RabbitMQ is message-queueing software, commonly labeled as a message broker or queue manager, that provides a method of exchanging data between processes, applications, and servers. It is a software where queues can be defined and gives your applications a common platform to send/receive messages and a secure place for them to reside until the message is received.

Here are a few highlight features of RabbitMQ:

* Robust messaging for applications
* Easy to use
* Runs on all major operating systems
* Supports a huge number of developer platforms
* Open source and commercially supported

## What You Can Monitor

Opsview’s RabbitMQ Opspack is designed to pull important information from the RabbitMQ management API using your RabbitMQ credentials. If you are looking to build a hierarchy of your RabbitMQ infrastructure, we can help.  Opsview can raise a notification, monitor using our time series graphs and even tie into a larger dashboard view for important services such as IO issues, memory utilization and more. Once downloaded, you can set notifications and thresholds within a few clicks.

## Service Checks

| Service Check | Description |
|:------------- |:----------- |
|check_mem_alarm | Returns critical if the memory alarm has gone off
|context_switches_details | Rate at which context switching takes place on this node during last statistics interval
|disk_free | Disk free space in bytes, will go critical if the disk_free_alarm is true
|fd_left | Number of file descriptors remaining
|fd_used | Percentage of file descriptors used
|io_read_avg_time | Average wait time (ms) for each disk read operation in the last statistics interval
|io_read_count | Rate of read operations by the persister in the last statistics interval
|io_seek_avg_time | Average wait time (ms) for each seek operation in the last statistics interval
|io_sync_avg_time | Average wait time (ms) for each fsync() operation in the last statistics interval
|io_write_avg_time | Average wait time (ms) for each disk write operation in the last statistics interval
|io_write_count | Rate of write operations by the persister in the last statistics interval
|mem_used | Total memory used
|running | Whether or not this node is up.
|sockets_left | File descriptors available for use as sockets remaining
|sockets_used | Percentage of file descriptors used as sockets

## Setup and Configuration

This plugin requires the node name of the node (stored in variable '%RABBITMQ_CREDENTIALS:3%') to monitor as well as the address (host).

Requires the node to be monitored to expose the management API, see here for details.

#### Configuring RabbitMQ
Enabling the management API is done using the below command on the host to monitor:

```rabbitmq-plugins enable rabbitmq_management```

A new user for RabbitMQ with permissions - "monitoring" should be created, as node level data is required.

Example basic rabbitmq.config file:
```
[
    {rabbit,
        [
            %% if want to test, using user/pass
            %% guest/guest - and monitor from a non-loopback interface
            %% uncomment next line to allow connections from any ip
            {loopback_users, []},
            %% In milliseconds
            {collect_statistics_interval,    300000}
        ]
    },
    {rabbitmq_management,
        [
            %% should be set to basic or detailed
            {rates_mode,    basic},
            %% change port to query management api here
            {listener,        [{port,    15672}]},
            {sample_retention_policies,
                %% list of {maxAgeSeconds,sampleEveryNSeconds}
                [
                    {global,    [{30000,300}]},
                    {basic,        [{30000,300}]}
                ]
            }
        ]
    }
].
```

#### Setting up Opsview to monitor RabbitMQ


Step 1: Add the host template
![Add host template](/docs/img/host-template.png?raw=true)

Step 2: Configure RABBITMQ_CREDENTIALS variable with username, password and 3rd argument as node name (eg rabbit@opsview-VM)

![Add host template](/docs/img/variables.png?raw=true)

Step 3: Reload and the system will now be monitored

![Add host template](/docs/img/output.png?raw=true)
