package browser

import (
	"context"
	"os"
	"time"

	"github.com/chromedp/chromedp"
	"github.com/sirupsen/logrus"
)

// Controller manages browser automation
type Controller struct {
	ctx    context.Context
	cancel context.CancelFunc
	log    *logrus.Logger
	debug  bool
}

// New creates a new browser controller
func New(headless bool, userDataDir string, viewportWidth, viewportHeight int, debug bool, log *logrus.Logger) (*Controller, error) {
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", headless),
		chromedp.Flag("disable-gpu", true),
		chromedp.Flag("no-sandbox", true),
		chromedp.Flag("disable-dev-shm-usage", true),
		chromedp.WindowSize(viewportWidth, viewportHeight),
	)

	if userDataDir != "" {
		opts = append(opts, chromedp.UserDataDir(userDataDir))
	}

	allocCtx, _ := chromedp.NewExecAllocator(context.Background(), opts...)
	
	var ctx context.Context
	if debug {
		ctx, _ = chromedp.NewContext(allocCtx, chromedp.WithDebugf(log.Debugf))
	} else {
		ctx, _ = chromedp.NewContext(allocCtx)
	}

	return &Controller{
		ctx:    ctx,
		cancel: nil,
		log:    log,
		debug:  debug,
	}, nil
}

// Start initializes the browser
func (c *Controller) Start() error {
	c.log.Info("Starting browser...")
	return chromedp.Run(c.ctx)
}

// Stop closes the browser
func (c *Controller) Stop() {
	c.log.Info("Stopping browser...")
	if c.cancel != nil {
		c.cancel()
	}
}

// Navigate navigates to a URL
func (c *Controller) Navigate(url string) error {
	c.log.WithField("url", url).Debug("Navigating...")
	return chromedp.Run(c.ctx,
		chromedp.Navigate(url),
	)
}

// WaitVisible waits for an element to be visible
func (c *Controller) WaitVisible(selector string, timeout time.Duration) error {
	ctx, cancel := context.WithTimeout(c.ctx, timeout)
	defer cancel()
	
	return chromedp.Run(ctx,
		chromedp.WaitVisible(selector, chromedp.ByQuery),
	)
}

// WaitReady waits for the page to be ready
func (c *Controller) WaitReady() error {
	return chromedp.Run(c.ctx,
		chromedp.WaitReady("body"),
	)
}

// Click clicks an element by selector
func (c *Controller) Click(selector string) error {
	c.log.WithField("selector", selector).Debug("Clicking...")
	return chromedp.Run(c.ctx,
		chromedp.Click(selector, chromedp.ByQuery),
	)
}

// Type types text into an element
func (c *Controller) Type(selector, text string) error {
	c.log.WithField("selector", selector).Debug("Typing...")
	return chromedp.Run(c.ctx,
		chromedp.SendKeys(selector, text, chromedp.ByQuery),
	)
}

// PressKey presses a keyboard key
func (c *Controller) PressKey(key string) error {
	c.log.WithField("key", key).Debug("Pressing key...")
	return chromedp.Run(c.ctx,
		chromedp.KeyEvent(key),
	)
}

// Screenshot takes a screenshot and saves to file
func (c *Controller) Screenshot(path string) error {
	c.log.WithField("path", path).Debug("Taking screenshot...")
	var buf []byte
	if err := chromedp.Run(c.ctx,
		chromedp.CaptureScreenshot(&buf),
	); err != nil {
		return err
	}
	
	return saveScreenshot(path, buf)
}

// ClickAt clicks at specific viewport coordinates
func (c *Controller) ClickAt(x, y float64) error {
	c.log.WithFields(logrus.Fields{
		"x": x,
		"y": y,
	}).Debug("Clicking at coordinates...")
	
	return chromedp.Run(c.ctx,
		chromedp.MouseClickXY(x, y),
	)
}

// Eval evaluates JavaScript and returns the result
func (c *Controller) Eval(script string, res interface{}) error {
	return chromedp.Run(c.ctx,
		chromedp.Evaluate(script, res),
	)
}

// EvalWithTimeout evaluates JavaScript with timeout
func (c *Controller) EvalWithTimeout(script string, res interface{}, timeout time.Duration) error {
	ctx, cancel := context.WithTimeout(c.ctx, timeout)
	defer cancel()
	
	return chromedp.Run(ctx,
		chromedp.Evaluate(script, res),
	)
}

// WaitForCondition waits for a JavaScript condition to be true
func (c *Controller) WaitForCondition(script string, timeout time.Duration) error {
	ctx, cancel := context.WithTimeout(c.ctx, timeout)
	defer cancel()
	
	return chromedp.Run(ctx,
		chromedp.Poll(script, nil, chromedp.WithPollingTimeout(timeout)),
	)
}

// GetContext returns the browser context
func (c *Controller) GetContext() context.Context {
	return c.ctx
}

func saveScreenshot(path string, buf []byte) error {
	return os.WriteFile(path, buf, 0644)
}
