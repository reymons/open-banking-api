package rest

import (
	"banking/core"
	"errors"
	"fmt"
	"log"
	"net/http/httptest"
	"testing"
)

type memLogger struct {
	dst []byte
}

func (l *memLogger) Write(src []byte) (int, error) {
	l.dst = make([]byte, len(src))
	copy(l.dst, src)
	return len(src), nil
}

func TestSendHttpError(t *testing.T) {
	type testCase struct {
		prfx              string
		err               error
		expectedCode      int
		expectedMsg       string
		expectedLoggerMsg string
	}

	var (
		reqMethod     = "GET"
		reqPath       = "/test"
		reqRemoteAddr = "127.0.0.1"
	)

	cases := []testCase{
		testCase{
			err:          core.ErrResourceNotFound,
			expectedCode: 404,
			expectedMsg:  "\n",
		},
		testCase{
			err:          core.ErrInvalidAccess,
			expectedCode: 403,
			expectedMsg:  "\n",
		},
		testCase{
			err:          core.ErrInvalidCredentials,
			expectedCode: 400,
			expectedMsg:  "Invalid credentials\n",
		},
		testCase{
			prfx:         "test",
			err:          errors.New("custom err"),
			expectedCode: 500,
			expectedMsg:  "\n",
			expectedLoggerMsg: fmt.Sprintf(
				"ERROR: %s - %s %s HTTP/1.1, status: 500, test: custom err\n",
				reqRemoteAddr,
				reqMethod,
				reqPath,
			),
		},
	}

	logger := memLogger{}
	log.SetOutput(&logger)
	log.SetFlags(0)

	for _, c := range cases {
		req := httptest.NewRequest(reqMethod, reqPath, nil)
		req.RemoteAddr = reqRemoteAddr
		w := httptest.NewRecorder()
		sendHttpError(w, req, c.prfx, c.err)

		if w.Code != c.expectedCode {
			t.Errorf("code mismatch: expected = %d, got = %d", c.expectedCode, w.Code)
		}

		msg := string(w.Body.Bytes())
		if msg != c.expectedMsg {
			t.Errorf("message mismatch: expected = %q, got = %q", c.expectedMsg, msg)
		}

		if c.expectedLoggerMsg == "" {
			continue
		}

		msg = string(logger.dst)
		if msg != c.expectedLoggerMsg {
			t.Errorf("logger message mismatch: expected = %q, got = %q", c.expectedLoggerMsg, msg)
		}
	}
}
