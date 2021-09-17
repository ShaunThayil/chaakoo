package cmd

import (
	"errors"
	"fmt"
	"strings"
	"tmuxt"
)

type Config struct {
	SessionName string             `mapstructure:"name"`
	Windows     map[string]*Window `mapstructure:"name"`
}

func (c *Config) Validate() error {
	if c == nil {
		return errors.New("config is nil")
	}
	if len(c.SessionName) == 0 {
		return errors.New("session name is required")
	}
	if len(c.Windows) == 0 {
		return fmt.Errorf("atleast 1 window is required for session - %s", c.SessionName)
	}
	for _, window := range c.Windows {
		if err := window.Validate(); err != nil {
			return err
		}
	}
	return nil
}

func (c *Config) Parse() error {
	for _, window := range c.Windows {
		if err := window.Parse(); err != nil {
			return fmt.Errorf("unable to parse grid for window - %s: %w", window.Name, err)
		}
	}
	return nil
}

type Window struct {
	Name      string `mapstructure:"name"`
	Grid      string `mapstructure:"grid"`
	FirstPane *tmuxt.Pane
	Commands  map[string]*Command `mapstructure:"commands"`
}

func (w *Window) Validate() error {
	if w == nil {
		return errors.New("window is nil")
	}
	if len(w.Name) == 0 {
		return errors.New("window name is required")
	}
	if len(strings.TrimSpace(w.Grid)) == 0 {
		return fmt.Errorf("grid for window, %s, is empty", w.Name)
	}
	return nil
}

func (w *Window) Parse() error {
	if w == nil {
		return errors.New("window is nil")
	}
	grid, err := tmuxt.PrepareGrid(w.Grid)
	if err != nil {
		return err
	}
	pane, err := tmuxt.PrepareGraph(grid)
	if err != nil {
		return err
	}
	w.FirstPane = pane
	return nil
}

type Command struct {
	CommandText      string `mapstructure:"command"`
	WorkingDirectory string `mapstructure:"working_dir"`
}
