package image

import (
	"bytes"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/docker/cli/cli/command/formatter"
	"github.com/docker/cli/internal/test"
	"github.com/moby/moby/api/types/image"
	"gotest.tools/v3/assert"
)

type historyCase struct {
	historyCtx historyContext
	expValue   string
	call       func() string
}

func TestHistoryContext_ID(t *testing.T) {
	id := test.RandomID()

	var ctx historyContext
	cases := []historyCase{
		{
			historyContext{
				h:     image.HistoryResponseItem{ID: id},
				trunc: false,
			}, id, ctx.ID,
		},
		{
			historyContext{
				h:     image.HistoryResponseItem{ID: id},
				trunc: true,
			}, formatter.TruncateID(id), ctx.ID,
		},
	}

	for _, c := range cases {
		ctx = c.historyCtx
		v := c.call()
		if strings.Contains(v, ",") {
			test.CompareMultipleValues(t, v, c.expValue)
		} else if v != c.expValue {
			t.Fatalf("Expected %s, was %s\n", c.expValue, v)
		}
	}
}

func TestHistoryContext_CreatedSince(t *testing.T) {
	longerAgo := time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC)
	dateStr := longerAgo.Local().Format(time.RFC3339)
	var ctx historyContext
	cases := []historyCase{
		{
			historyContext{
				h:     image.HistoryResponseItem{Created: time.Now().AddDate(0, 0, -7).Unix()},
				trunc: false,
				human: true,
			}, "7 days ago", ctx.CreatedSince,
		},
		{
			historyContext{
				h:     image.HistoryResponseItem{Created: longerAgo.Unix()},
				trunc: false,
				human: false,
			}, dateStr, ctx.CreatedSince,
		},
		{
			// The zero time is not displayed.
			historyContext{
				h:     image.HistoryResponseItem{Created: 0},
				trunc: false,
				human: true,
			}, "N/A", ctx.CreatedSince,
		},
		{
			// A time before the year 2000 is not displayed.
			historyContext{
				h:     image.HistoryResponseItem{Created: time.Date(1980, time.November, 10, 10, 23, 0, 0, time.UTC).Unix()},
				trunc: false,
				human: true,
			}, "N/A", ctx.CreatedSince,
		},
		{
			// A time after 2000 is displayed.
			historyContext{
				h:     image.HistoryResponseItem{Created: time.Now().AddDate(-11, 0, 0).Unix()},
				trunc: false,
				human: true,
			}, "11 years ago", ctx.CreatedSince,
		},
	}

	for _, c := range cases {
		ctx = c.historyCtx
		v := c.call()
		if strings.Contains(v, ",") {
			test.CompareMultipleValues(t, v, c.expValue)
		} else if v != c.expValue {
			t.Fatalf("Expected %q, was %q\n", c.expValue, v)
		}
	}
}

func TestHistoryContext_CreatedBy(t *testing.T) {
	const withTabs = `/bin/sh -c apt-key adv --keyserver hkp://pgp.mit.edu:80	--recv-keys 573BFD6B3D8FBC641079A6ABABF5BD827BD9BF62	&& echo "deb http://nginx.org/packages/mainline/debian/ jessie nginx" >> /etc/apt/sources.list  && apt-get update  && apt-get install --no-install-recommends --no-install-suggests -y       ca-certificates       nginx=${NGINX_VERSION}       nginx-module-xslt       nginx-module-geoip       nginx-module-image-filter       nginx-module-perl       nginx-module-njs       gettext-base  && rm -rf /var/lib/apt/lists/*` //nolint:revive // ignore line-length-limit
	const expected = `/bin/sh -c apt-key adv --keyserver hkp://pgp.mit.edu:80 --recv-keys 573BFD6B3D8FBC641079A6ABABF5BD827BD9BF62 && echo "deb http://nginx.org/packages/mainline/debian/ jessie nginx" >> /etc/apt/sources.list  && apt-get update  && apt-get install --no-install-recommends --no-install-suggests -y       ca-certificates       nginx=${NGINX_VERSION}       nginx-module-xslt       nginx-module-geoip       nginx-module-image-filter       nginx-module-perl       nginx-module-njs       gettext-base  && rm -rf /var/lib/apt/lists/*` //nolint:revive // ignore line-length-limit

	var ctx historyContext
	cases := []historyCase{
		{
			historyContext{
				h:     image.HistoryResponseItem{CreatedBy: withTabs},
				trunc: false,
			}, expected, ctx.CreatedBy,
		},
		{
			historyContext{
				h:     image.HistoryResponseItem{CreatedBy: withTabs},
				trunc: true,
			}, formatter.Ellipsis(expected, 45), ctx.CreatedBy,
		},
	}

	for _, c := range cases {
		ctx = c.historyCtx
		v := c.call()
		if strings.Contains(v, ",") {
			test.CompareMultipleValues(t, v, c.expValue)
		} else if v != c.expValue {
			t.Fatalf("Expected %s, was %s\n", c.expValue, v)
		}
	}
}

