package optimizeParser

import (
	"encoding/json"
	"io"
	"parser"
)

func ReadJSONAndSendOverChannel(reader io.Reader, c chan map[string]interface{}) error { //ds

	decoder := json.NewDecoder(reader)
	defer close(c)
	for {
		var mapObject map[string]interface{}
		if err := decoder.Decode(&mapObject); err == io.EOF {
			break
		} else if err != nil {

			return err
		}
		c <- mapObject
	}
	return nil
}

func ProcessjsonUsingConfig(config map[string]interface{}, receiverCh, writerCh chan map[string]interface{}) {

	for data := range receiverCh {
		parser.Filter(data, config)
		writerCh <- data
	}
	defer close(writerCh)
}

func GetDataFromChannelAndWriteJSON(receiver chan map[string]interface{}, writer io.Writer) error {

	enc := json.NewEncoder(writer)
	for data := range receiver {
		if err := enc.Encode(data); err != nil {
			return err
		}
	}
	return nil
}
