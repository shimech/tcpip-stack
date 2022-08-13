package ethernet_tap

import (
	"os"
	"syscall"

	"github.com/shimech/tcpip-stack/net"
	"github.com/shimech/tcpip-stack/platform/ioctl"
)

type Device struct {
	index     int
	name      string
	dtype     net.DeviceType
	mtu       uint16
	flags     net.Flag
	hlen      uint16
	alen      uint16
	addr      net.Address
	peer      net.Address
	broadcast net.Address
	ifaces    []net.Iface
	fd        uintptr
	irq       syscall.Signal
}

const CLONE_DEVICE = "/dev/net/tun"

func NewDevice(name string, address string) *Device {
	d := &Device{}
	return d
}

func (d *Device) Index() int {
	return d.index
}

func (d *Device) SetIndex(i int) {
	d.index = i
}

func (d *Device) Name() string {
	return d.name
}

func (d *Device) SetName(n string) {
	d.name = n
}

func (d *Device) Type() net.DeviceType {
	return d.dtype
}

func (d *Device) SetType(t net.DeviceType) {
	d.dtype = t
}

func (d *Device) MTU() uint16 {
	return d.mtu
}

func (d *Device) SetMTU(mtu uint16) {
	d.mtu = mtu
}

func (d *Device) Flags() net.Flag {
	return d.flags
}

func (d *Device) SetFlags(f net.Flag) {
	d.flags = f
}

func (d *Device) Hlen() uint16 {
	return d.hlen
}

func (d *Device) SetHlen(hlen uint16) {
	d.hlen = hlen
}

func (d *Device) Alen() uint16 {
	return d.alen
}

func (d *Device) SetAlen(alen uint16) {
	d.alen = alen
}

func (d *Device) Addr() net.Address {
	return d.addr
}

func (d *Device) Peer() net.Address {
	return d.peer
}

func (d *Device) Broadcast() net.Address {
	return d.broadcast
}

func (d *Device) SetBroadcast(b net.Address) {
	d.broadcast = b
}

func (d *Device) Ifaces() []net.Iface {
	return d.ifaces
}

func (d *Device) PrependIface(i net.Iface) {
	d.ifaces = append([]net.Iface{i}, d.ifaces...)
}

func (d *Device) Open() error {
	f, err := os.OpenFile(CLONE_DEVICE, os.O_RDWR, 0600)
	if err != nil {
		return err
	}
	defer f.Close()
	d.fd = f.Fd()
	n, err := ioctl.TUNSETIFF(d.fd, d.name, syscall.IFF_TAP|syscall.IFF_NO_PI)
	if err != nil {
		return err
	}
	d.name = n
	flags, err := ioctl.SIOCGIFFLAGS(d.name)
	if err != nil {
		return err
	}
	flags |= (syscall.IFF_UP | syscall.IFF_RUNNING)
	if err := ioctl.SIOCSIFFLAGS(d.name, flags); err != nil {
		return err
	}
	return nil
}

// static int
// ether_tap_open(struct net_device *dev)
// {
//     struct ether_tap *tap;
//     struct ifreq ifr = {};

//     tap = PRIV(dev);
//     tap->fd = open(CLONE_DEVICE, O_RDWR);
//     if (tap->fd == -1) {
//         errorf("open: %s, dev=%s", strerror(errno), dev->name);
//         return -1;
//     }
//     strncpy(ifr.ifr_name, tap->name, sizeof(ifr.ifr_name)-1);
//     ifr.ifr_flags = IFF_TAP | IFF_NO_PI;
//     if (ioctl(tap->fd, TUNSETIFF, &ifr) == -1) {
//         errorf("ioctl [TUNSETIFF]: %s, dev=%s", strerror(errno), dev->name);
//         close(tap->fd);
//         return -1;
//     }

// /* Set Asynchronous I/O signal delivery destination */
// if (fcntl(tap->fd, F_SETOWN, getpid()) == -1) {
//     errorf("fcntl(F_SETOWN): %s, dev=%s", strerror(errno), dev->name);
//     close(tap->fd);
//     return -1;
// }
// /* Enable Asynchronous I/O */
// if (fcntl(tap->fd, F_SETFL, O_ASYNC) == -1) {
//     errorf("fcntl(F_SETFL): %s, dev=%s", strerror(errno), dev->name);
//     close(tap->fd);
//     return -1;
// }
// /* Use other signal instead of SIGIO */
// if (fcntl(tap->fd, F_SETSIG, tap->irq) == -1) {
//     errorf("fcntl(F_SETSIG): %s, dev=%s", strerror(errno), dev->name);
//     close(tap->fd);
//     return -1;
// }

//     if (memcmp(dev->addr, ETHER_ADDR_ANY, ETHER_ADDR_LEN) == 0) {
//         if (ether_tap_addr(dev) == -1) {
//             errorf("ether_tap_addr() failure, dev=%s", dev->name);
//             close(tap->fd);
//             return -1;
//         }
//     }
//     return 0;
// };

func (d *Device) Close() error {
	return nil
}

func (d *Device) Transmit() error {
	return nil
}

// static int
// ether_tap_close(struct net_device *dev)
// {
//     close(PRIV(dev)->fd);
//     return 0;
// }