func TestHistoryContext_Size(t *testing.T) {
	const size = int64(182964289)
	const expected = "183MB"

	var ctx historyContext
	cases := []historyCase{
		{
			historyContext{
				h:     image.HistoryResponseItem{Size: size},
				trunc: false,
				human: true,
			}, expected, ctx.Size,
		}, {
			historyContext{
				h:     image.HistoryResponseItem{Size: size},
				trunc: false,
				human: false,
			}, strconv.Itoa(182964289), ctx.Size,
		},
	}

	for _, c := range cases {
		ctx = c.historyCtx
		v := c.call()
		if strings.Contains(v, ",") {
			test.CompareMultipleValues(t, v, c.expValue)
		} else if v != c.expValue {
			t.Fatalf("Expected %s, was %s\n", c.expValue, v)
		}
	}
}

func TestHistoryContext_Comment(t *testing.T) {
	const comment = "Some comment"

	var ctx historyContext
	cases := []historyCase{
		{
			historyContext{
				h:     image.HistoryResponseItem{Comment: comment},
				trunc: false,
			}, comment, ctx.Comment,
		},
	}

	for _, c := range cases {
		ctx = c.historyCtx
		v := c.call()
		if strings.Contains(v, ",") {
			test.CompareMultipleValues(t, v, c.expValue)
		} else if v != c.expValue {
			t.Fatalf("Expected %s, was %s\n", c.expValue, v)
		}
	}
}

func TestHistoryContext_Table(t *testing.T) {
	out := bytes.NewBufferString("")
	unixTime := time.Now().AddDate(0, 0, -1).Unix()
	oldDate := time.Now().AddDate(-17, 0, 0).Unix()
	histories := []image.HistoryResponseItem{
		{
			ID:        "imageID1",
			Created:   unixTime,
			CreatedBy: "/bin/bash ls && npm i && npm run test && karma -c karma.conf.js start && npm start && more commands here && the list goes on",
			Size:      int64(182964289),
			Comment:   "Hi",
			Tags:      []string{"image:tag2"},
		},
		{ID: "imageID2", Created: unixTime, CreatedBy: "/bin/bash echo", Size: int64(182964289), Comment: "Hi", Tags: []string{"image:tag2"}},
		{ID: "imageID3", Created: unixTime, CreatedBy: "/bin/bash ls", Size: int64(182964289), Comment: "Hi", Tags: []string{"image:tag2"}},
		{ID: "imageID4", Created: unixTime, CreatedBy: "/bin/bash grep", Size: int64(182964289), Comment: "Hi", Tags: []string{"image:tag2"}},
		{ID: "imageID5", Created: 0, CreatedBy: "/bin/bash echo", Size: int64(182964289), Comment: "Hi", Tags: []string{"image:tag2"}},
		{ID: "imageID6", Created: oldDate, CreatedBy: "/bin/bash echo", Size: int64(182964289), Comment: "Hi", Tags: []string{"image:tag2"}},
	}

	//nolint:dupword // ignore "Duplicate words (CREATED) found"
	const expectedNoTrunc = `IMAGE      CREATED        CREATED BY                                                                                                                     SIZE      COMMENT
imageID1   24 hours ago   /bin/bash ls && npm i && npm run test && karma -c karma.conf.js start && npm start && more commands here && the list goes on   183MB     Hi
imageID2   24 hours ago   /bin/bash echo                                                                                                                 183MB     Hi
imageID3   24 hours ago   /bin/bash ls                                                                                                                   183MB     Hi
imageID4   24 hours ago   /bin/bash grep                                                                                                                 183MB     Hi
imageID5   N/A            /bin/bash echo                                                                                                                 183MB     Hi
imageID6   17 years ago   /bin/bash echo                                                                                                                 183MB     Hi
`
	//nolint:dupword // ignore "Duplicate words (CREATED) found"
	const expectedTrunc = `IMAGE      CREATED        CREATED BY                                      SIZE      COMMENT
imageID1   24 hours ago   /bin/bash ls && npm i && npm run test && kar…   183MB     Hi
imageID2   24 hours ago   /bin/bash echo                                  183MB     Hi
imageID3   24 hours ago   /bin/bash ls                                    183MB     Hi
imageID4   24 hours ago   /bin/bash grep                                  183MB     Hi
imageID5   N/A            /bin/bash echo                                  183MB     Hi
imageID6   17 years ago   /bin/bash echo                                  183MB     Hi
`

	cases := []struct {
		context  formatter.Context
		expected string
	}{
		{
			formatter.Context{
				Format: NewHistoryFormat("table", false, true),
				Trunc:  true,
				Output: out,
			},
			expectedTrunc,
		},
		{
			formatter.Context{
				Format: NewHistoryFormat("table", false, true),
				Trunc:  false,
				Output: out,
			},
			expectedNoTrunc,
		},
	}

	for _, tc := range cases {
		t.Run(string(tc.context.Format), func(t *testing.T) {
			err := HistoryWrite(tc.context, true, histories)
			assert.NilError(t, err)
			assert.Equal(t, out.String(), tc.expected)
			// Clean buffer
			out.Reset()
		})
	}
}
