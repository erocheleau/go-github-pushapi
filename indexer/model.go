package indexer

type Organization struct {
	Login            string `json:"login,omitempty"`
	ID               int    `json:"id,omitempty"`
	AvatarURL        string `json:"avatar_url,omitempty"`
	URL              string `json:"html_url,omitempty"`
	Name             string `json:"name,omitempty"`
	Company          string `json:"company,omitempty"`
	Blog             string `json:"blog,omitempty"`
	Location         string `json:"location,omitempty"`
	Email            string `json:"email,omitempty"`
	Description      string `json:"description,omitempty"`
	PublicReposCount int    `json:"public_repos,omitempty"`
	PublicGistsCount int    `json:"public_gists,omitempty"`
	FollowersCount   int    `json:"followers,omitempty"`
	FollowingCount   int    `json:"following,omitempty"`
	CreatedAt        string `json:"created_at,omitempty"`
	UpdatedAt        string `json:"updated_at,omitempty"`
	Type             string `json:"type,omitempty"`
	FileType         string `json:"filetype"`
	Data             string `json:"data"`
}

type Repository struct {
	ID             float64 `json:"id"`
	Name           string  `json:"name"`
	FullName       string  `json:"full_name"`
	Description    string  `json:"description"`
	Data           string  `json:"data"`
	Private        bool    `json:"private"`
	URL            string  `json:"html_url"`
	CreatedDate    string  `json:"created_at"`
	UpdatedAt      string  `json:"updated_at"`
	PushedAt       string  `json:"pushed_at"`
	GitURL         string  `json:"git_url"`
	Homepage       string  `json:"homepage,omitempty"`
	Language       string  `json:"language"`
	HasIssues      bool    `json:"has_issues"`
	ForksCount     float64 `json:"forks_count,omitempty"`
	IssuesCount    float64 `json:"open_issues_count,omitempty"`
	Owner          Owner   `json:"owner"`
	OwnerName      string  `json:"ownerName"`
	OwnerID        float64 `json:"ownerID"`
	OwnerAvatarURL string  `json:"ownerAvatarURL"`
	OwnerType      string  `json:"ownerType"`
	FileType       string  `json:"filetype"`
}

type Owner struct {
	Login     string  `json:"login"`
	ID        float64 `json:"id"`
	AvatarURL string  `json:"avatar_url"`
	Type      string  `json:"type"`
}
