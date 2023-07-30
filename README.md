# flexclient

This proof of concept tools listens for FlexRadio Discovery packets on 14992/UDP and retransmits them as UDP broadcast packets to 4992/UDP.

This tool was designed to run on VPN-connected clients using OpenVPN tun (layer 3) interfaces where the traditional FlexRadio broadcast packet on 4992/UDP is not able to transverse the VPN connection. A seperate tool (flextool) must runs on the same subnet as the FlexRadio SDRs that listens for discovery packetds on 4992/udp and then retransmits them to VPN client ips on 14992/udp.

## Usage

Download and execute on your OpenVPN connected client.

By default it listens for discovery packets sent on 14992/udp and retransmits them to 255.255.255.255:4992.

The listening port (default 14992/udp) can be changed at runtime using the port (`-p`) flag.

### Example Usage

```
flexclient.exe

flexclient.exe -p 15992
```
