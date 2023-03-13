package websocket

// "encoding/json"
// "io"
// "time"

type Response struct {
	IsSuccessful bool     `json:"isSuccessful"`
	Message      []string `json:"message,omitempty"`
}

// type DateInterval struct {
// 	StartDate *time.Time `json:"startDate,omitempty"`
// 	EndDate   *time.Time `json:"endDate,omitempty"`
// }

// type Filter struct {
// 	Level   *[]string     `json:"level,omitempty"`
// 	Dt      *DateInterval `json:"dt,omitempty"`
// 	Source  *[]string     `json:"source,omitempty"`
// 	Message *string       `json:"message,omitempty"`
// }

// func (f *Filter) Bind(body io.ReadCloser) error {
// 	return json.NewDecoder(body).Decode(f)
// }
