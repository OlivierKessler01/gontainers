package process

import (
	"os"
	"fmt"
)




func GetNamespaces(pid int) (namespaces []string, err error) {
	err = nil
	namespaces = make([]string, 0)

	utsNS, _ := os.Readlink(fmt.Sprintf("/proc/%d/ns/uts", pid))
	namespaces = append(namespaces, utsNS)
	pidNS, _ := os.Readlink(fmt.Sprintf("/proc/%d/ns/pid", pid))
	namespaces = append(namespaces, pidNS)

	return 
}
