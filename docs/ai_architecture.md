# AI Email Summary Architecture

## Overview

This document describes the architecture and implementation of the AI-powered email summarization feature for Email Sentinel CLI.

## Features

- **Multi-Provider Support**: Claude (Anthropic), OpenAI, and Google Gemini
- **Structured Output**: Summary (max 500 chars), Questions, and Action Items
- **Smart Caching**: Avoid redundant API calls for same emails
- **Rate Limiting**: Control costs with hourly/daily limits
- **Priority Filtering**: Only summarize high-priority emails (optional)
- **Async Processing**: Non-blocking summary generation
- **Retry Logic**: Handle transient API failures
- **Database Storage**: Persistent caching of summaries

## Architecture

```
┌──────────────┐
│ Email Alert  │
└──────┬───────┘
       │
       v
┌──────────────────────┐
│ AI Summary Service   │  ← Configuration (ai-config.yaml)
│  - Check Cache       │
│  - Rate Limit Check  │
│  - Priority Filter   │
└──────┬───────────────┘
       │
       v
┌──────────────────────┐
│ Provider Interface   │
│  ┌────────────────┐  │
│  │ Claude         │  │
│  │ OpenAI         │  │
│  │ Gemini         │  │
│  └────────────────┘  │
└──────┬───────────────┘
       │
       v
┌──────────────────────┐
│ Storage Layer        │
│  - Cache Summaries   │
│  - Track Tokens      │
└──────────────────────┘
       │
       v
┌──────────────────────┐
│ Notification System  │
│  - Desktop Toast     │
│  - Mobile Push       │
│  - Tray UI           │
└──────────────────────┘
```

## Data Flow

1. **Email Received** → Matched by filter
2. **Check Priority** → Skip if priority_only=true and email is not urgent
3. **Check Cache** → Return cached summary if exists
4. **Check Rate Limit** → Verify within hourly/daily limits
5. **Call AI Provider** → Generate summary with retries
6. **Store in Database** → Cache for future use
7. **Update UI** → Show in notifications and tray

## File Structure

```
internal/ai/
├── config.go          # YAML configuration parser
├── types.go           # Data structures
├── provider.go        # Provider interface + implementations
└── service.go         # Main service with caching/rate limiting

ai-config.yaml         # User configuration file

internal/storage/
├── migrations.go      # Migration_002_AddAISummariesTable
└── db.go              # AI summary storage functions

cmd/
└── start.go           # Integration point (--ai-summary flag)
```

## Database Schema

```sql
CREATE TABLE ai_summaries (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    message_id TEXT NOT NULL UNIQUE,    -- Gmail message ID
    summary TEXT NOT NULL,                -- Brief overview (max 500 chars)
    questions TEXT,                       -- JSON array of questions
    action_items TEXT,                    -- JSON array of action items
    provider TEXT NOT NULL,               -- "claude", "openai", "gemini"
    model TEXT NOT NULL,                  -- Model name used
    generated_at INTEGER NOT NULL,        -- Unix timestamp
    tokens_used INTEGER DEFAULT 0         -- API tokens consumed
);

CREATE INDEX idx_summary_message_id ON ai_summaries(message_id);
CREATE INDEX idx_summary_generated_at ON ai_summaries(generated_at DESC);
```

## Configuration

### api-config.yaml

```yaml
ai_summary:
  enabled: false
  provider: "claude"  # or "openai" or "gemini"

  api:
    claude:
      api_key: ""  # or env: ANTHROPIC_API_KEY
      model: "claude-3-5-haiku-20241022"
      max_tokens: 1024
      temperature: 0.3
    # ... similar for openai and gemini

  behavior:
    max_summary_length: 500
    priority_only: false
    enable_cache: true
    timeout_seconds: 15
    retry_attempts: 2
    include_in_notifications: true
    show_ai_icon: true

  rate_limit:
    max_per_hour: 60
    max_per_day: 500
```

### Environment Variables

```bash
# Set API keys via environment variables (recommended)
export ANTHROPIC_API_KEY="your-key-here"
export OPENAI_API_KEY="your-key-here"
export GEMINI_API_KEY="your-key-here"
```

## Usage

### Starting with AI Summaries

```bash
# Enable AI summaries (uses ai-config.yaml settings)
email-sentinel start --ai-summary --tray

# AI summaries will be generated for matching emails
# and displayed in notifications and tray
```

### Cost Optimization

1. **Use Smaller Models**
   - Claude Haiku: $0.25/$1.25 per 1M tokens
   - GPT-4o-mini: $0.15/$0.60 per 1M tokens
   - Gemini Flash: Free tier available

2. **Enable Priority-Only Mode**
   ```yaml
   behavior:
     priority_only: true  # Only summarize urgent emails
   ```

3. **Set Rate Limits**
   ```yaml
   rate_limit:
     max_per_hour: 30
     max_per_day: 200
   ```

4. **Enable Caching**
   ```yaml
   behavior:
     enable_cache: true  # Avoid re-summarizing same emails
   ```

## Provider-Specific Notes

