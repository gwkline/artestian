package main

import (
	"flag"
	"fmt"
	"log"
	"log/slog"
	"os"
	"time"

	"github.com/gwkline/artestian/pkg/agent"
	"github.com/gwkline/artestian/pkg/config"
	"github.com/gwkline/artestian/pkg/finder"
	"github.com/gwkline/artestian/pkg/generator"
	"github.com/gwkline/artestian/pkg/languages"
	"github.com/gwkline/artestian/pkg/prompt_logger"
	"github.com/gwkline/artestian/types"

	"github.com/joho/godotenv"
)

var (
	dir        = flag.String("dir", "", "Path to project root")
	aiProvider = flag.String("ai", "anthropic", "AI provider to use (currently only anthropic is supported)")
	logLevel   = flag.String("log-level", "info", "Log level (debug, info, warn, error)")
)

func main() {
	slog.Info("starting Artestian - AI-Powered Test Generator")

	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	if err := run(); err != nil {
		slog.Error("application error", "error", err)
		os.Exit(1)
	}
}

func run() error {
	flag.Parse()
	setupLogger()

	if *dir == "" {
		return fmt.Errorf("config file is required. Use -config flag to specify the path")
	}

	cfg, err := loadConfiguration(*dir)
	if err != nil {
		return err
	}

	examples, contextFiles, err := loadTestResources(cfg)
	if err != nil {
		return err
	}

	lang, err := initializeLanguage(cfg)
	if err != nil {
		return err
	}

	agent, err := initializeAIProvider(*aiProvider)
	if err != nil {
		return err
	}

	return generateTests(cfg, lang, examples, contextFiles, agent)
}

func loadConfiguration(dirPath string) (*config.Config, error) {
	slog.Debug("loading configuration", "path", dirPath)
	cfg, err := config.LoadConfig(dirPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load configuration: %w", err)
	}
	return cfg, nil
}

func loadTestResources(cfg *config.Config) ([]types.TestExample, []types.ContextFile, error) {
	slog.Debug("loading test examples from config")
	examples, err := cfg.LoadExamples()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to load test examples: %w", err)
	}

	slog.Debug("loading context files")
	contextFiles, err := cfg.LoadContextFiles()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to load context files: %w", err)
	}
	if len(contextFiles) > 0 {
		slog.Info("loaded context files", "count", len(contextFiles))
	}

	return examples, contextFiles, nil
}

func initializeLanguage(cfg *config.Config) (types.ILanguage, error) {
	slog.Debug("initializing language support", "language", cfg.GetLanguage())
	switch cfg.GetLanguage() {
	case "typescript":
		return languages.NewTypeScriptSupport(), nil
	case "go":
		return languages.NewGoSupport(), nil
	default:
		return nil, fmt.Errorf("unsupported language: %s", cfg.GetLanguage())
	}
}

func initializeAIProvider(provider string) (types.IAgent, error) {
	logger, err := prompt_logger.Init(false)
	if err != nil {
		return nil, fmt.Errorf("failed to create prompt logger: %w", err)
	}

	slog.Debug("initializing AI provider", "provider", provider)
	switch provider {
	case "anthropic":
		return agent.NewAnthropicProvider(logger)
	default:
		return nil, fmt.Errorf("unknown AI provider: %s", provider)
	}
}

func generateTests(cfg *config.Config, lang types.ILanguage, examples []types.TestExample, contextFiles []types.ContextFile, aiClient types.IAgent) error {
	slog.Debug("initializing file finder")
	fileFinder := finder.NewFileFinder(lang)

	slog.Debug("initializing test generator")
	testGen := generator.NewTestGenerator(fileFinder, aiClient, lang, examples, contextFiles)

	rootDir := cfg.GetRootDir()
	excludedDirs := cfg.GetExcludedDirs()

	slog.Debug("generating next test", "rootDir", rootDir)
	if err := testGen.GenerateNextTest(*dir, rootDir, excludedDirs); err != nil {
		slog.Error("failed to generate test", "error", err)
		return fmt.Errorf("error generating test: %w", err)
	}

	return nil
}

func setupLogger() {
	var level slog.Level
	switch *logLevel {
	case "debug":
		level = slog.LevelDebug
	case "info":
		level = slog.LevelInfo
	case "warn":
		level = slog.LevelWarn
	case "error":
		level = slog.LevelError
	default:
		level = slog.LevelInfo
	}

	opts := &slog.HandlerOptions{
		Level: level,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			// Format time as HH:MM:SS
			if a.Key == slog.TimeKey {
				return slog.Attr{
					Key:   slog.TimeKey,
					Value: slog.StringValue(time.Now().Format("15:04:05")),
				}
			}
			return a
		},
	}

	handler := slog.NewTextHandler(os.Stdout, opts)
	logger := slog.New(handler)
	slog.SetDefault(logger)
}
