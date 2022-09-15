package parser

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseSql(t *testing.T) {
	sql := `CREATE TABLE t_person_info (
  age INT(11) unsigned NULL,
  id BIGINT(11) PRIMARY KEY AUTO_INCREMENT NOT NULL COMMENT '这是id',
  name VARCHAR(30) NOT NULL DEFAULT 'default_name' COMMENT '这是名字',
  created_at datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  sex VARCHAR(2) NULL,
  num INT(11) DEFAULT 3 NULL,
  comment TEXT
  ) COMMENT="person info";`
	codes, err := ParseSQL(sql, WithTablePrefix("t_"), WithJSONTag(0))
	assert.Nil(t, err)
	for k, v := range codes {
		t.Log(k, v)
	}
}

var testData = [][]string{
	{
		"CREATE TABLE information (age INT(11) NULL);",
		"Age int `gorm:\"column:age\"`", "",
	},
	{
		"CREATE TABLE information (age BIGINT(11) NULL COMMENT 'is age');",
		"Age int64 `gorm:\"column:age\"` // is age", "",
	},
	{
		"CREATE TABLE information (id BIGINT(11) PRIMARY KEY AUTO_INCREMENT);",
		"ID int64 `gorm:\"column:id;primary_key;AUTO_INCREMENT\"`", "",
	},
	{
		"CREATE TABLE information (user_ip varchar(20));",
		"UserIP string `gorm:\"column:user_ip\"`", "",
	},
	{
		"CREATE TABLE information (created_at datetime NOT NULL DEFAULT CURRENT_TIMESTAMP);",
		"CreatedAt time.Time `gorm:\"column:created_at;default:CURRENT_TIMESTAMP;NOT NULL\"`", "time",
	},
	{
		"CREATE TABLE information (num INT(11) DEFAULT 3 NULL);",
		"Num int `gorm:\"column:num;default:3\"`", "",
	},
	{
		"CREATE TABLE information (num double(5,6) DEFAULT 31.50 NULL);",
		"Num float64 `gorm:\"column:num;default:31.50\"`", "",
	},
	{
		"CREATE TABLE information (comment TEXT);",
		"Comment string `gorm:\"column:comment\"`", "",
	},
	{
		"CREATE TABLE information (comment TINYTEXT);",
		"Comment string `gorm:\"column:comment\"`", "",
	},
	{
		"CREATE TABLE information (comment LONGTEXT);",
		"Comment string `gorm:\"column:comment\"`", "",
	},
}

func TestParseSqls(t *testing.T) {
	for i, test := range testData {
		msg := fmt.Sprintf("sql-%d", i)
		codes, err := ParseSQL(test[0], WithNoNullType())
		if !assert.NoError(t, err, msg) {
			continue
		}
		for k, v := range codes {
			t.Log(i+1, k, v)
		}
	}
}

func Test_toCamel(t *testing.T) {
	str := "user_example"
	t.Log(toCamel(str))
}
