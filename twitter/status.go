package twitter

import "strconv"

type Status struct {
	ID       string    `json:"id_str"`
	Entities EntitySet `json:"entities"`
}

type EntitySet struct {
	URLs []URLEntity `json:"urls"`
}

type URLEntity struct {
	ExpandedURL string `json:"expanded_url"`
}

func (s *Status) PrevID() (string, error) {
	id, err := strconv.Atoi(s.ID)

	if err != nil {
		return "", err
	}

	id--

	return strconv.Itoa(id), nil
}
