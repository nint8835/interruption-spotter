package cmd

import "log/slog"

func checkErr(err error, msg string) {
	if err != nil {
		slog.Error(msg, "err", err)
	}
}
