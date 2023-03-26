# rtpost
rtpost is a free and open source analyze media rtp packet loss rate system, designed to analyse the udp packets of the specified port in the pcap file obtained by tcpdump,  decode them into rtp packets and then analyse them to obtain the packet loss rate.

## 说明

该包针对的是使用RTP协议的音视频服务，对使用tcpdump抓取的包，将其中的UDP包按RTP协议进行解码，使用端口和SSRC作为过滤，统计对应SSRC的丢包率。

解码的RTP结构:

```go
type RTP struct {
	// byte 0
	Version   uint8 // 2 bits
	Padding   uint8 // 1 bit
	Extension uint8 // 1 bit
	CSRC      uint8 // 4 bits, Contributing source

	// byte 1
	Marker       uint8 // 1 bit
	PlayloadType uint8 // 7 bit

	// byte 2-3
	SequenceNumber uint16 // 16 bits

	// byte 4-7
	Timestamp uint32 // 32 bits

	// byte 8-11
	SSRC uint32 // 32 bits, stream_id used here
}
```

这里过滤端口为：

```go
udpPorts []string = []string{"8089", "8090", "9600", "9601", "9602", "9603", "9604", "9605", "18089", "18090", "19600", "19601", "19602", "19603", "19604", "19605"}
```

## 参数

### --help

```shell
rtpost is a free and open source analyze live rtp packet loss rate system,
designed to analyse the udp packets of the specified port in the pcap file obtained by tcpdump, 
decode them into rtp packets and then analyse them to obtain the packet loss rate.

Usage:
  rtpost [flags]
  rtpost [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  decode      decode rtp packets
  delta       delta rtp packets
  help        Help about any command
  loss        loss show rtp lost packets
  stat        stat rtp packets's information

Flags:
  -c, --csrc string         rtp csrc is hexadecimal, prefixed with 0x
  -d, --deduplicate         rtp info deduplicate by ip, port and ssrc
  -f, --file string         pcap file path
  -g, --group               rtp info deduplicate in group
  -h, --help                help for rtpost
  -s, --ssrc string         rtp ssrc is hexadecimal, prefixed with 0x
  -t, --time-delta uint32   show bigger than delta's packet
  -v, --version             show rtpost version

Use "rtpost [command] --help" for more information about a command.
```

### -v/--version

显示rtpost的版本信息

### -f/--file

指定pcap的文件路径

### -s/--ssrc

指定rtp的ssrc,十六进制且以0x开头，忽略大小写，没有指定时包括所有的SSRC

### -c/--ssrc

指定rtp的csrc，十六进制且以0x开头，忽略大小写，通过CSRC找到对应的SSRC

### -d/--deduplicate

是否对rtp包进行去重，将重复seq的包进行过滤,没有指定时统计结果与wireshark一致

### -g/--group

是否对rtp包进行分组，分组方式按照srcIP:srcPort->dstIP:dstPort[SSRC]进行分组

### -t/--delta

在delta命令时，显示时延大于delta秒的RTP包

## 命令

### decode

展示指定pcap的rtp包的解码结果 ，可以指定-s,-c,-dg等参数

`rtpost -f "demo.pcap" -dg -s SSRC decode`

```shell
rtpost -f "..\pcap\2023-03-23\wifi\15-30.pcapng" -s 0x62CCF3F -dg  decode
[2023-03-23 14:15:19.194727] 10.148.60.126:8089->10.128.63.152:23459 RTP[708] {PT=DynamicRTP-116, SSRC=0x62CCF3F9, Seq=19993, Time=4145753973, CSRC=1, [0xECBD19D1]}
[2023-03-23 14:15:19.194727] 10.148.60.126:8089->10.128.63.152:23459 RTP[708] {PT=DynamicRTP-116, SSRC=0x62CCF3F9, Seq=19994, Time=4145753973, CSRC=1, [0xECBD19D1]}
[2023-03-23 14:15:19.194727] 10.148.60.126:8089->10.128.63.152:23459 RTP[708] {PT=DynamicRTP-116, SSRC=0x62CCF3F9, Seq=19995, Time=4145753973, CSRC=1, [0xECBD19D1]}
[2023-03-23 14:15:19.194727] 10.148.60.126:8089->10.128.63.152:23459 RTP[708] {PT=DynamicRTP-116, SSRC=0x62CCF3F9, Seq=19996, Time=4145753973, CSRC=1, [0xECBD19D1]}
[2023-03-23 14:15:19.194727] 10.148.60.126:8089->10.128.63.152:23459 RTP[708] {PT=DynamicRTP-116, SSRC=0x62CCF3F9, Seq=19997, Time=4145753973, CSRC=1, [0xECBD19D1]}
... ...
```

