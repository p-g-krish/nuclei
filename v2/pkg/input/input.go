package input

import (
	"net"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/projectdiscovery/hmap/store/hybrid"
	templateTypes "github.com/projectdiscovery/nuclei/v2/pkg/templates/types"
)

// Helper is a structure for helping with input transformation
type Helper struct {
	InputsHTTP *hybrid.HybridMap
}

// NewHelper returns a new inpt helper instance
func NewHelper() *Helper {
	helper := &Helper{}
	return helper
}

// Close closes the resources associated with input helper
func (h *Helper) Close() error {
	var err error
	if h.InputsHTTP != nil {
		err = h.InputsHTTP.Close()
	}
	return err
}

// Transform transforms an input based on protocol type and returns
// appropriate input based on it.
func (h *Helper) Transform(input string, protocol templateTypes.ProtocolType) string {
	switch protocol {
	case templateTypes.DNSProtocol, templateTypes.WHOISProtocol:
		return h.convertInputToType(input, typeHostOnly, "")
	case templateTypes.FileProtocol, templateTypes.OfflineHTTPProtocol:
		return h.convertInputToType(input, typeFilepath, "")
	case templateTypes.HTTPProtocol, templateTypes.HeadlessProtocol:
		return h.convertInputToType(input, typeURL, "")
	case templateTypes.NetworkProtocol:
		return h.convertInputToType(input, typeHostWithOptionalPort, "")
	case templateTypes.SSLProtocol:
		return h.convertInputToType(input, typeHostWithPort, "443")
	case templateTypes.WebsocketProtocol:
		return h.convertInputToType(input, typeWebsocket, "")
	}
	return input
}

type inputType int

const (
	typeHostOnly inputType = iota + 1
	typeHostWithPort
	typeHostWithOptionalPort
	typeURL
	typeFilepath
	typeWebsocket
)

// convertInputToType converts an input based on an inputType.
// Various formats are supported for inputs and their transformation
func (h *Helper) convertInputToType(input string, inputType inputType, defaultPort string) string {
	notURL := !strings.Contains(input, "://")
	parsed, _ := url.Parse(input)
	var host, port string
	if !notURL {
		host, port, _ = net.SplitHostPort(parsed.Host)
	} else {
		host, port, _ = net.SplitHostPort(input)
	}
	hasPort := port != ""

	if inputType == typeFilepath {
		// if it has ports most likely it's not a file
		if hasPort {
			return ""
		}
		if filepath.IsAbs(input) {
			return input
		}
		if absPath, _ := filepath.Abs(input); absPath != "" && fileOrFolderExists(absPath) {
			return input
		}
		if _, err := filepath.Match(input, ""); err != filepath.ErrBadPattern && notURL {
			return input
		}
	} else if inputType == typeHostOnly {
		if host != "" {
			return host
		}
		if !notURL {
			return parsed.Hostname()
		} else {
			return input
		}
	} else if inputType == typeURL {
		if parsed != nil && (parsed.Scheme == "http" || parsed.Scheme == "https") {
			return input
		}
		if h.InputsHTTP != nil {
			if probed, ok := h.InputsHTTP.Get(input); ok {
				return string(probed)
			}
		}
	} else if inputType == typeHostWithPort {
		if host != "" && port != "" {
			return net.JoinHostPort(host, port)
		}
		if parsed != nil && port == "" && parsed.Scheme == "https" {
			return net.JoinHostPort(parsed.Host, "443")
		}
		if defaultPort != "" {
			return net.JoinHostPort(input, defaultPort)
		}
	} else if inputType == typeHostWithOptionalPort {
		if host != "" && port != "" {
			return net.JoinHostPort(host, port)
		}
		if parsed != nil && port == "" && parsed.Scheme == "https" {
			return net.JoinHostPort(parsed.Host, "443")
		}
		if defaultPort != "" {
			return net.JoinHostPort(input, defaultPort)
		}
		return input
	} else if inputType == typeWebsocket {
		if parsed != nil && (parsed.Scheme == "ws" || parsed.Scheme == "wss") {
			return input
		}
	}
	return ""
}

func fileOrFolderExists(filename string) bool {
	_, err := os.Stat(filename)
	return !os.IsNotExist(err)
}
