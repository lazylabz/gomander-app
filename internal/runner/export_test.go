package runner

// WaitForCommand exists for the Runner's own tests against real processes:
// nothing in the app awaits a process, it reacts to the ProcessFinished event.
func (c *DefaultRunner) WaitForCommand(commandId string) {
	c.mutex.Lock()
	running, exists := c.runningCommands[commandId]
	c.mutex.Unlock()

	if exists {
		running.wg.Wait()
	}
}
