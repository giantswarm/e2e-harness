package release

type Version string

func NewStableVersion() Version {
	return Version(":stable")
}

func NewVersion(gitSHA string) Version {
	return Version("@1.0.0-" + gitSHA)
}

func (v Version) String() string {
	return string(v)
}
