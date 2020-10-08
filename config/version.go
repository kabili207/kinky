// Code generated by tools/version.go - DO NOT EDIT.

package config

import "time"

const (
	BuildTimestamp     = int64(1602193067)
	BuildDate          = "2020-10-08T23:37:47+02:00"
	Version            = "0.0.1"
	GitDirty           = true
	GitSha             = "366bcc8b3b8da88910a305d48e7e9a966c50398a"
	GitCommitTimestamp = int64(1602002424)
	GitCommitDate      = "2020-10-06T18:40:24+02:00"
	GitBranch          = "master"
)

func BuildTime() time.Time {
	return time.Unix(BuildTimestamp, 0)
}

func GitCommitTime() time.Time {
	return time.Unix(GitCommitTimestamp, 0)
}

func GitShaShort() string {
	return GitSha[:6]
}
