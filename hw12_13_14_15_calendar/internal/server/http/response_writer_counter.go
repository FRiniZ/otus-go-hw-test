package internalhttp

import (
	"bufio"
	"fmt"
	"net"
	"net/http"
	"strconv"
	"strings"
	"sync/atomic"
	"time"
)

// stolen idea from https://github.com/miolini/datacounter/blob/master/response_writer.go

type ResponseWriterCounter struct {
	http.ResponseWriter
	r          *http.Request
	count      uint64
	started    time.Time
	statusCode int
}

func NewResponseWriterCounter(rw http.ResponseWriter, r *http.Request) *ResponseWriterCounter {
	return &ResponseWriterCounter{
		ResponseWriter: rw,
		r:              r,
		count:          0,
		started:        time.Now(),
		statusCode:     0,
	}
}

func (rwl *ResponseWriterCounter) Write(buf []byte) (int, error) {
	n, err := rwl.ResponseWriter.Write(buf)
	atomic.AddUint64(&rwl.count, uint64(n))
	return n, err
}

func (rwl *ResponseWriterCounter) Header() http.Header {
	return rwl.ResponseWriter.Header()
}

func (rwl *ResponseWriterCounter) WriteHeader(statusCode int) {
	rwl.statusCode = statusCode
	rwl.Header().Set("X-Runtime", fmt.Sprintf("%.6f", time.Since(rwl.started).Seconds()))
	rwl.ResponseWriter.WriteHeader(statusCode)
}

func (rwl *ResponseWriterCounter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	return rwl.ResponseWriter.(http.Hijacker).Hijack()
}

func (rwl *ResponseWriterCounter) Count() uint64 {
	return atomic.LoadUint64(&rwl.count)
}

func (rwl *ResponseWriterCounter) Started() time.Time {
	return rwl.started
}

func (rwl *ResponseWriterCounter) StatusCode() int {
	return rwl.statusCode
}

func (rwl *ResponseWriterCounter) String() string {
	var b strings.Builder

	ip, _, err := net.SplitHostPort(rwl.r.RemoteAddr)
	if err != nil {
		return fmt.Sprintf("userip: %q is not IP\n", rwl.r.RemoteAddr)
	}

	//	userIP := net.ParseIP(ip)
	//	if userIP == nil {
	//		return fmt.Sprintf("userip: %q is not IP:port", rwl.r.RemoteAddr)
	//	}

	b.WriteString(ip)
	b.WriteString(" ")
	b.WriteString(rwl.started.Format("02/Jan/2006:15:04:05 -0700"))
	b.WriteString(" ")
	b.WriteString(rwl.r.Method)
	b.WriteString(" ")
	b.WriteString(rwl.r.RequestURI)
	b.WriteString(" ")
	b.WriteString(rwl.r.Proto)
	b.WriteString(" ")
	b.WriteString(strconv.Itoa(rwl.statusCode))
	b.WriteString(" ")
	b.WriteString(strconv.FormatUint(rwl.count, 10))
	b.WriteString(" \"")
	b.WriteString(rwl.r.Header.Get("User-Agent"))
	b.WriteString("\"\n")

	return b.String()
}
