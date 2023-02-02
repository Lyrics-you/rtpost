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

## 命令

### --help

```shell
rtpost is a free and open source analyze live rtp packet loss rate system,
designed to analyse the udp packets of the specified port in the pcap file obtained by tcpdump, 
decode them into rtp packets and then analyse them to obtain the packet loss rate.

Usage:
  rtpost [flags]

Flags:
  -f, --file string       pcap file path        
  -h, --help              help for rtpost       
  -s, --ssrc string       rtp ssrc
  -v, --visual            rtp info visualization
```

## -f/--file

指定pcap的文件路径

## -s/--ssrc

指定rtp的ssrc

## -v/--visual

是否可视化



可以使用 rtpost -f "file" -s "ssrc" 进行分析丢包率，和wireshark比对，结果一致

程序有待后续改进
