package err

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"
)

//CustomWrapF custom error wrap with flag
// example result : [<flag>] <args> : <err>
func CustomWrapF(err error, flag string, args ...string) error {
	var msg string
	var strPkg string
	if strPkg = strings.TrimSpace(flag); strPkg != "" {
		strPkg = fmt.Sprintf("[%s] ", strings.ToUpper(flag))
	}

	msg = strings.Join(args, " | ")
	msg = strPkg + msg

	if err == nil {
		err = errors.New(fmt.Sprintf("[FOR DEVELOPER] forget to set error in %s", flag))
	}

	return errors.Wrap(err, msg)
}

//CustomWrap custom error wrap without flag
// example result : <args> : <err>
func CustomWrap(err error, args ...string) error {
	msg := strings.Join(args, " | ")
	return errors.Wrap(err, msg)
}
