package jenkinsstatus

import "github.com/bndr/gojenkins"

// JenkinsClient represents Jenkins client high level interface.
type JenkinsClient struct {
	jenkins *gojenkins.Jenkins
}

// NewJenkinsClient returns initialized JenkinsClient object.
func NewJenkinsClient(url, username, password string) (*JenkinsClient, error) {
	jenkins := gojenkins.CreateJenkins(nil, url, username, password)
	_, err := jenkins.Init()
	if err != nil {
		return nil, err
	}
	return &JenkinsClient{
		jenkins: jenkins,
	}, nil
}

// GetStatus returns build status for the specified Jenkins job.
// Returned status can be ABORTED, FAILURE, NOT_BUILT, SUCCESS, UNSTABLE.
func (c *JenkinsClient) GetStatus(id string, parentIDs ...string) (string, error) {
	job, err := c.jenkins.GetJob(id, parentIDs...)
	if err != nil {
		return "", err
	}
	build, err := job.GetLastBuild()
	if err != nil {
		return "", err
	}
	return build.GetResult(), nil
}
