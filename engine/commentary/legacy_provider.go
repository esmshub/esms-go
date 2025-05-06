package commentary

import (
	"errors"
	"fmt"
	"math/rand"
	"strings"
	"sync"

	"github.com/esmshub/esms-go/engine/commentary/formatters"
	"github.com/esmshub/esms-go/engine/events"
	"github.com/esmshub/esms-go/engine/models"
	"github.com/esmshub/esms-go/engine/pkg/utils"
	"go.uber.org/zap"
)

var (
	DefaultCommentaryProviderAliases = map[string]string{
		events.InjuryTimeAddedEventName:             CommInjuryTimeBEvent,
		events.HalfTimeEventName:                    CommHalfTimeEvent,
		events.FullTimeEventName:                    CommFullTimeEvent,
		events.ChanceEventName:                      CommChanceEvent,
		events.ChanceBeatsDefenderEventName:         CommChanceBeatsDefenderEvent,
		events.AssistedChanceEventName:              CommAssistedChanceEvent,
		events.AssistedChanceBeatsDefenderEventName: CommAssistedChanceBeatsDefenderEvent,
		events.ShotOnTargetEventName:                CommShotEvent,
		events.ShotOffTargetEventName:               CommShotOffTargetEvent,
		events.ShotOffTargetDeflectionEventName:     CommShotOffTargetDeflectionEvent,
		events.ShotTackledEventName:                 CommShotTackledEvent,
		events.ShotTackledCornerEventName:           CommShotTackledCornerEvent,
		events.ShotSavedEventName:                   CommShotSavedEvent,
		events.ShotClearedEventName:                 CommShotClearedEvent,
		events.ShotSavedCornerEventName:             CommShotSavedCornerEvent,
		events.GoalScoredEventName:                  CommGoalEvent,
		events.GoalScoredCancelledEventName:         CommGoalCancelledEvent,
		events.OwnGoalScoredEventName:               CommOwnGoalEvent,
		events.CornerCaughtEventName:                CommCornerCaughtEvent,
		events.CornerClearedEventName:               CommCornerClearedEvent,
		events.CornerShotEventName:                  CommCornerShotEvent,
		// models.YellowCardEventName:      CommYellowCardEvent,
		// models.RedCardEventName:         CommRedCardEvent,
		// models.SecondYellowCardEventName: CommSecondYellowCardEvent,
		// models.FirstYellowCardEventName: CommYellowCardEvent,
		// models.ExtraTimeEventName:       CommExtraTimeEvent,
		// models.AttemptsEventName:        CommAttemptsEvent,
		// models.ShotsOnTargetEventName:   CommShotsOnTargetEvent,
		// models.ShotsOffTargetEventName:  CommShotsOffTargetEvent,

	}
)

type LegacyFileCommentaryProvider struct {
	mu         *sync.RWMutex
	eventMap   map[string][]string
	aliases    map[string]string
	commentary []string
}

