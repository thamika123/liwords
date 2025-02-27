package config

import (
	"context"
	"errors"
	"net/http"

	"github.com/domino14/liwords/pkg/stores/common"
	macondoconfig "github.com/domino14/macondo/config"
	"github.com/lithammer/shortuuid"
	"github.com/namsral/flag"
)

type ArgonConfig struct {
	Time    int
	Memory  int
	Threads int
	Keylen  int
}

type Config struct {
	MacondoConfig macondoconfig.Config
	ArgonConfig   ArgonConfig

	DBHost           string
	DBPort           string
	DBUser           string
	DBPassword       string
	DBSSLMode        string
	DBName           string
	DBMigrationsPath string
	DBConnUri        string
	DBConnDSN        string

	ListenAddr   string
	SecretKey    string
	NatsURL      string
	MailgunKey   string
	RedisURL     string
	DiscordToken string
	// Puzzles
	PuzzleGenerationSecretKey      string
	ECSClusterName                 string
	PuzzleGenerationTaskDefinition string
}

type CtxKey string

const CtxKeyword CtxKey = CtxKey("config")

// Load loads the configs from the given arguments
func (c *Config) Load(args []string) error {
	fs := flag.NewFlagSet("macondo", flag.ContinueOnError)

	fs.BoolVar(&c.MacondoConfig.Debug, "debug", false, "debug logging on")

	fs.StringVar(&c.MacondoConfig.LetterDistributionPath, "letter-distribution-path", "../macondo/data/letterdistributions", "directory holding letter distribution files")
	fs.StringVar(&c.MacondoConfig.StrategyParamsPath, "strategy-params-path", "../macondo/data/strategy", "directory holding strategy files")
	fs.StringVar(&c.MacondoConfig.LexiconPath, "lexicon-path", "../macondo/data/lexica", "directory holding lexicon files")
	fs.StringVar(&c.MacondoConfig.DefaultLexicon, "default-lexicon", "NWL20", "the default lexicon to use")
	fs.StringVar(&c.MacondoConfig.DefaultLetterDistribution, "default-letter-distribution", "English", "the default letter distribution to use. English, EnglishSuper, Spanish, Polish, etc.")
	fs.StringVar(&c.DBHost, "db-host", "", "the database host")
	fs.StringVar(&c.DBPort, "db-port", "", "the database port")
	fs.StringVar(&c.DBUser, "db-user", "", "the database user")
	fs.StringVar(&c.DBPassword, "db-password", "", "the database password")
	fs.StringVar(&c.DBSSLMode, "db-ssl-mode", "", "the database SSL mode")
	fs.StringVar(&c.DBName, "db-name", "", "the database name")
	fs.StringVar(&c.ListenAddr, "listen-addr", ":8001", "listen on this address")
	fs.StringVar(&c.SecretKey, "secret-key", "", "secret key must be a random unguessable string")
	fs.StringVar(&c.NatsURL, "nats-url", "nats://localhost:4222", "the NATS server URL")
	fs.StringVar(&c.MailgunKey, "mailgun-key", "", "the Mailgun secret key")
	fs.StringVar(&c.RedisURL, "redis-url", "", "the Redis URL")
	fs.StringVar(&c.DiscordToken, "discord-token", "", "the token used for moderator action discord notifications")
	fs.StringVar(&c.DBMigrationsPath, "db-migrations-path", "", "the path where migrations are stored")
	fs.StringVar(&c.PuzzleGenerationSecretKey, "puzzle-generation-secret-key", shortuuid.New(), "a secret key used for generating puzzles")
	fs.StringVar(&c.ECSClusterName, "ecs-cluster-name", "", "the ECS cluster this runs on")
	fs.StringVar(&c.PuzzleGenerationTaskDefinition, "puzzle-generation-task-definition", "", "the task definition for the puzzle generation ECS task")
	// For password hashing:
	fs.IntVar(&c.ArgonConfig.Keylen, "argon-key-len", 32, "the Argon key length")
	fs.IntVar(&c.ArgonConfig.Time, "argon-time", 1, "the Argon time")
	fs.IntVar(&c.ArgonConfig.Memory, "argon-memory", 64*1024, "the Argon memory (KB)")
	fs.IntVar(&c.ArgonConfig.Threads, "argon-threads", 4, "the Argon threads")
	err := fs.Parse(args)
	// build the DB conn string from the passed-in DB arguments
	c.DBConnUri = common.PostgresConnUri(c.DBHost, c.DBPort, c.DBName, c.DBUser, c.DBPassword, c.DBSSLMode)
	c.DBConnDSN = common.PostgresConnDSN(c.DBHost, c.DBPort, c.DBName, c.DBUser, c.DBPassword, c.DBSSLMode)
	return err
}

// Get the Macondo config from the context
func GetMacondoConfig(ctx context.Context) (*macondoconfig.Config, error) {
	ctxConfig, ok := ctx.Value(CtxKeyword).(*Config)
	if !ok {
		return nil, errors.New("config is not ok")
	}
	if ctxConfig == nil {
		return nil, errors.New("config is nil")
	}
	return &ctxConfig.MacondoConfig, nil
}

func CtxMiddlewareGenerator(config *Config) (mw func(http.Handler) http.Handler) {
	mw = func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := context.WithValue(r.Context(), CtxKeyword, config)
			r = r.WithContext(ctx)
			h.ServeHTTP(w, r)
		})
	}
	return
}
