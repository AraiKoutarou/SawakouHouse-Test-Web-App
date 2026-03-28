#!/bin/bash

# 環境変数の確認
if [ -z "$GEMINI_API_KEY" ] || [ -z "$GITHUB_TOKEN" ] || [ -z "$PR_NUMBER" ]; then
  echo "必要な環境変数が設定されていません"
  exit 1
fi

# PR差分を取得
DIFF_URL="https://api.github.com/repos/${REPO_OWNER}/${REPO_NAME}/pulls/${PR_NUMBER}/files"
FILES=$(curl -s -H "Authorization: token ${GITHUB_TOKEN}" "$DIFF_URL")

# 差分を整形
DIFF_TEXT="以下のプルリクエストのコードレビューを日本語で実施してください。\\n\\n重要: 必ず日本語でレビューを記述してください。\\n\\n以下の点に注目してください：\\n1. バグや潜在的なエラーの指摘\\n2. コードの可読性と保守性の向上案\\n3. セキュリティ上の懸念点\\n4. 良い実装については褒める\\n\\n変更内容：\\n${FILES}"

# Gemini APIに送信
RESPONSE=$(curl -s "https://generativelanguage.googleapis.com/v1beta/models/gemini-1.5-flash:generateContent?key=${GEMINI_API_KEY}" \
  -H 'Content-Type: application/json' \
  -X POST \
  -d "{\"contents\":[{\"parts\":[{\"text\":${DIFF_TEXT}}]}]}")

# レスポンスからテキストを抽出
REVIEW_TEXT=$(echo "$RESPONSE" | jq -r '.candidates[0].content.parts[0].text // "レビューの生成に失敗しました"')

# GitHubにコメント投稿
COMMENT_BODY=$(jq -n --arg text "## 🤖 Gemini Code Review\n\n$REVIEW_TEXT" '{body: $text}')
curl -s -X POST \
  -H "Authorization: token ${GITHUB_TOKEN}" \
  -H "Content-Type: application/json" \
  -d "$COMMENT_BODY" \
  "https://api.github.com/repos/${REPO_OWNER}/${REPO_NAME}/issues/${PR_NUMBER}/comments"

echo "レビューコメントを投稿しました"
