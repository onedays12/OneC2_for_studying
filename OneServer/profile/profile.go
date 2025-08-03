package profile

import (
	"errors"
	"fmt"
	"os"
	"regexp"
)

func NewProfile() *Profile {
	return new(Profile)
}

func (p *Profile) Validate() error {
	if err := p.validateTeamServer(); err != nil {
		return err
	}
	if err := p.validateServerResponse(); err != nil {
		return err
	}
	if err := p.validateZap(); err != nil {
		return err
	}
	return nil
}

func (p *Profile) validateTeamServer() error {
	c := &p.TeamServerConfig
	if c.Port < 1 || c.Port > 65535 {
		return fmt.Errorf("TeamServer.port must be between 1 and 65535 (current %d)", c.Port)
	}
	if !ValidUriString(c.Endpoint) {
		return fmt.Errorf("TeamServer.endpoint must be a valid URI (current %q)", c.Endpoint)
	}
	if c.Password == "" {
		return errors.New("TeamServer.password must be set")
	}
	if err := fileMustExist(c.Cert, "TeamServer.cert"); err != nil {
		return err
	}
	if err := fileMustExist(c.Key, "TeamServer.key"); err != nil {
		return err
	}
	if c.AccessTokenLiveHours < 1 {
		return errors.New("TeamServer.access_token_live_hours must be > 0")
	}
	if c.RefreshTokenLiveHours < 1 {
		return errors.New("TeamServer.refresh_token_live_hours must be > 0")
	}
	return nil
}

func (p *Profile) validateServerResponse() error {
	if p.ServerResponseConfig.Page != "" {
		if err := fileMustExist(p.ServerResponseConfig.Page, "ServerResponse.page"); err != nil {
			return err
		}
	}
	return nil
}

func (p *Profile) validateZap() error {
	c := &p.ZapConfig
	switch c.Level {
	case "debug", "info", "warn", "error", "dpanic", "panic", "fatal":
	default:
		return fmt.Errorf("Zap.level must be one of debug/info/warn/error/dpanic/panic/fatal (current %q)", c.Level)
	}
	if c.MaxSize <= 0 {
		return errors.New("Zap.max_size must be > 0")
	}
	if c.MaxBackups < 0 {
		return errors.New("Zap.max_backups must be >= 0")
	}
	if c.MaxAge < 0 {
		return errors.New("Zap.max_age must be >= 0")
	}
	return nil
}

// ---------- 辅助函数 ----------

func fileMustExist(path, name string) error {
	if path == "" {
		return fmt.Errorf("%s must be set", name)
	}
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("%s: file does not exist", name)
		}
		return fmt.Errorf("%s: %w", name, err)
	}
	return nil
}

func ValidUriString(s string) bool {
	re := regexp.MustCompile(`^/(?:[a-zA-Z0-9-_.]+(?:/[a-zA-Z0-9-_.]+)*)?$`)
	return re.MatchString(s)
}
