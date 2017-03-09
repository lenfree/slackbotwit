package awswit

import (
	"fmt"

	witai "github.com/lenfree/wit.ai-go"
)

type awsEntities struct {
	EntityName entityName
	Intent     intent
	Count      count
	EntityType entityType
	Property   property
}

type entityName struct {
	Name   string
	Entity witai.EntityRes
}

type intent struct {
	Name   string
	Entity witai.EntityRes
}

type count struct {
	Name   string
	Entity witai.EntityRes
}

type entityType struct {
	Name   string
	Entity witai.EntityRes
}

type property struct {
	Name   string
	Entity witai.EntityRes
}

func newWit(token string) *witai.Client {
	return witai.NewClient(token)
}

func query(c *witai.Client, txt string) witai.Message {
	c.Verbose = true

	fmt.Printf("Text: %s\n", txt)
	// message
	result, err := c.QueryMessage(txt, nil, "", "")
	if err != nil {
		fmt.Printf("%s\n", err.Error())
	}
	fmt.Printf("query message result: %+v\n", result)
	fmt.Printf("outcome: %+v\n", result.Outcomes)
	return result
}

func parse(m witai.Message) awsEntities {

	var intention intent
	var entity entityName
	var number count
	var enType entityType
	var enProperty property

	for key := range m.Outcomes {
		switch key {
		case "intent":
			{
				intention = intent{Name: key,
					Entity: witai.EntityRes{Confidence: m.Outcomes[key][0].Confidence,
						Type:  m.Outcomes[key][0].Type,
						Value: m.Outcomes[key][0].Value,
					},
				}
			}
		case "entity":
			{
				entity = entityName{Name: key,
					Entity: witai.EntityRes{Confidence: m.Outcomes[key][0].Confidence,
						Type:  m.Outcomes[key][0].Type,
						Value: m.Outcomes[key][0].Value,
					},
				}
			}
		case "count":
			{
				number = count{Name: key,
					Entity: witai.EntityRes{Confidence: m.Outcomes[key][0].Confidence,
						Type:  m.Outcomes[key][0].Type,
						Value: m.Outcomes[key][0].Value,
					},
				}
			}
		case "type":
			{
				enType = entityType{Name: key,
					Entity: witai.EntityRes{Confidence: m.Outcomes[key][0].Confidence,
						Type:  m.Outcomes[key][0].Type,
						Value: m.Outcomes[key][0].Value,
					},
				}
			}
		case "property":
			{
				enProperty = property{Name: key,
					Entity: witai.EntityRes{Confidence: m.Outcomes[key][0].Confidence,
						Type:  m.Outcomes[key][0].Type,
						Value: m.Outcomes[key][0].Value,
					},
				}
			}
		}
	}
	return awsEntities{
		Intent:     intention,
		EntityName: entity,
		Count:      number,
		EntityType: enType,
		Property:   enProperty,
	}
}
