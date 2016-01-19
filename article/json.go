package article

import (
//"encoding/json"
)

/*func (n *Article) MarshalJSON() ([]byte, error) {
	strBody := struct {
		ArticleModel
		Body string
	}{
		n.ArticleModel,
		n.BodyStr(),
	}

	bts, err := json.MarshalIndent(&strBody, " ", " ")
	if err != nil {
		return nil, err
	}
	return bts, err
}*/
