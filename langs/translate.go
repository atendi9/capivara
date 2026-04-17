package langs

type PkgName string

const (
	PkgAssert PkgName = "capivara::assert"
	PkgRunner PkgName = "capivara::runner"
)

func Translate(pkgName PkgName) Messages {
	return translations[pkgName]
}
