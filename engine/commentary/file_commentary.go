package commentary

import (
	_ "embed"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/esmshub/esms-go/engine/internal/rng"
	"github.com/esmshub/esms-go/engine/utils"
	// "github.com/jameshowe/esms-go/internal/rng"
	// "github.com/jameshowe/esms-go/pkg/config"
	// "github.com/jameshowe/esms-go/pkg/log"
	// "github.com/jameshowe/esms-go/pkg/util"
)

// // go:embed data/language.dat
var commentaryData []byte

const DEFAULT_COMMENTARY_CONFIG_FILENAME = "language.dat"

type CommentaryEventMap map[string][]string

type FileCommentaryProvider struct {
	tokenReplacements map[string]string
	eventMap          CommentaryEventMap
}

func (p *FileCommentaryProvider) Load(filePath string) error {
	p.tokenReplacements = map[string]string{}
	p.eventMap = CommentaryEventMap{}
	fmt.Println("Reading commentary...")

	_, err := utils.ReadFile(filePath, func(line string, row int) error {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" || strings.HasPrefix(trimmed, "|") {
			return nil
		}

		evt, evtErr := utils.Substring(trimmed, "[", "]")
		if evtErr != nil {
			fmt.Println(fmt.Errorf("unable to parse commentary event: %+v", evtErr))
			return nil
		}
		comm, commsErr := utils.Substring(trimmed, "{", "}")
		if commsErr != nil {
			fmt.Println(trimmed)
			fmt.Println(fmt.Errorf("unable to parse commentary event: %+v", commsErr))
			return nil
		}
		p.eventMap[evt] = append(p.eventMap[evt], comm)
		return nil
	})
	return err
}

func (p *FileCommentaryProvider) GetEventText(event string, args ...any) string {
	evts := p.eventMap[event]
	text := evts[0]
	if len(evts) > 0 {
		text = evts[rng.RandomRange(0, len(evts))]
	}
	str := fmt.Sprintf(text, args...)
	for token, value := range p.tokenReplacements {
		str = strings.ReplaceAll(str, token, value)
	}
	return strings.ReplaceAll(str, "\\n", "\n")
}

func (p *FileCommentaryProvider) AddTokenReplacement(token, value string) {
	p.tokenReplacements[token] = value
}

// func GetCommentaryFilePath() string {
// 	return filepath.Join(GetConfigDirectory(), DEFAULT_COMMENTARY_CONFIG_FILENAME)
// }

func WriteCommentaryConfig(targetDir string) {
	outputFile := filepath.Join(targetDir, DEFAULT_COMMENTARY_CONFIG_FILENAME)

	// log.Trace().Msg("Attempting to write Commentary file...")
	fmt.Println("Attempting to write Commentary file...")
	err := os.WriteFile(outputFile, commentaryData, 0644)
	if err != nil {
		fmt.Println(fmt.Errorf("unable to save Commentary file: %+v", err))
		return
	}

	fmt.Printf("Commentary file saved to %s", outputFile)
}

func NewFileCommentaryProvider() *FileCommentaryProvider {
	return &FileCommentaryProvider{
		tokenReplacements: make(map[string]string),
	}
}
