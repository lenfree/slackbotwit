package awswit

import (
	"fmt"

	witai "github.com/lenfree/wit.ai-go"
)

type entityName struct {
	Name   string
	Entity witai.EntityRes
}

type intent struct {
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

func parse(m witai.Message) (intent, entityName) {

	var intention intent
	var entity entityName
	for key := range m.Outcomes {
		if key == "intent" {
			intention = intent{Name: key,
				Entity: witai.EntityRes{Confidence: m.Outcomes[key][0].Confidence,
					Type:  m.Outcomes[key][0].Type,
					Value: m.Outcomes[key][0].Value,
				},
			}
		} else {
			entity = entityName{Name: key,
				Entity: witai.EntityRes{Confidence: m.Outcomes[key][0].Confidence,
					Type:  m.Outcomes[key][0].Type,
					Value: m.Outcomes[key][0].Value,
				},
			}
		}
	}
	return intention, entity
}
