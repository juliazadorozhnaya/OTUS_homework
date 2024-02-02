//go:build !bench
// +build !bench

package hw10programoptimization

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetDomainStat(t *testing.T) {
	data := `{"Id":1,"Name":"Howard Mendoza","Username":"0Oliver","Email":"aliquid_qui_ea@Browsedrive.gov","Phone":"6-866-899-36-79","Password":"InAQJvsq","Address":"Blackbird Place 25"}
{"Id":2,"Name":"Jesse Vasquez","Username":"qRichardson","Email":"mLynch@broWsecat.com","Phone":"9-373-949-64-00","Password":"SiZLeNSGn","Address":"Fulton Hill 80"}
{"Id":3,"Name":"Clarence Olson","Username":"RachelAdams","Email":"RoseSmith@Browsecat.com","Phone":"988-48-97","Password":"71kuz3gA5w","Address":"Monterey Park 39"}
{"Id":4,"Name":"Gregory Reid","Username":"tButler","Email":"5Moore@Teklist.net","Phone":"520-04-16","Password":"r639qLNu","Address":"Sunfield Park 20"}
{"Id":5,"Name":"Janice Rose","Username":"KeithHart","Email":"nulla@Linktype.com","Phone":"146-91-01","Password":"acSBF5","Address":"Russell Trail 61"}`

	t.Run("find 'com'", func(t *testing.T) {
		result, err := GetDomainStat(bytes.NewBufferString(data), "com")
		require.NoError(t, err)
		require.Equal(t, DomainStat{
			"browsecat.com": 2,
			"linktype.com":  1,
		}, result)
	})

	t.Run("find 'gov'", func(t *testing.T) {
		result, err := GetDomainStat(bytes.NewBufferString(data), "gov")
		require.NoError(t, err)
		require.Equal(t, DomainStat{"browsedrive.gov": 1}, result)
	})

	t.Run("find 'unknown'", func(t *testing.T) {
		result, err := GetDomainStat(bytes.NewBufferString(data), "unknown")
		require.NoError(t, err)
		require.Equal(t, DomainStat{}, result)
	})
}

func TestGetDomainStatEmptyData(t *testing.T) {
	data := `` // Пустые данные

	result, err := GetDomainStat(bytes.NewBufferString(data), "com")
	require.NoError(t, err)
	require.Equal(t, DomainStat{}, result, "Result should be empty for empty data")
}

func TestGetDomainStatCaseSensitivity(t *testing.T) {
	// Данные содержат домены в разном регистре
	data := `{"Id":1,"Name":"Alice Johnson","Username":"alice","Email":"alice@Example.COM","Phone":"123","Password":"pass","Address":"Some Address 1"}
{"Id":2,"Name":"Bob Smith","Username":"bob","Email":"bob@example.com","Phone":"456","Password":"pass2","Address":"Some Address 2"}`

	result, err := GetDomainStat(bytes.NewBufferString(data), "com")
	require.NoError(t, err)
	require.Equal(t, DomainStat{"example.com": 2}, result, "Domain count should be case insensitive")
}

func TestGetDomainStatInvalidJSON(t *testing.T) {
	// Некорректный JSON
	data := `{"Id":1, "Name":"Invalid JSON"`

	_, err := GetDomainStat(bytes.NewBufferString(data), "com")
	require.Error(t, err, "Should return an error for invalid JSON")
}

func TestGetDomainStatNoEmailsMatchingDomain(t *testing.T) {
	// Данные, в которых нет email, соответствующих домену
	data := `{"Id":1,"Name":"Charlie Brown","Username":"charlie","Email":"charlie@notmatching.org","Phone":"789","Password":"pass3","Address":"Some Address 3"}`

	result, err := GetDomainStat(bytes.NewBufferString(data), "com")
	require.NoError(t, err)
	require.Equal(t, DomainStat{}, result, "Result should be empty when no emails match the domain")
}

func TestGetDomainStatMultipleLinesAndDomains(t *testing.T) {
	// Данные с несколькими записями и доменами
	data := `{"Id":1,"Name":"Deborah Clark","Username":"deborah","Email":"deborah@domain1.com","Phone":"101","Password":"pass4","Address":"Address 4"}
{"Id":2,"Name":"Evan Wright","Username":"evan","Email":"evan@domain2.com","Phone":"202","Password":"pass5","Address":"Address 5"}
{"Id":3,"Name":"Fiona Graham","Username":"fiona","Email":"fiona@domain1.com","Phone":"303","Password":"pass6","Address":"Address 6"}`

	result, err := GetDomainStat(bytes.NewBufferString(data), "com")
	require.NoError(t, err)
	require.Equal(t, DomainStat{
		"domain1.com": 2,
		"domain2.com": 1,
	}, result, "Result should correctly count multiple domains")
}
