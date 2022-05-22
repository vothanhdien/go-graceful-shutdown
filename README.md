# Golang Graceful Shutdown

This project I'm using http server and cron jobs for demonstrate how to archive graceful shutdown.

Graceful shutdown achieved when:
* [ ] System waiting for Http finished process last request before exit
* [ ] System waiting for cron jobs process finished before exit

To run project exec `make run`
