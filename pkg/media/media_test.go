package media

import "testing"

func TestMedia_CleanURL(t *testing.T) {
	type fields struct {
		URL      string
		Filename string
		Parent   []string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "Should return the URL without query params",
			fields: fields{
				URL: "https://thumbs46.redgifs.com/some-video.mp4?expires=1690588800&signature=v2:some-signature&for=some-ip&hash=some-hash",
			},
			want: "https://thumbs46.redgifs.com/some-video.mp4",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := Media{
				URL:      tt.fields.URL,
				Filename: tt.fields.Filename,
				Parent:   tt.fields.Parent,
			}
			if got := m.CleanURL(); got != tt.want {
				t.Errorf("CleanURL() = %v, want %v", got, tt.want)
			}
		})
	}
}
