package loadbalancer

import "testing"

func TestNewLoadBalancer(t *testing.T) {
	ld := New(ROUNDED_ROBIN, []string{
		"http://localhost:5001",
		"http://localhost:5002",
		"http://localhost:5003",
	})

	chooseUrlTests := []string{
		"http://localhost:5001",
		"http://localhost:5002",
		"http://localhost:5003",
	}
	for _, u := range chooseUrlTests {
		s := ld.chooseServer()
		if s.Url != u {
			t.Errorf("expected url is %q. got %q\n", u, s.Url)
		}
	}
}

// TODO:
func TestTrafficCount(t *testing.T) {
	t.Skip("test is yet to be implemented")
}
