package config

import "testing"

func TestChannelEmail_Validate(t *testing.T) {
	type fields struct {
		Name         string
		From         string
		To           string
		ServerName   string
		AuthUsername string
		AuthPassword string
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
		errText string
	}{
		{
			name:    "empty name",
			fields:  fields{Name: "", From: "", To: "", ServerName: "", AuthUsername: "", AuthPassword: ""},
			wantErr: true,
			errText: "name must be not empty",
		},
		{
			name:    "empty from",
			fields:  fields{Name: "foo", From: "", To: "", ServerName: "", AuthUsername: "", AuthPassword: ""},
			wantErr: true,
			errText: "from must be not empty",
		},
		{
			name:    "empty to",
			fields:  fields{Name: "foo", From: "gopher@example.net", To: "", ServerName: "", AuthUsername: "", AuthPassword: ""},
			wantErr: true,
			errText: "to must be not empty",
		},
		{
			name:    "empty server_name",
			fields:  fields{Name: "foo", From: "gopher@example.net", To: "foo@example.com", ServerName: "", AuthUsername: "", AuthPassword: ""},
			wantErr: true,
			errText: "server_name must be not empty",
		},
		{
			name:    "ok",
			fields:  fields{Name: "foo", From: "gopher@example.net", To: "foo@example.com", ServerName: "mail.example.com", AuthUsername: "", AuthPassword: ""},
			wantErr: false,
			errText: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := ChannelEmail{
				Name:         tt.fields.Name,
				From:         tt.fields.From,
				To:           tt.fields.To,
				ServerName:   tt.fields.ServerName,
				AuthUsername: tt.fields.AuthUsername,
				AuthPassword: tt.fields.AuthPassword,
			}
			err := cfg.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err != nil && err.Error() != tt.errText {
				t.Errorf("Validate() error = '%s', wantErrText '%s'", err.Error(), tt.errText)
			}
		})
	}
}