### stat

展示rtp统计后的结果（丢包率，期望收到的rtp包的数量，丢失的包的数量，ssrc以及csrc信息）

`rtpost -f "demo.pcap" -dg -s SSRC stat`

```shell
rtpost -f "..\pcap\2023-03-23\wifi\15-30.pcapng" -dg stat
10.128.63.152:24560->10.148.60.126:8089 {SSRC:0x4DD38B09 CSRC:[]} (recv:43968, expect:43968, loss:0, rate:0.00%)
10.128.63.152:23459->10.148.60.126:8089 {SSRC:0x9C687F33 CSRC:[]} (recv:13209, expect:13209, loss:0, rate:0.00%)
10.128.63.152:23459->10.148.60.126:8089 {SSRC:0xE6992304 CSRC:[]} (recv:44292, expect:44292, loss:0, rate:0.00%)
10.148.60.126:8089->10.128.63.152:24560 {SSRC:0x4CCEB77F CSRC:[0x4DD38B09]} (recv:43676, expect:43965, loss:289, rate:0.66%)
10.148.60.126:8089->10.128.63.152:24560 {SSRC:0xD6E84261 CSRC:[0x7948C4E4]} (recv:43601, expect:43889, loss:288, rate:0.66%)   
10.148.60.126:8089->10.128.63.152:23459 {SSRC:0x62CCF3F9 CSRC:[0xECBD19D1]} (recv:206107, expect:207588, loss:1481, rate:0.71%)
```

### loss

展示rtp丢失包的结果，丢失包首位信息较全，中间丢失的包会进行标注

`rtpost -f "demo.pcap" -dg -s SSRC loss`

```shell
rtpost -f "..\pcap\2023-03-23\wifi\15-30.pcapng" -s 0x62CCF3F -dg  loss
... ...
----------------------------
[2023-03-23 14:28:57.748712] 10.148.60.126:8089->10.128.63.152:23459 RTP[703] {PT=DynamicRTP-116, SSRC=0x62CCF3F9, Seq=17977, Time=4219430223, CSRC=1, [0xECBD19D1]}
|- lost Seq=17978
|- lost Seq=17979
|- lost Seq=17980
[2023-03-23 14:28:57.748712] 10.148.60.126:8089->10.128.63.152:23459 RTP[703] {PT=DynamicRTP-116, SSRC=0x62CCF3F9, Seq=17981, Time=4219430223, CSRC=1, [0xECBD19D1]}
----------------------------
[2023-03-23 14:28:57.748712] 10.148.60.126:8089->10.128.63.152:23459 RTP[734] {PT=DynamicRTP-116, SSRC=0x62CCF3F9, Seq=17994, Time=4219432023, CSRC=1, [0xECBD19D1]}
|- lost Seq=17995
|- lost Seq=17996
|- lost Seq=17997
[2023-03-23 14:28:57.748712] 10.148.60.126:8089->10.128.63.152:23459 RTP[720] {PT=DynamicRTP-116, SSRC=0x62CCF3F9, Seq=17998, Time=4219437693, CSRC=1, [0xECBD19D1]}
```

### delta

展示相隔包间隔大于的包的信息

`rtpost -f "demo.pcap -dg -s SSRC delta -t 3"`

可以使用 rtpost -f "file" -s "ssrc" 进行分析丢包率，和wireshark比对，结果一致

```shell
rtpost -f "..\pcap\2023-03-23\wifi\15-30.pcapng" -s 0x62CCF3F -dg  delta -t 1
SSRC=0x62CCF3F9 delta(>=1s)  :  3
[2023-03-23 14:17:54.110683][1.243s] 10.148.60.126:8089->10.128.63.152:23459 RTP[694] {PT=DynamicRTP-116, SSRC=0x62CCF3F9, Seq=56455, Time=4159518933, CSRC=1, [0xECBD19D1]}
[2023-03-23 14:25:07.935139][1.143s] 10.148.60.126:8089->10.128.63.152:23459 RTP[718] {PT=DynamicRTP-116, SSRC=0x62CCF3F9, Seq=29088, Time=4198634283, CSRC=1, [0xECBD19D1]}
[2023-03-23 14:25:21.007613][1.384s] 10.148.60.126:8089->10.128.63.152:23459 RTP[730] {PT=DynamicRTP-116, SSRC=0x62CCF3F9, Seq=31684, Time=4199665143, CSRC=1, [0xECBD19D1]}
```

## 未来

针对窗口统计丢包率
