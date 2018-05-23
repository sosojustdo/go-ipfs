package commands

import "testing"

func TestParsePeerParam(t *testing.T) {
	// [arg, maddr, id, err]
	tests := [][]string{
		{"QmdfTbBqBPQ7VNxZEYEj14VmRuZBkqFbiwReogJgS1zR1n", "", "QmdfTbBqBPQ7VNxZEYEj14VmRuZBkqFbiwReogJgS1zR1n", ""},
		{"/ipfs/QmdfTbBqBPQ7VNxZEYEj14VmRuZBkqFbiwReogJgS1zR1n", "", "QmdfTbBqBPQ7VNxZEYEj14VmRuZBkqFbiwReogJgS1zR1n", ""},
		{"/ip4/127.0.0.1/tcp/222/ipfs/QmdfTbBqBPQ7VNxZEYEj14VmRuZBkqFbiwReogJgS1zR1n", "/ip4/127.0.0.1/tcp/222", "QmdfTbBqBPQ7VNxZEYEj14VmRuZBkqFbiwReogJgS1zR1n", ""},
		{"/ip4/127.0.0.1/tcp/222/QmdfTbBqBPQ7VNxZEYEj14VmRuZBkqFbiwReogJgS1zR1n", "/ip4/127.0.0.1/tcp/222", "QmdfTbBqBPQ7VNxZEYEj14VmRuZBkqFbiwReogJgS1zR1n", "expected peer address to contain /ipfs/"},
	}

	for i, test := range tests {
		ma, id, err := ParsePeerParam(test[0])
		if err != nil {
			if test[3] == err.Error() {
				continue
			}
			t.Fatalf("test %d unexpected error: '%s'", i, err)
		}
		if test[3] != "" {
			t.Errorf("test %d expected error", i)
		}
		if ma == nil {
			if test[1] != "" {
				t.Errorf("test %d got nil maddr", i)
			}
		} else {
			if ma.String() != test[1] {
				t.Errorf("test %d got unexpected maddr '%s'", i, ma.String())
			}
		}
		if id.Pretty() != test[2] {
			t.Errorf("test %d got unexpected peerid '%s'", i, id)
		}
	}
}
