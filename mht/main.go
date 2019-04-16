package mht

import (
	"encoding/hex"
	"fmt"
	"github.com/satori/go.uuid"
)

type Mht struct {
	Url         string
	ContentType string
	Boundary    string
	Delime      string
}

func UUID() string {
	for {
		uid, err := uuid.NewV4()
		if err == nil {
			return hex.EncodeToString(uid.Bytes())
		}
	}
	return ""
}

func DefaultMht() *Mht {
	return &Mht{
		Delime:   "--",
		Boundary: UUID(),
	}
}

func (m *Mht) Header() string {
	m.Boundary = fmt.Sprintf(`boundary="----MultipartBoundary--%s----"`, UUID())
	return fmt.Sprintf(`Content-Type: multipart/related; type="text/html"; %s`, m.Boundary)
}

/*
self._msg['MIME-Version'] = '1.0'
        self._msg.add_header('Content-Type', 'multipart/related', type='text/html')

MIME-Version: 1.0
Content-Type: multipart/related;
	type="text/html";
	boundary="----MultipartBoundary--F04fDaDlxkjtQjJtRNkNt6CRYacCpnAk6dXtNLgYrf----"
 */
