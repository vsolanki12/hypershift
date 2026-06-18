package util

import (
	"strconv"

	hyperv1 "github.com/openshift/hypershift/api/hypershift/v1beta1"
)

// LogLevelToKlogVerbosity maps a LogLevel enum value to a klog verbosity integer.
// Returns 2 (Normal) for unrecognized or empty values.
func LogLevelToKlogVerbosity(level hyperv1.LogLevel) int {
	switch level {
	case hyperv1.Debug:
		return 4
	case hyperv1.Trace:
		return 6
	case hyperv1.TraceAll:
		return 8
	default:
		return 2
	}
}

// KASVerbosityLevel returns the klog verbosity for the KAS container.
// Precedence: structured API field > deprecated annotation > default value (2)
func KASVerbosityLevel(hcp *hyperv1.HostedControlPlane) int {
	// Structured API field takes priority over deprecated annotation.
	if hcp.Spec.OperatorConfiguration != nil &&
		hcp.Spec.OperatorConfiguration.KubeAPIServer != nil &&
		hcp.Spec.OperatorConfiguration.KubeAPIServer.LogLevel != "" {
		return LogLevelToKlogVerbosity(hcp.Spec.OperatorConfiguration.KubeAPIServer.LogLevel)
	}

	// Deprecated annotation fallback for backward compatibility.
	if ann := hcp.Annotations[hyperv1.KubeAPIServerVerbosityLevelAnnotation]; ann != "" {
		if v, err := strconv.Atoi(ann); err == nil {
			return v
		}
	}

	// Default value for new clusters.
	return 2
}
