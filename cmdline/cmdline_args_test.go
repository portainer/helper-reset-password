package cmdline

import (
	"flag"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	helper_reset_password "github.com/portainer/helper-reset-password"
)

func TestParseArguments(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		args    []string
		want    Arguments
		wantErr bool
	}{
		{
			name: "default values without flags",
			args: []string{},
			want: Arguments{
				Password:     "",
				PasswordHash: "",
				DataPath:     helper_reset_password.DataStorePath,
			},
			wantErr: false,
		},
		// Data-path flag tests
		{
			name: "custom data-path specified",
			args: []string{"-data-path", "/custom/path"},
			want: Arguments{
				Password:     "",
				PasswordHash: "",
				DataPath:     "/custom/path",
			},
			wantErr: false,
		},
		{
			name: "custom data-path with spaces",
			args: []string{"-data-path", "/custom path/data"},
			want: Arguments{
				Password:     "",
				PasswordHash: "",
				DataPath:     "/custom path/data",
			},
			wantErr: false,
		},
		{
			name: "empty data-path value",
			args: []string{"-data-path", ""},
			want: Arguments{
				Password:     "",
				PasswordHash: "",
				DataPath:     helper_reset_password.DataStorePath,
			},
			wantErr: false,
		},
		// Password flag tests
		{
			name: "password flag specified",
			args: []string{"-password", "mySecurePassword123"},
			want: Arguments{
				Password:     "mySecurePassword123",
				PasswordHash: "",
				DataPath:     helper_reset_password.DataStorePath,
			},
			wantErr: false,
		},
		{
			name: "empty password value",
			args: []string{"-password", ""},
			want: Arguments{
				Password:     "",
				PasswordHash: "",
				DataPath:     helper_reset_password.DataStorePath,
			},
			wantErr: false,
		},
		// Password-hash flag tests
		{
			name: "password-hash flag specified",
			args: []string{"-password-hash", "$2a$10$abcdefghijklmnopqrstuv"},
			want: Arguments{
				Password:     "",
				PasswordHash: "$2a$10$abcdefghijklmnopqrstuv",
				DataPath:     helper_reset_password.DataStorePath,
			},
			wantErr: false,
		},
		{
			name: "empty password-hash value",
			args: []string{"-password-hash", ""},
			want: Arguments{
				Password:     "",
				PasswordHash: "",
				DataPath:     helper_reset_password.DataStorePath,
			},
			wantErr: false,
		},
		// Error cases
		{
			name:    "password and password-hash conflict",
			args:    []string{"-password", "test123", "-password-hash", "hash123"},
			want:    Arguments{},
			wantErr: true,
		},
		{
			name:    "password and password-hash conflict with data-path",
			args:    []string{"-password", "test123", "-password-hash", "hash123", "-data-path", "/custom"},
			want:    Arguments{},
			wantErr: true,
		},
		// Combined flags tests
		{
			name: "password with custom data-path",
			args: []string{"-password", "test123", "-data-path", "/custom/path"},
			want: Arguments{
				Password:     "test123",
				PasswordHash: "",
				DataPath:     "/custom/path",
			},
			wantErr: false,
		},
		{
			name: "password-hash with custom data-path",
			args: []string{"-password-hash", "hash123", "-data-path", "/custom/path"},
			want: Arguments{
				Password:     "",
				PasswordHash: "hash123",
				DataPath:     "/custom/path",
			},
			wantErr: false,
		},
		{
			name: "all valid flags combined",
			args: []string{"-password-hash", "$2a$10$xyz", "-data-path", "/opt/portainer"},
			want: Arguments{
				Password:     "",
				PasswordHash: "$2a$10$xyz",
				DataPath:     "/opt/portainer",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset flags for each test
			flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)

			// Set os.Args to simulate command line arguments
			oldArgs := os.Args
			os.Args = append([]string{"cmd"}, tt.args...)
			defer func() { os.Args = oldArgs }()

			got, err := ParseArguments()

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}
