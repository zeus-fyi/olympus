package client

import (
	"net/url"
	"testing"

	"github.com/google/go-querystring/query"
	"github.com/gorilla/schema"
	"github.com/stretchr/testify/suite"
)

type ClientTestSuite struct {
	suite.Suite
}

type Data struct {
	Indexes []string `url:"indexes,omitempty"`
}

type GData struct {
	Indexes []string `schema:"indexes,omitempty"`
}

func (c *ClientTestSuite) TestQueryParams() {
	expectedStr := "indexes=1&indexes=2"
	d := Data{
		Indexes: []string{"1", "2"},
	}
	values, err := query.Values(d)
	c.Assert().Nil(err)
	encodedString := values.Encode()
	c.Assert().Equal(expectedStr, encodedString)

	g := GData{
		Indexes: []string{"1", "2"},
	}
	var encoder = schema.NewEncoder()
	params := url.Values{}

	err = encoder.Encode(g, params)
	c.Assert().Nil(err)
	s := params.Encode()
	c.Assert().Equal(expectedStr, s)

}

func TestClientTestSuite(t *testing.T) {
	suite.Run(t, new(ClientTestSuite))
}
