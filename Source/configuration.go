package main



// get the WakeUp Time for the core poster
func (p *RssWatcherPlugin) getWakeUpTime() (uint, error) {
	//config := p.getConfiguration()

	var wakeUpTime uint = 15

	//var err error
	/*if len(config.Heartbeat) > 0 {
		heartbeatTime, err = strconv.Atoi(config.Heartbeat)
		if err != nil {
			return 15, err
		}
	}*/

	return wakeUpTime, nil
}