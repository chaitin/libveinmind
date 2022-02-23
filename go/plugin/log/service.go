package log

import (
	"context"
	"os"
	"time"

	"golang.org/x/sync/errgroup"
	"golang.org/x/xerrors"

	"github.com/chaitin/libveinmind/go/plugin/service"
)

const Namespace = "github.com/chaitin/libveinmind/logging"

type logConfig struct {
	Level Level         `json:"level"`
	Delay time.Duration `json:"delay"`
}

// clientCore is the core used when a plugin is hosted,
// and in which case the log should be written back to
// the host for further processing.
type clientCore struct {
	ctx     context.Context
	group   *errgroup.Group
	level   Level
	logCh   chan Log
	closeCh chan struct{}
	log     func([]Log) error
}

func (c *clientCore) Enabled(l Level) bool {
	return l <= c.level
}

func (c *clientCore) Do(l Log) {
	select {
	case <-c.ctx.Done():
		return
	case c.logCh <- l:
	}

	// Handle the level of fatal and panic to keep in sync
	// with the defined behaviour.
	if l.Level < ErrorLevel {
		c.CloseWait()
		if l.Level == PanicLevel {
			panic(l.Msg)
		} else {
			os.Exit(1)
		}
	}
}

func (c *clientCore) CloseWait() {
	select {
	case <-c.ctx.Done():
	case c.closeCh <- struct{}{}:
	}
	_ = c.group.Wait()
}

func (c *clientCore) runBufferThread(
	d time.Duration, bufferCh chan<- []Log,
) error {
	ticker := time.NewTicker(d)
	defer ticker.Stop()
	var buffer []Log
	flushBuffer := func() {
		b := buffer
		buffer = nil
		select {
		case <-c.ctx.Done():
		case bufferCh <- b:
		}
	}
	for {
		tickerCh := ticker.C
		if buffer == nil {
			tickerCh = nil
		}
		select {
		case <-c.ctx.Done():
			return nil
		case <-tickerCh:
			flushBuffer()
		case <-c.closeCh:
			flushBuffer()
			return xerrors.New("logger closed")
		case item := <-c.logCh:
			buffer = append(buffer, item)
		}
	}
}

func (c *clientCore) runDirectThread(bufferCh chan<- []Log) error {
	for {
		var l Log
		select {
		case <-c.ctx.Done():
			return nil
		case l = <-c.logCh:
		}
		select {
		case <-c.ctx.Done():
			return nil
		case bufferCh <- []Log{l}:
		}
	}
}

func (c *clientCore) runCallThread(bufferCh <-chan []Log) error {
	for {
		select {
		case <-c.ctx.Done():
			return nil
		case buffer := <-bufferCh:
			if err := c.log(buffer); err != nil {
				return err
			}
		}
	}
}

func newClientCore() (*clientCore, error) {
	if !service.Hosted() {
		return nil, xerrors.New("client is not hosted")
	}
	var manifest struct{}
	err := service.GetManifest(Namespace, &manifest)
	if err != nil {
		return nil, err
	}
	var getConfig func() (*logConfig, error)
	service.GetService(Namespace, "getConfig", &getConfig)
	cfg, err := getConfig()
	if err != nil {
		return nil, err
	}
	var log func([]Log) error
	service.GetService(Namespace, "log", &log)
	group, ctx := errgroup.WithContext(context.Background())
	result := &clientCore{
		ctx:     ctx,
		group:   group,
		level:   cfg.Level,
		logCh:   make(chan Log),
		closeCh: make(chan struct{}),
		log:     log,
	}
	bufferCh := make(chan []Log)
	group.Go(func() error {
		return result.runCallThread(bufferCh)
	})
	if cfg.Delay <= 0 {
		group.Go(func() error {
			return result.runDirectThread(bufferCh)
		})
	} else {
		group.Go(func() error {
			return result.runBufferThread(cfg.Delay, bufferCh)
		})
	}
	return result, nil
}

type loggerService struct {
	level  Level
	delay  time.Duration
	fields Fields
	core   Core
}

func (s *loggerService) getConfig() logConfig {
	return logConfig{
		Level: s.level,
		Delay: s.delay,
	}
}

func (s *loggerService) log(buffer []Log) {
	for _, item := range buffer {
		if s.fields != nil {
			if item.Fields == nil {
				item.Fields = make(Fields)
			}
			for k, v := range s.fields {
				// TODO: preserve the field from recursive
				// plugin invocation.
				item.Fields[k] = v
			}
		}
		if item.Level < ErrorLevel {
			// Clamp the log level, since the fatal and panic
			// in the plugin is not equivalent to the fatal
			// and panic of the host program.
			item.Level = ErrorLevel
		}
		s.core.Do(item)
	}
}

func (s *loggerService) Add(registry *service.Registry) {
	registry.Define(Namespace, struct{}{})
	registry.AddService(Namespace, "getConfig", s.getConfig)
	registry.AddService(Namespace, "log", s.log)
}

type serviceOption struct {
	level Level
	delay time.Duration
}

// ServiceOption are the options for creating logger services
// and provide it to plugins to execute.
type ServiceOption func(*serviceOption)

// WithMaxLevel sets the maximum level of log that will be
// transferred back to the host from plugin.
func WithMaxLevel(level Level) ServiceOption {
	return func(opt *serviceOption) {
		opt.level = level
	}
}

// WithBufferDelay sets timeout of log buffering before they
// could be sent back to the host.
func WithBufferDelay(d time.Duration) ServiceOption {
	return func(opt *serviceOption) {
		opt.delay = d
	}
}

// WithNoBuffer tells the plugin to send back logs to the host
// without internal buffering.
func WithNoBuffer() ServiceOption {
	return func(opt *serviceOption) {
		opt.delay = 0
	}
}

func newLoggerService(
	core Core, fields Fields, opts ...ServiceOption,
) *loggerService {
	opt := &serviceOption{
		level: InfoLevel,
		delay: time.Millisecond * 20,
	}
	for _, f := range opts {
		f(opt)
	}
	return &loggerService{
		level:  opt.level,
		delay:  opt.delay,
		fields: fields,
		core:   core,
	}
}

func (l *Logger) NewService(opts ...ServiceOption) service.Services {
	return newLoggerService(l.core, nil, opts...)
}

func (l *Logger) Add(r *service.Registry) {
	r.AddServices(l.NewService())
}

func (e *Entry) NewService(opts ...ServiceOption) service.Services {
	return newLoggerService(e.l.core, e.fields, opts...)
}

func (e *Entry) Add(r *service.Registry) {
	r.AddServices(e.NewService())
}
