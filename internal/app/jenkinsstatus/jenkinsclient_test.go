package jenkinsstatus

import "testing"

func TestJenkinsGetStatus(t *testing.T) {
	jenkins, err := NewJenkinsClient("http://127.0.0.1:8080", "jenkins", "jenkins.jenkins")
	if err != nil {
		t.Fatal(err)
	}
	sts, err := jenkins.GetStatus("webapp-develop-config-tests", "webapp-develop")
	if err != nil {
		t.Fatal(err)
	}
	if sts != "success" {
		t.Errorf("expected %s, got %s", "success", sts)
	}
}
