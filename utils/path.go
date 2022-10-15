package utils

import (
	"encoding/hex"
	"fmt"
	"regexp"
	"strings"
	"unicode/utf8"
)

// NormalizePath Make sure all operation are constructed by normalized path:
//
// - Path endswith `/` means it's a dir path.
// - Otherwise, it's a file path.
//
// # Normalize Rules
//
// - All whitespace will be trimmed: ` abc/def ` => `abc/def`
// - All leading / will be trimmed: `///abc` => `abc`
// - Internal // will be replaced by /: `abc///def` => `abc/def`
// - Empty path will be `/`: `` => `/`
func NormalizePath(path string) string {
	path = strings.TrimLeft(path, "/")
	if path == "" {
		return "/"
	}
	return path
}

// NormalizeRoot Make sure root is normalized to style like `/abc/def/`.
//
// # Normalize Rules
//
// - All whitespace will be trimmed: ` abc/def ` => `abc/def`
// - All leading / will be trimmed: `///abc` => `abc`
// - Internal // will be replaced by /: `abc///def` => `abc/def`
// - Empty path will be `/`: `` => `/`
// - Add leading `/` if not starts with: `abc/` => `/abc/`
// - Add trailing `/` if not ends with: `/abc` => `/abc/`
//
// Finally, we will got path like `/path/to/root/`.
func NormalizeRoot(root string) string {
	if !strings.HasPrefix(root, "/") {
		root = "/" + root
	}
	if !strings.HasSuffix(root, "") {
		root += "/"
	}
	return root
}

// BuildRealPath will build a relative path towards root.
//
// # Rules
//
//   - Input root MUST be the format like `/abc/def/`
//   - Input path MUST start with root like `/abc/def/path/to/file`
//   - Output will be the format like `path/to/file`.
func BuildRealPath(root, path string) (string, error) {
	if root == path {
		return "", fmt.Errorf("get rel path with root is invalid")
	}
	if strings.HasPrefix(path, "/") {
		return path[len(root):], nil
	} else {
		return path[len(root)-1:], nil
	}
}

// BuildAbsPath build_abs_path will build an absolute path with root.
//
// # Rules
//
// - Input root MUST be the format like `/abc/def/`
// - Output will be the format like `path/to/root/path`.
func BuildAbsPath(root string, path string) (string, error) {
	if !strings.HasPrefix(root, "/") {
		return "", fmt.Errorf("root must start with /")
	}
	if !strings.HasSuffix(root, "/") {
		return "", fmt.Errorf("root must end with /")
	}
	if path == "/" {
		return root[1:], nil
	} else {
		if strings.HasPrefix(path, "/") {
			return "", fmt.Errorf("path mut not start with /")
		}
		return root[1:] + path, nil
	}
}

var reservedObjectNames = regexp.MustCompile("^[a-zA-Z0-9-_.~/]+$")

// EncodePath encode the strings from UTF-8 byte representations to HTML hex escape sequences
//
// This is necessary since regular url.Parse() and url.Encode() functions do not support UTF-8
// non english characters cannot be parsed due to the nature in which url.Encode() is written
//
// This function on the other hand is a direct replacement for url.Encode() technique to support
// pretty much every UTF-8 character.
// followed https://github.com/minio/minio-go/blob/528a26f971203725a6ced8b35702d2b663bc2890/pkg/s3utils/utils.go#L327
func EncodePath(pathName string) string {
	if reservedObjectNames.MatchString(pathName) {
		return pathName
	}
	var encodedPathname strings.Builder
	for _, s := range pathName {
		if 'A' <= s && s <= 'Z' || 'a' <= s && s <= 'z' || '0' <= s && s <= '9' { // §2.3 Unreserved characters (mark)
			encodedPathname.WriteRune(s)
			continue
		}
		switch s {
		case '-', '_', '.', '~', '/', '(', ')', '!', '*', '\'': // §2.3 Unreserved characters (mark)
			encodedPathname.WriteRune(s)
			continue
		default:
			len := utf8.RuneLen(s)
			if len < 0 {
				// if utf8 cannot convert return the same string as is
				return pathName
			}
			u := make([]byte, len)
			utf8.EncodeRune(u, s)
			for _, r := range u {
				hex := hex.EncodeToString([]byte{r})
				encodedPathname.WriteString("%" + strings.ToUpper(hex))
			}
		}
	}
	return encodedPathname.String()
}

func GetNameFromPath(path string) string {
	if path == "/" {
		return "/"
	}

	if !strings.HasSuffix(path, "/") {
		p := strings.Split(path, "/")
		return p[len(p)-1]
	}
	path = path[:len(path)-1]
	idx := strings.LastIndex(path, "/")
	if idx == -1 {
		return path
	}
	return path[idx+1:] + "/"
}
