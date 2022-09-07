package middlewares

import (
	"github.com/jmoiron/sqlx"
	"net/http"
	"sync"
	"tie.prodigy9.co/config"
	"tie.prodigy9.co/controllers/render"
	"tie.prodigy9.co/data"
)

type dataContext struct {
	sync.RWMutex
	db  *sqlx.DB
	cfg *config.Config
}

func newDataContext(cfg *config.Config) *dataContext {
	return &dataContext{cfg: cfg}
}

func (c *dataContext) Get() (*sqlx.DB, error) {
	if db := c.tryGet(); db == nil {
		if err := c.tryInit(); err != nil {
			return nil, err
		} else {
			return c.Get()
		}
	} else {
		return db, nil
	}
}
func (c *dataContext) tryGet() *sqlx.DB {
	c.RLock()
	defer c.RUnlock()

	return c.db
}
func (c *dataContext) tryInit() error {
	c.Lock()
	defer c.Unlock()

	if c.db != nil {
		return nil
	}

	if db, err := data.Connect(c.cfg); err != nil {
		return err
	} else {
		c.db = db
		return nil
	}
}

func AddDataContext(handler http.Handler, cfg *config.Config) http.Handler {
	dc := newDataContext(cfg)
	return http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
		if db, err := dc.Get(); err != nil {
			render.Error(resp, req, 500, err)
		} else {
			handler.ServeHTTP(resp, req.WithContext(
				data.NewContext(req.Context(), db)))
		}
	})
}
