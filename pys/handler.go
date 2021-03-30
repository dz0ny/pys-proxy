package pys

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

type UserData struct {
	ClientIPAddress string `json:"client_ip_address"`
	ClientUserAgent string `json:"client_user_agent"`
}
type Event struct {
	EventName      string            `json:"event_name"`
	EventTime      int64             `json:"event_time"`
	EventSourceURL string            `json:"event_source_url"`
	EventID        string            `json:"event_id"`
	UserData       UserData          `json:"user_data"`
	CustomData     map[string]string `json:"custom_data"`
	ActionSource   string            `json:"action_source"`
}

func Handler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err == nil {

		ev := Event{
			ActionSource:   "website",
			EventSourceURL: r.FormValue("data[event_url]"),
			EventName:      r.FormValue("event"),
			EventID:        r.FormValue("eventID"),
			EventTime:      time.Now().Unix(),
			UserData: UserData{
				ClientIPAddress: r.Header.Get("X-Real-IP"),
				ClientUserAgent: r.Header.Get("User-Agent"),
			},
			CustomData: map[string]string{
				"domain":     r.Header.Get("Host"),
				"user_roles": r.FormValue("data[user_role]"),
				"user_role":  r.FormValue("data[user_role]"),
				"plugin":     r.FormValue("data[plugin]"),
				"page_title": r.FormValue("data[page_title]"),
				"post_type":  r.FormValue("data[post_type]"),
				"post_id":    r.FormValue("data[post_id]"),
				"event_url":  r.FormValue("data[event_url]"),
			},
		}
		evenJson, _ := json.Marshal(ev)
		data, err := PostToFB(r.FormValue("ids[]"), r.Header.Get("Access-Token"), evenJson, http.DefaultClient)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(w, string(data))
		} else {
			w.WriteHeader(http.StatusAccepted)
			fmt.Fprint(w, string(evenJson))
		}

	} else {
		fmt.Fprint(w, "Failed!")
	}

}

func PostToFB(id, access_token string, event []byte, c *http.Client) ([]byte, error) {
	payload := url.Values{}

	payload.Set("data", "["+string(event)+"]")
	payload.Set("access_token", access_token)

	res, err := c.PostForm("https://graph.facebook.com/v10.0/"+id+"/events", payload)
	if err != nil {
		panic(err)
	}

	return ioutil.ReadAll(res.Body)
}
