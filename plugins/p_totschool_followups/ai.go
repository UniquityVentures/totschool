package p_totschool_followups

import (
	"context"
	"log"
	"sync"
	"time"

	"google.golang.org/genai"
	"gorm.io/gorm"
)

type generationTask struct {
	FollowupID   uint
	Content      string
	SystemPrompt string
}

var (
	workCh    = make(chan generationTask, 64)
	cancelMu  sync.Mutex
	cancelled = map[uint]bool{}
)

func Generate(db *gorm.DB, followupID uint, content, systemPrompt string) {
	one := 1
	db.Model(&Followup{}).Where("id = ?", followupID).Updates(map[string]any{
		"generation_id":    &one,
		"generated_letter": "",
	})
	workCh <- generationTask{
		FollowupID:   followupID,
		Content:      content,
		SystemPrompt: systemPrompt,
	}
}

func CancelGeneration(db *gorm.DB, followupID uint) {
	cancelMu.Lock()
	cancelled[followupID] = true
	cancelMu.Unlock()
	db.Model(&Followup{}).Where("id = ?", followupID).Update("generation_id", nil)
}

func isCancelled(followupID uint) bool {
	cancelMu.Lock()
	defer cancelMu.Unlock()
	if cancelled[followupID] {
		delete(cancelled, followupID)
		return true
	}
	return false
}

func runWorker(db *gorm.DB) {
	clientConfig := &genai.ClientConfig{}
	if followupAIConfig.APIKey != "" {
		clientConfig.APIKey = followupAIConfig.APIKey
	}
	model := "gemini-2.5-flash"
	if followupAIConfig.Model != "" {
		model = followupAIConfig.Model
	}

	client, err := genai.NewClient(context.Background(), clientConfig)
	if err != nil {
		log.Printf("[followups] genai client not available: %v", err)
		for task := range workCh {
			db.Model(&Followup{}).Where("id = ?", task.FollowupID).Update("generation_id", nil)
		}
		return
	}

	for task := range workCh {
		if isCancelled(task.FollowupID) {
			continue
		}

		ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
		config := &genai.GenerateContentConfig{}
		if task.SystemPrompt != "" {
			config.SystemInstruction = genai.NewContentFromText(task.SystemPrompt, "user")
		}
		resp, err := client.Models.GenerateContent(ctx, model, genai.Text(task.Content), config)
		cancel()

		if isCancelled(task.FollowupID) {
			continue
		}

		if err != nil {
			log.Printf("[followups] generation failed for followup %d: %v", task.FollowupID, err)
			db.Model(&Followup{}).Where("id = ?", task.FollowupID).Update("generation_id", nil)
			continue
		}

		respText := resp.Text()
		if respText == "" {
			log.Printf("[followups] empty response for followup %d", task.FollowupID)
			db.Model(&Followup{}).Where("id = ?", task.FollowupID).Update("generation_id", nil)
			continue
		}

		db.Model(&Followup{}).Where("id = ?", task.FollowupID).Updates(map[string]any{
			"generated_letter": respText,
			"generation_id":    nil,
		})
	}
}

const followupLetterSystemPrompt = `You are an expert financial advisory letter writer. Generate an initial financial advisory follow-up letter using the client's proposal questionnaire answers.

Output requirements:
1. Output markdown only, with clear headings and readable paragraphs.
2. Do not wrap the output in a code block.
3. Personalize the letter using the client profile, advisor name, city, and proposal answers.
4. Use Indian financial planning context.
5. Maintain a professional, warm, advisory tone.
6. Always validate the client's current stance positively while still identifying improvement areas.
7. If a numeric field is missing or unclear, say that it should be confirmed instead of inventing a number.

Template structure:

1. Header and greeting
- Title: "Wealth Creation Strategy For [User_Name] and Family: Initial Advisory Letter"
- Date, To, City
- Subject: "Commencing Your Wealth Building Journey: Strategy and Next Steps"
- Dear [User_Name],

2. Executive summary and cash flow analysis
- Use Monthly_Income and Monthly_Expenditure from the questionnaire.
- Calculate Monthly_Surplus = Monthly_Income - Monthly_Expenditure when both values are available.
- Calculate Savings_Rate = (Monthly_Surplus / Monthly_Income) * 100 when possible.
- Present the savings rate positively and explain that the surplus can be optimized for faster goal achievement.

3. Phase 1: Foundation and Protection (Risk Mitigation)
- Explain that risk protection comes before aggressive investing.
- Emergency Fund:
  - Prefer Base_Survival_Expenditure if provided; otherwise use Monthly_Expenditure.
  - Emergency_Fund_Target = selected monthly expense x 6.
  - Timeline_Months = Emergency_Fund_Target / Monthly_Surplus when possible.
- Family Risk Protection:
  - ELV_Target = Monthly_Income x 200 when possible.
  - Recommend life cover around ELV_Target.
  - Recommend standalone family health insurance separate from employer cover.

4. Phase 2: Wealth Acceleration and Growth Strategy
- Use Children_Goals and Retirement_Goals from the questionnaire.
- For long horizons, recommend diversified growth assets and systematic investment discipline.
- If Savings_Rate is under 30%, set Target_Savings_Rate_Percentage to 30% and calculate Surplus_Shortfall.
- If heavy aspirational goals are present, use firmer language about widening the surplus gap.

5. Actionable roadmap and execution steps
- Immediate: secure life insurance and health insurance.
- Emergency fund: dedicate surplus until the target is reached.
- Long-term: deploy recurring investments after the emergency fund and protection base are secure.

6. Getting started: the insurance strategy and Why LIC
- Explain the first insurance step and how it protects the spouse and dependents.
- Mention spouse name if provided.
- Explain Why LIC:
  - Advantages: Government of India backed sovereign guarantee, brand trust, claim-settlement reliability.
  - Disadvantages: bureaucratic onboarding, documentation scrutiny, moral/medical underwriting, slower timelines.
- State that patience during approval is worthwhile because long-term reliability outweighs delays.

7. Sign-off
- End with a brief invitation to discuss the plan and next steps.
- Sign with advisor name and "Family Wealth Educator".

Core philosophy:
- Balance Risk Isolation (Phase 1) with Goal Funding (Phase 2).
- If the user has low savings and high liabilities, emphasize debt consolidation and minimal emergency benchmarks before growth assets.
- If aspirational capital goals are high, explicitly connect the goals to the need for stronger surplus creation.
`

const followupLetterEditorSystemPrompt = `You are an expert financial advisory letter editor. Edit or rewrite the given follow-up letter according to the user's instructions.

Rules:
1. Only output the edited letter content.
2. Preserve factual information unless the user specifically asks to modify it.
3. Preserve the professional advisory tone unless instructed otherwise.
4. Keep the output valid markdown.
5. Do not wrap your response in a code block.`
