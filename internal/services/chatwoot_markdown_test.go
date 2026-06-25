package services

import "testing"

func TestNormalizeChatwootMarkdownLinksForCodechat(t *testing.T) {
	const mapsURL = "https://maps.app.goo.gl/UYZvebsBMEzk7cqk8"

	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "plain text stays unchanged",
			input: "Ola, tudo bem?",
			want:  "Ola, tudo bem?",
		},
		{
			name:  "plain URL stays unchanged",
			input: mapsURL,
			want:  mapsURL,
		},
		{
			name:  "non-link bold stays unchanged",
			input: "**Atencao**",
			want:  "**Atencao**",
		},
		{
			name:  "whatsapp native bold URL stays unchanged",
			input: "*" + mapsURL + "*",
			want:  "*" + mapsURL + "*",
		},
		{
			name:  "bold autolink becomes URL",
			input: "**<" + mapsURL + ">**",
			want:  mapsURL,
		},
		{
			name:  "mixed bold autolink becomes URL",
			input: "<**" + mapsURL + ">**",
			want:  mapsURL,
		},
		{
			name:  "bold autolink in sentence becomes URL",
			input: "Segue o endereco: **<" + mapsURL + ">** ",
			want:  "Segue o endereco: " + mapsURL,
		},
		{
			name:  "autolink becomes URL",
			input: "<" + mapsURL + ">",
			want:  mapsURL,
		},
		{
			name:  "bold URL becomes URL",
			input: "**" + mapsURL + "**",
			want:  mapsURL,
		},
		{
			name:  "dangling bold marker after URL is removed",
			input: mapsURL + "**",
			want:  mapsURL,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			got := normalizeChatwootMarkdownLinksForCodechat(tc.input)
			if got != tc.want {
				t.Fatalf("normalized text mismatch\nwant: %q\n got: %q", tc.want, got)
			}
		})
	}
}
