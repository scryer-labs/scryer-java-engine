package engine

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"slices"
	"strings"

	"github.com/scryer-labs/scryer-java-engine/internal/worker"
)

const (
	MaxRequestBytes   = 1 << 20
	MaxRequestIDBytes = 256
)

type Request struct {
	ProtocolVersion int    `json:"protocolVersion"`
	RequestID       string `json:"requestId"`
	Stack           string `json:"stack"`
	Operation       string `json:"operation"`
	Scan            struct {
		Facts json.RawMessage `json:"facts"`
	} `json:"scan"`
	raw []byte
}

func ParseRequest(data []byte) (Request, error) {
	if len(data) > MaxRequestBytes {
		return Request{}, errors.New("request exceeds 1 MiB")
	}
	var request Request
	if err := json.Unmarshal(data, &request); err != nil || request.RequestID == "" || len(request.RequestID) > MaxRequestIDBytes {
		return Request{}, errors.New("invalid request JSON or requestId")
	}
	request.raw = data
	return request, nil
}

func (r Request) Execute(ctx context.Context, installation string, out, diagnostics io.Writer) error {
	identity := Describe()
	code, message := r.validate(identity)
	if code != "" {
		return json.NewEncoder(out).Encode(struct {
			ProtocolVersion int               `json:"protocolVersion"`
			RequestID       string            `json:"requestId"`
			Stack           string            `json:"stack"`
			Operation       string            `json:"operation"`
			Report          any               `json:"report"`
			Error           map[string]string `json:"error"`
		}{
			ProtocolVersion: identity.ProtocolVersion,
			RequestID:       r.RequestID,
			Stack:           identity.ID,
			Operation:       r.Operation,
			Error:           map[string]string{"code": code, "message": message},
		})
	}
	return worker.ExecuteAt(ctx, installation, bytes.NewReader(r.raw), out, diagnostics)
}

func (r Request) validate(identity Identity) (code, message string) {
	if strings.TrimSpace(r.Operation) == "" {
		return "JAVA_ANALYZER_ERROR", "Missing/non-string operation"
	}
	if r.ProtocolVersion != identity.ProtocolVersion || r.Stack != identity.ID {
		return "JAVA_ANALYZER_ERROR", "Unsupported protocol version or stack"
	}
	supported := capabilities()
	if !slices.Contains(supported.Operations, r.Operation) {
		return "UNSUPPORTED_OPERATION", "Java engine does not support operation " + r.Operation
	}
	if r.Operation != "scan" || r.Scan.Facts == nil {
		return "", ""
	}
	var facts []*string
	if err := json.Unmarshal(r.Scan.Facts, &facts); err != nil || facts == nil || slices.Contains(facts, nil) {
		return "JAVA_ANALYZER_ERROR", "facts must be an array of strings"
	}
	var unsupported []string
	for _, fact := range facts {
		if !slices.Contains(supported.ScanFacts, *fact) {
			unsupported = append(unsupported, *fact)
		}
	}
	if len(unsupported) > 0 {
		return "UNSUPPORTED_FACT", "Java scan does not support: " + strings.Join(unsupported, ", ")
	}
	return "", ""
}
