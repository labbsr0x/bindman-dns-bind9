package nsupdate

import (
	"os"
	"reflect"
	"testing"
)

const basePath = "./data"

func TestMain(m *testing.M) {
	exitCode := m.Run()
	os.Remove(basePath)
	os.Exit(exitCode)
}

func TestBuilder_New(t *testing.T) {
	type fields struct {
		Server   string
		Port     string
		KeyFile  string
		BasePath string
		Zone     string
		Debug    bool
	}
	type args struct {
		basePath string
	}
	tests := []struct {
		name       string
		fields     fields
		args       args
		wantResult *NSUpdate
		wantErr    bool
	}{
		{
			name: "success required fields",
			fields: fields{
				Server:  "server",
				KeyFile: "K<anything>.+157+<anything>.key",
				Zone:    "zone.com",
			},
			args: args{basePath},
			wantResult: &NSUpdate{Builder{
				Server:   "server",
				KeyFile:  "K<anything>.+157+<anything>.key",
				Zone:     "zone.com",
				BasePath: basePath,
			}},
			wantErr: false,
		},
		{
			name: "override basePath with function parameter basePath",
			fields: fields{
				Server:   "server",
				KeyFile:  "K<anything>.+157+<anything>.key",
				Zone:     "zone.com",
				BasePath: "base path to be override",
			},
			args: args{basePath},
			wantResult: &NSUpdate{Builder{
				Server:   "server",
				KeyFile:  "K<anything>.+157+<anything>.key",
				Zone:     "zone.com",
				BasePath: basePath,
			}},
			wantErr: false,
		},
		{
			name: "missing server",
			fields: fields{
				KeyFile: "K<anything>.+157+<anything>.key",
				Zone:    "zone.com",
			},
			args:    args{basePath},
			wantErr: true,
		},
		{
			name: "empty spaces server",
			fields: fields{
				Server:  " ",
				KeyFile: "K<anything>.+157+<anything>.key",
				Zone:    "zone.com",
			},
			args:    args{basePath},
			wantErr: true,
		},
		{
			name: "missing zone",
			fields: fields{
				Server:  "server",
				KeyFile: "K<anything>.+157+<anything>.key",
			},
			args:    args{basePath},
			wantErr: true,
		},
		{
			name: "empty spaces zone",
			fields: fields{
				Server:  "server",
				KeyFile: "K<anything>.+157+<anything>.key",
				Zone:    " ",
			},
			args:    args{basePath},
			wantErr: true,
		},
		{
			name: "missing keyFile",
			fields: fields{
				Server: "server",
				Zone:   "zone",
			},
			args:    args{basePath},
			wantErr: true,
		},
		{
			name: "invalid keyFile",
			fields: fields{
				Server:  "server",
				KeyFile: "file.key",
				Zone:    "zone",
			},
			args:    args{basePath},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &Builder{
				Server:   tt.fields.Server,
				Port:     tt.fields.Port,
				KeyFile:  tt.fields.KeyFile,
				BasePath: tt.fields.BasePath,
				Zone:     tt.fields.Zone,
				Debug:    tt.fields.Debug,
			}
			gotResult, err := b.New(tt.args.basePath)
			if (err != nil) != tt.wantErr {
				t.Errorf("Builder.New() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotResult, tt.wantResult) {
				t.Errorf("Builder.New() = %v, want %v", gotResult, tt.wantResult)
			}
		})
	}
}
