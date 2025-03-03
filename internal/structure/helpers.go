package structure

import "fmt"

func (c *Commands) Register(name string, f func(*State, Command) error) {
	c.Handlers[name] = f
}

func (c *Commands) Run(s *State, cmd Command) error {
	handler, exists := c.Handlers[cmd.Name]
	if !exists {
		return fmt.Errorf("command not found: %s", cmd.Name)
	}
	handler(s, cmd)
	return nil
}