func (p *LegacyFileCommentaryProvider) Load(filePath string) error {
	p.eventMap = map[string][]string{}
	zap.L().Debug("Loading commentary", zap.String("file", filePath))

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

func (p *LegacyFileCommentaryProvider) getCommentaryKey(event events.Event) string {
	data, ok := event.GetData().(map[string]any)
	if data != nil && !ok {
		zap.L().Warn("event data is not of type map[string]any", zap.Any("event", event))
	}

	eventName := event.GetName()
	switch eventName {
	case events.ShotOnTargetEventName:
		if oneOnOne, ok := data["one_on_one"].(bool); ok && oneOnOne {
			eventName = CommOneOnOneShotEvent
		}
	case events.ShotOffTargetEventName:
		if oneOnOne, ok := data["one_on_one"].(bool); ok && oneOnOne {
			eventName = CommOneOnOneShotOffTargetEvent
		}
	case events.ShotSavedEventName:
		if oneOnOne, ok := data["one_on_one"].(bool); ok && oneOnOne {
			eventName = CommOneOnOneShotSavedEvent
		}
	}

	return eventName
}

func (p *LegacyFileCommentaryProvider) WriteCommentary(event events.Event) error {
	commKey := p.getCommentaryKey(event)

	p.mu.RLock()
	if alias, ok := p.aliases[commKey]; ok {
		commKey = alias
	}
	p.mu.RUnlock()

	text := p.getCommentaryText(commKey)

	matchEvent, ok := event.(*models.MatchEvent)
	if !ok {
		return errors.New("event is not of type models.MatchEvent")
	}

	switch event.GetName() {
	case events.InjuryTimeAddedEventName:
		text = formatters.FormatInjuryTimeEvent(text, matchEvent)
	case events.KickOffEventName:
		text = formatters.FormatKickOffEvent(text, matchEvent)
	case events.HalfTimeEventName:
		comm := formatters.FormatScoreEvent(p.getCommentaryText("COMM_HTSCORE"), matchEvent)
		text = fmt.Sprintf("%s%s\n", formatters.FormatHalfTimeEvent(text, matchEvent), comm)
	case events.FullTimeEventName:
		text = formatters.FormatFullTimeEvent(text, matchEvent)
	case events.ChanceEventName, events.ChanceBeatsDefenderEventName:
		text = formatters.FormatChanceEvent(text, matchEvent)
	case events.AssistedChanceEventName, events.AssistedChanceBeatsDefenderEventName:
		text = formatters.FormatAssistedChanceEvent(text, matchEvent)
	case events.ShotTackledEventName, events.ShotTackledCornerEventName, events.ShotTackledRecoveryEventName, events.ShotClearedEventName:
		text = formatters.FormatShotTackledEvent(text, matchEvent)
	case events.ShotOnTargetEventName:
		text = formatters.FormatShotEvent(text, matchEvent)
	case events.ShotOffTargetEventName, events.ShotOffTargetDeflectionEventName:
		shotCommKey := CommShotEvent
		if data, ok := event.GetData().(map[string]any); ok {
			if val, exists := data["one_on_one"]; exists {
				if oneOnOne, ok := val.(bool); ok && oneOnOne {
					shotCommKey = CommOneOnOneShotEvent
				}
			}
		}
		text = strings.Join(
			[]string{
				formatters.FormatShotEvent(p.getCommentaryText(shotCommKey), matchEvent),
				formatters.FormatMatchEvent(text, matchEvent),
			},
			"",
		)
	case events.OwnGoalScoredEventName:
		text = formatters.FormatOwnGoalScoredEvent(text, matchEvent)
	case events.CornerCaughtEventName:
		comm := formatters.FormatCornerEvent(p.getCommentaryText(CommCornerTakenEvent), matchEvent)
		text = strings.Join(
			[]string{
				p.getCommentaryText(CommCornerEvent),
				comm,
				formatters.FormatCornerCaughtEvent(text, matchEvent),
			},
			"",
		)
	case events.CornerClearedEventName:
		comm := formatters.FormatCornerEvent(p.getCommentaryText(CommCornerTakenEvent), matchEvent)
		text = strings.Join(
			[]string{
				p.getCommentaryText(CommCornerEvent),
				comm,
				formatters.FormatCornerClearedEvent(text, matchEvent),
			},
			"",
		)
	case events.CornerShotEventName:
		comm := formatters.FormatCornerEvent(p.getCommentaryText(CommCornerTakenEvent), matchEvent)
		text = strings.Join(
			[]string{
				p.getCommentaryText(CommCornerEvent),
				comm,
				formatters.FormatCornerShotEvent(text, matchEvent),
			},
			"",
		)
	default:
		text = formatters.FormatMatchEvent(text, matchEvent)
	}

	// Format newlines
	p.write(strings.ReplaceAll(text, "\\n", "\n"))

	return nil
}

func (p *LegacyFileCommentaryProvider) getCommentaryText(key string) string {
	p.mu.RLock()
	defer p.mu.RUnlock()

	options := p.eventMap[key]
	if len(options) > 0 {
		return options[rand.Intn(len(options))]
	} else {
		zap.L().Warn("no commentary found for key", zap.String("key", key))
	}

	return ""
}

func (p *LegacyFileCommentaryProvider) write(commentary string) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.commentary = append(p.commentary, commentary)
}

func (p *LegacyFileCommentaryProvider) GetCommentary() []string {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.commentary
}

func (p *LegacyFileCommentaryProvider) SetAliases(aliases map[string]string) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.aliases = aliases
}

func NewLegacyFileCommentaryProvider() *LegacyFileCommentaryProvider {
	return &LegacyFileCommentaryProvider{
		aliases:    DefaultCommentaryProviderAliases,
		commentary: []string{},
		mu:         &sync.RWMutex{},
	}
}
