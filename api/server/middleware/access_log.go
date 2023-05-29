package middleware

import (
	"bytes"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/SyntSugar/ss-infra-go/consts"
	"github.com/SyntSugar/ss-infra-go/machine"
	"github.com/gin-gonic/gin"
	"github.com/valyala/fasttemplate"
)

const (
	hostHeader                  = "Host"
	defaultSlowRequestThreshold = 1500 * time.Millisecond
)

var (
	bufferPool = sync.Pool{
		New: func() any {
			return bytes.NewBuffer(make([]byte, 256))
		},
	}
	dash    = []byte("-")
	newLine = byte('\n')
)

type LogItem struct {
	TraceId        string
	ContentLength  int64
	URL            *url.URL
	RequestHeader  http.Header
	ResponseHeader http.Header
	RemoteAddr     string
	Method         string
	Proto          string
	ReceivedAt     time.Time
	FirstByteTime  time.Time
	Latency        time.Duration
	BytesSent      int
	StatusCode     int
}

type AccessLogger struct {
	enabled  bool
	IP       string
	pattern  string
	writer   io.Writer
	template *fasttemplate.Template

	slowRequestThreshold time.Duration
}

// To set the enabled field as true by default to avoid the need for manual enabling of access log recording.
func NewAccessLogger(writer io.Writer, pattern string) (*AccessLogger, error) {
	accessLog := &AccessLogger{
		enabled: true,
		pattern: pattern,
		writer:  writer,
	}
	if err := accessLog.buildTemplate(); err != nil {
		return nil, err
	}
	localIP, err := machine.GetLocalIP()
	if err != nil {
		return nil, err
	}
	accessLog.IP = localIP
	accessLog.slowRequestThreshold = defaultSlowRequestThreshold

	return accessLog, nil
}

func (accessLogger *AccessLogger) Disabled() {
	accessLogger.enabled = false
}

func (accessLogger *AccessLogger) Enabled() {
	accessLogger.enabled = true
}

func (accessLogger *AccessLogger) Status() string {
	if accessLogger.enabled {
		return "enabled"
	}
	return "disabled"
}

func (accessLogger *AccessLogger) SetSlowRequestThreshold(duration time.Duration) {
	accessLogger.slowRequestThreshold = duration
}

func (accessLogger *AccessLogger) buildTemplate() (err error) {
	templateText := strings.NewReplacer(
		"%a", "${RemoteIP}",
		"%A", "${LocalIP}",
		"%b", "${BytesSent|-}",
		"%B", "${BytesSent|0}",
		"%H", "${Proto}",
		"%m", "${Method}",
		"%q", "${QueryString}",
		"%r", "${Method} ${RequestURI}",
		"%s", "${StatusCode}",
		"%t", "${ReceivedAt|02/Jan/2006:15:04:05 -0700}",
		"%U", "${URLPath}",
		"%D", "${Latency|-}",
		"%T", "${Latency|s}",
		"%F", "${FirstByteTime|-}",
	).Replace(accessLogger.pattern)

	timeFormatPattern := regexp.MustCompile("%({([^}]+)})t")
	templateText = timeFormatPattern.ReplaceAllString(templateText, "${ReceivedAt|$2}")
	requestHeaderPattern := regexp.MustCompile("%({([^}]+)})i")
	templateText = requestHeaderPattern.ReplaceAllStringFunc(templateText, strings.ToLower)
	templateText = requestHeaderPattern.ReplaceAllString(templateText, "${RequestHeader|$2}")
	responseHeaderPattern := regexp.MustCompile("%({([^}]+)})o")
	templateText = responseHeaderPattern.ReplaceAllStringFunc(templateText, strings.ToLower)
	templateText = responseHeaderPattern.ReplaceAllString(templateText, "${ResponseHeader|$2}")
	accessLogger.template, err = fasttemplate.NewTemplate(templateText, "${", "}")
	return
}

