package whoami

const unknown = "UNKNOWN"

var (
	name       = unknown
	version    = unknown
	number     = unknown
	buildAt    = unknown
	commitHash = unknown
	gitBranch  = unknown
)

func Name() string {
	return name
}

func Version() string {
	return version
}

func Number() string {
	return number
}

func BuildAt() string {
	return buildAt
}

func CommitHash() string {
	return commitHash
}

func GitBranch() string {
	return gitBranch
}
