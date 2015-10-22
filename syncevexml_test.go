package syncevexml

import (
  "testing"
  "io/ioutil"

  "github.com/jovon/syncevexml/types"
)

func TestXMLToJSON(t *testing.T) {
    data, _ := ioutil.ReadFile("./xml_examples/skillQueue.xml")
    jsonResult := `{"currentTime":"2009-03-18 13:19:43","rows":[{"queuePosition":"1","typeID":"11441","level":"3","startSP":"7072","endSP":"40000","startTime":"2009-03-18 02:01:06","endTime":"2009-03-18 15:19:21"},{"queuePosition":"2","typeID":"20533","level":"4","startSP":"112000","endSP":"633542","startTime":"2009-03-18 15:19:21","endTime":"2009-03-30 03:16:14"}],"cachedUntil":"2009-03-18 13:34:43"}`
    result, _ := XMLToJSON(data, new(types.SkillQueueXML))
    if result != jsonResult {
        t.Errorf("XMLToJSON(%q) == %q, want %q", data, result, jsonResult)
    } 
}