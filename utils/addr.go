package utils

import (
	"fmt"
	"strings"
)

func NodeHostPrefix(nodeName, nodeID string) string {
	return fmt.Sprintf("%s%s", nodeName, nodeID)
}

func NodeHost(nodeName, nodeID string, domainSuffix string) string {
	suffix := strings.Trim(domainSuffix, ".")
	return fmt.Sprintf("%s.%s", NodeHostPrefix(nodeName, nodeID), suffix)
}

func WorkerHostPrefix(workerName string) string {
	return workerName
}

func WorkerHost(workerName, domainSuffix string) string {
	suffix := strings.Trim(domainSuffix, ".")
	return fmt.Sprintf("%s.%s", WorkerHostPrefix(workerName), suffix)
}
