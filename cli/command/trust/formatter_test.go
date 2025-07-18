package trust

import (
	"bytes"
	"testing"

	"github.com/docker/cli/cli/command/formatter"
	"github.com/docker/cli/internal/test"
	"gotest.tools/v3/assert"
	is "gotest.tools/v3/assert/cmp"
)

func TestTrustTag(t *testing.T) {
	digest := test.RandomID()
	trustedTag := "tag"

	var ctx trustTagContext

	cases := []struct {
		trustTagCtx trustTagContext
		expValue    string
		call        func() string
	}{
		{
			trustTagContext{
				s: SignedTagInfo{
					Name:    trustedTag,
					Digest:  digest,
					Signers: nil,
				},
			},
			digest,
			ctx.Digest,
		},
		{
			trustTagContext{
				s: SignedTagInfo{
					Name:    trustedTag,
					Digest:  digest,
					Signers: nil,
				},
			},
			trustedTag,
			ctx.SignedTag,
		},
		// Empty signers makes a row with empty string
		{
			trustTagContext{
				s: SignedTagInfo{
					Name:    trustedTag,
					Digest:  digest,
					Signers: nil,
				},
			},
			"",
			ctx.Signers,
		},
		{
			trustTagContext{
				s: SignedTagInfo{
					Name:    trustedTag,
					Digest:  digest,
					Signers: []string{"alice", "bob", "claire"},
				},
			},
			"alice, bob, claire",
			ctx.Signers,
		},
		// alphabetic signing on Signers
		{
			trustTagContext{
				s: SignedTagInfo{
					Name:    trustedTag,
					Digest:  digest,
					Signers: []string{"claire", "bob", "alice"},
				},
			},
			"alice, bob, claire",
			ctx.Signers,
		},
	}

	for _, c := range cases {
		ctx = c.trustTagCtx
		v := c.call()
		if v != c.expValue {
			t.Fatalf("Expected %s, was %s\n", c.expValue, v)
		}
	}
}

func TestTrustTagContextWrite(t *testing.T) {
	cases := []struct {
		context  formatter.Context
		expected string
	}{
		// Errors
		{
			formatter.Context{
				Format: "{{InvalidFunction}}",
			},
			`template parsing error: template: :1: function "InvalidFunction" not defined`,
		},
		{
			formatter.Context{
				Format: "{{nil}}",
			},
			`template parsing error: template: :1:2: executing "" at <nil>: nil is not a command`,
		},
		// Table Format
		{
			formatter.Context{
				Format: NewTrustTagFormat(),
			},
			`SIGNED TAG   DIGEST     SIGNERS
tag1         deadbeef   alice
tag2         aaaaaaaa   alice, bob
tag3         bbbbbbbb   
`,
		},
	}

	signedTags := []SignedTagInfo{
		{Name: "tag1", Digest: "deadbeef", Signers: []string{"alice"}},
		{Name: "tag2", Digest: "aaaaaaaa", Signers: []string{"alice", "bob"}},
		{Name: "tag3", Digest: "bbbbbbbb", Signers: []string{}},
	}

	for _, tc := range cases {
		t.Run(string(tc.context.Format), func(t *testing.T) {
			var out bytes.Buffer
			tc.context.Output = &out

			if err := TagWrite(tc.context, signedTags); err != nil {
				assert.Error(t, err, tc.expected)
			} else {
				assert.Equal(t, out.String(), tc.expected)
			}
		})
	}
}

// With no trust data, the TagWrite will print an empty table:
// it's up to the caller to decide whether or not to print this versus an error
func TestTrustTagContextEmptyWrite(t *testing.T) {
	emptyCase := struct {
		context  formatter.Context
		expected string
	}{
		formatter.Context{
			Format: NewTrustTagFormat(),
		},
		`SIGNED TAG   DIGEST    SIGNERS
`,
	}

	emptySignedTags := []SignedTagInfo{}
	out := bytes.NewBufferString("")
	emptyCase.context.Output = out
	err := TagWrite(emptyCase.context, emptySignedTags)
	assert.NilError(t, err)
	assert.Check(t, is.Equal(emptyCase.expected, out.String()))
}

func TestSignerInfoContextEmptyWrite(t *testing.T) {
	emptyCase := struct {
		context  formatter.Context
		expected string
	}{
		formatter.Context{
			Format: NewSignerInfoFormat(),
		},
		`SIGNER    KEYS
`,
	}
	emptySignerInfo := []SignerInfo{}
	out := bytes.NewBufferString("")
	emptyCase.context.Output = out
	err := SignerInfoWrite(emptyCase.context, emptySignerInfo)
	assert.NilError(t, err)
	assert.Check(t, is.Equal(emptyCase.expected, out.String()))
}

func TestSignerInfoContextWrite(t *testing.T) {
	cases := []struct {
		context  formatter.Context
		expected string
	}{
		// Errors
		{
			formatter.Context{
				Format: "{{InvalidFunction}}",
			},
			`template parsing error: template: :1: function "InvalidFunction" not defined`,
		},
		{
			formatter.Context{
				Format: "{{nil}}",
			},
			`template parsing error: template: :1:2: executing "" at <nil>: nil is not a command`,
		},
		// Table Format
		{
			formatter.Context{
				Format: NewSignerInfoFormat(),
				Trunc:  true,
			},
			`SIGNER    KEYS
alice     key11, key12
bob       key21
eve       foobarbazqux, key31, key32
`,
		},
		// No truncation
		{
			formatter.Context{
				Format: NewSignerInfoFormat(),
			},
			`SIGNER    KEYS
alice     key11, key12
bob       key21
eve       foobarbazquxquux, key31, key32
`,
		},
	}

	signerInfo := []SignerInfo{
		{Name: "alice", Keys: []string{"key11", "key12"}},
		{Name: "bob", Keys: []string{"key21"}},
		{Name: "eve", Keys: []string{"key31", "key32", "foobarbazquxquux"}},
	}
	for _, tc := range cases {
		t.Run(string(tc.context.Format), func(t *testing.T) {
			var out bytes.Buffer
			tc.context.Output = &out

			if err := SignerInfoWrite(tc.context, signerInfo); err != nil {
				assert.Error(t, err, tc.expected)
			} else {
				assert.Equal(t, out.String(), tc.expected)
			}
		})
	}
}
