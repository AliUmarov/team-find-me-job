package gigachat

import (
	"encoding/json"
	"fmt"
)

type improveRequest struct {
	Model    string        `json:"model"`
	Messages []chatMessage `json:"messages"`
}

type chatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type improveResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

type resumeAIResult struct {
	Improved string `json:"improved"`
	Score    int    `json:"score"`
}

func ImproveResume(fullText string, client *Client) (string, int, error) {

	prompt := fmt.Sprintf(`
Ты — помощник по улучшению резюме.

1) Улучши текст резюме, сохранив факты, но сделав формулировки более профессиональными и читаемыми.
2) Дай оценку резюме по шкале от 1 до 10, где 10 — идеальное резюме для сильного кандидата.

ОТВЕТ ВЕРНИ СТРОГО В ФОРМАТЕ JSON БЕЗ ОБЪЯСНЕНИЙ И ТЕКСТА ВОКРУГ. ПРИМЕР:
{
  "improved": "улучшенный текст",
  "score": 8
}

Вот текст резюме:

%s
`, fullText)

	req := improveRequest{
		Model: "GigaChat",
		Messages: []chatMessage{
			{
				Role:    "user",
				Content: prompt,
			},
		},
	}

	raw, err := client.send(req)
	if err != nil {
		return "", 0, err
	}

	var apiResp improveResponse
	if err := json.Unmarshal(raw, &apiResp); err != nil {
		return "", 0, fmt.Errorf("не удалось разобрать ответ GigaChat: %w", err)
	}

	if len(apiResp.Choices) == 0 {
		return "", 0, fmt.Errorf("ответ пустой (choices=0)")
	}

	content := apiResp.Choices[0].Message.Content
	if content == "" {
		return "", 0, fmt.Errorf("пустой content в ответе модели")
	}

	var result resumeAIResult
	if err := json.Unmarshal([]byte(content), &result); err != nil {
		return "", 0, fmt.Errorf("не удалось распарсить JSON из content: %w; raw content: %s", err, content)
	}

	if result.Improved == "" {
		return "", 0, fmt.Errorf("модель не вернула поле improved")
	}
	if result.Score <= 0 || result.Score > 100 {
		result.Score = 50
	}

	return result.Improved, result.Score, nil
}