### Claude (Anthropic)
- **Best for**: Structured output, following instructions
- **Model**: `claude-3-5-haiku-20241022` (fast, cost-effective)
- **API**: https://api.anthropic.com/v1/messages
- **Pricing**: $0.25 input / $1.25 output per 1M tokens

### OpenAI
- **Best for**: JSON output, reliability
- **Model**: `gpt-4o-mini` (affordable, good quality)
- **API**: https://api.openai.com/v1/chat/completions
- **Pricing**: $0.15 input / $0.60 output per 1M tokens

### Google Gemini
- **Best for**: Cost efficiency, generous free tier
- **Model**: `gemini-1.5-flash` (fast, free tier)
- **API**: https://generativelanguage.googleapis.com/v1beta/models
- **Pricing**: Free tier available

## Integration Points

### 1. Email Processing Pipeline

```go
// cmd/start.go - checkEmails()
if aiSummaryEnabled && aiService != nil {
    go func(alert storage.Alert) {
        summary, err := aiService.GenerateSummary(
            alert.MessageID,
            alert.Sender,
            alert.Subject,
            "", // body (not available in snippet API)
            alert.Snippet,
            alert.Priority,
        )
        if err != nil {
            log.Printf("⚠️  AI summary failed: %v", err)
            return
        }
        // Update alert with summary
        alert.AISummary = summary
        // Refresh notifications/tray
    }(alert)
}
```

### 2. Notification System

```go
// internal/notify/ - Enhanced notifications
if alert.AISummary != nil {
    body := fmt.Sprintf("%s\n\n📝 %s",
        alert.Subject,
        alert.AISummary.Summary)

    if len(alert.AISummary.Questions) > 0 {
        body += "\n\n❓ Questions:\n"
        for _, q := range alert.AISummary.Questions {
            body += fmt.Sprintf("  • %s\n", q)
        }
    }

    if len(alert.AISummary.ActionItems) > 0 {
        body += "\n\n✅ Action Items:\n"
        for _, item := range alert.AISummary.ActionItems {
            body += fmt.Sprintf("  • %s\n", item)
        }
    }
}
```

### 3. System Tray

```go
// internal/tray/tray.go - Enhanced tooltips
if alert.AISummary != nil {
    tooltip += fmt.Sprintf("\n\n🤖 AI Summary:\n%s", alert.AISummary.Summary)

    if len(alert.AISummary.Questions) > 0 {
        tooltip += fmt.Sprintf("\n\n❓ %d question(s)", len(alert.AISummary.Questions))
    }

    if len(alert.AISummary.ActionItems) > 0 {
        tooltip += fmt.Sprintf("\n✅ %d action item(s)", len(alert.AISummary.ActionItems))
    }
}
```

## Error Handling

1. **API Failures**: Retry with exponential backoff
2. **Rate Limits**: Skip summarization, log warning
3. **Invalid Responses**: Log error, continue without summary
4. **Network Timeouts**: Respect timeout_seconds configuration
5. **Cache Failures**: Log warning, proceed with API call

## Security Considerations

1. **API Keys**: Store in environment variables, not in code
2. **Email Content**: Sent to third-party AI providers
3. **Data Retention**: Summaries stored in local SQLite database
4. **Rate Limiting**: Prevents excessive API costs
5. **Timeout Protection**: Prevents hanging requests

## Performance

- **Async Processing**: Summaries generated in background goroutines
- **Caching**: Avoid duplicate API calls for same emails
- **Rate Limiting**: Prevent API overload
- **Timeouts**: Fast failure for slow API responses
- **Non-Blocking**: Email processing continues if summarization fails

## Testing

```bash
# Test AI configuration
email-sentinel start --ai-summary

# Monitor logs for AI activity
# Look for: 🤖 Generating AI summary...
#           ✅ AI summary generated (X tokens)
```

## Future Enhancements

1. **Sentiment Analysis**: Detect email tone (urgent, happy, angry)
2. **Category Detection**: Auto-categorize emails
3. **Smart Replies**: Generate suggested responses
4. **Meeting Detection**: Extract meeting details
5. **Contact Extraction**: Pull out names, emails, phone numbers
6. **Multi-Language**: Translate summaries
7. **Custom Prompts**: Per-filter custom prompts
8. **Batch Processing**: Summarize multiple emails at once
9. **Cost Tracking**: Dashboard for API usage/costs
10. **A/B Testing**: Compare provider quality/cost

## Troubleshooting

### "AI summary not enabled"
- Check `enabled: true` in ai-config.yaml
- Use `--ai-summary` flag

### "API key not set"
- Set environment variable (ANTHROPIC_API_KEY, etc.)
- Or add to ai-config.yaml (not recommended)

### "Rate limit exceeded"
- Increase max_per_hour/max_per_day
- Or wait for rate limit window to reset

### "API error 401"
- Invalid API key
- Check environment variable

### "Timeout"
- Increase timeout_seconds
- Check network connection

### "Failed to parse response"
- Provider returned non-JSON
- Check model supports JSON output
- Review prompt template

## Support

For issues or questions:
- GitHub: https://github.com/datateamsix/email-sentinel
- Documentation: See README.md
