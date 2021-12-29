package forward

import (
	"io"
	"log"
	"net"
	"strings"
	"sync"

	dc "github.com/digitalcircle-com-br/gw/lib/base"
	"github.com/digitalcircle-com-br/random"
)

type Listener struct {
	lconnstr string
	rconnstr string
	closed   bool
	listener net.Listener
	conns    sync.Map
}

func (l *Listener) Init(k string, v string) {
	prel, ok := listeners.Load(k)
	if ok {
		prelListenerPtr := prel.(*Listener)
		if prelListenerPtr.lconnstr != k || prelListenerPtr.rconnstr != v {
			prelListenerPtr.Close()
		} else {
			return
		}
	}

	l.lconnstr = k
	l.rconnstr = v
	l.conns = sync.Map{}
	kparts := strings.Split(k, ":")
	vparts := strings.Split(v, ":")

	go func() {
		listener, err := net.Listen(kparts[0], strings.Join(kparts[1:], ":"))
		l.listener = listener
		if err != nil {
			log.Printf("Error setting forward for: %s => %v: %s", k, v, err.Error())
			return
		}
		for {
			lconn, err := l.listener.Accept()
			lconnid := random.Str(63)
			rconnid := random.Str(63)
			if err != nil {
				log.Printf("Error setting forwad for: %s => %s: %s", k, v, err.Error())
				if lconn != nil {
					lconn.Close()
				}

				if l.closed {
					return
				}
				continue
			}
			rconn, err := net.Dial(vparts[0], strings.Join(vparts[1:], ":"))
			if err != nil {
				log.Printf("Error connecting forwad for: %s => %s: %s", k, v, err.Error())
				lconn.Close()
				if rconn != nil {
					rconn.Close()
				}
				if l.closed {
					return
				}
				continue
			}

			go func(lid, rid string, rc, lc *net.Conn) {
				io.Copy(*rc, *lc)
				(*lc).Close()
				(*rc).Close()
				l.conns.Delete(lid)
				l.conns.Delete(rid)

			}(lconnid, rconnid, &lconn, &rconn)

			go func(lid, rid string, rc, lc *net.Conn) {
				io.Copy(*lc, *rc)
				(*lc).Close()
				(*rc).Close()
				l.conns.Delete(lid)
				l.conns.Delete(rid)
			}(lconnid, rconnid, &lconn, &rconn)

			l.conns.Store(lconnid, &lconn)
			l.conns.Store(rconnid, &rconn)
		}
	}()

	dc.Log("Initiating Conn FWD to %s=>%s", k, v)
	listeners.Store(l.lconnstr, l)

}
func (l *Listener) Close() {
	if l.listener != nil {
		l.listener.Close()
	}

	l.conns.Range(func(key, value interface{}) bool {
		(*value.(*net.Conn)).Close()
		return true
	})
	l.conns = sync.Map{}
	l.closed = true
	listeners.Delete(l.lconnstr)
}

var listeners sync.Map

func Load(fwds map[string]string) {

	for k, v := range fwds {
		l := Listener{}
		l.Init(k, v)
	}
	listeners.Range(func(key, value interface{}) bool {
		kstr := key.(string)
		vlis := value.(*Listener)
		_, ok := fwds[kstr]
		if !ok {
			vlis.Close()
		}
		return true
	})

}