func (accessLogger *AccessLogger) record(item *LogItem) error {
	buf, _ := bufferPool.Get().(*bytes.Buffer)

	// buf need put after access writter write.
	buf.Reset()

	_, err := accessLogger.template.ExecuteFunc(buf,
		func(w io.Writer, tag string) (int, error) {
			tagLower := ""
			tagIndex := strings.Index(tag, "|")
			if tagIndex > 0 {
				tagLower = strings.ToLower(tag[:tagIndex])
			}

			switch tag {
			case "AM-Trace-ID":
				return w.Write([]byte(item.TraceId))
			case "Content-Length":
				return w.Write([]byte(strconv.FormatInt(item.ContentLength, 10)))
			case "RemoteIP":
				return w.Write([]byte(item.RemoteAddr))
			case "LocalIP":
				return w.Write([]byte(accessLogger.IP))
			case "BytesSent|-":
				if item.BytesSent == 0 {
					return w.Write(dash)
				}
				return w.Write([]byte(strconv.Itoa(item.BytesSent)))
			case "BytesSent|0":
				return w.Write([]byte(strconv.Itoa(item.BytesSent)))
			case "Proto":
				return w.Write([]byte(item.Proto))
			case "Method":
				return w.Write([]byte(item.Method))
			case "QueryString":
				return w.Write([]byte(item.URL.Query().Encode()))
			case "RequestURI":
				return w.Write([]byte(item.URL.RequestURI()))
			case "URLPath":
				return w.Write([]byte(item.URL.Path))
			case "StatusCode":
				return w.Write([]byte(strconv.Itoa(item.StatusCode)))
			case "Latency|-":
				return w.Write([]byte(strconv.FormatInt(item.Latency.Milliseconds(), 10)))
			case "Latency|s":
				return w.Write([]byte(strconv.FormatFloat(item.Latency.Seconds(), 'f', 3, 64)))
			case "FirstByteTime|-":
				if item.FirstByteTime.IsZero() {
					return w.Write(dash)
				}
				return w.Write([]byte(strconv.FormatInt(item.FirstByteTime.Sub(item.ReceivedAt).Milliseconds(), 10)))
			default:
				if tagIndex > 0 {
					switch tagLower {
					case "ReceivedAt":
						return w.Write([]byte(item.ReceivedAt.Format(string([]byte(tag)[tagIndex+1:]))))
					case "RequestHeader":
						hv := item.RequestHeader.Get(string([]byte(tag)[tagIndex+1:]))
						if hv == "" {
							return w.Write(dash)
						}
						return w.Write([]byte(hv))
					case "ResponseHeader":
						hv := item.ResponseHeader.Get(string([]byte(tag)[tagIndex+1:]))
						if hv == "" {
							return w.Write(dash)
						}
						return w.Write([]byte(hv))
					}
				}
			}
			return 0, nil
		})
	if err != nil {
		return err
	}

	buf.WriteByte(newLine)
	_, err = accessLogger.writer.Write(buf.Bytes())
	bufferPool.Put(buf)
	return err
}

func createLogItem(r *http.Request, w ResponseWriter, receivedAt time.Time, duration time.Duration) LogItem {
	logItem := LogItem{}

	r.Header.Set(hostHeader, r.Host)
	logItem.TraceId = r.Header.Get(consts.HeaderXCloudTraceContext)
	logItem.ContentLength = r.ContentLength
	logItem.URL = r.URL
	logItem.RequestHeader = r.Header
	if splitIndex := strings.Index(r.RemoteAddr, ":"); splitIndex != -1 {
		logItem.RemoteAddr = r.RemoteAddr[:splitIndex]
	}
	logItem.Method = r.Method
	logItem.Proto = r.Proto
	logItem.ReceivedAt = receivedAt
	logItem.Latency = duration
	logItem.ResponseHeader = w.Header()
	logItem.BytesSent = w.Size()
	logItem.StatusCode = w.Status()
	logItem.FirstByteTime = w.FirstByteTime()
	return logItem
}

func ignoreRequest(req *http.Request) bool {
	return req.Method == http.MethodOptions
}

func AccessLog(logger *AccessLogger) gin.HandlerFunc {
	return func(c *gin.Context) {
		if (!logger.enabled && logger.slowRequestThreshold <= 0) || ignoreRequest(c.Request) {
			c.Next()
			return
		}

		receivedAt := time.Now()
		original := c.Writer
		proxyWriter := newResponseWriter(c.Writer)
		if writer, ok := proxyWriter.(gin.ResponseWriter); ok {
			c.Writer = writer
		}

		c.Next()

		duration := time.Since(receivedAt)
		if !logger.enabled && duration < logger.slowRequestThreshold {
			return
		}
		logItem := createLogItem(c.Request, proxyWriter, receivedAt, duration)
		_ = logger.record(&logItem)
		c.Writer = original
	}
}
