module github.com/wagecloud/wagecloud-server

go 1.23.4

require github.com/libvirt/libvirt-go-xml v7.4.0+incompatible

require (
	github.com/google/uuid v1.6.0
	github.com/jackc/pgx/v5 v5.7.4
	github.com/kdomanski/iso9660 v0.4.0
	github.com/spf13/viper v1.20.1
	go.uber.org/zap v1.27.0
	golang.org/x/crypto v0.37.0
	gopkg.in/natefinch/lumberjack.v2 v2.2.1
	gopkg.in/yaml.v3 v3.0.1
)

require (
	github.com/golang-jwt/jwt/v5 v5.2.2
	libvirt.org/go/libvirt v1.11001.0
)

require (
	github.com/fsnotify/fsnotify v1.9.0 // indirect
	github.com/getsentry/sentry-go v0.31.1
	github.com/go-chi/chi/v5 v5.2.1
	github.com/go-viper/mapstructure/v2 v2.2.1 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20240606120523-5a60cdf6a761 // indirect
	github.com/jackc/puddle/v2 v2.2.2 // indirect
	github.com/pelletier/go-toml/v2 v2.2.4 // indirect
	github.com/sagikazarmark/locafero v0.9.0 // indirect
	github.com/sourcegraph/conc v0.3.0 // indirect
	github.com/spf13/afero v1.14.0 // indirect
	github.com/spf13/cast v1.7.1 // indirect
	github.com/spf13/pflag v1.0.6 // indirect
	github.com/subosito/gotenv v1.6.0 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	golang.org/x/sync v0.13.0 // indirect
	golang.org/x/sys v0.32.0 // indirect
	golang.org/x/text v0.24.0 // indirect
)
